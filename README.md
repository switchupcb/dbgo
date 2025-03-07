# Generate a database consumer module for your database based on domain types.

Use `dbgo` to stop wasting time developing optimized code (e.g, `Go`, `SQL`) for your database and domain models.

## What is dbgo?

`dbgo` generates a database consumer package for your database and domain models _(i.e., Go types)_.

## Why don't you use other database frameworks?

`dbgo` generates optimized Go code and SQL queries for your database and gives you the option to use domain models as a source of truth.

Other database frameworks generate generic Go code and SQL queries based on the database as a source of truth.

**Here is an example of the difference between `dbgo` and other frameworks.**

### What is your workflow with dbgo?

Your workflow with `dbgo` involves defining Go types _(e.g., domain models)_ and connecting to an existing database to generate:
1. a Repository Go package _(e.g., for a business domain)_ which transfers data from a datastore to your domain models.
2. a Datastore Go package _(e.g., for a `psql` database)_ without unnecessary "data access objects" or "data transfer functionality" _to reduce CPU usage and memory allocations_.
   1. Database Go models for Read (Select) operations which do not use reflection during runtime.
   2. Database Driver Go code to call **C**reate (Insert), **R**ead (Select), **U**pdate, **D**elete operations in a single or batch statement.
   3. Database Query Manager for SQL queries and Stored Procedures.
   4. Database Query Developer to develop custom type-safe SQL statements using Go type database models.
   5. Database Schema _(e.g., tables, views)_ Go type models.

**So, you can immediately use the domain type with the database once your Go types are defined and your database is set up.**

### What is your workflow with other database frameworks?

Your workflow with other database frameworks involves:
- Generators which generate Go and SQL code, but are impossible to customize when you need more than basic CRUD operations.
- ORMs which use reflection to perform CRUD operations that fetch EVERY FIELD FROM A TABLE instead of only fetching data you need.
- Query Builders which only help YOU CREATE SQL using Go types.

YOU WASTE TIME patching your repository on each database update using other database frameworks because no code is generated for your domain types.

## Table of Contents

| Topic                                | Category                                                                                                                                                                                                                                             |
| :----------------------------------- | :--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| [Using `dbgo`](#how-do-you-use-dbgo) |                                                                                                                                                                                                                                                      |  |
| with a domain                        | [1. Define Go types](#step-1-define-go-types-domain-models), [2. Deploy Database](#step-2-deploy-database), [3. Map Domain to Database](#step-3-map-domain-fields-to-database)                                                                       |
| with a database                      | [4. Configure setup file](#step-4-configure-the-setup-file), [5. Generate SQL](#step-5-generate-sql-statements-and-stored-procedures), [6. Generate Database Consumer](#generate-a-database-consumer-module-for-your-database-based-on-domain-types) |
|                                      |

## How do you use dbgo?

This demonstration generates a database consumer package for an `Account` domain. 

> _`dbgo` can generate a database consumer module without defining a domain. Skip to [step 4](#step-4-configure-the-setup-file) when this is your use case._

### Step 1. Define Go types (domain models)

Go types are defined in a file.

`./domain/domain.go`

```go
// Account represents a user's account.
type Account struct {
    ID int
    Username string
    Password string
    Name string
}
```

### Step 2. Deploy Database

You must connect to an existing database to run `dbgo`.

Here is the database diagram for the database used in this example.

```
TODO: add database diagram
```

### Step 3. Map Domain Fields to Database

Map the domain's fields to database schema _(e.g, table)_ fields.

`./domain/domain.go`

```go
// Account represents a user's account.
type Account struct {
    ID int
    Username string `dbgo:"users.name"`
    Password string `dbgo:"users.password"`
    Name string     `dbgo:"accounts.name"`
}
```

### Step 4. Configure the setup file

You set up `dbgo` with a YAML file.

`setup.yml`

```yml
generated:
    # Define the code generator inputs.
    input:
        # domain package containing Go types (relative to the setup file)
        dpkg: ./domain

        # database connection and schema (public by default).
        db: 
            # connection string or environment variable (e.g., `$VAR`)
            connection: postgresql://user:pass@localhost:5432/dbgo?sslmode=disable 
            schema: public

        # database query directory containing SQL query files (relative to the setup file).
        queries: datastore/psql/queries

    # Define where the code is generated (relative to the setup file).
    output:
        # domain repository package containing repository model functions.
        dpkg: datastore/domain

        # database package containing database model functions. 
        dbpkg: datastore/psql

        # Define the optional custom templates used to generate the file (.go supported).
        # template: ./generate.go

# Define custom options (which are passed to generator options) for customization.
custom:
  option: The possibilities are endless.
```

### Step 5. Generate SQL statements and stored procedures

Use the `dbgo query` manager to save customized type-safe SQL statements or generate them.

**1\)** Install the command line tool: `xstruct`.
```
go install github.com/switchupcb/xstruct@latest
```

**2\)** Install the command line tool: `sqlc`.
```
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

**3\)** Install the command line tool: `dbgo`.

```
go install github.com/switchupcb/dbgo@latest
```

**4\)** Run the executable with the following options to add SQL to the queries directory.

| Command Line                                | Description                                                                                                                                                                                                                           |
| :------------------------------------------ | :------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| `db query gen -y path/to/yml`             | Generates SQL statements for Read (Select) operations and adds Stored Procedures for Create (Insert), Update, Delete operations to the database.                                                                                      |
| `db query template <name> -y path/to/yml` | Adds a `name` template to the queries `templates` directory. The template contains Go type database models you can use to return a type-safe SQL statement from the `SQL()` function in `name.go` which is  called by `db query save`. |
| `db query save <name> -y path/to/yml`      | Saves an SQL file _(with the same name as the template \[e.g., `name.sql`\])_ containing an SQL statement _(returned from the `SQL()` function in `name.go`)_ to the queries directory.                                                |

_Here are additional usage notes._
- _`-y`, `--yml`: The path to the YML file must be specified in reference to the current working directory._
- _`db query template`: Every template is updated when this command is executed without a specified template._
- _`db query save`: Every template is saved when this command is executed without a specified template._
- _`db query save`: You are not required to initialize a `go.mod` file to run templates, but using `go get github.com/switchupcb/jet/v2@dbgo` in a `go.mod` related to the template files helps you identify compiler errors in your template files._

#### How do you develop type-safe SQL?

Running `db query template <name> -y path/to/yml` adds a `name.go` file with database models as Go types. Use these Go types with [`jet`](https://github.com/go-jet/jet) to return an `stmt.Sql()` from `SQL()`, which cannot be interpreted unless the Go code referencing struct fields can be compiled.

_Read <a href="https://github.com/go-jet/jet#how-quickly-bugs-are-found" target="_blank">"How quickly bugs are found"</a> for more information._

### Step 6. Generate the database consumer package

Install the command line tool when you haven't already: `dbgo`.

```
go install github.com/switchupcb/dbgo@latest
```

Run the executable with given options.
    
```bash
dbgo gen -y path/to/yml
```

_The path to the YML file must be specified in reference to the current working directory._

**You can view the output of this example [here](/examples/main/).**
