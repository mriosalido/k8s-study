FROM golang:1.17-alpine AS builder
WORKDIR /go/src/k8s-study
COPY . .
RUN go mod tidy
RUN go build
FROM alpine:3.14
COPY --from=builder /go/src/k8s-study/k8s-study /usr/local/bin
ENTRYPOINT ["/usr/local/bin/k8s-study"]
#CMD ["/usr/local/bin/server"]