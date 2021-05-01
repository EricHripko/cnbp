HOST_PORT?=8080
TEST_IMAGE_PREFIX?=sample

.PHONY: build
build:
	docker build --progress=plain -t erichripko/cnbp .

.PHONY: test
test:
	go test -v ./...

.PHONY: e2e-sample-ruby-bundler
e2e-sample-ruby-bundler: e2e-sample-ruby-bundler-build e2e-sample-ruby-bundler-run

.PHONY: e2e-sample-ruby-bundler-build
e2e-sample-ruby-bundler-build: export TEST_IMAGE_NAME:=$(TEST_IMAGE_PREFIX)-ruby-app
e2e-sample-ruby-bundler-build:
	@echo "> Building $(TEST_IMAGE_NAME)..."
	docker build --progress=plain -t $(TEST_IMAGE_NAME) -f samples/ruby-bundler/project.toml samples/ruby-bundler/

.PHONY: e2e-sample-ruby-bundler-run
e2e-sample-ruby-bundler-run: export TEST_IMAGE_NAME:=$(TEST_IMAGE_PREFIX)-ruby-app
e2e-sample-ruby-bundler-run:
	@echo "> Running $(TEST_IMAGE_NAME) on PORT $(HOST_PORT)..."
	@ID=`docker run -d --rm -it -p 8080:8080 $(TEST_IMAGE_NAME) /cnb/lifecycle/launcher`; \
	until CONTENTS=`curl -s http://localhost:8080`; do sleep 2; done; echo $$CONTENTS; \
	docker stop $$ID > /dev/null

.PHONY: clear-build-cache
clear-build-cache:
	docker builder prune -f
