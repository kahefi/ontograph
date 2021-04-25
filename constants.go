package ontograph

// Static URIS used in ontologies (RDF, RDFS and OWL)
const (
	OWLOntology                  string = "http://www.w3.org/2002/07/owl#Ontology"
	OWLVersionInfo               string = "http://www.w3.org/2002/07/owl#versionInfo"
	OWLImports                   string = "http://www.w3.org/2002/07/owl#imports"
	OWLInverseOf                 string = "http://www.w3.org/2002/07/owl#inverseOf"
	OWLClass                     string = "http://www.w3.org/2002/07/owl#Class"
	OWLEquivalentClass           string = "http://www.w3.org/2002/07/owl#equivalentClass"
	OWLDisjointWith              string = "http://www.w3.org/2002/07/owl#disjointWith"
	OWLObjectProperty            string = "http://www.w3.org/2002/07/owl#ObjectProperty"
	OWLFunctionalProperty        string = "http://www.w3.org/2002/07/owl#FunctionalProperty"
	OWLInverseFunctionalProperty string = "http://www.w3.org/2002/07/owl#InverseFunctionalProperty"
	OWLSymmetricProperty         string = "http://www.w3.org/2002/07/owl#SymmetricProperty"
	OWLAsymmetricProperty        string = "http://www.w3.org/2002/07/owl#AsymmetricProperty"
	OWLTransitiveProperty        string = "http://www.w3.org/2002/07/owl#TransitiveProperty"
	OWLReflexiveProperty         string = "http://www.w3.org/2002/07/owl#ReflexiveProperty"
	OWLIrreflexiveProperty       string = "http://www.w3.org/2002/07/owl#IrreflexiveProperty"
	OWLPropertyDisjointWith      string = "http://www.w3.org/2002/07/owl#propertyDisjointWith"
	OWLEquivalentProperty        string = "http://www.w3.org/2002/07/owl#equivalentProperty"
	OWLDatatypeProperty          string = "http://www.w3.org/2002/07/owl#DatatypeProperty"
	OWLNamedIndividual           string = "http://www.w3.org/2002/07/owl#NamedIndividual"
	OWLSameAs                    string = "http://www.w3.org/2002/07/owl#sameAs"

	RDFType string = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"

	RDFSComment       string = "http://www.w3.org/2000/01/rdf-schema#comment"
	RDFSLabel         string = "http://www.w3.org/2000/01/rdf-schema#label"
	RDFSSubClassOf    string = "http://www.w3.org/2000/01/rdf-schema#subClassOf"
	RDFSSubPropertyOf string = "http://www.w3.org/2000/01/rdf-schema#subPropertyOf"
	RDFSDomain        string = "http://www.w3.org/2000/01/rdf-schema#domain"
	RDFSRange         string = "http://www.w3.org/2000/01/rdf-schema#range"
	RDFSDatatype      string = "http://www.w3.org/2000/01/rdf-schema#Datatype"

	XSDString   string = "http://www.w3.org/2001/XMLSchema#string"
	XSDInteger  string = "http://www.w3.org/2001/XMLSchema#integer"
	XSDDouble   string = "http://www.w3.org/2001/XMLSchema#double"
	XSDFloat    string = "http://www.w3.org/2001/XMLSchema#float"
	XSDBoolean  string = "http://www.w3.org/2001/XMLSchema#boolean"
	XSDDate     string = "http://www.w3.org/2001/XMLSchema#date"
	XSDTime     string = "http://www.w3.org/2001/XMLSchema#time"
	XSDDateTime string = "http://www.w3.org/2001/XMLSchema#dateTime"
)