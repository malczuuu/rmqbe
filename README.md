# RMQ BE

RabbitMQ backend for HTTP authorization plugin based on MongoDB.

## Build

The project uses `Makefile` to simplify build command. Just call `make` and
project will compile into `rmqbe` binary.

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
