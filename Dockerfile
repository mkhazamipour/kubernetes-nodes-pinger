FROM golang:1.21.1 as builder

WORKDIR /app

COPY . /app

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/kubernetes-nodes-pinger

FROM alpine AS final
 
LABEL maintainer="Morteza Khazamipour"

WORKDIR /app

COPY --from=builder /app/kubernetes-nodes-pinger /app

COPY .env /app/

EXPOSE 8080

ENTRYPOINT [ "/app/kubernetes-nodes-pinger" ]
