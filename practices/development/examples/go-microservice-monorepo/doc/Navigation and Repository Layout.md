# Navigation and Repository Layout

The Go Microservices repository is a monorepo. It is designed to support multiple microservices relying on a common
library defining core architectural pieces each microservice can depend on - things such as middleware, configuration
code, code for connecting to the database, and the like. Each microservice then lives in its own folder and has a common
layout, making it easy to find different parts of the API/business logic/etc.

## Layout of the common library

* **commonlib** - Top-level folder containing foundational code which can be used as building blocks for new microservices
  * **auth** - Contains code for extracting authentication information from incoming requests, as well as data structures representing the contents of authentication information. For more info, see [Authentication.md](./Authentication.md).
  * **config** - Contains code for defining configuration options and retrieving them from a configuration registry. For more information see [Configuration.md](./Configuration.md)
    * **sharedoptions** - Contains common configuration options that may be used by all microservices
  * **database** - Contains database-related code, including functions for managing transactions and extracting the database connection from the current request context. Relevant information can be found in [Middleware.md](./Middleware.md), [Microservice Architecture.md](./Microservice Architecture.md), and [Testing.md](./Testing.md).
  * **logger** - Contains the global logger instance and functions for initializing it. For more information, see [Logging.md](./Logging.md).
  * **request** - Contains utilities for extracting information from requests, such as deserializing the request body or pulling out the request context. For more information, see [Microservice Architecture.md](./Microservice Architecture.md).
    * **testhelper** - Contains utilities for building HTTP requests in test code. See [Testing.md](./Testing.md) for more information.
  * **response** - Contains utilities for generating a standard error structure on HTTP responses. For more information, see [Microservice Architecture.md](./Microservice Architecture.md).
      * **dtos** - Contains code for common DTO types used in HTTP responses
      * **testhelper** - Contains utilities for deserializing HTTP response bodies into data structures in tests. See [Testing.md](./Testing.md) for more information.
  * **types** - Contains useful types for representing things in your code such as nullable values.
  * **sharedfeatures** - Contains common software features and controllers that can be used across microservices
    * **FEATURE NAME** - The name of the folder describes the microservice feature implemented by the business logic in this directory. See [Microservice Architecture.md](./Microservice Architecture.md) for more information.
      * **controller** - Contains REST controller definitions and DTOs which use and drive the business logic
      * **adapter** - Contains external access adapters used by the business logic to access or send data in external systems, such as message queues or databases.

## Layout of a microservice

For microservice layout, there are certain files at the root of the directory which are notable, so they will be mentioned in addition
to package-level stuff.

* **MICROSERVICE NAME** - This directory, such as "user", lives at the top level of the repository and defines the content of a microservice's code
  * **.env** - Contains configuration options ingested by [dotenv](https://github.com/joho/godotenv) on application startup so that you don't need to manually define configuration in your environment variables.
  * **main.go** - The entrypoint of the whole application. You should be able to follow startup logic from here.
  * **bootstrap.go** - Functions invoked by main.go to stand up the subsystems of the application, such as initializing the logger and connecting to the database. It also has functions for creating the HTTP router and attaching routes from all REST controllers in the app.
  * **options** - Contains the global configuration registry and initialization functions for it. See [Configuration.md](./Configuration.md) for more information.
  * **features** - Contains implementations for features that the microservice exposes
    * **FEATURE NAME** - The name of the folder describes the microservice feature implemented by the business logic in this directory. See [Microservice Architecture.md](./Microservice Architecture.md) for more information.
      * **controller** - Contains REST controller definitions which use and drive the business logic
      * **adapter** - Contains external access adapters used by the business logic to access or send data in external systems, such as message queues or databases
