# ontograph

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kahefi/ontograph/blob/main/LICENSE)

A Go package that contains utilities to work with RDF ontology graphs on different graph storages.

## Installation
Simply run 
```bash
go get -u github.com/kahefi/ontograph
```

## Usage
You can use two levels of abstraction: Either interact directly with a triple store graph through the `GraphStore` interface or use the higher abstraction `OntologyGraph` that is passed one of the triple stores as backend.

### Simple In-Memory Stores
TODO

### Blazegraph SPARQL Stores
TODO

### Ontology Graphs
TODO

## Linting & Testing
Check package for sanity by running the Ginkgo BDD test suites:
```bash
go test -cover
```
Also make sure that code linting is passed by fixing the issues indicated with the linters:
```bash
golangci-lint run && goreportcard-cli -v
``` 