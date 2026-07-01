repository := "ofwfurtado"
image := "proxmox-mcp-server:latest"


@build:
    echo "Building..."
    go build -o build/
    echo "Done!"

@docker-build image=image tag="latest":
    echo "Building Docker image {{image}}..."
    docker build -t {{repository}}/{{image}}:{{tag}} .
    echo "Done!"
