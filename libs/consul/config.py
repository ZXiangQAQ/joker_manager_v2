from typing import Optional

from pydantic import BaseModel, validator


class Config(BaseModel):
    host: str
    port: int
    username: Optional[str] = None
    password: Optional[str] = None

    @validator('username', pre=True)
    def password_match(cls, v, values, **kwargs):
        if v and not values["password"]:
            raise ValueError('username exists but no password provided')
        return v
