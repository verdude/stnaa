PROJECT := stnaa
BUILDROOT := build
EXE := $(PROJECT)
EXEPATH := $(BUILDROOT)/$(EXE)
CONFIG_PREFIX := .
CONFIG := config.toml

$(EXEPATH): $(wildcard *.go)
	go build -o $(EXEPATH) -ldflags='-X main.config=$(CONFIG_PREFIX)/$(CONFIG)'

.PHONY: clean
clean:
	rm -rf $(BUILDROOT)
