# Build the manager binary
FROM registry.access.redhat.com/ubi9/go-toolset:1.17.7 as builder

ARG OPERATOR_VERSION
WORKDIR /workspace
USER root
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY api/ api/
COPY controllers/ controllers/
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o manager \
    -ldflags="-X 'main.Version=${OPERATOR_VERSION}'" main.go

FROM registry.access.redhat.com/ubi9/ubi-micro
WORKDIR /
COPY --from=builder /workspace/manager .
USER 65532:65532
# Define GOTRACEBACK to mark this container as using the Go language runtime
# for `skaffold debug` (https://skaffold.dev/docs/workflows/debug/).
ENV GOTRACEBACK=single

ENTRYPOINT ["/manager"]
