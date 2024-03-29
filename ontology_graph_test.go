package ontograph_test

import (
    "fmt"

    "github.com/lithammer/shortuuid"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"

    . "github.com/kahefi/ontograph"
)

var _ = Describe("OntologyGraph", func() {
    var testUri string
    var graph GraphStore
    var ont *OntologyGraph

    checkIndividuals := func(indiv1, indiv2 OntologyIndividual) {
        Expect(indiv1.URI).To(Equal(indiv2.URI))
        Expect(indiv1.Types).To(ConsistOf(indiv2.Types))
        Expect(indiv1.SameIndividualAs).To(ConsistOf(indiv2.SameIndividualAs))
        // Check object properties
        for uri := range indiv1.ObjectProperties {
            Expect(indiv1.ObjectProperties[uri]).To(ConsistOf(indiv2.ObjectProperties[uri]))
        }
        for uri := range indiv2.ObjectProperties {
            Expect(indiv2.ObjectProperties[uri]).To(ConsistOf(indiv1.ObjectProperties[uri]))
        }
        // Check data properties
        for uri := range indiv1.DataProperties {
            Expect(indiv1.DataProperties[uri]).To(ConsistOf(indiv2.DataProperties[uri]))
        }
        for uri := range indiv2.DataProperties {
            Expect(indiv2.DataProperties[uri]).To(ConsistOf(indiv1.DataProperties[uri]))
        }
        // Check labels and comments
        Expect(indiv1.Label).To(Equal(indiv2.Label))
        Expect(indiv1.Comment).To(Equal(indiv2.Comment))
    }

    BeforeEach(func() {
        // Setup ontology
        testUri = fmt.Sprintf("https://www.ontograph.com/test-%s", shortuuid.New())
        graph = NewMemoryStore(testUri)
        var err error
        ont, err = InitOntologyGraph(graph)
        Expect(err).NotTo(HaveOccurred())
    })

    AfterEach(func() {
    })

    Describe("Loading the ontology graph", func() {
        It("should match the initialisation", func() {
            ont, err := LoadOntologyGraph(graph)
            Expect(err).NotTo(HaveOccurred())
            Expect(ont.GetURI()).To(Equal(testUri))
        })
    })

    Describe("Setting ontology labels and comments", func() {
        It("should have added the expected labels", func() {
            err := ont.SetLabel("label", "en")
            Expect(err).NotTo(HaveOccurred())
            err = ont.SetLabel("should not appear", "de")
            Expect(err).NotTo(HaveOccurred())
            err = ont.SetLabel("titel", "de")
            Expect(err).NotTo(HaveOccurred())
            err = ont.SetLabel("42", "")
            Expect(err).NotTo(HaveOccurred())
            // Check that labels were set correctly
            Expect(ont.GetLabel("de")).To(Equal("titel"))
            Expect(ont.GetLabel("en")).To(Equal("label"))
            Expect(ont.GetLabel("")).To(Equal("42"))

            // Reload ontology and check labels again to make sure they were stored correctly
            ont, err = LoadOntologyGraph(graph)
            Expect(err).NotTo(HaveOccurred())
            Expect(ont.GetURI()).To(Equal(testUri))
            Expect(ont.GetLabel("de")).To(Equal("titel"))
            Expect(ont.GetLabel("en")).To(Equal("label"))
            Expect(ont.GetLabel("")).To(Equal("42"))
        })
        It("should have added the expected comments", func() {
            err := ont.SetComment("comment", "en")
            Expect(err).NotTo(HaveOccurred())
            err = ont.SetComment("should not appear", "de")
            Expect(err).NotTo(HaveOccurred())
            err = ont.SetComment("kommentar", "de")
            Expect(err).NotTo(HaveOccurred())
            err = ont.SetComment("42", "")
            Expect(err).NotTo(HaveOccurred())
            // Check that comments were set correctly
            Expect(ont.GetComment("de")).To(Equal("kommentar"))
            Expect(ont.GetComment("en")).To(Equal("comment"))
            Expect(ont.GetComment("")).To(Equal("42"))

            // Reload ontology and check comments again to make sure they were stored correctly
            ont, err = LoadOntologyGraph(graph)
            Expect(err).NotTo(HaveOccurred())
            Expect(ont.GetURI()).To(Equal(testUri))
            Expect(ont.GetComment("de")).To(Equal("kommentar"))
            Expect(ont.GetComment("en")).To(Equal("comment"))
            Expect(ont.GetComment("")).To(Equal("42"))
        })
    })

    Describe("Retrieving the version of the ontology", func() {
        When("a version was set", func() {
            BeforeEach(func() {
                err := ont.SetVersion("0.42.1-get")
                Expect(err).NotTo(HaveOccurred())
            })
            It("should return the expected version", func() {
                version, err := ont.GetVersion()
                Expect(err).NotTo(HaveOccurred())
                Expect(version).To(Equal("0.42.1-get"))
            })
        })
        When("no version was set", func() {
            It("should return an empty string", func() {
                version, err := ont.GetVersion()
                Expect(err).NotTo(HaveOccurred())
                Expect(version).To(Equal(""))
            })
        })
    })

    Describe("Setting the version of the ontology", func() {
        It("should have added the version information to the ontology", func() {
            By("not returning an error")
            err := ont.SetVersion("0.42.1-set")
            Expect(err).NotTo(HaveOccurred())
            By("containing the expected version")
            version, err := ont.GetVersion()
            Expect(err).NotTo(HaveOccurred())
            Expect(version).To(Equal("0.42.1-set"))
        })
    })

    Describe("Retrieving the imported ontologies", func() {
        When("imports have been defined", func() {
            var testImports []string
            BeforeEach(func() {
                testImports = []string{"http://abc-1.com", "https://abc-2.com", "http://test.de/42"}
                for _, uri := range testImports {
                    err := ont.AddImport(uri)
                    Expect(err).NotTo(HaveOccurred())
                }
            })
            It("should return the expected list of URIs", func() {
                uris, err := ont.GetImports()
                Expect(err).NotTo(HaveOccurred())
                Expect(uris).To(ConsistOf(testImports))
            })
        })
        When("no imports have been defined", func() {
            It("should return an empty list", func() {
                uris, err := ont.GetImports()
                Expect(err).NotTo(HaveOccurred())
                Expect(uris).To(BeEmpty())
            })
        })
    })

    Describe("Adding an import to the ontology", func() {
        It("should have added the URI to the list of imports in the ontology", func() {
            err := ont.AddImport("http://abc-1.com")
            Expect(err).NotTo(HaveOccurred())
            uris, err := ont.GetImports()
            Expect(err).NotTo(HaveOccurred())
            Expect(uris).To(ContainElement("http://abc-1.com"))
        })
    })

    Describe("Adding and retrieving an ontology class", func() {
        When("the class belongs to the graph", func() {
            It("should successfully add the class to the store", func() {
                class := OntologyClass{
                    URI:          testUri + "#class",
                    EquivalentTo: []string{"http://abc.com#class2", "http://abc.com#class3"},
                    SubClassOf:   []string{"http://abc.com#parent1", "http://abc.com#parent2"},
                    DisjointWith: []string{"http://abc.com#notclass"},
                    Label:        map[string]string{"": "a label", "de": "ein title", "en": "a label"},
                    Comment:      map[string]string{"": "some comment", "de": "ein kommentar"},
                }
                err := ont.UpsertResource(&class)
                By("not raising an error")
                Expect(err).NotTo(HaveOccurred())
                By("having stored the expected class")
                retClass, err := ont.GetClass(class.URI)
                Expect(err).NotTo(HaveOccurred())
                Expect(retClass.URI).To(Equal(class.URI))
                Expect(retClass.EquivalentTo).To(ConsistOf(class.EquivalentTo))
                Expect(retClass.SubClassOf).To(ConsistOf(class.SubClassOf))
                Expect(retClass.DisjointWith).To(ConsistOf(class.DisjointWith))
                Expect(retClass.Label).To(Equal(class.Label))
                Expect(retClass.Comment).To(Equal(class.Comment))
            })
        })
        When("the class does not belong to the graph", func() {
            It("should reject the class", func() {
                class := OntologyClass{
                    URI: testUri + "x" + "#class",
                }
                err := ont.UpsertResource(&class)
                By("raising the expected error")
                Expect(err).To(Equal(ErrResourceDoesNotBelongToGraph))
                By("not having stored the class")
                _, err = ont.GetClass(class.URI)
                Expect(err).To(Equal(ErrResourceNotFound))
            })
        })
    })

    Describe("Adding and retrieving an ontology object property", func() {
        When("the object property belongs to the graph", func() {
            It("should successfully add the object property to the store", func() {
                prop := OntologyObjectProperty{
                    URI:                 testUri + "#objectprop",
                    EquivalentTo:        []string{"http://abc.com#prop2", "http://abc.com#prop3"},
                    SubPropertyOf:       []string{"http://abc.com#parent1", "http://abc.com#parent2"},
                    InverseOf:           []string{"http://abc.com#inv"},
                    Domains:             []string{"http://abc.com#class1", "http://abc.com#class2"},
                    Ranges:              []string{"http://abc.com#class3"},
                    DisjointWith:        []string{"http://abc.com#prop3"},
                    IsFunctional:        true,
                    IsInverseFunctional: true,
                    IsTransitive:        true,
                    IsSymmetric:         true,
                    IsAsymmetric:        true,
                    IsReflexive:         true,
                    IsIrreflexive:       true,
                    Label:               map[string]string{"": "a label", "de": "ein title", "en": "a label"},
                    Comment:             map[string]string{"": "some comment", "de": "ein kommentar"},
                }
                err := ont.UpsertResource(&prop)
                By("not raising an error")
                Expect(err).NotTo(HaveOccurred())
                By("having stored the expected object property")
                retProp, err := ont.GetObjectProperty(prop.URI)
                Expect(err).NotTo(HaveOccurred())
                Expect(retProp.URI).To(Equal(prop.URI))
                Expect(retProp.EquivalentTo).To(ConsistOf(prop.EquivalentTo))
                Expect(retProp.SubPropertyOf).To(ConsistOf(prop.SubPropertyOf))
                Expect(retProp.InverseOf).To(ConsistOf(prop.InverseOf))
                Expect(retProp.Domains).To(ConsistOf(prop.Domains))
                Expect(retProp.Ranges).To(ConsistOf(prop.Ranges))
                Expect(retProp.DisjointWith).To(ConsistOf(prop.DisjointWith))
                Expect(retProp.IsFunctional).To(Equal(prop.IsFunctional))
                Expect(retProp.IsInverseFunctional).To(Equal(prop.IsInverseFunctional))
                Expect(retProp.IsTransitive).To(Equal(prop.IsTransitive))
                Expect(retProp.IsSymmetric).To(Equal(prop.IsSymmetric))
                Expect(retProp.IsAsymmetric).To(Equal(prop.IsAsymmetric))
                Expect(retProp.IsReflexive).To(Equal(prop.IsReflexive))
                Expect(retProp.IsIrreflexive).To(Equal(prop.IsIrreflexive))
                Expect(retProp.Label).To(Equal(prop.Label))
                Expect(retProp.Comment).To(Equal(prop.Comment))
            })
        })
        When("the object property does not belong to the graph", func() {
            It("should reject the object property", func() {
                prop := OntologyObjectProperty{
                    URI: testUri + "x" + "#objectprop",
                }
                err := ont.UpsertResource(&prop)
                By("raising the expected error")
                Expect(err).To(Equal(ErrResourceDoesNotBelongToGraph))
                By("not having stored the object property")
                _, err = ont.GetObjectProperty(prop.URI)
                Expect(err).To(Equal(ErrResourceNotFound))
            })
        })
    })

    Describe("Adding and retrieving an ontology data property", func() {
        When("the data property belongs to the graph", func() {
            It("should successfully add the data property to the store", func() {
                prop := OntologyDataProperty{
                    URI:           testUri + "#dataprop",
                    EquivalentTo:  []string{"http://abc.com#prop2", "http://abc.com#prop3"},
                    SubPropertyOf: []string{"http://abc.com#parent1", "http://abc.com#parent2"},
                    Domains:       []string{"http://abc.com#class1"},
                    Ranges:        []string{"http://abc.com#datatype1", "http://abc.com#datatype2"},
                    DisjointWith:  []string{"http://abc.com#prop3"},
                    IsFunctional:  true,
                    Label:         map[string]string{"": "a label", "de": "ein title", "en": "a label"},
                    Comment:       map[string]string{"": "some comment", "de": "ein kommentar"},
                }
                err := ont.UpsertResource(&prop)
                By("not raising an error")
                Expect(err).NotTo(HaveOccurred())
                By("having stored the expected data property")
                retProp, err := ont.GetDataProperty(prop.URI)
                Expect(err).NotTo(HaveOccurred())
                Expect(retProp.URI).To(Equal(prop.URI))
                Expect(retProp.EquivalentTo).To(ConsistOf(prop.EquivalentTo))
                Expect(retProp.SubPropertyOf).To(ConsistOf(prop.SubPropertyOf))
                Expect(retProp.Domains).To(ConsistOf(prop.Domains))
                Expect(retProp.Ranges).To(ConsistOf(prop.Ranges))
                Expect(retProp.DisjointWith).To(ConsistOf(prop.DisjointWith))
                Expect(retProp.IsFunctional).To(Equal(prop.IsFunctional))
                Expect(retProp.Label).To(Equal(prop.Label))
                Expect(retProp.Comment).To(Equal(prop.Comment))
            })
        })
        When("the data property does not belong to the graph", func() {
            It("should reject the object property", func() {
                prop := OntologyDataProperty{
                    URI: testUri + "x" + "#dataprop",
                }
                err := ont.UpsertResource(&prop)
                By("raising the expected error")
                Expect(err).To(Equal(ErrResourceDoesNotBelongToGraph))
                By("not having stored the object property")
                _, err = ont.GetObjectProperty(prop.URI)
                Expect(err).To(Equal(ErrResourceNotFound))
            })
        })
    })

    Describe("Adding and retrieving an ontology datatype property", func() {
        When("the data property belongs to the graph", func() {
            It("should successfully add the data property to the store", func() {
                prop := OntologyDataProperty{
                    URI:           testUri + "#dataprop",
                    EquivalentTo:  []string{"http://abc.com#prop2", "http://abc.com#prop3"},
                    SubPropertyOf: []string{"http://abc.com#parent1", "http://abc.com#parent2"},
                    Domains:       []string{"http://abc.com#class1"},
                    Ranges:        []string{"http://abc.com#datatype1", "http://abc.com#datatype2"},
                    DisjointWith:  []string{"http://abc.com#prop3"},
                    IsFunctional:  true,
                    Label:         map[string]string{"": "a label", "de": "ein title", "en": "a label"},
                    Comment:       map[string]string{"": "some comment", "de": "ein kommentar"},
                }
                err := ont.UpsertResource(&prop)
                By("not raising an error")
                Expect(err).NotTo(HaveOccurred())
                By("having stored the expected data property")
                retProp, err := ont.GetDataProperty(prop.URI)
                Expect(err).NotTo(HaveOccurred())
                Expect(retProp.URI).To(Equal(prop.URI))
                Expect(retProp.EquivalentTo).To(ConsistOf(prop.EquivalentTo))
                Expect(retProp.SubPropertyOf).To(ConsistOf(prop.SubPropertyOf))
                Expect(retProp.Domains).To(ConsistOf(prop.Domains))
                Expect(retProp.Ranges).To(ConsistOf(prop.Ranges))
                Expect(retProp.DisjointWith).To(ConsistOf(prop.DisjointWith))
                Expect(retProp.IsFunctional).To(Equal(prop.IsFunctional))
                Expect(retProp.Label).To(Equal(prop.Label))
                Expect(retProp.Comment).To(Equal(prop.Comment))
            })
        })
        When("the data property does not belong to the graph", func() {
            It("should reject the object property", func() {
                prop := OntologyDataProperty{
                    URI: testUri + "x" + "#dataprop",
                }
                err := ont.UpsertResource(&prop)
                By("raising the expected error")
                Expect(err).To(Equal(ErrResourceDoesNotBelongToGraph))
                By("not having stored the data property")
                _, err = ont.GetObjectProperty(prop.URI)
                Expect(err).To(Equal(ErrResourceNotFound))
            })
        })
    })

    Describe("Adding and retrieving an ontology datatype", func() {
        When("the datatype belongs to the graph", func() {
            It("should successfully add the datatype to the store", func() {
                datatype := OntologyDatatype{
                    URI:     testUri + "#datatype",
                    Label:   map[string]string{"": "a label", "de": "ein title", "en": "a label"},
                    Comment: map[string]string{"": "some comment", "de": "ein kommentar"},
                }
                err := ont.UpsertResource(&datatype)
                By("not raising an error")
                Expect(err).NotTo(HaveOccurred())
                By("having stored the expected datatype")
                retDatatype, err := ont.GetDatatype(datatype.URI)
                Expect(err).NotTo(HaveOccurred())
                Expect(retDatatype.URI).To(Equal(datatype.URI))
                Expect(retDatatype.Label).To(Equal(datatype.Label))
                Expect(retDatatype.Comment).To(Equal(datatype.Comment))
            })
        })
        When("the data property does not belong to the graph", func() {
            It("should reject the object property", func() {
                datatype := OntologyDatatype{
                    URI: testUri + "x" + "#datatype",
                }
                err := ont.UpsertResource(&datatype)
                By("raising the expected error")
                Expect(err).To(Equal(ErrResourceDoesNotBelongToGraph))
                By("not having stored the datatype")
                _, err = ont.GetDatatype(datatype.URI)
                Expect(err).To(Equal(ErrResourceNotFound))
            })
        })
    })

    Describe("Adding and retrieving an ontology individual", func() {
        When("the individual belongs to the graph", func() {
            It("should successfully add the individual to the store", func() {
                indiv := OntologyIndividual{
                    URI:              testUri + "#indiv",
                    Types:            []string{"http://abc.com#type1", "http://abc.com#type2", "http://abc.com#type3"},
                    SameIndividualAs: []string{"http://abc.com#indiv2"},
                    Label:            map[string]string{"": "a label", "de": "ein title", "en": "a label"},
                    Comment:          map[string]string{"": "some comment", "de": "ein kommentar"},
                }
                indiv.AddObjectProperty("http://abc.com#prop1", "http://abc.com#indiv3")
                indiv.AddObjectProperty("http://abc.com#prop1", "http://abc.com#indiv4")
                indiv.AddObjectProperty("http://abc.com#prop3", "http://abc.com#indiv4")
                indiv.AddDataProperty("http://abc.com#dataprop1", XSDStringLiteral("Some string literal").Generic())
                indiv.AddDataProperty("http://abc.com#dataprop2", XSDIntegerLiteral(42).Generic())
                err := ont.UpsertResource(&indiv)
                By("not raising an error")
                Expect(err).NotTo(HaveOccurred())
                By("having stored the expected individual")
                retIndiv, err := ont.GetIndividual(indiv.URI)
                Expect(err).NotTo(HaveOccurred())
                checkIndividuals(retIndiv, indiv)
            })
        })
        When("the individual does not belong to the graph", func() {
            It("should reject the individual", func() {
                indiv := OntologyIndividual{
                    URI: testUri + "x" + "#indiv",
                }
                err := ont.UpsertResource(&indiv)
                By("raising the expected error")
                Expect(err).To(Equal(ErrResourceDoesNotBelongToGraph))
                By("not having stored the individual")
                _, err = ont.GetIndividual(indiv.URI)
                Expect(err).To(Equal(ErrResourceNotFound))
            })
        })
    })

    Describe("Retrieving ontology individuals", func() {
        var indiv1, indiv2, indiv3, indiv4 OntologyIndividual
        var filter TripleFilter
        BeforeEach(func() {
            // Setup a bunch of individuals
            indiv1 = OntologyIndividual{
                URI:              testUri + "#indiv1",
                Types:            []string{"http://abc.com#type1"},
                Label:            map[string]string{},
                Comment:          map[string]string{},
                SameIndividualAs: []string{},
            }
            indiv2 = OntologyIndividual{
                URI:              testUri + "#indiv2",
                Types:            []string{"http://abc.com#type2"},
                Label:            map[string]string{},
                Comment:          map[string]string{},
                SameIndividualAs: []string{},
            }
            indiv3 = OntologyIndividual{
                URI:              testUri + "#indiv3",
                Types:            []string{"http://abc.com#type1", "http://abc.com#type2", "http://abc.com#type3"},
                Label:            map[string]string{},
                Comment:          map[string]string{},
                SameIndividualAs: []string{},
            }
            indiv4 = OntologyIndividual{
                URI:              testUri + "#indiv4",
                Types:            []string{"http://abc.com#type2", "http://abc.com#type3"},
                Label:            map[string]string{},
                Comment:          map[string]string{},
                SameIndividualAs: []string{},
            }
            // Add object properties
            indiv1.AddObjectProperty("http://abc.com#prop1", "http://abc.com#indiv2")
            indiv1.AddObjectProperty("http://abc.com#prop1", "http://abc.com#indiv3")
            indiv2.AddObjectProperty("http://abc.com#prop2", "http://abc.com#indiv1")
            // indiv3 does not have any object proerties
            indiv1.AddDataProperty("http://abc.com#dataprop1", XSDStringLiteral("Some string literal").Generic())
            indiv3.AddDataProperty("http://abc.com#dataprop2", XSDIntegerLiteral(42).Generic())
            // Add all individuals
            err := ont.UpsertResource(&indiv1)
            Expect(err).NotTo(HaveOccurred())
            err = ont.UpsertResource(&indiv2)
            Expect(err).NotTo(HaveOccurred())
            err = ont.UpsertResource(&indiv3)
            Expect(err).NotTo(HaveOccurred())
            err = ont.UpsertResource(&indiv4)
            Expect(err).NotTo(HaveOccurred())
            // Initialize filter
            filter = TripleFilter{}
        })

        When("not supplying any filter", func() {
            It("should return all individuals in the ontology", func() {
                indivs, err := ont.GetIndividuals(nil)
                Expect(err).NotTo(HaveOccurred())
                found1, found2, found3, found4 := false, false, false, false
                for _, indiv := range indivs {
                    if indiv.URI == indiv1.URI {
                        checkIndividuals(indiv, indiv1)
                        found1 = true
                    } else if indiv.URI == indiv2.URI {
                        checkIndividuals(indiv, indiv2)
                        found2 = true
                    } else if indiv.URI == indiv3.URI {
                        checkIndividuals(indiv, indiv3)
                        found3 = true
                    } else if indiv.URI == indiv4.URI {
                        checkIndividuals(indiv, indiv4)
                        found4 = true
                    }
                }
                Expect(found1).To(BeTrue())
                Expect(found2).To(BeTrue())
                Expect(found3).To(BeTrue())
                Expect(found4).To(BeTrue())
            })
        })
        When("filtered by a single class", func() {
            It("should return the individuals of the specified class only", func() {
                filter = filter.OrWithClass("http://abc.com#type1")
                indivs, err := ont.GetIndividuals(filter)
                Expect(err).NotTo(HaveOccurred())
                Expect(len(indivs)).To(Equal(2))
                found1, found3 := false, false
                for _, indiv := range indivs {
                    if indiv.URI == indiv1.URI {
                        checkIndividuals(indiv, indiv1)
                        found1 = true
                    } else if indiv.URI == indiv3.URI {
                        checkIndividuals(indiv, indiv3)
                        found3 = true
                    }
                }
                Expect(found1).To(BeTrue())
                Expect(found3).To(BeTrue())
            })
        })
        When("filtered by all given classes", func() {
            It("should return the individuals that match all the specified classes", func() {
                filter = filter.AndWithClass("http://abc.com#type2")
                filter = filter.AndWithClass("http://abc.com#type3")
                indivs, err := ont.GetIndividuals(filter)
                Expect(err).NotTo(HaveOccurred())
                Expect(len(indivs)).To(Equal(2))
                found3, found4 := false, false
                for _, indiv := range indivs {
                    if indiv.URI == indiv3.URI {
                        checkIndividuals(indiv, indiv3)
                        found3 = true
                    } else if indiv.URI == indiv4.URI {
                        checkIndividuals(indiv, indiv4)
                        found4 = true
                    }
                }
                Expect(found3).To(BeTrue())
                Expect(found4).To(BeTrue())
            })
        })
        When("filtered by any given class", func() {
            It("should return the individuals that match any of the specified classes", func() {
                filter = filter.OrWithClass("http://abc.com#type1")
                filter = filter.OrWithClass("http://abc.com#type3")
                indivs, err := ont.GetIndividuals(filter)
                Expect(err).NotTo(HaveOccurred())
                Expect(len(indivs)).To(Equal(3))
                found1, found3, found4 := false, false, false
                for _, indiv := range indivs {
                    if indiv.URI == indiv1.URI {
                        checkIndividuals(indiv, indiv1)
                        found1 = true
                    } else if indiv.URI == indiv3.URI {
                        checkIndividuals(indiv, indiv3)
                        found3 = true
                    } else if indiv.URI == indiv4.URI {
                        checkIndividuals(indiv, indiv4)
                        found4 = true
                    }
                }
                Expect(found1).To(BeTrue())
                Expect(found3).To(BeTrue())
                Expect(found4).To(BeTrue())
            })
        })
        When("filtered by an object property", func() {
            It("should return the individuals with the specified property only", func() {
                filter = filter.AndWithObjectProperty("http://abc.com#prop2", "http://abc.com#indiv1")
                indivs, err := ont.GetIndividuals(filter)
                Expect(err).NotTo(HaveOccurred())
                Expect(len(indivs)).To(Equal(1))
                Expect(indivs[0].URI).To(Equal(indiv2.URI))
                checkIndividuals(indivs[0], indiv2)
            })
        })
        When("filtered by a data property", func() {
            It("should return the individuals with the specified property only", func() {
                filter := filter.AndWithDataProperty("http://abc.com#dataprop2", XSDIntegerLiteral(42).Generic())
                indivs, err := ont.GetIndividuals(filter)
                Expect(err).NotTo(HaveOccurred())
                Expect(len(indivs)).To(Equal(1))
                Expect(indivs[0].URI).To(Equal(indiv3.URI))
                checkIndividuals(indivs[0], indiv3)
            })
        })
        When("filtered by a chain of classes and properties", func() {
            It("should return the expected individuals only", func() {
                filter = filter.AndWithClass("http://abc.com#type2")
                filter = filter.AndWithObjectProperty("http://abc.com#prop2", "http://abc.com#indiv1")
                filter = filter.OrWithClass("http://abc.com#type3")
                filter = filter.AndWithDataProperty("http://abc.com#dataprop2", XSDIntegerLiteral(42).Generic())
                indivs, err := ont.GetIndividuals(filter)
                Expect(err).NotTo(HaveOccurred())
                Expect(len(indivs)).To(Equal(2))
                found2, found3 := false, false
                for _, indiv := range indivs {
                    if indiv.URI == indiv2.URI {
                        checkIndividuals(indiv, indiv2)
                        found2 = true
                    } else if indiv.URI == indiv3.URI {
                        checkIndividuals(indiv, indiv3)
                        found3 = true
                    }
                }
                Expect(found2).To(BeTrue())
                Expect(found3).To(BeTrue())
            })
        })
    })
})
