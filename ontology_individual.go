package ontograph

// An OntologyIndividual represents an individual from an ontology.
type OntologyIndividual struct {
    uri string
}

// NewIndividual creates a new ontologcial individual.
func NewIndividual(uri string) *OntologyIndividual {
    indiv := OntologyIndividual{
        uri: uri,
    }
    return &indiv
}

// GetURI retrieves the URI of the individual.
func (indiv *OntologyIndividual) GetURI() string {
    return indiv.uri
}
