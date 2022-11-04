# Microservice Template Repo

This repo comes pre-populated with the base settings for a new golang microservice:

* Provides a Buildkite pipeline with build, lint, tests, Trivy, and Cosign built-in
* Base golang http server application with audit + metrics built-in
* Working Docker compose setup and Makefile
* golangci lint config

## Usage

Search this repo for `TODO` to find the places where code would need to be updated which should be limited do:

```bash
.buildkite/pipeline.yml

cmd/serve.go
cmd/root.go
cmd/errors.go

srv/doc.go
srv/errors.go
```

## Development

This repo includes a `docker-compose.yml` and a `Makefile` to make getting started easy.

`make docker-up` will start a basic https server with an audit container
