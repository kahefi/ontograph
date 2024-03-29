package ontograph

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// BlazegraphEndpoint is the SPARQL endpoint for a Blazegraph database
type BlazegraphEndpoint struct {
	host   string
	client *http.Client
}

// NewBlazegraphEndpoint creates a new endpoint on the specified host address of the Blazegraph database.
func NewBlazegraphEndpoint(hostAddr string) *BlazegraphEndpoint {
	ep := BlazegraphEndpoint{
		host:   hostAddr,
		client: http.DefaultClient,
	}
	return &ep
}

// NewBlazegraphStore creates a new store associated with a graph URI in the specified namespace. Operations will be conducted through the specified endpoint. This constructor does neither check if the namespace or graph exist nor if the endpoint is online.
func (ep *BlazegraphEndpoint) NewBlazegraphStore(uri, namespace string) *BlazegraphStore {
	store := BlazegraphStore{
		uri:       uri,
		namespace: namespace,
		endpoint:  ep,
	}
	return &store
}

// IsOnline checks if the Blazegraph endpoint is online (i.e. if it responds with HTTP 200 on its status endpoint).
func (ep *BlazegraphEndpoint) IsOnline() (bool, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s/bigdata/status", ep.host), nil)
	if err != nil {
		return false, err
	}
	code, _, err := ep.doHTTP(req)
	if err != nil {
		return false, err
	}
	if code != http.StatusOK {
		return false, fmt.Errorf("Unexpected status response: %d (Expected 200)", code)
	}
	return true, nil
}

// GetNamespaces retrieves a list of namespaces in the database.
func (ep *BlazegraphEndpoint) GetNamespaces() ([]string, error) {
	// Create request
	path := fmt.Sprintf("%s/bigdata/namespace?describe-each-named-graph=false", ep.host)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, err
	}

	// Execute request
	statusCode, data, err := ep.doHTTP(req)
	if err != nil {
		return nil, err
	}
	// Check response status
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("Failed to query namespaces from database (HTTP %d)", statusCode)
	}
	var rex = regexp.MustCompile("/bigdata/namespace/(.+)/sparql")
	matches := rex.FindAllStringSubmatch(string(data), -1)
	namespaces := []string{}
	for _, m := range matches {
		namespaces = append(namespaces, m[1])
	}

	// // Nothing found
	return namespaces, nil
}

// CreateNamespace creates a new namespace with the given ID in the database.
// The namespace must not contain special characters or `.`.
func (ep *BlazegraphEndpoint) CreateNamespace(id string) error {
	payload := fmt.Sprintf(`
	com.bigdata.rdf.store.AbstractTripleStore.vocabularyClass=com.bigdata.rdf.vocab.core.BigdataCoreVocabulary_v20160317
	com.bigdata.rdf.store.AbstractTripleStore.textIndex=false
	com.bigdata.rdf.store.AbstractTripleStore.axiomsClass=com.bigdata.rdf.axioms.NoAxioms
	com.bigdata.rdf.sail.isolatableIndices=false
	com.bigdata.rdf.store.AbstractTripleStore.justify=false
	com.bigdata.rdf.sail.truthMaintenance=false
	com.bigdata.namespace.%s.spo.com.bigdata.btree.BTree.branchingFactor=1024
	com.bigdata.rdf.sail.namespace=%s
	com.bigdata.rdf.store.AbstractTripleStore.quads=true
	com.bigdata.namespace.%s.lex.com.bigdata.btree.BTree.branchingFactor=400
	com.bigdata.rdf.store.AbstractTripleStore.geoSpatial=false
	com.bigdata.rdf.store.AbstractTripleStore.statementIdentifiers=false`, id, id, id)

	// Create request
	req, err := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/bigdata/namespace", ep.host), strings.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "text/plain")

	// Execute request
	statusCode, _, err := ep.doHTTP(req)
	if err != nil {
		return err
	}

	if statusCode != http.StatusCreated {
		return fmt.Errorf("Failed to create blazegraph namespace '%s' (HTTP %d)", id, statusCode)
	}
	return nil
}

// DropNamespace removes the namespace with the given ID from the database.
// If the namespace does not exist in the first place, no error is returned (use `NamespaceExists` to check specifically for existence).
func (ep *BlazegraphEndpoint) DropNamespace(id string) error {
	// Delete request
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("%s/bigdata/namespace/%s", ep.host, url.PathEscape(id)), nil)
	if err != nil {
		return err
	}

	// Execute request
	statusCode, _, err := ep.doHTTP(req)
	if err != nil {
		return err
	}

	if statusCode != http.StatusOK {
		return fmt.Errorf("Failed to delete blazegraph namespace '%s' (HTTP %d)", id, statusCode)
	}
	return nil
}

// NamespaceExists checks if a namespace with the given ID exists.
func (ep *BlazegraphEndpoint) NamespaceExists(id string) (bool, error) {
	// Retrieve list of namespaces
	namespaces, err := ep.GetNamespaces()
	if err != nil {
		return false, err
	}
	// Lookup namespace id
	for _, s := range namespaces {
		if s == id {
			return true, nil
		}
	}
	// Nothing found
	return false, nil
}

// GetGraphs retrieves a list of graphs within the namespace of the database.
func (ep *BlazegraphEndpoint) GetGraphs(namespace string) ([]string, error) {
	sparqlReq := fmt.Sprintf(`SELECT DISTINCT ?g WHERE { GRAPH ?g { ?s ?p ?o } }`)
	resSet, code, err := ep.DoSparqlJSONQuery(namespace, sparqlReq)
	if err != nil {
		return []string{}, err
	}
	if code != http.StatusOK {
		return []string{}, fmt.Errorf("Unexpected response when retrieving graphs (HTTP %d)", code)
	}
	// Retrieve graph URIs from result set
	graphUris := []string{}
	for _, binding := range resSet.Results.Bindings {
		graphUris = append(graphUris, binding["g"].Value)
	}
	return graphUris, nil
}

// DoSparqlTurtleQuery queries the database for data in Turtle (ttl) format.
func (ep *BlazegraphEndpoint) DoSparqlTurtleQuery(namespace, sparqlQuery string) ([]byte, int, error) {
	// Setup request payload
	encQuery := fmt.Sprintf("query=%s", url.QueryEscape(sparqlQuery))
	// Create request
	path := fmt.Sprintf("%s/bigdata/namespace/%s/sparql", ep.host, url.PathEscape(namespace))
	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(encQuery))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/x-turtle")

	// Execute request
	code, data, err := ep.doHTTP(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	return data, code, nil
}

// DoSparqlJSONQuery queries the database for data in JSON Result Set format.
func (ep *BlazegraphEndpoint) DoSparqlJSONQuery(namespace, sparqlQuery string) (JSONResultSet, int, error) {
	var resSet JSONResultSet
	// Setup request payload
	encQuery := fmt.Sprintf("query=%s", url.QueryEscape(sparqlQuery))

	// Create request
	path := fmt.Sprintf("%s/bigdata/namespace/%s/sparql", ep.host, url.PathEscape(namespace))
	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(encQuery))
	if err != nil {
		return resSet, http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/sparql-results+json")

	// Execute request
	code, data, err := ep.doHTTP(req)
	if err != nil {
		return resSet, http.StatusInternalServerError, err
	}
	if code != http.StatusOK {
		return resSet, code, nil
	}

	// Decode response body
	err = json.Unmarshal(data, &resSet)
	return resSet, code, err
}

// DoSparqlUpdate performs a SPARQL update on the database
func (ep *BlazegraphEndpoint) DoSparqlUpdate(namespace, sparqlUpdate string) (int, error) {
	// Setup request payload
	encUpdate := fmt.Sprintf("update=%s", url.QueryEscape(sparqlUpdate))
	// Create request
	path := fmt.Sprintf("%s/bigdata/namespace/%s/sparql", ep.host, url.PathEscape(namespace))
	req, err := http.NewRequest(http.MethodPost, path, strings.NewReader(encUpdate))
	if err != nil {
		return http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")

	// Execute request
	code, _, err := ep.doHTTP(req)
	if err != nil {
		return http.StatusInternalServerError, err
	}
	// Return status
	return code, nil
}

// doHTTP executes the given request and returns HTTP status code, result data and error.
// In case that the returned status code is -1, there was an error with the request itself.
// If the status code is a valid HTTP code and error is not nil, there was an error with
// decoding the response body.
func (ep *BlazegraphEndpoint) doHTTP(req *http.Request) (int, []byte, error) {
	res, err := ep.client.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer res.Body.Close()
	// Read body data
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return res.StatusCode, nil, err
	}
	return res.StatusCode, data, nil
}

// A JSONResultSet represents the result set for SPARQL queries in JSON format (see https://www.w3.org/TR/sparql11-results-json for details)
type JSONResultSet struct {
	Head struct {
		Link []string `json:"link,omitempty"`
		Vars []string `json:"vars,omitempty"`
	} `json:"head,omitempty"`
	Results struct {
		Distinct bool                              `json:"distinct,omitempty"` // Deprecated
		Ordered  bool                              `json:"ordered,omitempty"`  // Deprecated
		Bindings []map[string]JSONResultSetBinding `json:"bindings,omitempty"`
	} `json:"results,omitempty"`
	Boolean bool `json:"boolean,omitempty"`
}

// A JSONResultSetBinding represents a binding results from the JSONResultSet.Results.Bindings slice.
type JSONResultSetBinding struct {
	Type     string `json:"type,omitempty"` // "uri", "literal", "typed-literal" or "bnode"
	Value    string `json:"value,omitempty"`
	Lang     string `json:"xml:lang,omitempty"`
	DataType string `json:"datatype,omitempty"`
}
