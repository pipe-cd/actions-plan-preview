FROM golang:1.16.5-alpine3.14

ARG PIPECTL_VERSION=v0.10.2
ARG PIPECTL_URL=https://github.com/pipe-cd/pipe/releases/download/${PIPECTL_VERSION}/pipectl_${PIPECTL_VERSION}_linux_amd64

RUN apk update && apk add curl && \
  curl -LO ${PIPECTL_URL} && \
  chmod +x pipectl_${PIPECTL_VERSION}_linux_amd64 && mv pipectl_${PIPECTL_VERSION}_linux_amd64 /bin/pipectl

COPY . /app

RUN cd /app && \
  go build -o /plan-preview . && \
  chmod +x /plan-preview

ENTRYPOINT ["/plan-preview"]
