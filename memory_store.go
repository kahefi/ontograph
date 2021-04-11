package ontograph

import (
	"errors"
	"io"

	"bytes"
	"regexp"
	"strings"

	"fmt"

	"github.com/deiu/rdf2go"
)

// MemoryStore is an in-memory implementation of the graph store. It uses the rdf2go backend to implement the methods and is suitable for smaller ontologies that fit into the working memory. While fast, avoid big graphs and consider using a database store for them instead.
type MemoryStore struct {
	uri   string
	graph *rdf2go.Graph
}

// NewMemoryStore creates a new in-memory graph store.
func NewMemoryStore(uri string) *MemoryStore {
	store := MemoryStore{
		uri:   uri,
		graph: rdf2go.NewGraph(""),
	}
	return &store
}

// ParseFromTurtle creates a new memory store from the parsed TTL data given in the reader.
func ParseFromTurtle(reader io.Reader) (*MemoryStore, error) {
	// Create a new graph
	g := rdf2go.NewGraph("")
	// Parse graph
	if err := g.Parse(reader, "text/turtle"); err != nil {
		return nil, err
	}
	// Find base URI
	const RDFType string = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"
	const OWLOntology string = "http://www.w3.org/2002/07/owl#Ontology"
	triple := g.One(nil, rdf2go.NewResource(RDFType), rdf2go.NewResource(OWLOntology))
	if triple == nil {
		// Use prefix from first triple as URI
		triple = <-g.IterTriples()
		if triple == nil {
			return nil, errors.New("No triple found in reader data")
		}
	}
	if triple == nil {
		return nil, errors.New("invalid ontology: missing owl:Ontology object ")
	}
	subj := triple.Subject.String()
	// Return new hive ontology
	store := MemoryStore{
		uri:   subj[1 : len(subj)-1],
		graph: g,
	}
	return &store, nil
}

// GetURI returns the named graph URI.
func (store *MemoryStore) GetURI() string {
	return store.uri
}

// GetFirstMatch retrieves the first triple that matches the pattern. Empty strings in subject, predicate or object are treated as wildcards.
func (store *MemoryStore) GetFirstMatch(subj, pred, obj string) (*Triple, error) {
	trp := store.graph.One(store.toTerm(subj), store.toTerm(pred), store.toTerm(obj))
	if trp == nil {
		return nil, nil
	}
	triple := Triple{
		Subject:   Term(trp.Subject.String()),
		Predicate: Term(trp.Predicate.String()),
		Object:    Term(trp.Object.String()),
	}
	return &triple, nil
}

// GetAllMatches retrieves all triples that match the pattern. Empty strings in subject, predicate or object are treated as wildcards.
func (store *MemoryStore) GetAllMatches(subj, pred, obj string) ([]Triple, error) {
	triples := []Triple{}
	// If the triple pattern is a complete wildcard, return all triples
	if subj == "" && pred == "" && obj == "" {
		for trp := range store.graph.IterTriples() {
			triples = append(triples, Triple{
				Subject:   Term(trp.Subject.String()),
				Predicate: Term(trp.Predicate.String()),
				Object:    Term(trp.Object.String()),
			})
		}
		return triples, nil
	}

	// Otherwise, find all occurrences using the `All` method
	for _, trp := range store.graph.All(store.toTerm(subj), store.toTerm(pred), store.toTerm(obj)) {
		triples = append(triples, Triple{
			Subject:   Term(trp.Subject.String()),
			Predicate: Term(trp.Predicate.String()),
			Object:    Term(trp.Object.String()),
		})
	}
	return triples, nil
}

// DeleteAllMatches removes all triples that match the pattern. Empty strings in subject, predicate or object are treated as wildcards.
func (store *MemoryStore) DeleteAllMatches(subj, pred, obj string) error {
	// Find all matching triples
	triples, err := store.GetAllMatches(subj, pred, obj)
	if err != nil {
		return err
	}
	// Delete all triples
	err = store.DeleteTriplesUnchecked(triples)
	return err
}

// GetAllTriples returns all triples in the store. The operation is equivalent to GetAllMatches("", "", "").
func (store *MemoryStore) GetAllTriples() ([]Triple, error) {
	return store.GetAllMatches("", "", "")
}

// AddTriple adds the given triple to the store. If the triple already exists, it errors with `ErrTripleAlreadyExists`.
func (store *MemoryStore) AddTriple(trp Triple) error {
	// Check if triple already exists
	foundTrp := store.graph.One(store.toTerm(trp.Subject.String()), store.toTerm(trp.Predicate.String()), store.toTerm(trp.Object.String()))
	if foundTrp != nil {
		return ErrTripleAlreadyExists
	}
	// Otherwise, add triple to store
	store.graph.AddTriple(store.toTerm(trp.Subject.String()), store.toTerm(trp.Predicate.String()), store.toTerm(trp.Object.String()))
	return nil
}

// AddTriples adds all the given triples to the store. If one of the triples already exist, it errors with `ErrTripleAlreadyExists`.
func (store *MemoryStore) AddTriples(trps []Triple) error {
	addedTrps := []Triple{}
	// Add all triples in sequence
	var err error
	for _, trp := range trps {
		err = store.AddTriple(trp)
		// Stop loop if there was an error
		if err != nil {
			break
		}
		// Otherwise, remember added triple
		addedTrps = append(addedTrps, trp)
	}
	// If there was an error, revoke the adding and return
	if err != nil {
		_ = store.DeleteTriplesUnchecked(addedTrps)
		return err
	}
	// All fine
	return nil
}

// AddTripleUnchecked adds the given triple to the store. It does not error if the triple already exists.
func (store *MemoryStore) AddTripleUnchecked(trp Triple) error {
	// Rdf2go will just add dplicate triples, so we have to check for existence eitherway and catch the conflict error
	err := store.AddTriple(trp)
	if err == ErrTripleAlreadyExists {
		return nil
	}
	return err
}

// AddTriplesUnchecked adds all the given triples to the store. It does not error if any of the triples already exists.
func (store *MemoryStore) AddTriplesUnchecked(trps []Triple) error {
	for _, trp := range trps {
		err := store.AddTripleUnchecked(trp)
		if err != nil {
			return err
		}
	}
	return nil
}

// DeleteTriple removes the given triple from the store. If the triple does not exist, it errors with `ErrTripleDoesNotExist`.
func (store *MemoryStore) DeleteTriple(trp Triple) error {
	// Check if triple exists
	foundTrp := store.graph.One(store.toTerm(trp.Subject.String()), store.toTerm(trp.Predicate.String()), store.toTerm(trp.Object.String()))
	if foundTrp == nil {
		return ErrTripleDoesNotExist
	}
	// Delete triple from store
	store.graph.Remove(foundTrp)
	return nil
}

// DeleteTriples remove all the given triples from the store. If one of the triples do not exist, it errors with `ErrTripleDoesNotExist` and no triple is deleted.
func (store *MemoryStore) DeleteTriples(trps []Triple) error {
	deletedTrps := []Triple{}
	// Delete all triples in sequence
	var err error
	for _, trp := range trps {
		err = store.DeleteTriple(trp)
		// Stop loop if there was an error
		if err != nil {
			break
		}
		// Otherwise, remember deleted triple
		deletedTrps = append(deletedTrps, trp)
	}
	// If there was an error, revoke the deletion and return
	if err != nil {
		_ = store.AddTriplesUnchecked(deletedTrps)
		return err
	}
	// All fine
	return nil
}

// DeleteTripleUnchecked removes the given triple from the store. It does not error if the triple does not exist.
func (store *MemoryStore) DeleteTripleUnchecked(trp Triple) error {
	// We need to get the exact triple object from the store to remove it...
	rdfTrp := store.graph.One(store.toTerm(trp.Subject.String()), store.toTerm(trp.Predicate.String()), store.toTerm(trp.Object.String()))
	store.graph.Remove(rdfTrp)
	return nil
}

// DeleteTriplesUnchecked removes all the given triples from the store. It does not error if any of the triples do not exist.
func (store *MemoryStore) DeleteTriplesUnchecked(trps []Triple) error {
	for _, trp := range trps {
		err := store.DeleteTripleUnchecked(trp)
		if err != nil {
			return err
		}
	}
	return nil
}

// Drop cleas the store and renders it unusable.
func (store *MemoryStore) Drop() error {
	store.uri = ""
	store.graph = nil
	return nil
}

// SerializeToTurtle writes the entire store into the writer in Turtle (TTL) format. If pretty is set to true, the TTL is pretty printed.
func (store *MemoryStore) SerializeToTurtle(w io.Writer, pretty bool) error {
	// Use native serializer if not pretty
	if !pretty {
		return store.graph.Serialize(w, "text/turtle")
	}

	// Setup base prefix map
	prefixMap := map[string]string{
		"":     store.uri + "#",
		"rdf":  "http://www.w3.org/1999/02/22-rdf-syntax-ns#",
		"rdfs": "http://www.w3.org/2000/01/rdf-schema#",
		"owl":  "http://www.w3.org/2002/07/owl#",
		"xsd":  "http://www.w3.org/2001/XMLSchema#",
	}
	// Find all imports
	const OWLImports string = "http://www.w3.org/2002/07/owl#imports"
	trps, err := store.GetAllMatches(NewResourceTerm(store.uri).String(), NewResourceTerm(OWLImports).String(), "")
	if err != nil {
		return err
	}
	importURIs := []string{}
	for _, trp := range trps {
		importURIs = append(importURIs, trp.Object.Value())
	}
	// Add imports to prefix map
	for _, importURI := range importURIs {
		abbr := importURI[strings.LastIndex(importURI, "/")+1:]
		prefixMap[abbr] = importURI + "#"
	}

	// Serialize ontology into buffer
	ttlBytes := new(bytes.Buffer)
	err = store.graph.Serialize(ttlBytes, "text/turtle")
	if err != nil {
		return err
	}
	// Convert result to string
	ttlContent := ttlBytes.String()

	// Setup Prefix block
	ttlPrefixes := ""
	for abbr, prefix := range prefixMap {
		// Setup prefix entry
		ttlPrefixes = fmt.Sprintf("%s@prefix %s: <%s> .\n", ttlPrefixes, abbr, prefix)
		// Apply prefixes
		var re = regexp.MustCompile(fmt.Sprintf(`\<%s(.+?)\>`, prefix))
		ttlContent = re.ReplaceAllString(ttlContent, fmt.Sprintf(`%s:$1`, abbr))
	}
	// Pretty format triples
	ttlContent = strings.Replace(ttlContent, " .", " .\n\n", -1)

	// Append prefix block and base path
	ttlContent = fmt.Sprintf("%s@base <%s> .\n\n%s", ttlPrefixes, store.uri, ttlContent)

	// Write result
	_, err = io.WriteString(w, ttlContent)
	return err
}

// Size returns the total number of triples in the store.
func (store *MemoryStore) Size() (int, error) {
	return store.graph.Len(), nil

}

// Helper functions

// toTerm converts the given string term in NTriple format into a rdf2go term.
func (store *MemoryStore) toTerm(term string) rdf2go.Term {
	if term == "" {
		return nil
	}
	t := Term(term)
	if t.IsResource() {
		return rdf2go.NewResource(t.Value())
	}
	if t.IsLiteral() {
		if t.Language() != "" {
			return rdf2go.NewLiteralWithLanguage(t.Value(), t.Language())
		}
		if t.Datatype() != "" {
			return rdf2go.NewLiteralWithDatatype(t.Value(), store.toTerm(NewResourceTerm(t.Datatype()).String()))
		}
		return rdf2go.NewLiteral(t.Value())
	}
	panic(fmt.Sprintf("Invalid term '%s'", term))
}
