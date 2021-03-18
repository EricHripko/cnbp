package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/EricHripko/cnbp/pkg/cib"
	"github.com/EricHripko/cnbp/pkg/cnbp2llb"
	"github.com/containerd/containerd/platforms"

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
				builder := strings.TrimSpace(string(dtMetadata))
				if strings.Contains(builder, "\n") {
					// Strip BuildKit syntax comment
					lines := strings.Split(builder, "\n")
					builder = lines[len(lines)-1]
				}

				// Prepare build environment
				env, err := cnbp2llb.BuildEnvironment(ctx, svc, tp, builder)
				if err != nil {
					return errors.Wrap(err, "cannot prepare environment")
				}

				// Detect and build
				detected := cnbp2llb.Detect(ctx, svc, env)
				built := cnbp2llb.Build(ctx, svc, env, detected)

				// Export
				ref, img, err := cnbp2llb.Export(ctx, svc, built)
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
