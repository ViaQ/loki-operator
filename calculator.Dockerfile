# Build the manager binary
FROM golang:1.16 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY cmd/size-calculator/main.go main.go
COPY api/ api/
COPY internal/ internal/

#@follow_tag(registry-proxy.engineering.redhat.com/rh-osbs/openshift-ose-cli:v4.7)
FROM registry-proxy.engineering.redhat.com/rh-osbs/openshift-ose-cli:v4.7.0-202104250659.p0 AS origincli
COPY --from=origincli /usr/bin/oc /usr/bin

# Build
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=on go build -a -o size-calculator main.go

# Use distroless as minimal base image to package the size-calculator binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/size-calculator .
USER 65532:65532

ENTRYPOINT ["/size-calculator"]
