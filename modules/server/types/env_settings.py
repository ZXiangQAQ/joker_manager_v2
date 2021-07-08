from pydantic import BaseSettings, AnyUrl


class MySQLDsn(AnyUrl):
    allowed_schemes = {'mysql', 'sqlite'}
    user_required = True


class DotenvSettings(BaseSettings):
    # Config will be read from environment variables and/or ".env" files.

    class Config:
        env_file = '.env'
        env_file_encoding = 'utf-8'
