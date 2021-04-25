package ontograph

// An OntologyDataProperty represents a data property from an ontology.
type OntologyDataProperty struct {
	URI           string
	EquivalentTo  []string
	SubPropertyOf []string
	Domains       []string
	Ranges        []string
	DisjointWith  []string
	IsFunctional  bool
	Label         map[string]string
	Comment       map[string]string
}

// GetURI returns the URI of the data property.
func (prop *OntologyDataProperty) GetURI() string {
	return prop.URI
}

// ToTriples converts the data property into a set of triples.
func (prop *OntologyDataProperty) ToTriples() []Triple {
	trps := []Triple{}
	subj := NewResourceTerm(prop.URI)

	// Define property definition triple
	trps = append(trps, Triple{
		Subject:   subj,
		Predicate: NewResourceTerm(RDFType),
		Object:    NewResourceTerm(OWLDatatypeProperty),
	})
	// Add equivalentTo triples
	for _, uri := range prop.EquivalentTo {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(OWLEquivalentProperty),
			Object:    NewResourceTerm(uri),
		})
	}
	// Add subPropertyOf triples
	for _, uri := range prop.SubPropertyOf {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSSubPropertyOf),
			Object:    NewResourceTerm(uri),
		})
	}
	// Add domain triples
	for _, uri := range prop.Domains {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSDomain),
			Object:    NewResourceTerm(uri),
		})
	}
	// Add range triples
	for _, uri := range prop.Ranges {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSRange),
			Object:    NewResourceTerm(uri),
		})
	}
	// Add disjointWith triples
	for _, uri := range prop.DisjointWith {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(OWLPropertyDisjointWith),
			Object:    NewResourceTerm(uri),
		})
	}

	// Add logical property triples
	if prop.IsFunctional {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFType),
			Object:    NewResourceTerm(OWLFunctionalProperty),
		})
	}

	// Add labels
	for lang, label := range prop.Label {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSLabel),
			Object:    NewLiteralTerm(label, lang, ""),
		})
	}
	// Add comments
	for lang, comment := range prop.Comment {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSComment),
			Object:    NewLiteralTerm(comment, lang, ""),
		})
	}
	// Done, return triples
	return trps
}
