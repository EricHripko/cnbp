# Overview

This repository houses a barebones BuildKit frontend for
[Cloud Native Buildpacks](https://buildpacks.io/). Note that it's Linux only
since BuildKit [doesn't support Windows yet](https://github.com/moby/buildkit/issues/616).
Frontend aims to implement [Platform API 0.5](https://github.com/buildpacks/spec/),
but is not there yet.

## How to build it

- Clone the repository
- Build the frontend image: `docker build -t erichripko/cnbp .`

## How to try it

- Make sure that BuildKit is enabled in your daemon
- Clone CNBP [samples repository](https://github.com/buildpacks/samples)
- Add the following lines to the `Dockerfile`. As you can see, the only
  content apart from the syntax stanza should be the name of the builder
  to be used. This corresponds to `<image-name>` of
  `pack build <image-name> [flags]` invocation.

```dockerfile
# syntax = erichripko/cnbp
paketobuildpacks/builder:full
```

- Build it

```shell
$ docker build -t sample .
[+] Building 80.8s (28/28) FINISHED
 => [internal] load build definition from Dockerfile                                                               0.1s
 => => transferring dockerfile: 173B                                                                               0.0s
 => [internal] load .dockerignore                                                                                  0.0s
 => => transferring context: 2B                                                                                    0.0s
 => resolve image config for docker.io/erichripko/cnbp:latest                                                      0.0s
 => CACHED docker-image://docker.io/erichripko/cnbp:latest                                                         0.0s
 => [internal] load build definition from Dockerfile                                                               0.0s
 => => transferring dockerfile: 173B                                                                               0.0s
 => load metadata for docker.io/paketobuildpacks/builder:full                                                      0.0s
 => [internal] load context                                                                                        0.1s
 => => transferring context: 87.94kB                                                                               0.0s
 => Builder is paketobuildpacks/builder:full                                                                       0.6s
 => Load sources                                                                                                   0.2s
 => Detection                                                                                                      1.2s
 => Load group definition                                                                                          0.2s
 => Load plan definition                                                                                           0.2s
 => Build                                                                                                         41.2s
 => load metadata for docker.io/paketobuildpacks/run:full-cnb                                                      1.4s
 => Run image is index.docker.io/paketobuildpacks/run:full-cnb                                                    30.0s
 => Exporting launcher                                                                                             1.0s
 => Exporting buildpack layer /layers/paketo-buildpacks_ca-certificates/helper                                     0.1s
 => Exporting buildpack layer /layers/paketo-buildpacks_bellsoft-liberica/helper                                   0.2s
 => Exporting buildpack layer /layers/paketo-buildpacks_bellsoft-liberica/java-security-properties                 0.1s
 => Exporting buildpack layer /layers/paketo-buildpacks_bellsoft-liberica/jre                                      0.5s
 => Exporting buildpack layer /layers/paketo-buildpacks_bellsoft-liberica/jvmkill                                  0.1s
 => Exporting buildpack layer /layers/paketo-buildpacks_executable-jar/classpath                                   0.1s
 => Exporting buildpack layer /layers/paketo-buildpacks_spring-boot/helper                                         0.1s
 => Exporting buildpack layer /layers/paketo-buildpacks_spring-boot/spring-cloud-bindings                          0.1s
 => Exporting buildpack layer /layers/paketo-buildpacks_spring-boot/web-application-type                           0.1s
 => Exporting app layer                                                                                            0.2s
 => Exporting build metadata                                                                                       0.2s
 => exporting to image                                                                                             0.9s
 => => exporting layers                                                                                            0.8s
 => => writing image sha256:321aa7926a9e38422f53a1ebcfff5e32af9dcff4788e825199140eeb5d5e38a6                       0.0s
 => => naming to docker.io/library/sample                                                                          0.0s
  ```

- Run the application: `docker run --rm -d -p 8080:8080 sample /cnb/lifecycle/launcher`
- Observe that it works as expected

```shell
$ curl http://localhost:8080
...
        <h1>Hello, Buildpacker!</h1>
    </div>
</body>
</html>
```
