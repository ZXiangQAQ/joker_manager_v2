# joker_manager_v2

重构项目, Python Web + Golang Agent + Vue Frontend

cmd没写，暂时直接拉起来
```python
uvicorn modules.server.main:app
```

目录结构

```shell
.
├── docs                  # 文档相关
├── libs                  # 封装的依赖包
├── models                # 数据库模型
├── modules               # 模块
│   ├── agent             # agent
│   └── server            # server
│       ├── api           # controller
│       ├── core          # core
│       ├── middleware    # middlewares
│       ├── service       # services
│       └── types         # schemas
└── version               # 版本信息

```