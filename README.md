# RMQ BE

RabbitMQ backend for HTTP authorization plugin based on MongoDB. See
[`rabbitmq-auth-backend-http`][1] for more details.

## Build

The project uses `Makefile` to simplify build command. Just call `make` and
project will compile into `rmqbe` binary.

Run `go test ./...` to run all tests.

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

[1]: https://github.com/rabbitmq/rabbitmq-auth-backend-http
