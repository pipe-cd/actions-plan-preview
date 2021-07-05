FROM golang:1.16.5-alpine3.14 AS builder
COPY . /app
WORKDIR /app
RUN go build -o /plan-preview .

FROM gcr.io/pipecd/pipectl:v0.10.2
COPY --from=builder /plan-preview ./
RUN chmod +x ./plan-preview
ENTRYPOINT ["./plan-preview"]
