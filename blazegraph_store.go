package ontograph

import (
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

// BlazegraphStore is a SPARQL endpoint implementation of the graph store. It uses a Blazegraph database to implement the methods and is suitable for larger ontologies that might not fit into memory.
type BlazegraphStore struct {
	uri       string
	namespace string
	endpoint  *BlazegraphEndpoint
}

// GetURI returns the named graph URI.
func (store *BlazegraphStore) GetURI() string {
	return store.uri
}

// GetFirstMatch retrieves the first triple that matches the pattern. Empty strings in subject, predicate or object are treated as wildcards.
func (store *BlazegraphStore) GetFirstMatch(subj, pred, obj string) (*Triple, error) {
	// TODO: might be implemented more efficiently?
	matches, err := store.GetAllMatches(subj, pred, obj)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, nil
	}
	return &matches[0], nil
}

// GetAllMatches retrieves all triples that match the pattern. Empty strings in subject, predicate or object are treated as wildcards.
func (store *BlazegraphStore) GetAllMatches(subj, pred, obj string) ([]Triple, error) {
	// Parse pattern to query parameters
	s := "?s"
	p := "?p"
	o := "?o"
	if subj != "" {
		s = Term(subj).String()
	}
	if pred != "" {
		p = Term(pred).String()
	}
	if obj != "" {
		o = Term(obj).String()
	}
	// Construct SPARQL query
	sparqlReq := fmt.Sprintf(`SELECT ?s ?p ?o WHERE { GRAPH <%s> { %s %s %s. } }`, store.uri, s, p, o)

	// Execute SPARQL query
	resSet, code, err := store.endpoint.DoSparqlJSONQuery(store.namespace, sparqlReq)
	if err != nil {
		return nil, err
	}
	if code != http.StatusOK {
		return nil, fmt.Errorf("Received unexpected status code from SPARQL query (HTTP %d): %s", code, sparqlReq)
	}
	// We got a result set, iterate through bindings and parse corresponding triples
	resTrps := []Triple{}
	for _, trpBinding := range resSet.Results.Bindings {
		sTerm := Term(subj)
		if subj == "" {
			sTerm = binding2Term(trpBinding["s"])
		}
		pTerm := Term(pred)
		if pred == "" {
			pTerm = binding2Term(trpBinding["p"])
		}
		oTerm := Term(obj)
		if obj == "" {
			oTerm = binding2Term(trpBinding["o"])
		}
		// Return result triple
		resTrps = append(resTrps, Triple{
			Subject:   sTerm,
			Predicate: pTerm,
			Object:    oTerm,
		})
	}
	return resTrps, nil
}

// DeleteAllMatches removes all triples that match the pattern. Empty strings in subject, predicate or object are treated as wildcards.
func (store *BlazegraphStore) DeleteAllMatches(subj, pred, obj string) error {
	// Parse pattern to query parameters
	s := "?s"
	p := "?p"
	o := "?o"
	if subj != "" {
		s = Term(subj).String()
	}
	if pred != "" {
		p = Term(pred).String()
	}
	if obj != "" {
		o = Term(obj).String()
	}
	// Setup SPARQL query for deletion
	sparqlReq := fmt.Sprintf(`DELETE WHERE { GRAPH <%s> { %s %s %s . } }`, store.uri, s, p, o)
	code, err := store.endpoint.DoSparqlUpdate(store.namespace, sparqlReq)
	// Check response status
	if err != nil {
		return err
	}
	if code == http.StatusNotFound {
		return nil
	}
	if code != http.StatusOK {
		return fmt.Errorf("Failed to delete triples from graph '%s' on namespace '%s' (HTTP %d)", store.namespace, store.uri, code)
	}
	// We succeeded
	return nil
}

// GetAllTriples returns all triples in the store. The operation is equivalent to GetAllMatches("", "", "").
func (store *BlazegraphStore) GetAllTriples() ([]Triple, error) {
	return store.GetAllMatches("", "", "")
}

// AddTriple adds the given triple to the store. If the triple already exists, it errors with `ErrTripleAlreadyExists`.
func (store *BlazegraphStore) AddTriple(trp Triple) error {
	// Check if triple already exists
	foundTrp, err := store.tripleExists(trp)
	if err != nil {
		return err
	}
	if foundTrp {
		return ErrTripleAlreadyExists
	}
	// Otherwise, add triple to store
	return store.AddTripleUnchecked(trp)
}

// AddTriples adds all the given triples to the store. If one of the triples already exist, it errors with `ErrTripleAlreadyExists`.
func (store *BlazegraphStore) AddTriples(trps []Triple) error {
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
func (store *BlazegraphStore) AddTripleUnchecked(trp Triple) error {
	// Setup SPARQL insert query
	ttlData := fmt.Sprintf("%s %s %s .", trp.Subject.String(), trp.Predicate.String(), trp.Object.String())
	sparqlReq := fmt.Sprintf("INSERT DATA { GRAPH <%s> { %s } }", store.uri, ttlData)
	code, err := store.endpoint.DoSparqlUpdate(store.namespace, sparqlReq)
	// Check response status
	if err != nil {
		return err
	}
	if code == http.StatusNotFound {
		return fmt.Errorf("Namespace '%s' does not exist (HTTP %d)", store.namespace, http.StatusNotFound)
	}
	if code != http.StatusOK {
		return fmt.Errorf("Failed to insert triple into graph '%s' on namespace '%s' (HTTP %d)", store.namespace, store.uri, code)
	}
	// We succeeded
	return nil
}

// AddTriplesUnchecked adds all the given triples to the store. It does not error if any of the triples already exists.
func (store *BlazegraphStore) AddTriplesUnchecked(trps []Triple) error {
	// Convert triples to TTL
	var ttlDataBuffer strings.Builder
	for _, trp := range trps {
		ttlDataBuffer.WriteString(fmt.Sprintf("%s %s %s .", trp.Subject.String(), trp.Predicate.String(), trp.Object.String()))
	}

	sparqlReq := fmt.Sprintf("INSERT DATA { GRAPH <%s> { %s } }", store.uri, ttlDataBuffer.String())
	code, err := store.endpoint.DoSparqlUpdate(store.namespace, sparqlReq)
	// Check response status
	if err != nil {
		return err
	}
	if code == http.StatusNotFound {
		return fmt.Errorf("Namespace '%s' does not exist (HTTP %d)", store.namespace, http.StatusNotFound)
	}
	if code != http.StatusOK {
		return fmt.Errorf("Failed to insert triples into graph '%s' on namespace '%s' (HTTP %d)", store.namespace, store.uri, code)
	}
	// We succeeded
	return nil
}

// DeleteTriple removes the given triple from the store. If the triple does not exist, it errors with `ErrTripleDoesNotExist`.
func (store *BlazegraphStore) DeleteTriple(trp Triple) error {
	// Check if triple already exists
	foundTrp, err := store.tripleExists(trp)
	if err != nil {
		return err
	}
	if !foundTrp {
		return ErrTripleDoesNotExist
	}

	// Otherwise, delete triple from store
	return store.DeleteTripleUnchecked(trp)
}

// DeleteTriples remove all the given triples from the store. If one of the triples do not exist, it errors with `ErrTripleDoesNotExist` and no triple is deleted.
func (store *BlazegraphStore) DeleteTriples(trps []Triple) error {
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
func (store *BlazegraphStore) DeleteTripleUnchecked(trp Triple) error {
	// Setup SPARQL deletion query
	ttlData := fmt.Sprintf("%s %s %s .", trp.Subject.String(), trp.Predicate.String(), trp.Object.String())
	sparqlReq := fmt.Sprintf("DELETE DATA { GRAPH <%s> { %s } }", store.uri, ttlData)
	code, err := store.endpoint.DoSparqlUpdate(store.namespace, sparqlReq)
	// Check response status
	if err != nil {
		return err
	}
	if code == http.StatusNotFound {
		return nil
	}
	if code != http.StatusOK {
		return fmt.Errorf("Failed to delete triple from graph '%s' on namespace '%s' (HTTP %d)", store.namespace, store.uri, code)
	}
	// We succeeded
	return nil
}

// DeleteTriplesUnchecked removes all the given triples from the store. It does not error if any of the triples do not exist.
func (store *BlazegraphStore) DeleteTriplesUnchecked(trps []Triple) error {
	// Convert triples to TTL
	var ttlDataBuffer strings.Builder
	for _, trp := range trps {
		ttlDataBuffer.WriteString(fmt.Sprintf("%s %s %s .", trp.Subject.String(), trp.Predicate.String(), trp.Object.String()))
	}
	// Fire SPARQL delete query for triples
	sparqlReq := fmt.Sprintf("DELETE DATA { GRAPH <%s> { %s } }", store.uri, ttlDataBuffer.String())
	code, err := store.endpoint.DoSparqlUpdate(store.namespace, sparqlReq)
	// Check response status
	if err != nil {
		return err
	}
	if code == http.StatusNotFound {
		return nil
	}
	if code != http.StatusOK {
		return fmt.Errorf("Failed to delete triples from graph '%s' on namespace '%s' (HTTP %d)", store.namespace, store.uri, code)
	}
	// We succeeded
	return nil
}

// Drop clears the store and renders it unusable.
func (store *BlazegraphStore) Drop() error {
	// Check if graph exists in the first place
	if store.endpoint == nil {
		return fmt.Errorf("Store was already dropped")
	}
	sparqlReq := fmt.Sprintf("ASK WHERE { GRAPH <%s> { ?s ?p ?o } }", store.uri)
	resSet, code, err := store.endpoint.DoSparqlJSONQuery(store.namespace, sparqlReq)
	// Check response status
	if err != nil {
		return err
	}
	if code == http.StatusNotFound || (code == http.StatusOK && !resSet.Boolean) {
		return fmt.Errorf("Graph '%s' does not exist on '%s", store.uri, store.namespace)
	}
	if code != http.StatusOK {
		return fmt.Errorf("Failed to query for existence of '%s' on namespace '%s' (HTTP %d)", store.uri, store.namespace, code)
	}

	// Drop graph
	sparqlReq = fmt.Sprintf("DROP GRAPH <%s>", store.uri)
	code, err = store.endpoint.DoSparqlUpdate(store.namespace, sparqlReq)
	// Check response status
	if err != nil {
		return err
	}
	if code == http.StatusNotFound {
		return fmt.Errorf("Namespace '%s' does not exist (HTTP %d)", store.namespace, code)
	}
	if code != http.StatusOK {
		return fmt.Errorf("Failed to delete graph '%s' on '%s' (HTTP %d)", store.uri, store.namespace, code)
	}
	store.uri = ""
	store.namespace = ""
	store.endpoint = nil
	return nil
}

// SerializeToTurtle writes the entire store into the writer in Turtle (TTL) format. If pretty is set to true, the TTL is pretty printed.
func (store *BlazegraphStore) SerializeToTurtle(w io.Writer, pretty bool) error {
	// Compile SPARQL construct query
	sparqlReq := fmt.Sprintf("CONSTRUCT { ?s ?p ?o } FROM <%s> WHERE {  ?s ?p ?o . }", store.uri)
	ttlBytes, code, err := store.endpoint.DoSparqlTurtleQuery(store.namespace, sparqlReq)
	// Check response status
	if err != nil {
		return err
	}
	if code == http.StatusNotFound {
		return fmt.Errorf("Namspace '%s' does not exist (HTTP %d)", store.namespace, http.StatusNotFound)
	}
	if code != http.StatusOK {
		return fmt.Errorf("Failed to query for graph '%s' (HTTP %d)", store.uri, code)
	}

	// Write out returned TTL if we do not need to prettify it
	if !pretty {
		_, err := w.Write(ttlBytes)
		return err
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

	// Convert TTL to string
	ttlContent := string(ttlBytes)

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
func (store *BlazegraphStore) Size() (int, error) {
	// Setup SPARQL query
	sparqlReq := fmt.Sprintf("SELECT (COUNT(*) as ?n) FROM <%s> WHERE { ?s ?p ?o } ", store.uri)
	resSet, code, err := store.endpoint.DoSparqlJSONQuery(store.namespace, sparqlReq)
	// Check response status
	if err != nil {
		return 0, err
	}
	if code == http.StatusNotFound {
		return 0, fmt.Errorf("Namspace '%s' does not exist (HTTP %d)", store.namespace, http.StatusNotFound)
	}
	if code != http.StatusOK {
		return 0, fmt.Errorf("Failed to execute SELECT query on namespace '%s' (HTTP %d)", store.namespace, code)
	}
	return strconv.Atoi(resSet.Results.Bindings[0]["n"].Value)
}

// ********************
// * Helper functions *
// ********************

func (store *BlazegraphStore) tripleExists(trp Triple) (bool, error) {
	// Make query
	sparqlReq := fmt.Sprintf("ASK WHERE { GRAPH <%s> { %s %s %s } }", store.uri, trp.Subject.String(), trp.Predicate.String(), trp.Object.String())
	resSet, code, err := store.endpoint.DoSparqlJSONQuery(store.namespace, sparqlReq)
	// Check response status
	if err != nil {
		return false, err
	}
	if code == http.StatusNotFound {
		return false, nil
	}
	if code != http.StatusOK {
		return false, fmt.Errorf("Failed to execute ASK query on namespace '%s' (HTTP %d)", store.namespace, code)
	}
	return resSet.Boolean, nil
}

func binding2Term(binding JSONResultSetBinding) Term {
	switch binding.Type {
	case "uri":
		return NewResourceTerm(binding.Value)
	case "literal":
		return NewLiteralTerm(binding.Value, binding.Lang, binding.DataType)
	case "typed-literal":
		return NewLiteralTerm(binding.Value, binding.Lang, binding.DataType)
	default:
		panic(fmt.Sprintf("Unknown JSON Result Set binding type '%s'", binding.Type))
	}
}
