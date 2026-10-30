# Etapa 1: Compilar a aplicação
FROM --platform=$BUILDPLATFORM golang:1.27.3-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH \
    go build -trimpath -o evaluation-service .

# Etapa 2: Executar a aplicação
FROM alpine:3.24

RUN apk add --no-cache ca-certificates \
    && addgroup -S app \
    && adduser -S -G app app

WORKDIR /app

COPY --from=builder /app/evaluation-service .

USER app

EXPOSE 8004

CMD ["./evaluation-service"]
