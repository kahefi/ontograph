package ontograph

// An OntologyClass represents a class from an ontology.
type OntologyClass struct {
	URI          string
	EquivalentTo []string
	SubClassOf   []string
	DisjointWith []string
	Label        map[string]string
	Comment      map[string]string
}

// GetURI returns the URI of the class.
func (class *OntologyClass) GetURI() string {
	return class.URI
}

// ToTriples converts the class into a set of triples.
func (class *OntologyClass) ToTriples() []Triple {
	trps := []Triple{}
	subj := NewResourceTerm(class.URI)
	// Define class definition triple
	trps = append(trps, Triple{
		Subject:   subj,
		Predicate: NewResourceTerm(RDFType),
		Object:    NewResourceTerm(OWLClass),
	})
	// Add equivalentTo triples
	for _, uri := range class.EquivalentTo {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(OWLEquivalentClass),
			Object:    NewResourceTerm(uri),
		})
	}
	// Add subclassOf triples
	for _, uri := range class.SubClassOf {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSSubClassOf),
			Object:    NewResourceTerm(uri),
		})
	}
	// Add subclassOf triples
	for _, uri := range class.DisjointWith {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(OWLDisjointWith),
			Object:    NewResourceTerm(uri),
		})
	}
	// Add labels
	for lang, label := range class.Label {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSLabel),
			Object:    NewLiteralTerm(label, lang, ""),
		})
	}
	// Add comments
	for lang, comment := range class.Comment {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSComment),
			Object:    NewLiteralTerm(comment, lang, ""),
		})
	}
	// Done, return triples
	return trps
}
