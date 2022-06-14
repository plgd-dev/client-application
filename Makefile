SHELL = /bin/bash
SERVICE_NAME = client-application
VERSION_TAG ?= $(shell git rev-parse --short=7 --verify HEAD)
TMP_PATH = $(shell pwd)/.tmp
GOPATH ?= $(shell go env GOPATH)
WORKING_DIRECTORY := $(shell pwd)
SIMULATOR_NAME_SUFFIX ?= $(shell hostname)
BUILD_PATH ?=$(TMP_PATH)/build
WWW_PATH=$(TMP_PATH)/www
CLIENT_APPLICATION_VERSION_PATH_VARIABLE = main.Version
CLIENT_APPLICATION_UI_SEPARATOR_PATH_VARIABLE = main.UISeparator

CERT_TOOL_IMAGE ?= ghcr.io/plgd-dev/hub/cert-tool:vnext
DEVSIM_IMAGE ?= ghcr.io/iotivity/iotivity-lite/cloud-server-discovery-resource-observable-debug:master
DEVSIM_PATH = $(shell pwd)/.tmp/devsim
CERT_PATH = $(TMP_PATH)/certs
UI_SEPARATOR = "--------UI--------"

certificates:
	mkdir -p $(CERT_PATH)
	docker pull $(CERT_TOOL_IMAGE)
	docker run --rm --user $(USER_ID):$(GROUP_ID) -v $(CERT_PATH):/out $(CERT_TOOL_IMAGE) --outCert=/out/rootcacrt.pem --outKey=/out/rootcakey.pem --cert.subject.cn="ca" --cmd.generateRootCA
	docker run --rm --user $(USER_ID):$(GROUP_ID) -v $(CERT_PATH):/out $(CERT_TOOL_IMAGE) --signerCert=/out/rootcacrt.pem --signerKey=/out/rootcakey.pem --outCert=/out/httpcrt.pem --outKey=/out/httpkey.pem --cert.san.domain=localhost --cert.san.ip=127.0.0.1 --cert.subject.cn="mfg" --cmd.generateCertificate
.PHONY: certificates

env: clean certificates
	mkdir -p $(DEVSIM_PATH)
	docker pull $(DEVSIM_IMAGE)
	docker run -d \
		--privileged \
		--network=host \
		--name devsim \
		-v $(DEVSIM_PATH):/tmp \
		$(DEVSIM_IMAGE) devsim-$(SIMULATOR_NAME_SUFFIX)
.PHONY: env

clean:
	docker rm -f devsim || true
	rm -rf $(TMP_PATH) || true
.PHONY: clean


build-web:
	mkdir -p $(WWW_PATH)
	docker build --tag build-web -f ./web/Dockerfile ./web
	$(eval container_id=$(shell docker create build-web))
	docker cp $(container_id):/web/build/ $(WWW_PATH)/
	docker rm -f $(container_id)
	cd $(WWW_PATH)/build && tar -czf $(TMP_PATH)/ui.tar.gz .
.PHONY: build-web

define build-binary
	GOOS=$(1) GOARCH=$(2) go build -ldflags="-X '$(CLIENT_APPLICATION_VERSION_PATH_VARIABLE)=$(VERSION_TAG)'" -ldflags="-X '$(CLIENT_APPLICATION_UI_SEPARATOR_PATH_VARIABLE)=$(UI_SEPARATOR)'" -o $(3) $(WORKING_DIRECTORY)/cmd
	printf "\n$(UI_SEPARATOR)\n" >> $(3)
	wc -c $(3)
	cat $(TMP_PATH)/ui.tar.gz >> $(3)
endef

build: clean build-web
	mkdir -p $(BUILD_PATH)
	$(call build-binary,linux,386,$(BUILD_PATH)/$(SERVICE_NAME).linux.386)
	$(call build-binary,linux,amd64,$(BUILD_PATH)/$(SERVICE_NAME).linux.amd64)
	$(call build-binary,linux,arm,$(BUILD_PATH)/$(SERVICE_NAME).linux.arm)
	$(call build-binary,linux,arm64,$(BUILD_PATH)/$(SERVICE_NAME).linux.arm64)
	$(call build-binary,windows,386,$(BUILD_PATH)/$(SERVICE_NAME).windows.386.exe)
	$(call build-binary,windows,amd64,$(BUILD_PATH)/$(SERVICE_NAME).windows.amd64.exe)
	$(call build-binary,windows,arm64,$(BUILD_PATH)/$(SERVICE_NAME).windows.arm64.exe)
	$(call build-binary,darwin,amd64,$(BUILD_PATH)/$(SERVICE_NAME).macos.amd64)
	$(call build-binary,darwin,arm64,$(BUILD_PATH)/$(SERVICE_NAME).macos.arm64)
.PHONY: build

test: env
	LISTEN_FILE_CA_POOL=$(WORKING_DIRECTORY)/.tmp/certs/rootcacrt.pem \
	LISTEN_FILE_CERT_DIR_PATH=$(WORKING_DIRECTORY)/.tmp/certs \
	LISTEN_FILE_CERT_NAME=httpcrt.pem \
	LISTEN_FILE_CERT_KEY_NAME=httpkey.pem \
	go test -v --race -p 1 -covermode=atomic -coverpkg=./... -coverprofile=$(WORKING_DIRECTORY)/.tmp/coverage.txt  ./...
.PHONY: test


GOOGLEAPIS_PATH := $(WORKING_DIRECTORY)/dependency/googleapis
GRPCGATEWAY_MODULE_PATH := $(shell go list -m -f '{{.Dir}}' github.com/grpc-ecosystem/grpc-gateway/v2 | head -1)
PLGDHUB_MODULE_PATH := $(shell go list -m -f '{{.Dir}}' github.com/plgd-dev/hub/v2 | head -1)

proto/generate:
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_device.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_devices.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_resource.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_device_resource_links.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/own_device.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/disown_device.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/clear_cache.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) -I=$(GOOGLEAPIS_PATH) -I=$(GRPCGATEWAY_MODULE_PATH) --go-grpc_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/service.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) -I=$(GOOGLEAPIS_PATH) -I=$(GRPCGATEWAY_MODULE_PATH) --openapiv2_out=$(GOPATH)/src \
		--openapiv2_opt logtostderr=true \
		$(WORKING_DIRECTORY)/pb/service.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) -I=$(GOOGLEAPIS_PATH) -I=$(GRPCGATEWAY_MODULE_PATH) --grpc-gateway_out $(GOPATH)/src \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		$(WORKING_DIRECTORY)/pb/service.proto
.PHONY: proto/generate

