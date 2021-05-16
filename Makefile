.PHONY: build

build:
	sam build
openapi-gen:
	oapi-codegen -generate "types" -package openapi ./openapi.yaml > openapi/openapi.gen.go
