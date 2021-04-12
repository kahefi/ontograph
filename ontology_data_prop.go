package ontograph

// An OntologyDataProperty represents an data property from an ontology.
type OntologyDataProperty struct {
    uri string
}

// NewDataProperty creates a new ontological data property.
func NewDataProperty(uri string) *OntologyDataProperty {
    prop := OntologyDataProperty{
        uri: uri,
    }
    return &prop
}

// GetURI retrieves the URI of the data property.
func (prop *OntologyDataProperty) GetURI() string {
    return prop.uri
}
