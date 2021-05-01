package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/EricHripko/buildkit-fdk/pkg/cib"
	"github.com/EricHripko/cnbp/pkg/cnbp2llb"
	"github.com/EricHripko/cnbp/pkg/config"
	"github.com/containerd/containerd/platforms"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend/gateway/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

const (
	keyMultiPlatform = "multi-platform"
)

// Build the image with this frontend.
func Build(ctx context.Context, c client.Client) (*client.Result, error) {
	return BuildWithService(ctx, c, cib.NewService(ctx, c))
}

// BuildWithService uses the provided container image build service to
// perform the build.
//nolint:gocyclo // Frontends are complex
func BuildWithService(ctx context.Context, c client.Client, svc cib.Service) (*client.Result, error) {
	opts := svc.GetOpts()

	// Identify target platforms
	targetPlatforms, err := svc.GetTargetPlatforms()
	if err != nil {
		return nil, err
	}
	exportMap := len(targetPlatforms) > 1
	if v := opts[keyMultiPlatform]; v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, errors.Errorf("invalid boolean value %s", v)
		}
		if !b && exportMap {
			return nil, errors.Errorf("returning multiple target plaforms is not allowed")
		}
		exportMap = b
	}
	expPlatforms := &exptypes.Platforms{
		Platforms: make([]exptypes.Platform, len(targetPlatforms)),
	}

	// Build an image for each platform
	res := client.NewResult()
	eg, ctx := errgroup.WithContext(ctx)
	for i, tp := range targetPlatforms {
		func(i int, tp *specs.Platform) {
			eg.Go(func() error {
				// Fetch the builder name
				dtMetadata, err := svc.GetMetadata()
				if err != nil {
					return err
				}

				buildCfg, err := config.FromProjectTOML(string(dtMetadata))
				if err != nil {
					return err
				}

				// Prepare build environment
				env, err := cnbp2llb.BuildEnvironment(ctx, svc, tp, buildCfg.Builder())
				if err != nil {
					return errors.Wrap(err, "cannot prepare environment")
				}

				// Prepare cache
				uid, err := FetchUID(ctx, env)
				if err != nil {
					return err
				}
				gid, err := FetchGID(ctx, env)
				if err != nil {
					return err
				}
				cache := llb.AddMount(
					cnbp2llb.CacheDir,
					llb.Scratch().File(
						llb.Mkdir(
							cnbp2llb.CacheDir,
							os.FileMode(0o755),
							llb.WithUIDGID(uid, gid),
						),
						llb.WithCustomName("Setting cache mount permissions"),
					),
					llb.SourcePath(cnbp2llb.CacheDir),
					llb.AsPersistentCacheDir("buildpacks-cache", llb.CacheMountPrivate),
				)

				// Detect and build
				detected := cnbp2llb.Detect(ctx, svc, env)
				analyzed := cnbp2llb.Analyze(ctx, svc, env, detected, cache)
				restored := cnbp2llb.Restore(ctx, svc, analyzed, cache)
				built := cnbp2llb.Build(ctx, svc, restored)

				// Export
				ref, img, err := cnbp2llb.Export(ctx, svc, built, cache)
				if err != nil {
					return errors.Wrap(err, "cannot export")
				}
				config, err := json.Marshal(img)
				if err != nil {
					return errors.Wrap(err, "failed to marshal image config")
				}
				if !exportMap {
					res.AddMeta(exptypes.ExporterImageConfigKey, config)
					res.SetRef(ref)
				} else {
					p := platforms.DefaultSpec()
					if tp != nil {
						p = *tp
					}

					k := platforms.Format(p)
					res.AddMeta(fmt.Sprintf("%s/%s", exptypes.ExporterImageConfigKey, k), config)
					res.AddRef(k, ref)
					expPlatforms.Platforms[i] = exptypes.Platform{
						ID:       k,
						Platform: p,
					}
				}
				return nil
			})
		}(i, tp)
	}
	if err := eg.Wait(); err != nil {
		return nil, err
	}

	// Export image(s)
	if exportMap {
		dt, err := json.Marshal(expPlatforms)
		if err != nil {
			return nil, err
		}
		res.AddMeta(exptypes.ExporterPlatformsKey, dt)
	}
	return res, nil
}
