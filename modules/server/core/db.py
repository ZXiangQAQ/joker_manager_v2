from contextlib import contextmanager
from functools import wraps
from inspect import signature
from typing import TypeVar, Callable

from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from modules.server.core.config import app_config as config

engine = create_engine(
    config.db.dsn,
    pool_pre_ping=config.db.pool_pre_ping,
    pool_size=config.db.pool_size,
    pool_recycle=config.db.pool_recycle,
    echo=config.db.show_sql
)
SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

RT = TypeVar("RT")


def find_session_idx(func: Callable[..., RT]) -> int:
    """
    获取 session 在参数列表的位置
    :param func:
    :return:
    """
    func_params = signature(func).parameters
    try:
        session_args_idx = tuple(func_params).index("session")
    except ValueError:
        raise ValueError(f"Function {func.__qualname__} has no `session` argument") from None
    return session_args_idx


@contextmanager
def create_session():
    """
    https://docs.python.org/zh-cn/3.8/library/contextlib.html
    创建一个 session 会话实例
    :return:
    """
    session = SessionLocal()
    try:
        yield session
        session.commit()
    except Exception:
        session.rollback()
        raise
    finally:
        session.close()


def provider_session(func: Callable[..., RT]) -> Callable[..., RT]:
    """
    装饰器管理 session 连接
    :param func:
    :return:
    """
    session_args_idx = find_session_idx(func)

    @wraps(func)
    def wrapper(*args, **kwargs) -> RT:
        if "session" in kwargs or session_args_idx < len(args):
            # 不需要实例化`session`
            return func(*args, **kwargs)
        else:
            # 通过上下文管理器, 实例化`session`, 并传递给`function`
            with create_session() as session:
                return func(*args, session=session, **kwargs)

    return wrapper
