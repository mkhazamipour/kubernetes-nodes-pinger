FROM golang:1.21.1 as builder

WORKDIR /app

COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/kubernetes-network-test

FROM alpine AS final
 
LABEL maintainer="Morteza Khazamipour"

WORKDIR /app

COPY --from=builder /app/kubernetes-network-test /app

COPY .env main.go /app/

EXPOSE 2112

ENTRYPOINT [ "/app/kubernetes-network-test" ]