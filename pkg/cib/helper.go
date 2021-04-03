package cib

import (
	"context"
	"errors"
	"strconv"

	"github.com/moby/buildkit/client/llb"
)

// FetchUID returns the UID of the user for the build.
func FetchUID(ctx context.Context, state llb.State) (int, error) {
	env, err := state.Env(ctx)
	if err != nil {
		return 0, err
	}

	for _, kv := range env {
		k, v := parseKeyValue(kv)
		if k == "CNB_USER_ID" {
			return strconv.Atoi(v)
		}
	}
	return 0, errors.New("user not found")
}

// FetchGID returns the GID of the user for the build.
func FetchGID(ctx context.Context, state llb.State) (int, error) {
	env, err := state.Env(ctx)
	if err != nil {
		return 0, err
	}

	for _, kv := range env {
		k, v := parseKeyValue(kv)
		if k == "CNB_GROUP_ID" {
			return strconv.Atoi(v)
		}
	}
	return 0, errors.New("group not found")
}
