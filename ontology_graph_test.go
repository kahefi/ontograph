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
    var ont *OntologyGraph

    BeforeEach(func() {
        testUri = fmt.Sprintf("https://www.ontograph.com/test-%s", shortuuid.New())
        graph := NewMemoryStore(testUri)
        ont = NewOntologyGraph(graph)
    })

    AfterEach(func() {
    })

    Describe("Retrieving the ontology URI", func() {
        It("should match the initialisation", func() {
            Expect(ont.GetURI()).To(Equal(testUri))
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
})
