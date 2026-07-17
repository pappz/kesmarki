BINARY := kesmarki

GOOS   := linux
GOARCH := arm
GOARM  := 5

.PHONY: build clean

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -o $(BINARY) .

clean:
	rm -f $(BINARY)
