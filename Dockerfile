FROM golang:1.26.4-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/proxmox-mcp-server .

FROM alpine:3.22

RUN apk add --no-cache ca-certificates \
	&& addgroup -S app \
	&& adduser -S -G app app

COPY --from=build /out/proxmox-mcp-server /usr/local/bin/proxmox-mcp-server

USER app

ENTRYPOINT ["proxmox-mcp-server"]
