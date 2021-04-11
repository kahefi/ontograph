# ontograph

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/kahefi/ontograph/blob/main/LICENSE)

A Go package that contains utilities to work with triple stores and ontology graphs through an unified API.

## Installation
Simply run 
```bash
go get -u github.com/kahefi/ontograph
```

## Usage
There are two levels of abstraction available for use: Either interact directly with a triple store graph through the `GraphStore` interface or use the higher abstraction `OntologyGraph` that is passed one of the triple stores as backend.

### The GraphStore Interface
The graphstore interface receives and outputs RDF graph data for ontologies in form of term triples (subject-predicate-object).

```golang
// Initialize a new in-memory graph store
graphUri = "https://www.your-ontology-graph.com/abc"
graph = NewMemoryStore(graphUri)

// Define a triple using the NTriple format
s1 := Term("<https://www.your-ontology-graph.com/abc#subject-1>")
p1 := Term("<https://www.your-ontology-graph.com/abc#predicate-1>")
o1 := Term("\"some value\"")
trp1 := NewTriple(s1, p1, o1)

// Define a triple using the Term wrappers
p2 := NewResourceTerm("https://www.your-ontology-graph.com/abc#predicate-2")
o2 := NewLiteralTerm("another value", "en", "xsd:string")
trp2 := NewTriple(s1, p2, o2)

// Print some components of a term
fmt.Println(o2.Value())
fmt.Println(o2.DataType())
fmt.Println(o2.Language())

// Add some triples
graph.AddTriplesUnchecked([]Triple{trp1, trp2})

// Delete all triples which have s1 as subject
graph.DeleteAllMatches(s1.String(), "", "")

```

#### Simple In-Memory Store
```golang
// Initialize a new in-memory graph store
graphUri = "https://www.your-ontology-graph.com/abc"
memGraph = NewMemoryStore(graphUri)

// Add some triples
memGraph.AddTriplesUnchecked([]Triple{trp1, trp2})
```

#### Blazegraph Database Store
```golang
// Initialize a new Blazegraph database endpoint
endpoint = NewBlazegraphEndpoint("http://blazeraph-host:8080")

// Create a new namespace (or use the default `kb`)
namespace := "tenant-1"
endpoint.CreateNamespace(namespace)

// Connect a graph store to the Blazegraph endpoint
graphUri = "https://www.your-ontology-graph.com/abc"
dbGraph = endpoint.NewBlazegraphStore(graphUri, namespace)

// Add some triples
dbGraph.AddTriplesUnchecked([]Triple{trp1, trp2})
```

### Ontology Graphs
TODO

## Linting & Testing
Check for sanity by running the Ginkgo BDD test suite for the package:
```bash
go test -cover -coverprofile=./cover.out && go tool cover -func=./cover.out
```
Also make sure that code linting is passed by fixing the issues indicated with the linters:
```bash
golangci-lint run && goreportcard-cli -v
```
