# ---------- Build stage ----------
FROM golang:1.23-bookworm AS builder
WORKDIR /app
ENV CGO_ENABLED=0 GO111MODULE=on
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# main อยู่ที่ cmd/main.go → build ที่โฟลเดอร์ ./cmd
RUN go build -ldflags="-s -w" -o /bin/app ./cmd

# ---------- Runtime stage ----------
FROM alpine:3.20
RUN addgroup -S app && adduser -S app -G app \
 && apk add --no-cache ca-certificates tzdata \
 && ln -sf /usr/share/zoneinfo/Asia/Bangkok /etc/localtime
USER app
WORKDIR /home/app
COPY --from=builder /bin/app /home/app/app
EXPOSE 3000
ENTRYPOINT ["/home/app/app"]
