# Delta-trade Core Service

The boilerplace of this project comes from: 
https://github.com/amitshekhariitbhu/go-backend-clean-architecture


## Major Packages used in this project

- **gin**: Gin is an HTTP web framework written in Go (Golang). It features a Martini-like API with much better performance -- up to 40 times faster. If you need a smashing performance, get yourself some Gin.
- **mongo go driver**: The Official Golang driver for MongoDB.
- **jwt**: JSON Web Tokens are an open, industry-standard RFC 7519 method for representing claims securely between two parties. Used for Access Token and Refresh Token.
- **viper**: For loading configuration from the `.env` file. Go configuration with fangs. Find, load, and unmarshal a configuration file in JSON, TOML, YAML, HCL, INI, envfile, or Java properties formats.
- **bcrypt**: Package bcrypt implements Provos and Mazières's bcrypt adaptive hashing algorithm.
- **testify**: A toolkit with common assertions and mocks that plays nicely with the standard library.
- **mockery**: A mock code autogenerator for Golang used in testing.
- Check more packages in `go.mod`.

### How to run this project?

It is suggested to run with docker-compose, however, you can setup local environment to run without too.


#### Run with Docker

- Install Docker and Docker Compose.
- Go to `docker` folder
- Run `docker-compose up -d`.
- Access API using `http://localhost:8080` and swagger page via `http://localhost:8080/swagger/index.html`

#### Run without Docker (Not tested)

- Install `go` if not installed on your machine.
- Install `MongoDB` if not installed on your machine.
- Important: Change the `DB_HOST` to `localhost` (`DB_HOST=localhost`) in `.env` configuration file. `DB_HOST=mongodb` is needed only when you run with Docker.
- Run `go run cmd/main.go`.
- Access API using `http://localhost:8080`

### How to update Swagger page

- Update comments tags on API functions or main app function (the updating bit)
- Intall swag via `go install github.com/go-openapi/swag`
- Go to `src`
- Run `swag init` and files in `docs` will be updated
- Reboot the app

### How to run the test?

```bash
# Run all tests
go test ./...
```

### How to generate the mock code?

In this project, to test, we need to generate mock code for the use-case, repository, and database.

```bash
# Generate mock code for the usecase and repository
mockery --dir=domain --output=domain/mocks --outpkg=mocks --all

# Generate mock code for the database
mockery --dir=mongo --output=mongo/mocks --outpkg=mocks --all
```

Whenever you make changes in the interfaces of these use-cases, repositories, or databases, you need to run the corresponding command to regenerate the mock code for testing.
