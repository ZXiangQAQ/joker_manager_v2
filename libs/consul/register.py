from dataclasses import dataclass
from typing import Optional

import consul

from libs.consul.config import Config


@dataclass
class ServiceCheck:
    interval: str
    deregister_critical_service_after: str
    http: str


class Register:

    def __init__(
            self,
            cfg: Config,
            address: str,
            port: int,
            service: str,
            deregister_critical_service_after: Optional[int] = 60
    ) -> None:
        self.cfg = cfg
        self.address = address
        self.port = port
        self.service = service
        self.deregister_critical_service_after = deregister_critical_service_after

    def register(self):
        service_check = ServiceCheck(
            interval="30s",
            deregister_critical_service_after="60s",
            http=self.address
        )
        client = consul.Consul(host=self.cfg.host, port=self.cfg.port)
        client.agent.service.register(
            name=self.service,
            service_id=self.identity(),
            address=self.address,
            port=self.port,
            tags=[],
            check=None
        )

    def deregister(self):
        pass

    def identity(self) -> str:
        return ''
