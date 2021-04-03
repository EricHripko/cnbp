package cnbp2llb

import (
	"context"

	"github.com/EricHripko/cnbp/pkg/cib"

	"github.com/moby/buildkit/client/llb"
)

// Restore desired cached artefacts from the previous build.
func Restore(ctx context.Context, build cib.Service, analyzed llb.State, cache llb.RunOption) llb.State {
	// Execute restorer
	// See https://github.com/buildpacks/spec/blob/main/platform.md#restorer
	return analyzed.Run(
		llb.Args([]string{"/cnb/lifecycle/restorer"}),
		llb.WithCustomName("Restore"),
		cache,
	).Root()
}
