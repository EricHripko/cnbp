// BuildKit frontend that implements CloudNative Buildpacks.
package main

import (
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/sirupsen/logrus"
)

func main() {
	if err := grpcclient.RunFromEnvironment(appcontext.Context(), Build); err != nil {
		logrus.Errorf("fatal error: %+v", err)
		panic(err)
	}
}
