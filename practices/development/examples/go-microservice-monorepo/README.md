# Go Microservice Template

This repository is a monorepo containing both a common library with reusable components that can be used across microservices
and individual microservices that build off the common library.

The common library is located under the `commonlib` directory, and the microservice that depends on it is located in `microsvc`.

To start a new microservice, you can create a copy of the `microsvc` directory and track down the TODO comments which will help
you customize the template for your needs.

## Getting Started
To get started on the Go microservices, you'll need the following prerequisites.

### Prerequisites
- [ ] [Docker](https://docs.rancherdesktop.io/getting-started/installation/) installed
- [ ] [Go](https://go.dev/dl/) (using version **Go 1.21**) installed
    - [ ] Install GO
    - [ ] Add `export PATH=$PATH:$GOPATH/bin` to rc file (i.e. .bashrc/.zshrc/etc.)
- [ ] [Install Dbamte](https://github.com/amacneil/dbmate/blob/main/README.md#installation)
  - For macOS: `brew install dbmate`
  - For Windows `scoop install dbmate`

### Start services
Once you have everything installed, you can start the microservice and the database by running the following:

1. `docker compose up -d` - This should build and start the example microservice container and its database
2. `cd microsvc && dbmate up` - Once the database is started, this will switch to the microservice directory and run the database migrations against the running database

You can also run the microservice template locally by changing into the `microsvc` directory and running `go run` with the database running and provisioned from the docker-compose file.

## Dbmate Migrations
Dbmate is a database migration tool that will run and keep track of migrations for the microservices in this template. Migrations are from the application and can be run from the command line. 
To run the migrations, `cd` into a microservice's directory and run `dbmate up`.

NOTE: `dbmate migrate` will only run pending migrations while `dbmate up` will create the database schema (if it does not already exist) and run any pending migrations.

### New Migrations
Ensure dbmate is installed locally by running `which dbmate`. If Dbmate is installed, it will return with a path to Dbmate, otherwise it will return with `dbmate not found`.

[Creating Migrations](https://github.com/amacneil/dbmate/blob/main/README.md#creating-migrations)
- In the microservice repo that the migration is for, run `dbmate new migration_name`
- Open newly create file `db/migrations/XXXXXXXXXXXXXX_migration_name.sql`
- Add the SQL for migration in the `migrate:up` section
- Add the SQL to undo the migration in the `migrate:down` section

## API Documentation

You can access the swagger documentation for the started microservice at [this location](http://localhost/api/swagger/index.html).

## Further reading
Documentation topics can be found in the [doc](doc) folder. A glossary is provided at the root in [index.md](doc/index.md).