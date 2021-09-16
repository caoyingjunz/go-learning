FROM golang:1.15-alpine3.12 as builder
ARG GOPROXY
ARG APP
ENV GOPROXY=${GOPROXY}
WORKDIR /go/pixiu
COPY http-learning.go .
RUN CGO_ENABLED=0 go build -a -o httpd http-learning.go

FROM jacky06/static:nonroot
ARG APP
WORKDIR /
COPY --from=builder /go/pixiu/httpd /usr/local/bin/httpd
USER root:root