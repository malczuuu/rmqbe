# RMQ BE

RabbitMQ backend for HTTP authorization plugin based on MongoDB. See
[`rabbitmq-auth-backend-http`][rabbitmq-auth-backend-http] for more details.

The main purpose of this project was to learn a bit of Go programming language.

## Build

It uses [Taskfile](https://taskfile.dev/docs/getting-started) tool.

```bash
git clone https://github.com/malczuuu/rmqbe.git
cd rmqbe

task test
task build

./dist/rmqbe
```

Then browse `http://localhost:8000`.

## API

| endpoint         | parameters                                                                          |
| ---------------- | ----------------------------------------------------------------------------------- |
| `POST /user`     | `username`<br/>`password`                                                           |
| `POST /vhost`    | `username`<br/>`vhost`<br/>`ip`                                                     |
| `POST /resource` | `username`<br/>`vhost`<br/>`resource`<br/>`name`<br/>`permission`                   |
| `POST /topic`    | `username`<br/>`vhost`<br/>`resource`<br/>`name`<br/>`permission`<br/>`routing_key` |

## MongoDB Model

```
Permission {
    resource:    String,
    name:        String,
    permission:  String,
    routing_key: String?
}
```

```
User {
    _id:         ObjectId,
    username:    String,
    password:    BCrypt,
    vhosts:      String[],
    permissions: Permission[]
}
```

[rabbitmq-auth-backend-http]: https://github.com/rabbitmq/rabbitmq-auth-backend-http
