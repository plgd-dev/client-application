FROM goreleaser/goreleaser AS builder
ARG UI_FILE

RUN apk add upx
WORKDIR /go/src/github.com/plgd-dev/client-application
COPY go.mod go.sum ./
COPY ${UI_FILE} ${UI_FILE}
COPY . .
RUN make build BUILD_UI=false UI_FILE=${UI_FILE}

FROM ghcr.io/iotivity/iotivity-lite/cloud-server-discovery-resource-observable:latest AS service
COPY ./docker/run.sh /plgd/run.sh
COPY --from=builder /go/src/github.com/plgd-dev/client-application/.tmp/dist/*/client-application /plgd/client-application
WORKDIR /plgd
ENTRYPOINT [ "/plgd/run.sh" ]