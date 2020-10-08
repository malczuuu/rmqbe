# RMQ BE

RabbitMQ backend for HTTP authorization plugin based on MongoDB. See
[`rabbitmq-auth-backend-http`][1] for more details.

The main purpose of this project was to learn a bit of Go programming
language.

## Build

Run `go build ./cmd/rmqbe` to build the application. Binary file will
be named `./rmqbe`.

Run `go test ./...` to run all tests.

Run application with `./rmqbe` and browse `http://localhost:8000`.

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
