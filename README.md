# Run fakeCAS with OSF

Please follow [README-docker-compose.md](https://github.com/CenterForOpenScience/osf.io/blob/develop/README-docker-compose.md) to run fakeCAS with OSF.

## Change the Image

By default, OSF uses the `master` image of fakeCAS, as shown below in [docker-compose.yml](https://github.com/CenterForOpenScience/osf.io/blob/develop/docker-compose.yml).

```yml
##################################
# Central Authentication Service #
##################################

fakecas:
  image: quay.io/centerforopenscience/fakecas:master
  command: fakecas -host=0.0.0.0:8080 -osfhost=localhost:5000 -dbaddress=postgres://postgres@postgres:5432/osf?sslmode=disable
  restart: unless-stopped
  ports:
    - 8080:8080
  depends_on:
    - postgres
  stdin_open: true
```

If you need the `develop` one, use `quay.io/centerforopenscience/fakecas:develop` instead. Run `docker-compose pull fakecas` to pull the latest image before starting fakeCAS.

## Pre-docker-compose

Starting [19.0.0](https://github.com/CenterForOpenScience/fakecas/milestone/1), fakeCAS no longer provides downloadable binary executables. Here is the last version [0.11.1](https://github.com/CenterForOpenScience/fakecas/releases/tag/0.11.1) that provides such a binary. However, the binary does not work with OSF starting [19.23.0](https://github.com/CenterForOpenScience/osf.io/releases/tag/19.23.0).

# Develop fakeCAS

Please take a look at the [Dockerfile](https://github.com/cslzchen/fakecas/blob/develop/Dockerfile) for how to develop fakeCAS locally. On macOS, use [`brew`](https://github.com/Homebrew/brew) to install [`go`](https://github.com/golang/go).
After installing go, pull dependencies using `go install` and run fakeCAS locally unsing `go run ./fakeCas.go`
