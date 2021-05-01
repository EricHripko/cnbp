package cnbp2llb

import (
	"context"

	"github.com/EricHripko/buildkit-fdk/pkg/cib"
	"github.com/moby/buildkit/client/llb"
)

// Build the project with the provided parameters. Transforms application
// source code into runnable artifacts that can be packaged into a container.
func Build(ctx context.Context, build cib.Service, restored llb.State) llb.State {
	// Execute builder
	// See https://github.com/buildpacks/spec/blob/main/platform.md#builder
	return restored.Run(
		llb.Args([]string{"/cnb/lifecycle/builder"}),
		llb.WithCustomName("Build"),
	).Root()
}
