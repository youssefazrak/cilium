FROM docker.io/library/golang:1.13.12 as builder
LABEL maintainer="maintainer@cilium.io"
ADD . /go/src/github.com/cilium/cilium
WORKDIR /go/src/github.com/cilium/cilium/operator
ARG LOCKDEBUG
ARG V
RUN make CGO_ENABLED=0 GOOS=linux LOCKDEBUG=$LOCKDEBUG PKG_BUILD=1 EXTRA_GOBUILD_FLAGS="-a -installsuffix cgo"
RUN strip cilium-operator

FROM docker.io/library/alpine:3.9.3 as certs
RUN apk --update add ca-certificates

FROM docker.io/library/alpine:3.9.3 as planer
ARG PLANER_VERSION=0.14.0
ADD https://artifactory.palantir.build/artifactory/internal-dist/com/palantir/deployability/planer/$PLANER_VERSION/planer-$PLANER_VERSION-linux-amd64.tgz!/planer /usr/local/bin
ADD planer /etc/planer
RUN chmod +x /usr/local/bin/*

FROM scratch
LABEL maintainer="maintainer@cilium.io"
COPY --from=builder /go/src/github.com/cilium/cilium/operator/cilium-operator /usr/bin/cilium-operator
COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=planer /usr/local/bin/planer /usr/local/bin/planer
COPY --from=planer /etc/planer /etc/planer
WORKDIR /
CMD ["/usr/bin/cilium-operator"]
