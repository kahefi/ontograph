package ontograph

// An OntologyDatatype represents an ontological data type (e.g. strings, integers, ...).
type OntologyDatatype struct {
	URI     string
	Label   map[string]string
	Comment map[string]string
}

// GetURI returns the URI of the data type.
func (dt *OntologyDatatype) GetURI() string {
	return dt.URI
}

// ToTriples converts the datatype into a set of triples.
func (dt *OntologyDatatype) ToTriples() []Triple {
	trps := []Triple{}
	subj := NewResourceTerm(dt.URI)

	// Define datatype definition triple
	trps = append(trps, Triple{
		Subject:   subj,
		Predicate: NewResourceTerm(RDFType),
		Object:    NewResourceTerm(RDFSDatatype),
	})

	// Add labels
	for lang, label := range dt.Label {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSLabel),
			Object:    NewLiteralTerm(label, lang, ""),
		})
	}
	// Add comments
	for lang, comment := range dt.Comment {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSComment),
			Object:    NewLiteralTerm(comment, lang, ""),
		})
	}
	// Done, return triples
	return trps
}
