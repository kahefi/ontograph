package ontograph

// An OntologyGraph represents an ontology backed by a grapg store using a higher abstraction level.
type OntologyGraph struct {
	graph GraphStore
}

// NewOntologyGraph creates a new ontology using the given graph store as backend.
func NewOntologyGraph(graph GraphStore) *OntologyGraph {
	ont := OntologyGraph{
		graph: graph,
	}
	return &ont
}

// GetURI returns the URI of the ontology
func (ont *OntologyGraph) GetURI() string {
	return ont.graph.GetURI()
}

// GetVersion returns the version set for this ontology. If not version is set, the empty string is returned.
func (ont *OntologyGraph) GetVersion() (string, error) {
	trp, err := ont.graph.GetFirstMatch(
		NewResourceTerm(ont.GetURI()).String(),
		NewResourceTerm(OWLVersionInfo).String(),
		"",
	)
	if err != nil {
		return "", err
	}
	// Check if a version was found at all
	if trp == nil {
		return "", nil
	}
	// Return version value
	return trp.Object.Value(), nil
}

// SetVersion sets a version for this ontology. All previous versions will be deleted when a new one is set!
func (ont *OntologyGraph) SetVersion(version string) error {
	// First delete all previous versions
	if err := ont.graph.DeleteAllMatches(NewResourceTerm(ont.GetURI()).String(), NewResourceTerm(OWLVersionInfo).String(), ""); err != nil {
		return err
	}
	// Set new version
	err := ont.graph.AddTripleUnchecked(Triple{
		Subject:   NewResourceTerm(ont.GetURI()),
		Predicate: NewResourceTerm(OWLVersionInfo),
		Object:    NewLiteralTerm(version, "", ""),
	})
	if err != nil {
		return err
	}
	// Success
	return nil
}

// GetImports returns a list of URIs for the imported ontologies.
func (ont *OntologyGraph) GetImports() ([]string, error) {
	// Get triples with import predicate
	trps, err := ont.graph.GetAllMatches(
		NewResourceTerm(ont.GetURI()).String(),
		NewResourceTerm(OWLImports).String(),
		"",
	)
	if err != nil {
		return nil, err
	}
	// Extract imported URIs
	importUris := []string{}
	for _, trp := range trps {
		importUris = append(importUris, trp.Object.Value())
	}
	return importUris, nil
}

// AddImport adds an ontology to the list of imports in the ontology.
func (ont *OntologyGraph) AddImport(uri string) error {
	// Get triples with import predicate
	return ont.graph.AddTriple(Triple{
		Subject:   NewResourceTerm(ont.GetURI()),
		Predicate: NewResourceTerm(OWLImports),
		Object:    NewResourceTerm(uri),
	})
}

func (ont *OntologyGraph) CreateIndividual(id string) error {
	return nil
}

func (ont *OntologyGraph) CreateObjectProperty(id string) error {
	return nil
}

func (ont *OntologyGraph) CreateDataProperty(id string) error {
	return nil
}

const (
	OWLOntology        string = "http://www.w3.org/2002/07/owl#Ontology"
	OWLVersionInfo     string = "http://www.w3.org/2002/07/owl#versionInfo"
	OWLImports         string = "http://www.w3.org/2002/07/owl#imports"
	OWLNamedIndividual string = "http://www.w3.org/2002/07/owl#NamedIndividual"
	OWLInverseOf       string = "http://www.w3.org/2002/07/owl#inverseOf"

	RDFType string = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"

	RDFSComment       string = "http://www.w3.org/2000/01/rdf-schema#comment"
	RDFSLabel         string = "http://www.w3.org/2000/01/rdf-schema#label"
	RDFSSubClassOf    string = "http://www.w3.org/2000/01/rdf-schema#subClassOf"
	RDFSSubPropertyOf string = "http://www.w3.org/2000/01/rdf-schema#subPropertyOf"
)
