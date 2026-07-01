# proxmox-mcp-server

MCP server for Proxmox VE.

Current capabilities:
- read-only inventory and inspection tools
- optional read-write power management tools
- `stdio` and `streamable-http` transports
- per-mode tool allow lists

## Local Tooling

This repository uses:
- [`mise`](https://mise.jdx.dev/) for local tool version management

Install project tools with:

```bash
mise install
```

Available tool versions are declared in [mise.toml](/home/fwfurtado/projects/fwfurtado/proxmox-mcp-server/mise.toml).

## Build

With Go directly:

```bash
go build -o build/proxmox-mcp-server .
```

With `just`:

```bash
just build
```

## Configuration

Environment variables:

```bash
PROXMOX_URL=https://pve.example.com:8006/api2/json
PROXMOX_TOKEN_ID=mcp@pve!readonly
PROXMOX_TOKEN_SECRET=...
PROXMOX_INSECURE_TLS=false

MCP_TRANSPORT=stdio
MCP_HTTP_ADDR=:8080
MCP_ALLOW_WRITE=false

MCP_READONLY_TOOLS=
MCP_READWRITE_TOOLS=
```

Notes:
- `PROXMOX_URL` must include `/api2/json`
- `MCP_READONLY_TOOLS` and `MCP_READWRITE_TOOLS` are comma-separated allow lists
- empty allow list means "register all tools in that mode"
- `*` is also accepted and means "register all tools in that mode"

## Run

If you use `mise`, you can run commands inside the configured tool environment with `mise exec -- ...`.

Example build with `mise` and `just`:

```bash
mise exec -- just build
```

`stdio`:

```bash
build/proxmox-mcp-server
```

`streamable-http`:

```bash
build/proxmox-mcp-server --transport streamable-http --http-addr :8080
```

Using a local `.env` file:

```bash
set -a
. ./.env
set +a
build/proxmox-mcp-server
```

HTTP transport:

```bash
set -a
. ./.env
set +a
build/proxmox-mcp-server --transport streamable-http --http-addr :8080
```

Enable write tools:

```bash
set -a
. ./.env
set +a
build/proxmox-mcp-server --allow-write
```

Restrict read-only tools:

```bash
set -a
. ./.env
set +a
build/proxmox-mcp-server --readonly-tools list_vms,get_vm,list_tasks
```

Restrict write tools:

```bash
set -a
. ./.env
set +a
build/proxmox-mcp-server --allow-write --readwrite-tools start_vm,shutdown_vm,reboot_container
```

## Tools

Read-only tools:
- `list_nodes`
- `list_vms`
- `get_vm`
- `get_vm_config`
- `list_containers`
- `get_container`
- `list_storage`
- `list_tasks`
- `get_task`
- `list_snapshots`
- `list_networks`
- `list_node_networks`
- `list_cluster_resources`

Read-write tools:
- `start_vm`
- `stop_vm`
- `shutdown_vm`
- `reboot_vm`
- `reset_vm`
- `start_container`
- `stop_container`
- `shutdown_container`
- `reboot_container`

## Proxmox Setup

The examples below assume:
- user: `mcp@pve`
- read-only token: `mcp@pve!readonly`
- read-write token: `mcp@pve!power`
- custom role: `MCPFull`

Run these commands on a Proxmox VE node as `root`.

### 1. Create the user

```bash
pveum user add mcp@pve
```

Optional comment:

```bash
pveum user modify mcp@pve --comment "MCP service user"
```

### 2. Create the read-write role

`MCPFull` is intended for this server's current write scope: read access plus VM/container power management.

```bash
pveum role add MCPFull -privs "Datastore.Audit Pool.Audit SDN.Audit Sys.Audit VM.Audit VM.PowerMgmt"
```

If the role already exists, update it:

```bash
pveum role modify MCPFull -privs "Datastore.Audit Pool.Audit SDN.Audit Sys.Audit VM.Audit VM.PowerMgmt"
```

### 3. Grant the base user the superset role

With API token privilege separation enabled, effective token permissions are the intersection of:
- user permissions
- token permissions

Because of that, the base user must have the superset of permissions needed by any token.

```bash
pveum aclmod / -user mcp@pve -role MCPFull
```

### 4. Create the read-only token

```bash
pveum user token add mcp@pve readonly --privsep 1
```

Grant `PVEAuditor` to the token:

```bash
pveum aclmod / -token 'mcp@pve!readonly' -role PVEAuditor
```

### 5. Create the read-write token

```bash
pveum user token add mcp@pve power --privsep 1
```

Grant `MCPFull` to the token:

```bash
pveum aclmod / -token 'mcp@pve!power' -role MCPFull
```

## Permission Model

With `--privsep 1`, token permissions are:

```text
effective = user permissions ∩ token permissions
```

That means:
- the user `mcp@pve` should have `MCPFull` on `/`
- the token `mcp@pve!readonly` can have `PVEAuditor` on `/`
- the token `mcp@pve!power` can have `MCPFull` on `/`

This keeps one service user with multiple tokens of different scope.

## Verify Effective Permissions

Read-only token:

```bash
curl -ksS \
  -H 'Authorization: PVEAPIToken=mcp@pve!readonly=TOKEN_SECRET' \
  'https://pve.example.com:8006/api2/json/access/permissions?path=/' | jq
```

Write token:

```bash
curl -ksS \
  -H 'Authorization: PVEAPIToken=mcp@pve!power=TOKEN_SECRET' \
  'https://pve.example.com:8006/api2/json/access/permissions?path=/' | jq
```

Expected for the write token:
- `Sys.Audit: 1`
- `VM.Audit: 1`
- `VM.PowerMgmt: 1`

Using a local `.env` file:

```bash
set -a
. ./.env
set +a
curl -ksS \
  -H "Authorization: PVEAPIToken=${PROXMOX_TOKEN_ID}=${PROXMOX_TOKEN_SECRET}" \
  "${PROXMOX_URL}/access/permissions?path=/" | jq
```

## Example `.env`

Read-only token:

```bash
PROXMOX_URL=https://pve.example.com:8006/api2/json
PROXMOX_TOKEN_ID=mcp@pve!readonly
PROXMOX_TOKEN_SECRET=your-readonly-token-secret
```

Write token:

```bash
PROXMOX_URL=https://pve.example.com:8006/api2/json
PROXMOX_TOKEN_ID=mcp@pve!power
PROXMOX_TOKEN_SECRET=your-power-token-secret
MCP_ALLOW_WRITE=true
```

## Operational Notes

- `get_vm` may need access to both VM status and VM config endpoints
- `list_storage` is node-oriented, so shared storage can appear once per node
- `get_task` returns task metadata plus task log lines
- on `streamable-http`, `CTRL+C` attempts graceful shutdown first and then forces close if needed
