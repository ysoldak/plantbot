SIZE   ?= full
TARGET ?= nano-rp2040


.PHONY: clean build flash

# --- Common targets ---

VERSION := $(shell git describe --tags --always)
LD_FLAGS := -ldflags="-X 'main.Version=$(VERSION)'"

clean:
	@rm -rf build

build:
	@mkdir -p build
	tinygo build $(LD_FLAGS) -target=$(TARGET) -size=$(SIZE) -opt=z -print-allocs=main -o ./build/debug.elf ./src

flash:
	tinygo flash $(LD_FLAGS) -target=$(TARGET) -size=$(SIZE) -opt=z -print-allocs=main ./src
