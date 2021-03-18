package cnbp2llb

import (
	"context"
	"path"

	"github.com/EricHripko/cnbp/pkg/cib"
	"github.com/moby/buildkit/client/llb"
)

// Build the project with the provided parameters. Transforms application
// source code into runnable artifacts that can be packaged into a container.
func Build(ctx context.Context, build cib.Service, env llb.State, detected llb.State) llb.State {
	// Provide detected inputs
	groupPath := path.Join(LayersDir, GroupPath)
	state := env.File(
		llb.Copy(
			detected,
			groupPath,
			groupPath,
		),
		llb.WithCustomName("Load group definition"),
	)
	planPath := path.Join(LayersDir, PlanPath)
	state = state.File(
		llb.Copy(
			detected,
			planPath,
			planPath,
		),
		llb.WithCustomName("Load plan definition"),
	)

	// Execute builder
	// See https://github.com/buildpacks/spec/blob/main/platform.md#builder
	return state.Run(
		llb.Args([]string{"/cnb/lifecycle/builder"}),
		llb.WithCustomName("Build"),
	).Root()
}
