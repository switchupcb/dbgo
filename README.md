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
