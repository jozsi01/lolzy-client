# 1. fázis: Build
FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS builder

# Munkakönyvtár beállítása
WORKDIR /app

# Függőségek letöltése (gyorsítja a későbbi buildelést)
COPY go.mod go.sum ./
RUN go mod download

# Forráskód másolása
COPY . .

# Az argumentumokat a Docker Buildx automatikusan kitölti
ARG TARGETOS
ARG TARGETARCH

# Fordítás a célarchitektúrára
RUN CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o main .

# 2. fázis: Futtató környezet (kicsi és biztonságos)
FROM alpine:latest

WORKDIR /root/

# Csak a lefordított binárist másoljuk át a builder fázisból
COPY --from=builder /app/main .

# Alkalmazás futtatása
CMD ["./main"]