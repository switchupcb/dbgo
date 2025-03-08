# Roadmap

Here is the feature release schedule for this software.

## Implemented

You can use the `dbgo query` manager to manage your SQL statements or generate them.
- `dbgo query gen`
- `dbgo query template`
- `dbgo query save`

## March 10, 2025: `dbgo gen` (without domain)

You can use `dbgo gen` to generate Database Driver Go code which calls your SQL queries based on the database as a single source of truth.

## March 17, 2025: `dbgo gen` (with domain)

You can use `dbgo gen` to generate Database Driver Go code which calls your SQL queries based on the domain as a single source of truth.

## March 24, 2025: `dbgo gen` (with .go templates)

You can use `dbgo gen` with `.go` files to customize the code generation algorithm, which will be updated to — by default — generate SQL statement-calling Go code for
- use of an SQL statement in a single network request.
- use of multiple SQL statements in an SQL batch statement _(i.e., a single network request)_.

## March 31, 2025: `dbgo gen` (with optimization)

Generated SQL statement-calling Go code is optimized for CPU usage and memory allocations beyond structural optimizations provided by the program.

An example of a structural optimization provided by the program is mapping SQL query results directly to domain models (Go types) to avoid memory allocations of a "data access object" and "domain transfer objects".

Structural optimizations provided by the program are already implemented by March 24, 2025. So, the March 31, 2025 update includes optimizations which have not been implemented _(e.g., in generated code)_.
