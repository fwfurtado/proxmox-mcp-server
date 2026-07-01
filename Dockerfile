FROM golang:1.26-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/proxmox-mcp-server .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates

COPY --from=build /out/proxmox-mcp-server /usr/local/bin/proxmox-mcp-server

ENTRYPOINT ["proxmox-mcp-server"]
