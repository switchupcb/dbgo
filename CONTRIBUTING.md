# Contributing

## Contributor License Agreement
 
Contributions to this project must be accompanied by a **Contributor License Agreement**. 

You or your employer retain the copyright to your contribution: Accepting this agreement gives us permission to use and redistribute your contributions as part of the project.

## Pull Requests

Pull requests must pass all [CI/CD](#cicd) measures and follow the [code specification](#specification).

## Project Structure

The repository consists of a detailed [README](README.md), [examples](/examples/), and [**command-line program**](/cmd/).

### Command Line Interface

The command-line interface program _(cmd)_ consists of 4 packages.

| Package              | Description                                                                                             |
| :------------------- | :------------------------------------------------------------------------------------------------------ |
| cmd `/cmd`           | Contains command-line interface program logic implemented with [cobra](https://github.com/spf13/cobra). |
| config `/cmd/config` | Contains loaders (e.g., `.yml`) used to configure the program's settings.                               |
| query `/cmd/query`   | Contains `dbgo query` command logic.                                                                    |
| gen `/cmd/gen`       | Contains `dbgo gen` command logic.                                                                      |

## Specification

### CI/CD

#### Static Code Analysis

`dbgo` uses [golangci-lint](https://github.com/golangci/golangci-lint) to statically analyze code. 

You can install golangci-lint with `go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.64.5` and run it using `golangci-lint run`. 
- you must add a `diff` tool in your PATH when you receive a `diff` error: There is a `diff` tool located in the `Git` bin.
- Use `golangci-lint run --disable-all --no-config -Egofmt --fix` when you receive `File is not ... with -...`.

## Roadmap

You can read the roadmap [here](/ROADMAP.md).