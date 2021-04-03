// Tiny kludge that reuses the functionality of default Buildpacks lifecycle
// to enable caching via BuildKit mounts.
package main

import (
	"os"
	"path"
	"strconv"

	"github.com/EricHripko/cnbp/pkg/cnbp2llb"
	"github.com/buildpacks/lifecycle"
	"github.com/buildpacks/lifecycle/api"
	"github.com/buildpacks/lifecycle/cache"
	"github.com/buildpacks/lifecycle/cmd"
	"github.com/buildpacks/lifecycle/layers"
)

func main() {
	// Read metadata and get user identity
	group, err := lifecycle.ReadGroup(path.Join(cnbp2llb.LayersDir, cnbp2llb.GroupPath))
	if err != nil {
		cmd.DefaultLogger.Fatalf("Failted to read Buildpacks group: %v\n", err)
	}
	uid, err := strconv.Atoi(os.Getenv("CNB_USER_ID"))
	if err != nil {
		cmd.DefaultLogger.Fatalf("Failted to parse user: %v\n", err)
	}
	gid, err := strconv.Atoi(os.Getenv("CNB_GROUP_ID"))
	if err != nil {
		cmd.DefaultLogger.Fatalf("Failted to parse group: %v\n", err)
	}

	// Create an exporter
	exporter := lifecycle.Exporter{
		Buildpacks: group.Group,
		LayerFactory: &layers.Factory{
			ArtifactsDir: cnbp2llb.LayersDir,
			UID:          uid,
			GID:          gid,
			Logger:       cmd.DefaultLogger,
		},
		Logger:      cmd.DefaultLogger,
		PlatformAPI: api.MustParse(os.Getenv("CNB_PLATFORM_API")),
	}

	// Attempt to export cache
	cacheStore, err := cache.NewVolumeCache(cnbp2llb.CacheDir)
	if err != nil {
		cmd.DefaultLogger.Fatalf("Failted to create cache store: %v\n", err)
	}
	if err := exporter.Cache(cnbp2llb.LayersDir, cacheStore); err != nil {
		cmd.DefaultLogger.Warnf("Failed to export cache: %v\n", err)
	}
}
