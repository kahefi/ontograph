package ontograph

// An OntologyObjectProperty represents an object property from an ontology.
type OntologyObjectProperty struct {
	URI                 string
	EquivalentTo        []string
	SubPropertyOf       []string
	InverseOf           []string
	Domains             []string
	Ranges              []string
	DisjointWith        []string
	IsFunctional        bool
	IsInverseFunctional bool
	IsTransitive        bool
	IsSymmetric         bool
	IsAsymmetric        bool
	IsReflexive         bool
	IsIrreflexive       bool
	Label               map[string]string
	Comment             map[string]string
}

// GetURI returns the URI of the object property.
func (prop *OntologyObjectProperty) GetURI() string {
	return prop.URI
}

// ToTriples converts the object property into a set of triples.
func (prop *OntologyObjectProperty) ToTriples() []Triple {
	trps := []Triple{}
	subj := NewResourceTerm(prop.URI)

	// Define property definition triple
	trps = append(trps, Triple{
		Subject:   subj,
		Predicate: NewResourceTerm(RDFType),
		Object:    NewResourceTerm(OWLObjectProperty),
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
	// Add inverseOf triples
	for _, uri := range prop.InverseOf {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(OWLInverseOf),
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
	if prop.IsInverseFunctional {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFType),
			Object:    NewResourceTerm(OWLInverseFunctionalProperty),
		})
	}
	if prop.IsTransitive {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFType),
			Object:    NewResourceTerm(OWLTransitiveProperty),
		})
	}
	if prop.IsSymmetric {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFType),
			Object:    NewResourceTerm(OWLSymmetricProperty),
		})
	}
	if prop.IsAsymmetric {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFType),
			Object:    NewResourceTerm(OWLAsymmetricProperty),
		})
	}
	if prop.IsReflexive {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFType),
			Object:    NewResourceTerm(OWLReflexiveProperty),
		})
	}
	if prop.IsIrreflexive {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFType),
			Object:    NewResourceTerm(OWLIrreflexiveProperty),
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
