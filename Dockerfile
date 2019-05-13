FROM golang:1.12-alpine as builder
RUN apk update && apk add git
WORKDIR /go/src/build
COPY . .

ENV GOPROXY=https://proxy.golang.org
ENV GO111MODULE=on
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-extldflags "-static"' -o server *.go

FROM busybox:latest as runtime
COPY --from=builder /go/src/build/server /usr/bin/server
EXPOSE 5555
ENTRYPOINT ["/usr/bin/server"]
