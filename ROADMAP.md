# Roadmap

Here is the feature release schedule for this software.

## Implemented

You can use `dbgo` with a PostreSQL database.

You can use the `dbgo query` manager to manage your SQL statements or generate them.
- `dbgo query gen`
- `dbgo query template`
- `dbgo query save`
  
You can use `dbgo gen` to generate Database Driver Go code which calls your SQL queries based on the database as a single source of truth.


## Delayed

The following features are delayed due to a change in the maintainer's priorities.

You can implement them with a pull request or by [sponsoring SwitchUpCB](https://github.com/sponsors/switchupcb?frequency=one-time).

### March 17, 2025: `dbgo gen` (with domain)

You can use `dbgo gen` to generate Database Driver Go code which calls your SQL queries based on the domain as a single source of truth.

You can place a `_dbgo.sql` file in your queries directory to generate Go code without an SQL file combination operation when `dbgo gen` is called.

You can use options to only generate a `combined.SQL` file (currently output from the implemented `--keep` option) or only generate Database Driver Go code (from the `_dbgo.sql` file).


### March 24, 2025: `dbgo gen` (with .go templates)

You can use `dbgo gen` with `.go` files to customize the code generation algorithm, which will be updated to — by default — generate SQL statement-calling Go code for
- use of an SQL statement in a single network request.
- use of multiple SQL statements in an SQL batch statement _(i.e., a single network request)_.

### March 31, 2025: `dbgo gen` (with optimization)

Generated SQL statement-calling Go code is optimized for CPU usage and memory allocations beyond structural optimizations provided by the program.

An example of a structural optimization provided by the program is mapping SQL query results directly to domain models (Go types) to avoid memory allocations of a "data access object" and "domain transfer objects".

Structural optimizations provided by the program are already implemented by March 24, 2025. So, the March 31, 2025 update includes optimizations which have not been implemented _(e.g., in generated code)_.

### Future: Stored Procedures

You can use `dbgo` to add Stored Procedures to your database for Create (Insert), Read, and Update operations.

### Future: Automatic sqlc Query Annotation Developer

You can use a customizable [sqlc Query Annotation](https://docs.sqlc.dev/en/stable/reference/query-annotations.html) developer to automatically add query annotations to your `dbgo gen` SQL file before Go code generation occurs. 

### Future: Database Support

You can use `dbgo` with SQLITE3 and MySQL.
