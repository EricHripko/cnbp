package cnbp2llb

import (
	"context"
	"path"

	"github.com/EricHripko/cnbp/pkg/cib"

	"github.com/moby/buildkit/client/llb"
)

// Analyze if there're any artefacts to reuse from the previous build.
func Analyze(ctx context.Context, build cib.Service, env llb.State, detected llb.State, cache llb.RunOption) llb.State {
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

	// Execute analyzer
	// See https://github.com/buildpacks/spec/blob/main/platform.md#analyzer
	state = state.AddEnv("CNB_CACHE_DIR", CacheDir)
	return state.Run(
		llb.Args([]string{"/cnb/lifecycle/analyzer", "this-image-definitely-does-not-exist"}),
		llb.WithCustomName("Analyze"),
		cache,
	).Root()
}
