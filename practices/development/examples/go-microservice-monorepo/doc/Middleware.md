# Middleware

The common library contains a set of ready-to-use middlewares that you can attach to a microservice. Many of these
are just customized versions of existing [echo middleware](https://echo.labstack.com/docs/category/middleware)
which have been tweaked to work with the rest of the systems in the common library. If you don't want to pick and
choose, you can get all the middlewares in a bundle with the `middleware.StandardMiddleware()` function.

## Auth middleware

The auth middleware extracts the content of JWTs and makes them easily available from the `echo.Context` on the implementation
of a route. Since Istio handles verification and checking for JWTs in production, this middleware simply serves as an
extractor and allows unauthenticated requests to come through, extracting the token claims if a token is present.

See the [authentication docs](Authentication.md#extracting-authentication-information-in-a-rest-controller) 
for more information on extracting JWT information from a request.

## Database connection middleware

The database connection middleware holds a `sqlx.DB` database connection and attaches it to the context of incoming requests.
This way, database-based driven ports can extract the database connection from the request context which is passed through
every layer.

See [this section in the microservice architecture docs](Microservice%20Architecture.md#acquiring-a-database-connection)
for more information on accessing the database connection.

## Logging middleware

The logging middleware automatically logs data about incoming HTTP requests to the server using the global logger. It
assumes the global logger is already initialized. See [the logging documentation](Logging.md#instantiating-the-logger) 
for information on initializing the logger.

## CORS middleware

The CORS middleware automatically handles CORS preflight requests made by the browser, expressing which domains are allowed
to access its APIs. The `sharedoptions.AllowedOrigins` shared option controls which hostnames are allowed to initiate HTTP
requests from the browser, allowing microservices using this middleware to prevent cross-site scripting attacks.

See [the configuration docs](Configuration.md) for more information on configuration registries.

# Recovery middleware

The recovery middleware causes the server to recover from panics occurring in route handlers. This is just the one
provided [by the Echo framework](https://echo.labstack.com/docs/middleware/recover).
