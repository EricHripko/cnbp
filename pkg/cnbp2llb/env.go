package cnbp2llb

import (
	"context"
	"fmt"
	"os"
	"path"

	"github.com/EricHripko/cnbp/pkg/cib"

	"github.com/moby/buildkit/client/llb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// BuildEnvironment sets up the build environment by loading the builder
// image, source code and setting appropriate environment variables.
func BuildEnvironment(ctx context.Context, build cib.Service, platform *specs.Platform, builder string) (state llb.State, err error) {
	// Load the builder image
	state, img, err := build.From(
		builder,
		platform,
		fmt.Sprintf("Builder is %s", builder),
	)
	if err != nil {
		err = errors.Wrap(err, "failed to load builder")
		return
	}
	// Copy source code
	src, err := build.SrcState()
	if err != nil {
		err = errors.Wrap(err, "failed to load sources")
		return
	}
	state = state.File(
		llb.Copy(
			src,
			"/",
			AppDir,
			&llb.CopyInfo{CopyDirContentsOnly: true},
			llb.WithUser(img.Config.User),
		),
		llb.WithCustomName("Load sources"),
	)
	// Setup environment
	state = state.AddEnv("CNB_PLATFORM_API", PlatformAPI)
	// Inject user-provided variables
	// See https://github.com/buildpacks/spec/blob/main/platform.md#user-provided-variables
	for name, value := range build.GetBuildArgs() {
		path := path.Join(PlatformDir, "env", name)
		state = state.File(
			llb.Mkfile(
				path,
				os.FileMode(0o644),
				[]byte(value),
			),
			llb.WithCustomNamef("Set %s=%s", name, value),
		)
	}
	return
}
