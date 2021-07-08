from fastapi import FastAPI
from starlette.middleware.cors import CORSMiddleware

from modules.server.api.v1.system import system
from modules.server.core.config import app_config
from modules.server.middleware.process_time import ProcessTimeMiddleware
from version import version


def get_applications() -> FastAPI:
    application = FastAPI(title=version.Project)

    # middlewares
    application.add_middleware(
        CORSMiddleware,
        allow_origins=app_config.server.allowed_hosts or ["*"],
        allow_credentials=True,
        allow_methods=["*"],
        allow_headers=["*"],
    )
    application.add_middleware(ProcessTimeMiddleware)

    # routers
    application.include_router(system, tags=["system"])
    return application


# singleton
app = get_applications()
