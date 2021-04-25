package ontograph

// An OntologyIndividual represents an individual from an ontology.
type OntologyIndividual struct {
	URI              string
	Types            []string
	SameIndividualAs []string
	ObjectProperties map[string][]string
	DataProperties   map[string][]GenericLiteral
	Label            map[string]string
	Comment          map[string]string
}

// GetURI returns the URI of the individual.
func (indiv *OntologyIndividual) GetURI() string {
	return indiv.URI
}

func (indiv *OntologyIndividual) AddObjectProperty(prop, target string) {
	if indiv.ObjectProperties == nil {
		indiv.ObjectProperties = map[string][]string{}
	}
	indiv.ObjectProperties[prop] = append(indiv.ObjectProperties[prop], target)
}

func (indiv *OntologyIndividual) AddDataProperty(prop string, target GenericLiteral) {
	if indiv.DataProperties == nil {
		indiv.DataProperties = map[string][]GenericLiteral{}
	}
	indiv.DataProperties[prop] = append(indiv.DataProperties[prop], target)
}

// ToTriples converts the individual into a set of triples.
func (indiv *OntologyIndividual) ToTriples() []Triple {
	trps := []Triple{}
	subj := NewResourceTerm(indiv.URI)

	// Define individual definition triple
	trps = append(trps, Triple{
		Subject:   subj,
		Predicate: NewResourceTerm(RDFType),
		Object:    NewResourceTerm(OWLNamedIndividual),
	})

	// Add type triples
	for _, uri := range indiv.Types {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFType),
			Object:    NewResourceTerm(uri),
		})
	}
	// Add SameIndividualAs triples
	for _, uri := range indiv.SameIndividualAs {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(OWLSameAs),
			Object:    NewResourceTerm(uri),
		})
	}

	// Add object property relations
	for propUri, targets := range indiv.ObjectProperties {
		for _, uri := range targets {
			trps = append(trps, Triple{
				Subject:   subj,
				Predicate: NewResourceTerm(propUri),
				Object:    NewResourceTerm(uri),
			})
		}
	}
	// Add data property relations
	for propUri, targets := range indiv.DataProperties {
		for _, lit := range targets {
			trps = append(trps, Triple{
				Subject:   subj,
				Predicate: NewResourceTerm(propUri),
				Object:    lit.Term(),
			})
		}
	}

	// Add labels
	for lang, label := range indiv.Label {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSLabel),
			Object:    NewLiteralTerm(label, lang, ""),
		})
	}
	// Add comments
	for lang, comment := range indiv.Comment {
		trps = append(trps, Triple{
			Subject:   subj,
			Predicate: NewResourceTerm(RDFSComment),
			Object:    NewLiteralTerm(comment, lang, ""),
		})
	}
	// Done, return triples
	return trps
}
