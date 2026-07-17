BINARY := kesmarki
BINDIR := bin

GOOS   := linux
GOARCH := arm
GOARM  := 5

.PHONY: build clean

build:
	GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) go build -o $(BINDIR)/$(BINARY) .

clean:
	rm -f $(BINDIR)/$(BINARY)
