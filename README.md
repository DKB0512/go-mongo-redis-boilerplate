# Go Boilerplate

A great starting point for building RESTful APIs in Go using Gin framework, connecting to a Mongo database and Redis as cache.

### Features

- Implements the Clean Architecture pattern for a scalable and maintainable
- Uses the Gin framework for efficient and fast handling of HTTP requests
- Integrates with MongoDB for powerful and flexible database operations
- Integrates with Redis database for efficient caching
- Uses Go Swag for Generating Rest Docs and Swagger panel (https://github.com/swaggo/swag)
- Uses Air for live reload app (https://github.com/cosmtrek/air)

##### Authentication

- Supports JWT authentication with configurable expiration and issuer, allowing for flexible and secure authentication processes.

### Getting Started

##### Prerequisites

- Go version 1.18.1 or higher

To get up and running with the Go-Boilerplate, follow these simple steps:

```
$ git clone https://github.com/DKB0512/go-mongo-redis-boilerplate.git
$ cd go-mongo-redis-boilerplate
$ cp internal/config/.env.example internal/config/.env # create a copy of the example environment file, and also follow configuration steps on the difference section below
$ go src/main.go Or air
```

#### Configuration

Generate Swagger doc files

```
$ swag init -d src/
```

### File Structure

    ..
    ├── docs                                            # Document for swagger.
    ├── src                                             #
    │   ├── common                                      # Common Types And Struct.
    |   │        └── controller.go                      # Base Controller Structure Type.
    │   │        └── model.go                           # Base Model Structure Type.
    |   ├── config                                      # Configs
    |   |        └── config.go                          # Base Config Module and Env Init.
    │   └── controllers                                 # Controllers
    │   │         └── articles.controllers.go           # Article Controller (example).
    │   │         └── base.go                           # Base Controller Structure.
    │   │         └── products.controllers.go           # Products Controller (example).
    │   │         └── swagger.controllers.go            # Swagger Controller
    │   │         └── users.controllers.go              # Users Controller (example).
    │   │         └── auth.controllers.go               # Auth Controller.
    │   │                                               #
    │   └── middleware                                  # Middlewares
    │   │         └── jwt.middleware.go                 # jwt Middlewares.
    |   |                                               #
    │   ├── core                                        # Core Configures
    │   │   └── db                                      # Db Configures
    │   │       └── mongo.go                            # MongoDb Configure File
    │   │       └── redis.go                            # Redis Configure File
    │   │                                               #
    │   ├── models                                      # Models
    │   │   └── base.go                                 # Base Model Structure.
    │   │   └── article.model.go                        # Article Model (example).
    │   │   └── cache.model.go                          # Cache Model (example).
    │   │   └── product.model.go                        # Product Model (example).
    │   │   └── user.model.go                           # User Model (example).
    │   │   └── auth.model.go                           # Auth Model.
    │   │                                               #
    │   ├── utils                                       # Utils.
    │   │   └── http.go                                 # Http Utils
    │   │   └── token.go                                # Token Utils
    │   │                                               #
    │   ├── main.go                                     # Main File.
    │   │                                               #
    ├── .env.example                                    # Enviroment Example File
    ├── Dockerfile                                      # Dockerfile
    ├── docker-compose.yml                              # docker compose file
    ├── .air.toml                                       # air configure
    └── ...

#### Special Thanks

@Mehdikarimian - Inital Boilerplate with PostgresDB & Gorm TypeOrm
