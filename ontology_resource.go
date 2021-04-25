package ontograph

// An OntologyResource abstracts a class, object property, data property, datatype property or an individual to a general resource.
type OntologyResource interface {
	GetURI() string
	ToTriples() []Triple
}
