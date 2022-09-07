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
CLIENT_APPLICATION_BINARY_PATH ?= 

CERT_TOOL_IMAGE ?= ghcr.io/plgd-dev/hub/cert-tool:vnext
DEVSIM_IMAGE ?= ghcr.io/iotivity/iotivity-lite/cloud-server-discovery-resource-observable-debug:master
DEVSIM_PATH = $(shell pwd)/.tmp/devsim
CERT_PATH = $(TMP_PATH)/certs
MFG_CERT_PATH = $(DEVSIM_PATH)/pki_certs
MFG_ROOT_CA_CRT = $(MFG_CERT_PATH)/cloudca.pem
MFG_CLIENT_APPLICATION_CRT = $(MFG_CERT_PATH)/mfgclientapplicationcrt.pem
MFG_CLIENT_APPLICATION_KEY = $(MFG_CERT_PATH)/mfgclientapplicationkey.pem
UI_SEPARATOR ?= "--------UI--------"

OAUTH_SERVER_PATH = $(shell pwd)/.tmp/oauth-server
OAUTH_SERVER_ID_TOKEN_PRIVATE_KEY = $(OAUTH_SERVER_PATH)/idTokenKey.pem
OAUTH_SERVER_ACCESS_TOKEN_PRIVATE_KEY = $(OAUTH_SERVER_PATH)/accessTokenKey.pem
CLOUD_SID = adebc667-1f2b-41e3-bf5c-6d6eabc68cc6

certificates:
	mkdir -p $(CERT_PATH)
	docker pull $(CERT_TOOL_IMAGE)
	docker run --rm --user $(USER_ID):$(GROUP_ID) -v $(CERT_PATH):/out $(CERT_TOOL_IMAGE) --outCert=/out/rootcacrt.pem --outKey=/out/rootcakey.pem --cert.subject.cn="ca" --cmd.generateRootCA
	docker run --rm --user $(USER_ID):$(GROUP_ID) -v $(CERT_PATH):/out $(CERT_TOOL_IMAGE) --signerCert=/out/rootcacrt.pem --signerKey=/out/rootcakey.pem --outCert=/out/httpcrt.pem --outKey=/out/httpkey.pem --cert.san.domain=localhost --cert.san.ip=127.0.0.1 --cert.subject.cn="http-server" --cmd.generateCertificate
	mkdir -p $(MFG_CERT_PATH)
	docker run --rm --user $(USER_ID):$(GROUP_ID) -v $(MFG_CERT_PATH):/out $(CERT_TOOL_IMAGE) --outCert=/out/cloudca.pem --outKey=/out/cloudcakey.pem --cert.subject.cn=MfgRootCA --cmd.generateRootCA
	docker run --rm --user $(USER_ID):$(GROUP_ID) -v $(MFG_CERT_PATH):/out $(CERT_TOOL_IMAGE) --signerCert=/out/cloudca.pem --signerKey=/out/cloudcakey.pem --outCert=/out/mfgcrt.pem --outKey=/out/mfgkey.pem --cert.subject.cn="mfg-device-cert" --cmd.generateCertificate
	docker run --rm --user $(USER_ID):$(GROUP_ID) -v $(MFG_CERT_PATH):/out $(CERT_TOOL_IMAGE) --signerCert=/out/cloudca.pem --signerKey=/out/cloudcakey.pem --outCert=/out/mfgclientapplicationcrt.pem --outKey=/out/mfgclientapplicationkey.pem --cert.subject.cn="mfg-client-application-cert" --cmd.generateCertificate
.PHONY: certificates

privateKeys:
	mkdir -p $(OAUTH_SERVER_PATH)
	openssl genrsa -out $(OAUTH_SERVER_ID_TOKEN_PRIVATE_KEY) 4096
	openssl ecparam -name prime256v1 -genkey -noout -out $(OAUTH_SERVER_ACCESS_TOKEN_PRIVATE_KEY)
.PHONY: privateKeys

env: clean certificates privateKeys
	mkdir -p $(DEVSIM_PATH)
	docker pull $(DEVSIM_IMAGE)
	docker run -d \
		--privileged \
		--network=host \
		--name devsim \
		-v $(DEVSIM_PATH):/tmp \
		-v $(MFG_CERT_PATH):/pki_certs \
		$(DEVSIM_IMAGE) devsim-$(SIMULATOR_NAME_SUFFIX)
.PHONY: env

clean:
	docker rm -f devsim || true
	rm -rf $(TMP_PATH) || true
.PHONY: clean


build-web:
	mkdir -p $(WWW_PATH)
	docker build --tag build-web:latest --target build-web -f ./web/Dockerfile ./web
	CONTAINER_ID=`docker create build-web:latest` && docker cp $$CONTAINER_ID:/web/build/ $(WWW_PATH)/ && docker rm -f $$CONTAINER_ID
	cd $(WWW_PATH)/build && tar -czf $(TMP_PATH)/ui.tar.gz .
.PHONY: build-web

inject-web: $(CLIENT_APPLICATION_BINARY_PATH)
	printf "\n$(UI_SEPARATOR)\n" >> $(CLIENT_APPLICATION_BINARY_PATH)
	wc -c $(CLIENT_APPLICATION_BINARY_PATH)
	cat $(TMP_PATH)/ui.tar.gz >> $(CLIENT_APPLICATION_BINARY_PATH)
.PHONY: inject-web

build: clean
	UI_SEPARATOR=$(UI_SEPARATOR) goreleaser build --rm-dist
.PHONY: build

test: env
	COVERAGE_FILE=$(WORKING_DIRECTORY)/.tmp/coverage.txt; \
	JSON_REPORT_FILE=$(WORKING_DIRECTORY)/.tmp/report.json; \
	export LISTEN_FILE_CA_POOL=$(WORKING_DIRECTORY)/.tmp/certs/rootcacrt.pem; \
	export LISTEN_FILE_CERT_DIR_PATH=$(WORKING_DIRECTORY)/.tmp/certs; \
	export LISTEN_FILE_CERT_NAME=httpcrt.pem; \
	export LISTEN_FILE_CERT_KEY_NAME=httpkey.pem; \
	export MFG_ROOT_CA_CRT=$(MFG_ROOT_CA_CRT); \
	export MFG_CLIENT_APPLICATION_CRT=$(MFG_CLIENT_APPLICATION_CRT); \
	export MFG_CLIENT_APPLICATION_KEY=$(MFG_CLIENT_APPLICATION_KEY); \
	export TEST_OAUTH_SERVER_ID_TOKEN_PRIVATE_KEY=$(OAUTH_SERVER_ID_TOKEN_PRIVATE_KEY); \
	export TEST_OAUTH_SERVER_ACCESS_TOKEN_PRIVATE_KEY=$(OAUTH_SERVER_ACCESS_TOKEN_PRIVATE_KEY); \
	export TEST_ROOT_CA_KEY=$(WORKING_DIRECTORY)/.tmp/certs/rootcakey.pem; \
	export TEST_ROOT_CA_CERT=$(WORKING_DIRECTORY)/.tmp/certs/rootcacrt.pem; \
	export TEST_CLOUD_SID=$(CLOUD_SID); \
	if [ -n "$${JSON_REPORT}" ]; then \
		go test -v --race -p 1 -covermode=atomic -coverpkg=./... -coverprofile=$${COVERAGE_FILE} -json ./... > "$${JSON_REPORT_FILE}" ; \
	else \
		go test -v --race -p 1 -covermode=atomic -coverpkg=./... -coverprofile=$${COVERAGE_FILE} ./... ; \
	fi ; \
	EXIT_STATUS=$$? ; \
	if [ $${EXIT_STATUS} -ne 0 ]; then \
		exit $${EXIT_STATUS}; \
	fi ;
.PHONY: test


GOOGLEAPIS_PATH := $(WORKING_DIRECTORY)/dependency/googleapis
GRPCGATEWAY_MODULE_PATH := $(shell go list -m -f '{{.Dir}}' github.com/grpc-ecosystem/grpc-gateway/v2 | head -1)
PLGDHUB_MODULE_PATH := $(shell go list -m -f '{{.Dir}}' github.com/plgd-dev/hub/v2 | head -1)

proto/generate:
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_device.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_devices.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_resource.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/update_resource.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/create_resource.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/delete_resource.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_device_resource_links.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/own_device.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/disown_device.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/clear_cache.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_information.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_identity_csr.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/update_identity_certificate.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_identity_certificate.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/get_json_web_keys.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) --go_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/update_json_web_keys.proto

	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) -I=$(GOOGLEAPIS_PATH) -I=$(GRPCGATEWAY_MODULE_PATH) --go-grpc_out=$(GOPATH)/src $(WORKING_DIRECTORY)/pb/service.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) -I=$(GOOGLEAPIS_PATH) -I=$(GRPCGATEWAY_MODULE_PATH) --openapiv2_out=$(GOPATH)/src \
		--openapiv2_opt logtostderr=true \
		$(WORKING_DIRECTORY)/pb/service.proto
	protoc -I=. -I=$(GOPATH)/src -I=$(PLGDHUB_MODULE_PATH) -I=$(GOOGLEAPIS_PATH) -I=$(GRPCGATEWAY_MODULE_PATH) --grpc-gateway_out $(GOPATH)/src \
		--grpc-gateway_opt logtostderr=true \
		--grpc-gateway_opt paths=source_relative \
		$(WORKING_DIRECTORY)/pb/service.proto
.PHONY: proto/generate

