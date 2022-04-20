## Basic Stratum Server

- Prerequisites
  - [ ] Docker & Docker compose installed on machine
  - [ ] Go modules enable on machine

#### Running the server

##### _start postgres docker container_

```sh
docker-compose up -d
```

##### _build binaries_

```sh
make
```

- #### Required Server Flags
  - `dsn` default : _postgresql://localhost:5432/luxor?user=luxor&password=luxor&sslmode=disable_
  - `port` default : 8080
