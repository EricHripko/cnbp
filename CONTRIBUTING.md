# Contributing Guide

You're welcome to contribute in whatever way you can. To make sure that
contributing is a good experience for everyone, please submit an issue to
discuss the change ahead of implementing (if appropriate and possible).

## 📋 Prerequisites

- Go of version `1.15` or later (get it [here](https://golang.org/doc/install))
- Docker with BuildKit support enabled (see [docs](https://docs.docker.com/develop/develop-images/build_enhancements/))

## 💻 Develop

- Clone the repository: `git clone https://github.com/EricHripko/cnbp.git`
- Edit away in your favourite environment

## 🏗️ Build

- Check that the code builds: `go build ./...`
- Build the frontend: `docker build -t erichripko/cnbp .`

## ✅ Verify

- Verify functionality with tests: `go test ./...`
  - This will run unit tests against the components
  - This will run end-to-end tests against the frontend

## 🔬 Analyse

- Code is formatted with standard Go tooling: `go fmt ./...`
- Code is linted with [golangci-lint](https://golangci-lint.run/): `golangci-lint run`

## 📢 Publish

Once the PR is landed, [GitHub Actions](https://github.com/features/actions)
will automatically publish the release of the frontend to
[Docker Hub](https://hub.docker.com/).

## 😀 Release

The frontend is referenced directly from _Docker Hub_, so it's immediately
released for everyone once published. `docker pull erichripko/cnbp` might be
necessary to get the latest version of the frontend if an older version is
already stored on the daemon.
