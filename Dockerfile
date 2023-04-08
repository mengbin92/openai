FROM golang:1.20.3-alpine AS builder
LABEL maintainer="mengbin1992@outlook.com"

COPY . /go/src/openai/
ENV GOPROXY="https://goproxy.cn,direct"
ENV CGO_ENABLED=0

# build openai client
WORKDIR /go/src/openai
RUN cd /go/src/openai && go build -ldflags "-s -w" -o openai

FROM alpine:3.17
LABEL maintainer="mengbin1992@outlook.com"

RUN mkdir /app

COPY --from=builder /go/src/openai/openai /app
COPY conf /app/conf

WORKDIR /app

ENTRYPOINT ["/app/openai"]