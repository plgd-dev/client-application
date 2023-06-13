FROM goreleaser/goreleaser AS builder
ARG UI_FILE

RUN apk add upx
WORKDIR /go/src/github.com/plgd-dev/client-application
COPY go.mod go.sum ./
COPY ${UI_FILE} ${UI_FILE}
COPY . .
RUN make build BUILD_UI=false UI_FILE=${UI_FILE}

FROM ghcr.io/iotivity/iotivity-lite/cloud-server-discovery-resource-observable:latest AS service
COPY ./docker/run.sh /usr/local/bin/run.sh
COPY --from=builder /go/src/github.com/plgd-dev/client-application/.tmp/dist/*/client-application /usr/local/bin/client-application
RUN mkdir -p /plgd
WORKDIR /plgd
ENTRYPOINT [ "/usr/local/bin/run.sh" ]