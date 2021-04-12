package ontograph

// An OntologyObjectProperty represents an object property from an ontology.
type OntologyObjectProperty struct {
    uri string
}

// NewObjectProperty creates a new ontological object property.
func NewObjectProperty(uri string) *OntologyObjectProperty {
    prop := OntologyObjectProperty{
        uri: uri,
    }
    return &prop
}

// GetURI retrieves the URI of the object property.
func (prop *OntologyObjectProperty) GetURI() string {
    return prop.uri
}
