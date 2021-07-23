package ontograph

import (
	"errors"
	"strings"
)

// An OntologyGraph represents an ontology backed by a grapg store using a higher abstraction level.
type OntologyGraph struct {
	graph   GraphStore
	label   map[string]string
	comment map[string]string
}

// InitOntologyGraph initializes a new ontology on the given graph store as backend and adds
// the appropriate definition triples to the store.
// This constructor only works when the ontology definitions do not exist (use LoadOntologyGraph
// otherwise).
func InitOntologyGraph(graph GraphStore) (*OntologyGraph, error) {
	// Check if ontology already exists
	trp, err := graph.GetFirstMatch(
		NewResourceTerm(graph.GetURI()).String(),
		NewResourceTerm(RDFType).String(),
		NewResourceTerm(OWLOntology).String(),
	)
	if err != nil {
		return nil, err
	}
	if trp != nil {
		return nil, ErrOntologyAlreadyExists
	}
	// Add ontology definition triples
	err = graph.AddTripleUnchecked(Triple{
		Subject:   NewResourceTerm(graph.GetURI()),
		Predicate: NewResourceTerm(RDFType),
		Object:    NewResourceTerm(OWLOntology),
	})
	if err != nil {
		return nil, err
	}
	// Success
	ont := OntologyGraph{
		graph:   graph,
		label:   map[string]string{},
		comment: map[string]string{},
	}
	return &ont, nil
}

// LoadOntologyGraph loads the ontology using the given graph store as backend.
// This constructor only works when the ontology definitions already exist (use InitOntologyGraph
// otherwise to create them).
func LoadOntologyGraph(graph GraphStore) (*OntologyGraph, error) {
	// Check if ontology does exists
	trp, err := graph.GetFirstMatch(
		NewResourceTerm(graph.GetURI()).String(),
		NewResourceTerm(RDFType).String(),
		NewResourceTerm(OWLOntology).String(),
	)
	if err != nil {
		return nil, err
	}
	if trp == nil {
		return nil, ErrOntologyNotFound
	}
	// Success
	ont := OntologyGraph{
		graph:   graph,
		label:   map[string]string{},
		comment: map[string]string{},
	}

	// Retrieve labels (if available)
	trps, err := ont.graph.GetAllMatches(
		NewResourceTerm(ont.GetURI()).String(),
		NewResourceTerm(RDFSLabel).String(),
		"",
	)
	if err != nil {
		return nil, err
	}
	for _, trp := range trps {
		ont.label[trp.Object.Language()] = trp.Object.Value()
	}
	// Retrieve comments (if available)
	trps, err = ont.graph.GetAllMatches(
		NewResourceTerm(ont.GetURI()).String(),
		NewResourceTerm(RDFSComment).String(),
		"",
	)
	if err != nil {
		return nil, err
	}
	for _, trp := range trps {
		ont.comment[trp.Object.Language()] = trp.Object.Value()
	}

	return &ont, nil
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

// SetLabel sets the ontology label for the specified language code.
// Any previous set label for the language will be removed.
// If `label` is empty, the label for the language code will be removed.
func (ont *OntologyGraph) SetLabel(label, lang string) error {
	// Check if previous label must be removed
	if val, ok := ont.label[lang]; ok {
		if err := ont.graph.DeleteTripleUnchecked(Triple{
			Subject:   NewResourceTerm(ont.GetURI()),
			Predicate: NewResourceTerm(RDFSLabel),
			Object:    NewLiteralTerm(val, lang, ""),
		}); err != nil {
			return err
		}
	}
	// We are done if a label is to be removed
	if label == "" {
		return nil
	}
	// We can add the new label triple
	if err := ont.graph.AddTripleUnchecked(Triple{
		Subject:   NewResourceTerm(ont.GetURI()),
		Predicate: NewResourceTerm(RDFSLabel),
		Object:    NewLiteralTerm(label, lang, ""),
	}); err != nil {
		return err
	}
	// Sync the local label map
	ont.label[lang] = label
	return nil
}

// GetLabel retrieves the ontology label for the specified language code.
func (ont *OntologyGraph) GetLabel(lang string) string {
	return ont.label[lang]
}

// SetComment sets the ontology comment for the specified language code.
// Any previous set comment for the language will be removed.
// If `comment` is empty, the comment for the language code will be removed.
func (ont *OntologyGraph) SetComment(comment, lang string) error {
	// Check if previous comment must be removed
	if val, ok := ont.comment[lang]; ok {
		if err := ont.graph.DeleteTripleUnchecked(Triple{
			Subject:   NewResourceTerm(ont.GetURI()),
			Predicate: NewResourceTerm(RDFSComment),
			Object:    NewLiteralTerm(val, lang, ""),
		}); err != nil {
			return err
		}
	}
	// We are done if a comment is to be removed
	if comment == "" {
		return nil
	}
	// We can add the new comment triple
	if err := ont.graph.AddTripleUnchecked(Triple{
		Subject:   NewResourceTerm(ont.GetURI()),
		Predicate: NewResourceTerm(RDFSComment),
		Object:    NewLiteralTerm(comment, lang, ""),
	}); err != nil {
		return err
	}
	// Sync the local comment map
	ont.comment[lang] = comment
	return nil
}

// GetComment retrieves the ontology comment for the specified language code.
func (ont *OntologyGraph) GetComment(lang string) string {
	return ont.comment[lang]
}

// // AddClass adds the given class to the ontology.
// func (ont *OntologyGraph) AddClass(class OntologyClass) error {
// 	return ont.graph.AddTriplesUnchecked(class.ToTriples())
// }

// func (ont *OntologyGraph) AddObjectProperty(prop OntologyObjectProperty) error {
// 	return ont.graph.AddTriplesUnchecked(prop.ToTriples())
// }

// func (ont *OntologyGraph) AddDataProperty(prop OntologyDataProperty) error {
// 	return ont.graph.AddTriplesUnchecked(prop.ToTriples())
// }

// func (ont *OntologyGraph) AddDataTypeProperty(prop OntologyDatatype) error {
// 	return ont.graph.AddTriplesUnchecked(indiv.ToTriples())
// }

// func (ont *OntologyGraph) AddIndividual(indiv OntologyIndividual) error {
// 	return ont.graph.AddTriplesUnchecked(indiv.ToTriples())
// }

// UpsertResource stores the given resource into the graph.
// Any already stored version of the resources will be deleted.
func (ont *OntologyGraph) UpsertResource(resource OntologyResource) error {
	uri := resource.GetURI()
	if uri[:strings.LastIndex(uri, "#")] != ont.graph.GetURI() {
		return ErrResourceDoesNotBelongToGraph
	}
	if err := ont.DeleteResource(resource.GetURI()); err != nil {
		return err
	}
	return ont.graph.AddTriplesUnchecked(resource.ToTriples())
}

// DeleteResource removes the resource and all its references from the graph.
func (ont *OntologyGraph) DeleteResource(uri string) error {
	// First delete all triples which have the URI as subject
	err := ont.graph.DeleteAllMatches(NewResourceTerm(uri).String(), "", "")
	if err != nil {
		return err
	}
	// Second delete all triples that reference the URI in their object
	return ont.graph.DeleteAllMatches("", "", NewResourceTerm(uri).String())
}

// GetClass retrieves the class with the specified URI from the graph.
func (ont *OntologyGraph) GetClass(uri string) (OntologyClass, error) {
	// Retrieve all relevant triples
	trps, err := ont.graph.GetAllMatches(NewResourceTerm(uri).String(), "", "")
	if err != nil {
		return OntologyClass{}, err
	}
	// Parse triples into the class structure
	class := OntologyClass{
		URI:          "",
		EquivalentTo: []string{},
		SubClassOf:   []string{},
		DisjointWith: []string{},
		Label:        map[string]string{},
		Comment:      map[string]string{},
	}
	for _, trp := range trps {
		if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLClass) {
			class.URI = uri
		} else if trp.Predicate == NewResourceTerm(OWLEquivalentClass) {
			class.EquivalentTo = append(class.EquivalentTo, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFSSubClassOf) {
			class.SubClassOf = append(class.SubClassOf, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(OWLDisjointWith) {
			class.DisjointWith = append(class.DisjointWith, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFSLabel) {
			class.Label[trp.Object.Language()] = trp.Object.Value()
		} else if trp.Predicate == NewResourceTerm(RDFSComment) {
			class.Comment[trp.Object.Language()] = trp.Object.Value()
		}
	}
	// If no URI was set, the requested URI is not a class
	if class.URI == "" {
		return OntologyClass{}, ErrResourceNotFound
	}
	return class, nil
}

// GetObjectProperty retrieves the object property with the specified URI from the graph.
func (ont *OntologyGraph) GetObjectProperty(uri string) (OntologyObjectProperty, error) {
	// Retrieve all relevant triples
	trps, err := ont.graph.GetAllMatches(NewResourceTerm(uri).String(), "", "")
	if err != nil {
		return OntologyObjectProperty{}, err
	}
	// Parse triples into the object property structure
	prop := OntologyObjectProperty{
		URI:                 "",
		EquivalentTo:        []string{},
		SubPropertyOf:       []string{},
		InverseOf:           []string{},
		Domains:             []string{},
		Ranges:              []string{},
		DisjointWith:        []string{},
		IsFunctional:        false,
		IsInverseFunctional: false,
		IsTransitive:        false,
		IsSymmetric:         false,
		IsAsymmetric:        false,
		IsReflexive:         false,
		IsIrreflexive:       false,
		Label:               map[string]string{},
		Comment:             map[string]string{},
	}
	for _, trp := range trps {
		if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLObjectProperty) {
			prop.URI = uri
		} else if trp.Predicate == NewResourceTerm(OWLEquivalentProperty) {
			prop.EquivalentTo = append(prop.EquivalentTo, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFSSubPropertyOf) {
			prop.SubPropertyOf = append(prop.SubPropertyOf, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(OWLInverseOf) {
			prop.InverseOf = append(prop.InverseOf, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFSDomain) {
			prop.Domains = append(prop.Domains, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFSRange) {
			prop.Ranges = append(prop.Ranges, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(OWLPropertyDisjointWith) {
			prop.DisjointWith = append(prop.DisjointWith, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLFunctionalProperty) {
			prop.IsFunctional = true
		} else if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLInverseFunctionalProperty) {
			prop.IsInverseFunctional = true
		} else if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLTransitiveProperty) {
			prop.IsTransitive = true
		} else if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLSymmetricProperty) {
			prop.IsSymmetric = true
		} else if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLAsymmetricProperty) {
			prop.IsAsymmetric = true
		} else if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLReflexiveProperty) {
			prop.IsReflexive = true
		} else if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLIrreflexiveProperty) {
			prop.IsIrreflexive = true
		} else if trp.Predicate == NewResourceTerm(RDFSLabel) {
			prop.Label[trp.Object.Language()] = trp.Object.Value()
		} else if trp.Predicate == NewResourceTerm(RDFSComment) {
			prop.Comment[trp.Object.Language()] = trp.Object.Value()
		}
	}
	// If no URI was set, the requested URI is not an object property
	if prop.URI == "" {
		return OntologyObjectProperty{}, ErrResourceNotFound
	}
	return prop, nil
}

// GetDataProperty retrieves the data property with the specified URI from the graph.
func (ont *OntologyGraph) GetDataProperty(uri string) (OntologyDataProperty, error) {
	// Retrieve all relevant triples
	trps, err := ont.graph.GetAllMatches(NewResourceTerm(uri).String(), "", "")
	if err != nil {
		return OntologyDataProperty{}, err
	}
	// Parse triples into the object property structure
	prop := OntologyDataProperty{
		URI:           "",
		EquivalentTo:  []string{},
		SubPropertyOf: []string{},
		Domains:       []string{},
		Ranges:        []string{},
		DisjointWith:  []string{},
		IsFunctional:  false,
		Label:         map[string]string{},
		Comment:       map[string]string{},
	}
	for _, trp := range trps {
		if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLDatatypeProperty) {
			prop.URI = uri
		} else if trp.Predicate == NewResourceTerm(OWLEquivalentProperty) {
			prop.EquivalentTo = append(prop.EquivalentTo, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFSSubPropertyOf) {
			prop.SubPropertyOf = append(prop.SubPropertyOf, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFSDomain) {
			prop.Domains = append(prop.Domains, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFSRange) {
			prop.Ranges = append(prop.Ranges, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(OWLPropertyDisjointWith) {
			prop.DisjointWith = append(prop.DisjointWith, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLFunctionalProperty) {
			prop.IsFunctional = true
		} else if trp.Predicate == NewResourceTerm(RDFSLabel) {
			prop.Label[trp.Object.Language()] = trp.Object.Value()
		} else if trp.Predicate == NewResourceTerm(RDFSComment) {
			prop.Comment[trp.Object.Language()] = trp.Object.Value()
		}
	}
	// If no URI was set, the requested URI is not an object property
	if prop.URI == "" {
		return OntologyDataProperty{}, ErrResourceNotFound
	}
	return prop, nil
}

// GetDatatype retrieves the datatype with the specified URI from the graph.
func (ont *OntologyGraph) GetDatatype(uri string) (OntologyDatatype, error) {
	// Retrieve all relevant triples
	trps, err := ont.graph.GetAllMatches(NewResourceTerm(uri).String(), "", "")
	if err != nil {
		return OntologyDatatype{}, err
	}
	// Parse triples into the object property structure
	prop := OntologyDatatype{
		URI:     "",
		Label:   map[string]string{},
		Comment: map[string]string{},
	}
	for _, trp := range trps {
		if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(RDFSDatatype) {
			prop.URI = uri
		} else if trp.Predicate == NewResourceTerm(RDFSLabel) {
			prop.Label[trp.Object.Language()] = trp.Object.Value()
		} else if trp.Predicate == NewResourceTerm(RDFSComment) {
			prop.Comment[trp.Object.Language()] = trp.Object.Value()
		}
	}
	// If no URI was set, the requested URI is not an object property
	if prop.URI == "" {
		return OntologyDatatype{}, ErrResourceNotFound
	}
	return prop, nil
}

// GetIndividual retrieves the individual with the specified URI from the graph.
func (ont *OntologyGraph) GetIndividual(uri string) (OntologyIndividual, error) {
	// Retrieve all relevant triples
	trps, err := ont.graph.GetAllMatches(NewResourceTerm(uri).String(), "", "")
	if err != nil {
		return OntologyIndividual{}, err
	}
	// Parse triples into the individual structure
	indiv := OntologyIndividual{
		URI:              "",
		Types:            []string{},
		SameIndividualAs: []string{},
		ObjectProperties: map[string][]string{},
		DataProperties:   map[string][]GenericLiteral{},
		Label:            map[string]string{},
		Comment:          map[string]string{},
	}
	for _, trp := range trps {
		if trp.Predicate == NewResourceTerm(RDFType) && trp.Object == NewResourceTerm(OWLNamedIndividual) {
			indiv.URI = uri
		} else if trp.Predicate == NewResourceTerm(RDFType) {
			indiv.Types = append(indiv.Types, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(OWLSameAs) {
			indiv.SameIndividualAs = append(indiv.SameIndividualAs, trp.Object.Value())
		} else if trp.Predicate == NewResourceTerm(RDFSLabel) {
			indiv.Label[trp.Object.Language()] = trp.Object.Value()
		} else if trp.Predicate == NewResourceTerm(RDFSComment) {
			indiv.Comment[trp.Object.Language()] = trp.Object.Value()
		} else {
			obj := trp.Object
			prop := trp.Predicate.Value()
			if obj.IsResource() {
				indiv.ObjectProperties[prop] = append(indiv.ObjectProperties[prop], obj.Value())
			} else if obj.IsLiteral() {
				indiv.DataProperties[prop] = append(indiv.DataProperties[prop], *NewGenericLiteral(obj))
			}
		}
	}
	// If no URI was set, the requested URI is not an object property
	if indiv.URI == "" {
		return OntologyIndividual{}, ErrResourceNotFound
	}
	return indiv, nil
}

// *****************
// * Shared Errors *
// *****************

// ErrOntologyNotFound is raised when an ontology does not exist.
var ErrOntologyNotFound error = errors.New("The requested ontology does not exist")

// ErrOntologyAlreadyExists is raised when an ontology already exists.
var ErrOntologyAlreadyExists error = errors.New("The requested ontology already exists")

// ErrResourceNotFound is raised on conflict errors when a triple already exists (i.e. adding triples).
var ErrResourceNotFound error = errors.New("The requested ontology resource does not exist in the graph")

// ErrResourceDoesNotBelongToGraph is raised when a resource is attempted to be added to the graph, but their base URIs do not match.
var ErrResourceDoesNotBelongToGraph error = errors.New("The URI of the resource does not match the URI of the graph")
