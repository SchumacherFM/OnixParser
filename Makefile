GOFMT=gofmt -s

GOFILES=$(wildcard *.go **/*.go)

format:
	${GOFMT} -w ${GOFILES}
