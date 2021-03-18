package cnbp2llb

import (
	"context"

	"github.com/EricHripko/cnbp/pkg/cib"

	"github.com/moby/buildkit/client/llb"
)

// Detect if given builder supports the project. Finds an ordered group of
// buildpacks to use during the build phase.
func Detect(ctx context.Context, build cib.Service, env llb.State) llb.State {
	// Execute detector
	// See https://github.com/buildpacks/spec/blob/main/platform.md#detector
	return env.Run(
		llb.Args([]string{"/cnb/lifecycle/detector"}),
		llb.WithCustomName("Detection"),
	).Root()
}
