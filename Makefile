all: build_native

build_native:
	go build -o kubeaudit .

test:
	go test -cover ./cmd/... .

clean:
	/bin/rm -v kubeaudit

.PHONY: build_native clean test
