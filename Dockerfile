FROM golang:1.11
RUN go get github.com/mitchellh/gox
ADD . /go/src/github.com/skpr/eks-map-groups
WORKDIR /go/src/github.com/skpr/eks-map-groups
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/skpr/eks-map-groups/bin/eks-map-groups_linux_amd64 /usr/local/bin/eks-map-groups
CMD ["eks-map-groups"]