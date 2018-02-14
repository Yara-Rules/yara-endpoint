
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
	xgo -image $(DOCKER_IMAGE) -dest build -out yara-endpoint-s -deps $(DEPS) -targets="linux/*,darwin/*,windows/*" .

linux:
	xgo -image $(DOCKER_IMAGE) -dest build -out yara-endpoint-s -deps $(DEPS) -targets="linux/*" .

linux-x86:
	xgo -image $(DOCKER_IMAGE) -dest build -out yara-endpoint-s -deps $(DEPS) -targets="linux/386" .

linux-x64:
	xgo -image $(DOCKER_IMAGE) -dest build -out yara-endpoint-s -deps $(DEPS) -targets="linux/amd64" .

darwin:
	xgo -image $(DOCKER_IMAGE) -dest build -out yara-endpoint-s -deps $(DEPS) -targets="darwin/*" .

darwin-x86:
	xgo -image $(DOCKER_IMAGE) -dest build -out yara-endpoint-s -deps $(DEPS) -targets="darwin/386" .

darwin-x64:
	xgo -image $(DOCKER_IMAGE) -dest build -out yara-endpoint-s -deps $(DEPS) -targets="darwin/amd64" .

windows:
	xgo -image $(DOCKER_IMAGE) -dest build -out yara-endpoint-s -deps $(DEPS) -targets="windows/*" .

windows-x86:
	xgo -image $(DOCKER_IMAGE) -dest build -out yara-endpoint-s -deps $(DEPS) -targets="windows/386" .

windows-x64:
	xgo -image $(DOCKER_IMAGE) -dest build -out yara-endpoint-s -deps $(DEPS) -targets="windows/amd64" .

clean:
	rm -rf $(TARGET_DIR)
