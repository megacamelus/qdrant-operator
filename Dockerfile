#
# Builder
#

FROM golang:1.22 as builder

ARG TARGETOS
ARG TARGETARCH
ARG LDFLAGS

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY cmd/ cmd/
COPY api/ api/
COPY internal/ internal/
COPY pkg/ pkg/
COPY hack/ hack/

RUN CGO_ENABLED=0 GOLDFLAGS=${LDFLAGS} GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} hack/scripts/build.sh .

#
# Image
#

FROM registry.access.redhat.com/ubi9/ubi-minimal:latest
WORKDIR /
COPY --from=builder /workspace/qdrant .
USER 65532:65532

ENTRYPOINT ["/qdrant"]
CMD ["run"]
