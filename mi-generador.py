import argparse
import yaml

TESTING_NETWORK_NAME = "testing_net"

def create_docker_compose_base() -> dict[str]:
    return {"name": "tp0"}

def create_networks() -> dict[str]:
    return {
        "networks": {
            TESTING_NETWORK_NAME: {
                "ipam": {
                    "driver": "default",
                    "config": [
                        {
                            "subnet": "172.25.125.0/24"
                        }
                    ]
                }
            }
        }
    }

def create_service_server() -> dict[str]:
    return {
        "server": {
            "container_name": "server",
            "image": "server:latest",
            "entrypoint": "python3 /main.py",
            "environment": [
                "PYTHONUNBUFFERED=1",
                "LOGGING_LEVEL=DEBUG"
            ],
            "networks": [ TESTING_NETWORK_NAME ]
        }
    }

def create_client(client_id: int) -> dict[str]:
    client_name = f"client{client_id}"

    return {
        client_name: {
            "container_name": client_name,
            "image": "client:latest",
            "entrypoint": "/client",
            "environment": [
                f"CLI_ID={client_id}",
                "CLI_LOG_LEVEL=DEBUG"
            ],
            "networks": [TESTING_NETWORK_NAME],
            "depends_on": ["server"]
        }
    }

def create_services(n_clients: int) -> dict[str]:
    services = {}
    server = create_service_server()
    services.update(server)

    for client_id in range(1, n_clients + 1):
        client = create_client(client_id)

        services.update(client)

    return { "services": services }

def create_docker_compose_data(n_clients: int) -> dict[str]:
    base = create_docker_compose_base()
    services = create_services(n_clients)
    networks = create_networks()
    
    base.update(services)
    base.update(networks)

    return base


def parse_args():
    parser = argparse.ArgumentParser(
        prog="mi-generador",
        description="Generador de docker compose"
    )

    parser.add_argument("filename", type=str)
    parser.add_argument("n_clients", type=int)

    return parser.parse_args()

def main():
    args = parse_args()

    docker_compose_content = create_docker_compose_data(args.n_clients)
    
    with open(args.filename, "w", encoding="utf-8") as docker_compose_file:
        yaml.safe_dump(docker_compose_content, docker_compose_file, encoding="utf-8", sort_keys=False)

if __name__ == "__main__":
    main()