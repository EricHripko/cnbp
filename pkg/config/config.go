package config

import (
	"github.com/pelletier/go-toml"
	"github.com/pelletier/go-toml/query"
	"github.com/pkg/errors"
)

type Config interface {
	Builder() string
	PreviousImage() string
}

type config struct {
	builder   string
	prevImage string
}

func (c *config) Builder() string {
	return c.builder
}

func (c *config) PreviousImage() string {
	return c.prevImage
}

type buildKitOpts struct {
	PreviousImage string `toml:"previous-image"`
}

func FromProjectTOML(data string) (Config, error) {
	tree, err := toml.Load(data)
	if err != nil {
		return &config{}, err
	}

	builder, err := extractBuilder(tree)
	if err != nil {
		return &config{}, err
	}

	buildKitOpts, err := extractBuildKitOptions(tree)
	if err != nil {
		return &config{}, err
	}

	return &config{
		builder:   builder,
		prevImage: buildKitOpts.PreviousImage,
	}, nil
}

func extractBuilder(tree *toml.Tree) (string, error) {
	builder, err := extractString(tree, "$.io.buildpacks.build.builder")
	if err != nil {
		return "", err
	}

	if builder == "" {
		return "", errors.New("no builder provided")
	}

	return builder, err
}

func extractBuildKitOptions(tree *toml.Tree) (*buildKitOpts, error) {
	opts := &buildKitOpts{}
	results, err := query.CompileAndExecute("$.io.buildpacks.ext.buildkit", tree)
	if err != nil {
		return opts, errors.Wrap(err, "failed finding '$.io.buildpacks.ext.buildkit'")
	}

	if len(results.Values()) == 0 {
		return opts, nil
	}

	optsTree := results.Values()[0].(*toml.Tree)

	err = optsTree.Unmarshal(opts)
	if err != nil {
		return opts, errors.Wrap(err, "failed to parse '$.io.buildpacks.ext.buildkit' value")
	}

	return opts, nil
}

func extractString(tree *toml.Tree, path string) (string, error) {
	results, err := query.CompileAndExecute(path, tree)
	if err != nil {
		return "", errors.Wrapf(err, "failed finding '%s'", path)
	}

	if len(results.Values()) == 0 {
		return "", nil
	}

	return results.Values()[0].(string), nil
}
