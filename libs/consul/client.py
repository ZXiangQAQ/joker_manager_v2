from dataclasses import dataclass
from typing import List, Dict
from urllib import parse

import consul


@dataclass
class ServiceInstance:
    id: str
    name: str
    version: str
    metadata: Dict[str, str]
    tags: List[str]
    endpoints: List[str]


class Client:

    def __init__(self, cli: consul.Consul) -> None:
        self.cli = cli

    def register(self, svc: ServiceInstance) -> None:
        address, port = None, None
        addresses = {}
        for endpoint in svc.endpoints:
            raw = parse.urlparse(endpoint)
            address = raw.hostname
            port = raw.port
            addresses[raw.scheme] = {"address": endpoint, "port": port}

        self.cli.agent.service.register(
            name=svc.name,
            service_id=svc.id,
            address=address,
            port=port,
            tags=svc.tags,
            check={
                "tcp": f"{address}:{port}",
                "interval": "10s",
                "deregisterCriticalServiceAfter": "180s"
            }
        )

    def deregister(self, service_id: str) -> None:
        self.cli.agent.service.deregister(service_id)

    def service(self, service: str, index: int, is_passing: bool):
        pass
