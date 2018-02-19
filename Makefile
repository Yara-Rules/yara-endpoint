
YARA_VERSION := v3.7.1


GO := $(shell type go)
XGO := $(shell type xgo)
DOCKER := $(shell type docker)
TARGET_DIR := "build"

# "yara-rules/x-go-yara"
DOCKER_IMAGE := "x-go-yara"
# DOCKER_BUILD := $(shell docker images | grep $(DOCKER_IMAGE))
DEPS := "https://github.com/VirusTotal/yara/archive/$(YARA_VERSION).tar.gz"


.PHONY: all clean linux linux-x86 linux-x64 darwin darwin-x86 darwin-x64 windows windows-x86 windows-x64

ifeq ($(GO), "")
	@echo "You must install Go first. Please visit https://golang.org/ and follow the instructions"
	exit 1
endif
ifeq ($(DOCKER), "")
	@echo "You must install Docker first. Please visit https://www.docker.com/ and follow the instructions"
	exit 1
endif
# ifeq ($(DOCKER_BUILD), "")
# 	@echo "Building a docker image to compile Yara-Endpoint"
# 	$(shell docker build "https://github.com/Xumeiquer/xgo.git#:docker/go-latest" -t $(DOCKER_IMAGE))
# endif

all:
	make -C client
	make -C server

linux:
	make -C client linux
	make -C server linux

linux-x86:
	make -C client linux-x86
	make -C server linux-x86

linux-x64:
	make -C client linux-x64
	make -C server linux-x64

darwin:
	make -C client darwin
	make -C server darwin

darwin-x86:
	make -C client darwin-x86
	make -C server darwin-x86

darwin-x64:
	make -C client darwin-x64
	make -C server darwin-x64

windows:
	make -C client windows
	make -C server windows

windows-x86:
	make -C client windows-x86
	make -C server windows-x86

windows-x64:
	make -C client windows-x64
	make -C server windows-x64

clean:
	make -C client clean
	make -C server clean
