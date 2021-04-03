// Package cnbp2llb provides a BuildKit-based platform for running Cloud
// Native Buildpacks.
// See https://github.com/buildpacks/spec/blob/main/platform.md
package cnbp2llb

// PlatformAPI specifies the version of API that this platform supports.
// See https://github.com/buildpacks/spec/blob/main/platform.md#platform-api-compatibility
const PlatformAPI = "0.5"

// AppDir specifies the default value for the path to application directory.
// See https://github.com/buildpacks/spec/blob/main/platform.md#inputs
const AppDir = "/workspace"

// PlatformDir specifies the default value for the path to platform directory.
// See https://github.com/buildpacks/spec/blob/main/platform.md#inputs
const PlatformDir = "/platform"

// LayersDir specifies the default value for the path to layers directory.
// See https://github.com/buildpacks/spec/blob/main/platform.md#inputs
const LayersDir = "/layers"

// CacheDir specifies the default path to the cache directory.
// This value isn't defined by the spec, but is commonly used in platform
// implementations.
const CacheDir = "/cache"

// GroupPath specifies the default value for the output group definition.
// See https://github.com/buildpacks/spec/blob/main/platform.md#inputs
const GroupPath = "group.toml"

// PlanPath specifies the default value for the output resolved build plan.
// See https://github.com/buildpacks/spec/blob/main/platform.md#inputs
const PlanPath = "plan.toml"

// MetadataPath specifies the path for the build metadata.
// See https://github.com/buildpacks/spec/blob/main/platform.md#outputs-4
const MetadataPath = "config/metadata.toml"

// StackPath specifies the default value for the path to stack file.
// See https://github.com/buildpacks/spec/blob/main/platform.md#inputs-4
const StackPath = "/cnb/stack.toml"

// LauncherPath specifies the path to launcher executable.
// See https://github.com/buildpacks/spec/blob/main/platform.md#inputs-4
const LauncherPath = "/cnb/lifecycle/launcher"
