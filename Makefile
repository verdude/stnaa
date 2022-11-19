PROJECT := stnaa
BUILDROOT := build
EXE := $(PROJECT)
EXEPATH := $(BUILDROOT)/$(EXE)
CONFIG_PREFIX := .
CONFIG := config.toml

.PHONY: all
all: $(wildcard *.go)
	go build -o $(EXEPATH) -ldflags='-X main.config=$(CONFIG_PREFIX)/$(CONFIG)'

.PHONY: clean
clean:
	rm -rf $(BUILDROOT)
