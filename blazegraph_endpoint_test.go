package ontograph_test

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/teris-io/shortid"

	. "github.com/kahefi/ontograph"
)

var _ = Describe("BlazegraphEndpoint", func() {
	var endpoint *BlazegraphEndpoint

	BeforeEach(func() {
		endpoint = NewBlazegraphEndpoint("http://127.0.0.1:5060")
	})

	AfterEach(func() {
	})

	Describe("Checking if an endpoint is online", func() {
		When("the endpoint exists", func() {
			It("should return true without error", func() {
				endpoint = NewBlazegraphEndpoint("http://127.0.0.1:5060")
				isOnline, err := endpoint.IsOnline()
				Expect(err).NotTo(HaveOccurred())
				Expect(isOnline).To(BeTrue())
			})
		})
		When("the URL is malformed", func() {
			It("return false with error", func() {
				endpoint = NewBlazegraphEndpoint("127.0.0.1:5060")
				isOnline, err := endpoint.IsOnline()
				Expect(err).To(HaveOccurred())
				Expect(isOnline).To(BeFalse())
			})
		})
		When("the endpoint does not exist", func() {
			It("return false with error", func() {
				endpoint = NewBlazegraphEndpoint("http://127.0.0.1:5061")
				isOnline, err := endpoint.IsOnline()
				Expect(err).To(HaveOccurred())
				Expect(isOnline).To(BeFalse())
			})
		})
	})

	Describe("Creating a new namespace", func() {
		var testNs string
		BeforeEach(func() {
			testNs = fmt.Sprintf("ns-%s", shortid.MustGenerate())
		})
		AfterEach(func() {
			_ = endpoint.DropNamespace(testNs)
		})
		When("the namespace does not exist yet", func() {
			It("should successfully create a new namespace", func() {
				By("not returning an error")
				err := endpoint.CreateNamespace(testNs)
				Expect(err).NotTo(HaveOccurred())
				By("having created the new namespace")
				namespaces, err := endpoint.GetNamespaces()
				Expect(err).NotTo(HaveOccurred())
				Expect(namespaces).To(ContainElement(testNs))
				exists, err := endpoint.NamespaceExists(testNs)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())
			})
		})
		When("the namespace already exists", func() {
			BeforeEach(func() {
				// Create namespace for existing conflict
				err := endpoint.CreateNamespace(testNs)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should reject the creation of a new namespace", func() {
				By("returning an error")
				err := endpoint.CreateNamespace(testNs)
				Expect(err).To(HaveOccurred())
				By("having left the namespace intact")
				exists, err := endpoint.NamespaceExists(testNs)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())
			})
		})
	})

	Describe("Deleting a namespace", func() {
		var testNs string
		BeforeEach(func() {
			testNs = fmt.Sprintf("ns-%s", shortid.MustGenerate())
		})
		AfterEach(func() {
			_ = endpoint.DropNamespace(testNs)
		})
		When("the namespace exist", func() {
			BeforeEach(func() {
				err := endpoint.CreateNamespace(testNs)
				Expect(err).NotTo(HaveOccurred())
			})
			It("should successfully delete the namespace", func() {
				By("not returning an error")
				err := endpoint.DropNamespace(testNs)
				Expect(err).NotTo(HaveOccurred())
				By("having having deleted the namespace")
				namespaces, err := endpoint.GetNamespaces()
				Expect(err).NotTo(HaveOccurred())
				Expect(namespaces).NotTo(ContainElement(testNs))
				exists, err := endpoint.NamespaceExists(testNs)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeFalse())
			})
		})
		When("the namespace does not exists", func() {
			It("should return a succeed message", func() {
				By("not returning an error")
				err := endpoint.CreateNamespace(testNs)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})

	Describe("Retrieving a list of namespace", func() {
		var testNs []string
		BeforeEach(func() {
			for _, i := range []int{0, 1, 2, 3, 4, 5} {
				testNs = append(testNs, fmt.Sprintf("ns-%s", shortid.MustGenerate()))
				err := endpoint.CreateNamespace(testNs[i])
				Expect(err).NotTo(HaveOccurred())
			}
		})
		AfterEach(func() {
			for _, ns := range testNs {
				_ = endpoint.DropNamespace(ns)
			}
		})
		It("should return a list containing the expected namespaces", func() {
			nsList, err := endpoint.GetNamespaces()
			By("not returning an error")
			Expect(err).NotTo(HaveOccurred())
			By("Containing each of the expected namespaces")
			Expect(nsList).To(ContainElements(testNs))
		})
	})

	Describe("Retrieving a list of graph URIs in a namespace", func() {
		var testGraphs []string
		var testNs string
		BeforeEach(func() {
			testNs = fmt.Sprintf("ns-%s", shortid.MustGenerate())
			err := endpoint.CreateNamespace(testNs)
			Expect(err).NotTo(HaveOccurred())
			for range []int{0, 1, 2, 3, 4, 5} {
				// Init new ontology for the graph
				testGraph := fmt.Sprintf("http://test.com/graph-%s", shortid.MustGenerate())
				_, err = InitOntologyGraph(endpoint.NewBlazegraphStore(testGraph, testNs))
				Expect(err).NotTo(HaveOccurred())
				// Register graph for testing later
				testGraphs = append(testGraphs, testGraph)
			}
		})
		AfterEach(func() {
			_ = endpoint.DropNamespace(testNs)
		})
		It("should return a list containing the expected graph URIs", func() {
			graphList, err := endpoint.GetGraphs(testNs)
			By("not returning an error")
			Expect(err).NotTo(HaveOccurred())
			By("Containing each of the expected namespaces")
			Expect(graphList).To(ContainElements(testGraphs))
		})
	})

	// DoSparqlTurtleQuery covered by BlazegraphStore tests

	// DoSparqlJsonQuery covered by BlazegraphStore tests

	// DoSparqlUpdate covered by BlazegraphStore tests

})
