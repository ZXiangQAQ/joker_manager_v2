from typing import Set

from pydantic import BaseSettings

from modules.server.types.env_settings import DotenvSettings, MySQLDsn


class ServerConfig(DotenvSettings):
    port: int
    log_level: str
    allowed_hosts: Set[str] = ["*"]


class DBConfig(DotenvSettings):
    dsn: MySQLDsn = "mysql://root@192.168.1.157:13306/test?charset=utf8mb4"
    pool_size: int = 100
    pool_recycle: int = 3600
    pool_pre_ping: bool = True
    show_sql: bool = False
    page_size: int = 20


class AppConfig(BaseSettings):
    server: ServerConfig = ServerConfig()
    db: DBConfig = DBConfig()


# singleton
app_config = AppConfig()
