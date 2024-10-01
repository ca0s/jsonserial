.PHONY: all

all: jsonserial

jsonserial:
	go build -o bin/jsonserial ./cmd/jsonserial