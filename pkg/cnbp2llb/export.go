package cnbp2llb

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/EricHripko/buildkit-fdk/pkg/cib"
	cnbp "github.com/buildpacks/lifecycle"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/pkg/errors"
	fsutil "github.com/tonistiigi/fsutil/types"
)

// Export the produced layers into an OCI image. Unlike other high-level
// functions in this package, we have to manage the export manually (without
// lifecycle) to fit in with the BuildKit model of the world.
func Export(ctx context.Context, build cib.Service, built llb.State, cache llb.RunOption) (ref client.Reference, img *dockerfile2llb.Image, err error) {
	// Shall write the contents of all cached layers to the cache
	// Shall record the diffID and layer content metadata of all cached layers
	// in the cache
	// Relying on an external image since BuildKit doesn't provide a way to
	// inject the binary from the frontend.
	// See: https://github.com/moby/buildkit/issues/2063
	built = built.Run(
		llb.Args([]string{"/frontend/go/bin/cacher"}),
		llb.WithCustomName("Populating cache"),
		// Mount frontend for the cacher binary
		llb.AddMount(
			"/frontend",
			// TODO: can we make this dynamic
			llb.Image("erichripko/cnbp"),
		),
		cache,
	).Root()
	ref, err = build.Solve(ctx, built)
	if err != nil {
		return
	}

	// Read the stack and group
	var groups cnbp.BuildpackGroup
	err = readToml(ctx, ref, path.Join(LayersDir, GroupPath), &groups)
	if err != nil {
		return
	}

	// Find launch layers
	var launchLayers []string
	for _, group := range groups.Group {
		id := strings.ReplaceAll(group.ID, "/", "_")
		groupPath := path.Join(LayersDir, id)

		var files []*fsutil.Stat
		files, err = ref.ReadDir(ctx, client.ReadDirRequest{Path: groupPath})
		if err != nil {
			return
		}
		for _, file := range files {
			mode := os.FileMode(file.Mode)
			if !mode.IsDir() {
				continue
			}

			// Maybe found a layer, attempt to read its metadata
			var metadata cnbp.BuildpackLayerMetadata
			err = readToml(
				ctx,
				ref,
				path.Join(groupPath, path.Base(file.Path)+".toml"),
				&metadata,
			)
			if err == nil && metadata.Launch {
				// Found a launch layer
				launchLayers = append(
					launchLayers,
					path.Join(groupPath, file.Path),
				)
			}
		}
		if err != nil {
			return
		}
	}

	// Produce the end OCI image
	var stack cnbp.StackMetadata
	err = readToml(ctx, ref, StackPath, &stack)
	if err != nil {
		return
	}
	platform, err := built.GetPlatform(ctx)
	if err != nil {
		return
	}
	// Must be an extension of the <run-image>
	state, img, err := build.From(
		stack.RunImage.Image,
		platform,
		fmt.Sprintf("Run image is %s", stack.RunImage.Image),
	)
	if err != nil {
		return
	}
	// Must contain one or more launcher layers
	state = state.File(
		llb.Copy(
			built,
			LauncherPath,
			LauncherPath,
			&llb.CopyInfo{CreateDestPath: true},
		),
		llb.WithCustomName("Exporting launcher"),
	)
	// Must contain all buildpack-provided launch layers
	for _, layer := range launchLayers {
		state = state.File(
			llb.Copy(
				built,
				layer,
				layer,
				&llb.CopyInfo{CreateDestPath: true},
			),
			llb.WithCustomNamef("Exporting buildpack layer %s", layer),
		)
	}

	// Must contain one or more app layers
	state = state.File(
		llb.Copy(
			built,
			AppDir,
			AppDir,
			&llb.CopyInfo{CopyDirContentsOnly: true},
		),
		llb.WithCustomName("Exporting app layer"),
	)
	// Must contain a layer that includes metadata.toml
	metadata := path.Join(LayersDir, MetadataPath)
	state = state.File(
		llb.Copy(
			built,
			metadata,
			metadata,
			&llb.CopyInfo{CreateDestPath: true},
		),
		llb.WithCustomName("Exporting build metadata"),
	)

	ref, err = build.Solve(ctx, state)
	return
}

func readToml(ctx context.Context, ref client.Reference, path string, v interface{}) error {
	data, err := ref.ReadFile(ctx, client.ReadRequest{Filename: path})
	if err != nil {
		return errors.Wrapf(err, "failed to read %s", path)
	}
	_, err = toml.DecodeReader(bytes.NewReader(data), v)
	if err != nil {
		return errors.Wrapf(err, "failed to decode %s", path)
	}
	return nil
}
