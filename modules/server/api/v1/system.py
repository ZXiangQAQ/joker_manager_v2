from fastapi import APIRouter
from typing import Dict, Any
from modules.server.core.config import app_config


system = APIRouter()


@system.get("/config")
def config():
    return app_config
