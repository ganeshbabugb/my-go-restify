# README

## Getting started

Before running the application you will need a working PostgreSQL installation and a valid DSN (data source name) for connecting to the database.

Please open the `cmd/api/main.go` file and edit it to include your valid DSN as the default value.

```
cfg.db.dsn = env.GetString("DB_DSN", "YOUR DEFAULT DSN GOES HERE")
```

Note that this DSN must be in the format `user:pass@localhost:port/db` and **not** be prefixed with `postgres://`.

Make sure that you're in the root of the project directory, fetch the dependencies with `go mod tidy`, then run the application using `go run ./cmd/api`:

```
$ go mod tidy
$ go run ./cmd/api
```

If you make a request to the `GET /status` endpoint using `curl` you should get a response like this:

```
$ curl -i localhost:8600/status
HTTP/1.1 200 OK
Content-Type: application/json
Date: Mon, 09 May 2022 20:46:37 GMT
Content-Length: 23

{
    "Status": "OK",
}
```

You can also start the application with live reload support by using the `run` task in the `Makefile`:

```
$ make run
```

## Project structure

Everything in the codebase is designed to be editable. Feel free to change and adapt it to meet your needs.

|     |     |
| --- | --- |
| **`assets`** | Contains the non-code assets for the application. |
| `↳ assets/migrations/` | Contains SQL migrations. |
| `↳ assets/efs.go` | Declares an embedded filesystem containing all the assets. |

|     |     |
| --- | --- |
| **`cmd/api`** | Your application-specific code (handlers, routing, middleware, helpers) for dealing with HTTP requests and responses. |
| `↳ cmd/api/context.go` | Contains helpers for working with request context. |
| `↳ cmd/api/errors.go` | Contains helpers for managing and responding to error conditions. |
| `↳ cmd/api/handlers.go` | Contains your application HTTP handlers. |
| `↳ cmd/api/helpers.go` | Contains helper functions for common tasks. |
| `↳ cmd/api/main.go` | The entry point for the application. Responsible for parsing configuration settings initializing dependencies and running the server. Start here when you're looking through the code. |
| `↳ cmd/api/middleware.go` | Contains your application middleware. |
| `↳ cmd/api/routes.go` | Contains your application route mappings. |
| `↳ cmd/api/server.go` | Contains a helper functions for starting and gracefully shutting down the server. |

|     |     |
| --- | --- |
| **`internal`** | Contains various helper packages used by the application. |
| `↳ internal/database/` | Contains your database-related code (setup, connection and queries). |
| `↳ internal/env` | Contains helper functions for reading configuration settings from environment variables. |
| `↳ internal/password/` | Contains helper functions for hashing and verifying passwords. |
| `↳ internal/request/` | Contains helper functions for decoding JSON requests. |
| `↳ internal/response/` | Contains helper functions for sending JSON responses. |
| `↳ internal/validator/` | Contains validation helpers. |
| `↳ internal/version/` | Contains the application version number definition. |

## Configuration settings

Configuration settings are managed via environment variables, with the environment variables read into your application in the `run()` function in the `main.go` file.

You can try this out by setting a `HTTP_PORT` environment variable to configure the network port that the server is listening on:

```
$ export HTTP_PORT="9999"
$ go run ./cmd/api
```

Feel free to adapt the `run()` function to parse additional environment variables and store their values in the `config` struct. The application uses helper functions in the `internal/env` package to parse environment variable values or return a default value if no matching environment variable is set. It includes `env.GetString()`, `env.GetInt()` and `env.GetBool()` functions for reading string, integer and bool values from environment variables. Again, you can add any additional helper functions that you need.

## Creating new handlers

Handlers are defined as `http.HandlerFunc` methods on the `application` struct. They take the pattern:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    // Your handler logic...
}
```

Handlers are defined in the `cmd/api/handlers.go` file. For small applications, it's fine for all handlers to live in this file. For larger applications (10+ handlers) you may wish to break them out into separate files.

## Handler dependencies

Any dependencies that your handlers have should be initialized in the `run()` function `cmd/api/main.go` and added to the `application` struct. All of your handlers, helpers and middleware that are defined as methods on `application` will then have access to them.

You can see an example of this in the `cmd/api/main.go` file where we initialize a new `logger` instance and add it to the `application` struct.

## Creating new routes

[chi](https://github.com/go-chi/chi) version 5 is used for routing. Routes are defined in the `routes()` method in the `cmd/api/routes.go` file. For example:

```
func (app *application) routes() http.Handler {
    mux := chi.NewRouter()

    mux.Get("/your/path", app.yourHandler)

    return mux
}
```

For more information about chi and example usage, please see the [official documentation](https://github.com/go-chi/chi).

## Adding middleware

Middleware is defined as methods on the `application` struct in the `cmd/api/middleware.go` file. Feel free to add your own. They take the pattern:

```
func (app *application) yourMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Your middleware logic...
        next.ServeHTTP(w, r)
    })
}
```

You can then register this middleware with the router using the `Use()` method:

```
func (app *application) routes() http.Handler {
    mux := chi.NewRouter()
    mux.Use(app.yourMiddleware)

    mux.Get("/your/path", app.yourHandler)

    return mux
}
```

It's possible to use middleware on specific routes only by creating route 'groups':

```
func (app *application) routes() http.Handler {
    mux := chi.NewRouter()
    mux.Use(app.yourMiddleware)

    mux.Get("/your/path", app.yourHandler)

    mux.Group(func(mux chi.Router) {
        mux.Use(app.yourOtherMiddleware)

        mux.Get("/your/other/path", app.yourOtherHandler)
    })

    return mux
}
```

Note: Route 'groups' can also be nested.

## Sending JSON responses

JSON responses and a specific HTTP status code can be sent using the `response.JSON()` function. The `data` parameter can be any JSON-marshalable type.

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"hello":  "world"}

    err := response.JSON(w, http.StatusOK, data)
    if err != nil {
        app.serverError(w, r, err)
    }
}
```

Specific HTTP headers can optionally be sent with the response too:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    data := map[string]string{"hello":  "world"}

    headers := make(http.Header)
    headers.Set("X-Server", "Go")

    err := response.JSONWithHeaders(w, http.StatusOK, data, headers)
    if err != nil {
        app.serverError(w, r, err)
    }
}
```

## Parsing JSON requests

HTTP requests containing a JSON body can be decoded using the `request.DecodeJSON()` function. For example, to decode JSON into an `input` struct:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name string `json:"Name"`
        Age  int    `json:"Age"`
    }

    err := request.DecodeJSON(w, r, &input)
    if err != nil {
        app.badRequest(w, r, err)
        return
    }

    ...
}
```

Note: The target decode destination passed to `request.DecodeJSON()` (which in the example above is `&input`) must be a non-nil pointer.

The `request.DecodeJSON()` function returns friendly, well-formed, error messages that are suitable to be sent directly to the client using the `app.badRequest()` helper.

There is also a `request.DecodeJSONStrict()` function, which works in the same way as `request.DecodeJSON()` except it will return an error if the request contains any JSON fields that do not match a name in the the target decode destination.

## Validating JSON requests

The `internal/validator` package includes a simple (but powerful) `validator.Validator` type that you can use to carry out validation checks.

Extending the example above:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    var input struct {
        Name      string              `json:"Name"`
        Age       int                 `json:"Age"`
        Validator validator.Validator `json:"-"`
    }

    err := request.DecodeJSON(w, r, &input)
    if err != nil {
        app.badRequest(w, r, err)
        return
    }

    input.Validator.CheckField(input.Name != "", "Name", "Name is required")
    input.Validator.CheckField(input.Age != 0, "Age", "Age is required")
    input.Validator.CheckField(input.Age >= 21, "Age", "Age must be 21 or over")

    if input.Validator.HasErrors() {
        app.failedValidation(w, r, input.Validator)
        return
    }

    ...
}
```

The `app.failedValidation()` helper will send a `422` status code along with any validation error messages. For the example above, the JSON response will look like this:

```
{
    "FieldErrors": {
        "Age": "Age must be 21 or over",
        "Name": "Name is required"
    }
}
```

In the example above we use the `CheckField()` method to carry out validation checks for specific fields. You can also use the `Check()` method to carry out a validation check that is _not related to a specific field_. For example:

```
input.Validator.Check(input.Password == input.ConfirmPassword, "Passwords do not match")
```

The `validator.AddError()` and `validator.AddFieldError()` methods also let you add validation errors directly:

```
input.Validator.AddFieldError("Email", "This email address is already taken")
input.Validator.AddError("Passwords do not match")
```

The `internal/validator/helpers.go` file also contains some helper functions to simplify validations that are not simple comparison operations.

|     |     |
| --- | --- |
| `NotBlank(value string)` | Check that the value contains at least one non-whitespace character. |
| `MinRunes(value string, n int)` | Check that the value contains at least n runes. |
| `MaxRunes(value string, n int)` | Check that the value contains no more than n runes. |
| `Between(value, min, max T)` | Check that the value is between the min and max values inclusive. |
| `Matches(value string, rx *regexp.Regexp)` | Check that the value matches a specific regular expression. |
| `In(value T, safelist ...T)` | Check that a value is in a 'safelist' of specific values. |
| `AllIn(values []T, safelist ...T)` | Check that all values in a slice are in a 'safelist' of specific values. |
| `NotIn(value T, blocklist ...T)` | Check that the value is not in a 'blocklist' of specific values. |
| `NoDuplicates(values []T)` | Check that a slice does not contain any duplicate (repeated) values. |
| `IsEmail(value string)` | Check that the value has the formatting of a valid email address. |
| `IsURL(value string)` | Check that the value has the formatting of a valid URL. |

For example, to use the `Between` check your code would look similar to this:

```
input.Validator.CheckField(validator.Between(input.Age, 18, 30), "Age", "Age must between 18 and 30")
```

Feel free to add your own helper functions to the `internal/validator/helpers.go` file as necessary for your application.

## Working with the database

This codebase is set up to use PostgreSQL with the [lib/pq](https://github.com/lib/pq) driver. You can control which database you connect to using the `DB_DSN` environment variable to pass in a DSN, or by adapting the default value in `run()`.

The codebase is also configured to use [jmoiron/sqlx](https://github.com/jmoiron/sqlx), so you have access to the whole range of sqlx extensions as well as the standard library `Exec()`, `Query()` and `QueryRow()` methods .

The database is available to your handlers, middleware and helpers via the `application` struct. If you want, you can access the database and carry out queries directly. For example:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    _, err := app.db.Exec("INSERT INTO people (name, age) VALUES ($1, $2)", "Alice", 28)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    ...
}
```

Generally though, it's recommended to isolate your database logic in the `internal/database` package and extend the `DB` type to include your own methods. For example, you could create a `internal/database/people.go` file containing code like:

```
type Person struct {
    ID    int    `db:"id"`
    Name  string `db:"name"`
    Age   int    `db:"age"`
}

func (db *DB) NewPerson(name string, age int) error {
    _, err := db.Exec("INSERT INTO people (name, age) VALUES ($1, $2)", name, age)
    return err
}

func (db *DB) GetPerson(id int) (Person, error) {
    var person Person
    err := db.Get(&person, "SELECT * FROM people WHERE id = $1", id)
    return person, err
}
```

And then call this from your handlers:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    err := app.db.NewPerson("Alice", 28)
    if err != nil {
        app.serverError(w, r, err)
        return
    }

    ...
}
```

## Managing SQL migrations

The `Makefile` in the project root contains commands to easily create and work with database migrations:

|     |     |
| --- | --- |
| `$ make migrations/new name=add_example_table` | Create a new database migration in the `assets/migrations` folder. |
| `$ make migrations/up` | Apply all up migrations. |
| `$ make migrations/down` | Apply all down migrations. |
| `$ make migrations/goto version=N` | Migrate up or down to a specific migration (where N is the migration version number). |
| `$ make migrations/force version=N` | Force the database to be specific version without running any migrations. |
| `$ make migrations/version` | Display the currently in-use migration version. |

Hint: You can run `$ make help` at any time for a reminder of these commands.

These `Makefile` tasks are simply wrappers around calls to the `github.com/golang-migrate/migrate/v4/cmd/migrate` tool. For more information, please see the [official documentation](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate).

By default all 'up' migrations are automatically run on application startup using embeded files from the `assets/migrations` directory. You can disable this by setting the `DB_AUTOMIGRATE` environment variable to `false`.

## Logging

Leveled logging is supported using the [slog](https://pkg.go.dev/log/slog) and [tint](https://github.com/lmittmann/tint) packages.

By default, a logger is initialized in the `main()` function. This logger writes all log messages above `Debug` level to `os.Stdout`.

```
logger := slog.New(tint.NewHandler(os.Stdout, &tint.Options{Level: slog.LevelDebug}))
```

Feel free to customize this further as necessary.

Also note: Any messages that are automatically logged by the Go `http.Server` are output at the `Warn` level.

## Using Basic Authentication

The `cmd/api/middleware.go` file contains a `basicAuth` middleware that you can use to protect your application — or specific application routes — with HTTP basic authentication.

You can try this out by visiting the [https://localhost:8600/basic-auth-protected](https://localhost:8600/basic-auth-protected) endpoint in any web browser and entering the default user name and password:

```
User name: admin
Password:  pa55word
```

You can change the user name and password by setting the `BASIC_AUTH_USERNAME` environment variable and `BASIC_AUTH_HASHED_PASSWORD` environment variable. For example:

```
$ export BASIC_AUTH_USERNAME='alice'
$ export BASIC_AUTH_HASHED_PASSWORD='$2a$10$xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
$ go run ./cmd/api
```

Note: You will probably need to wrap the username and password in `'` quotes to prevent your shell interpreting dollar and slash symbols as special characters.

The value for the `BASIC_AUTH_HASHED_PASSWORD` environment variable should be a bcrypt hash of the password, not the plaintext password itself. An easy way to generate the bcrypt hash for a password is to use the `gophers.dev/cmds/bcrypt-tool` package like so:

```
$ go run gophers.dev/cmds/bcrypt-tool@latest hash 'your_pa55word'
```

If you want to change the default values for username and password you can do so by editing the default command-line flag values in the `cmd/api/main.go` file.

## User accounts

The application is configured to support user accounts with fully-functional signup and authentication workflows.

A `User` struct describing the data for a user is defined in `internal/database/users.go`.

```
type User struct {
    ID             int       `db:"id"`
    Created        time.Time `db:"created"`
    Email          string    `db:"email"`
    HashedPassword string    `db:"hashed_password"`
}
```

Feel free to add additional fields to this struct (don't forget to also update the SQL queries, migrations, and handler code as necessary!).

A new user account can be created by sending a request to the `POST /users` endpoint:

```
$ curl -i -d '{"Email": "alice@example.com", "Password": "sectr3t_pa55word"}' localhost:8600/users
HTTP/1.1 204 No Content
Vary: Authorization
Date: Wed, 17 Aug 2022 05:18:12 GMT
```

Authentication is managed using stateless tokens. When running the application you should use your own secret key for signing the tokens. This key should be a random 32-character string generated using a CSPRNG which you pass to the application using the `JWT_SECRET` environment variable:

```
$ export JWT_SECRET_KEY="a1uiBXkmY03pxXok3OkFV39saE8Cn574"
$ go run ./cmd/api
```

A new authentication token can be created by sending the user's email and password to the `POST /authentication-tokens` endpoint.

```
$ curl -i -d '{"Email": "alice@example.com", "Password": "sectr3t_pa55word"}' localhost:8600/authentication-tokens
HTTP/1.1 200 OK
Content-Type: application/json
Vary: Authorization
Date: Wed, 17 Aug 2022 05:26:02 GMT
Content-Length: 353

{
    "AuthenticationToken": "eyJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOi8vbG9jYWxob3N0OjQ0NDQiLCJzdWIiOiIxIiwiYXVkIjpbImh0dHA6Ly9sb2NhbGhvc3Q6NDQ0NCJdLCJleHAiOjE2NjA4MDAzNjIuMjc0MDU2MiwibmJmIjoxNjYwNzEzOTYyLjI3NDA1NjcsImlhdCI6MTY2MDcxMzk2Mi4yNzQwNTY0fQ.t469-8hrwyZUN8gWmK5TeelXgstFnwBaoW977F2JbrE",
    "AuthenticationTokenExpiry": "2022-08-18T07:26:02+02:00"
}
```

The authentication token is a JWT containing the user's ID. By default authentication tokens are valid for 24 hours. You can change this by editing the code in the `createAuthenticationToken` handler.

Subsequent requests to the API should include the authentication token in a HTTP `Authorization` header in the following format:

```
Authorization: Bearer <authentication token>
```

The `authenticate` middleware is used to check for the presence of an `Authorization` header. If the token is valid, the token is decoded and the user information is fetched from the database. You can retrieve the details of the current user in your application handlers by calling the `contextGetAuthenticatedUser()` helper.

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    authenticatedUser, found := contextGetAuthenticatedUser(r)

    ...
}
```

If an `Authorization` header is provided with a request but it is invalid or expired, then the `authenticate` middleware will return a `401 Unauthorized` response and an error message to the client.

If no `Authorization` header is provided, then the request is coming from an unauthenticated client. In this case, the `authenticate` middleware _will not_ return an error, but subsequent calls to the `contextGetAuthenticatedUser()` helper function will return a boolean with the value `false`.

You can restrict access to specific handlers based on whether a request is coming from an authenticated client by using the `requireAuthenticatedUser` middleware. An example of using this can be seen in the `cmd/app/routes.go` file.

Important: You should only call the `requireAuthenticatedUser` middleware _after_ the `authenticate` middleware.

## Admin tasks

The `Makefile` in the project root contains commands to easily run common admin tasks:

|     |     |
| --- | --- |
| `$ make tidy` | Format all code using `go fmt` and tidy the `go.mod` file. |
| `$ make audit` | Run `go vet`, `staticheck`, `govulncheck`, execute all tests and verify required modules. |
| `$ make test` | Run all tests. |
| `$ make test/cover` | Run all tests and outputs a coverage report in HTML format. |
| `$ make build` | Build a binary for the `cmd/api` application and store it in the `/tmp/bin` folder. |
| `$ make run` | Build and then run a binary for the `cmd/api` application. |
| `$ make run/live` | Build and then run a binary for the `cmd/api` application (uses live reloading). |

## Live reload

When you use `make run/live` to run the application, the application will automatically be rebuilt and restarted whenever you make changes to any files with the following extensions:

```
.go
.tpl, .tmpl, .html
.css, .js, .sql
.jpeg, .jpg, .gif, .png, .bmp, .svg, .webp, .ico
```

Behind the scenes the live reload functionality uses the [cosmtrek/air](https://github.com/cosmtrek/air) tool. You can configure how it works (including which file extensions and folders are watched for changes) by editing the `Makefile` file.

## Running background tasks

A `backgroundTask()` helper is included in the `cmd/api/helpers.go` file. You can call this in your handlers, helpers and middleware to run any logic in a separate background goroutine. This useful for things like sending emails, or completing slow-running jobs.

You can call it like so:

```
func (app *application) yourHandler(w http.ResponseWriter, r *http.Request) {
    ...

    app.backgroundTask(r, func() error {
        // The logic you want to execute in a background task goes here.
        // It should return an error, or nil.
        err := doSomething()
        if err != nil {
            return err
        }

        return nil
    })

    ...
}
```

Using the `backgroundTask()` helper will automatically recover any panics in the background task logic, and when performing a graceful shutdown the application will wait for any background tasks to finish running before it exits.

## Application version

The application version number is generated automatically by Go based on your version control information, and will be either a version tag (e.g. `v1.2.3`) or a pseudo-version (e.g. `v0.0.0-20250219190134-59bdb76fda0c`). It can be retrieved by calling the `version.Get()` function from the `internal/version` package.

Alternatively, you can retrieve your latest version control system revision number (e.g. Git commit hash) by calling the `version.GetRevision()` function.

Important: The version control information will only be available after you have initialized version control for your repository (e.g. run `$ git init`) AND application is built using `go build`. If you run the application using `go run` then `version.Get()` will return the string `"(devel)"` or `"unavailable"`.

## Changing the module path

The module path is currently set to `github.com/ganeshbabugb/my-go-restify`. If you want to change this please find and replace all instances of `github.com/ganeshbabugb/my-go-restify` in the codebase with your own module path.