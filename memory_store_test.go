package ontograph_test

import (
	"fmt"
	"strings"

	"github.com/lithammer/shortuuid/v3"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/kahefi/ontograph"
)

var _ = Describe("MemoryStore", func() {
	var graph *MemoryStore
	var graphUri string
	var testTriples []Triple

	BeforeEach(func() {
		// Initialize graph store
		graphUri = fmt.Sprintf("https://www.ontograph.com/test-%s", shortuuid.New())
		graph = NewMemoryStore(graphUri)
		// Add relation 'rel-1' from base with targets to 3 resources
		trp1, err := NewTriple(NewResourceTerm(graphUri), NewResourceTerm(graphUri+"#rel-1"), NewResourceTerm(graphUri+"#a"))
		Expect(err).NotTo(HaveOccurred())
		trp2, err := NewTriple(NewResourceTerm(graphUri), NewResourceTerm(graphUri+"#rel-1"), NewResourceTerm(graphUri+"#b"))
		Expect(err).NotTo(HaveOccurred())
		trp3, err := NewTriple(NewResourceTerm(graphUri), NewResourceTerm(graphUri+"#rel-1"), NewResourceTerm(graphUri+"#c"))
		Expect(err).NotTo(HaveOccurred())
		// Add relation 'rel-2' from 'a' with a single target to a resources
		trp4, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#b"))
		Expect(err).NotTo(HaveOccurred())
		// Add relation 'rel-3' to 'rel-5' from 'c' with target to different literals
		trp5, err := NewTriple(NewResourceTerm(graphUri+"#c"), NewResourceTerm(graphUri+"#rel-3"), NewLiteralTerm("lit1", "", ""))
		Expect(err).NotTo(HaveOccurred())
		trp6, err := NewTriple(NewResourceTerm(graphUri+"#c"), NewResourceTerm(graphUri+"#rel-4"), NewLiteralTerm("lit2", "de", ""))
		Expect(err).NotTo(HaveOccurred())
		trp7, err := NewTriple(NewResourceTerm(graphUri+"#c"), NewResourceTerm(graphUri+"#rel-5"), NewLiteralTerm("lit3", "", graphUri+"#datatype"))
		Expect(err).NotTo(HaveOccurred())
		// Add all triples to test triples
		testTriples = []Triple{*trp1, *trp2, *trp3, *trp4, *trp5, *trp6, *trp7}
		err = graph.AddTriples(testTriples)
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		_ = graph.Drop()
	})

	Describe("Retrieving the graph URI", func() {
		It("should match the initialisation", func() {
			Expect(graph.GetURI()).To(Equal(graphUri))
		})
	})

	Describe("Retrieving a single triple match", func() {
		Context("when there is exactly one match", func() {
			It("should return the expected matches", func() {
				trp, err := graph.GetFirstMatch("", fmt.Sprintf("<%s#rel-3>", graphUri), "\"lit1\"")
				Expect(err).NotTo(HaveOccurred())
				Expect(*trp).To(Equal(testTriples[4]))
			})
		})
		Context("when there are multiple matches", func() {
			It("should return one of the matches", func() {
				trp, err := graph.GetFirstMatch(fmt.Sprintf("<%s>", graphUri), fmt.Sprintf("<%s#rel-1>", graphUri), "")
				Expect(err).NotTo(HaveOccurred())
				Expect(testTriples[0:3]).To(ContainElement(*trp))
			})
		})
		Context("when there is no match", func() {
			It("should return nil", func() {
				trp, err := graph.GetFirstMatch("", fmt.Sprintf("<%s#rel-42>", graphUri), "")
				Expect(err).NotTo(HaveOccurred())
				Expect(trp).To(BeNil())
			})
		})
	})

	Describe("Retrieving all triple matches", func() {
		Context("when there are matches", func() {
			It("should return all expected matches from the store", func() {
				trps, err := graph.GetAllMatches("", fmt.Sprintf("<%s#rel-1>", graphUri), "")
				Expect(err).NotTo(HaveOccurred())
				Expect(trps).To(ConsistOf(testTriples[0:3]))
			})
		})
		Context("when there is no match", func() {
			It("should return an empty slice", func() {
				trps, err := graph.GetAllMatches("", fmt.Sprintf("<%s#rel-1>", graphUri), "\"lit1\"")
				Expect(err).NotTo(HaveOccurred())
				Expect(trps).To(BeEmpty())
			})
		})
		Context("when all triples are matched", func() {
			It("should return all all triples in the store", func() {
				trps, err := graph.GetAllMatches("", "", "")
				Expect(err).NotTo(HaveOccurred())
				Expect(trps).To(ConsistOf(testTriples))
			})
		})
	})

	Describe("Removing all triple matches", func() {
		Context("when there are multiple matches", func() {
			It("should remove all the matches from the store", func() {
				err := graph.DeleteAllMatches("", "", fmt.Sprintf("<%s#b>", graphUri))
				Expect(err).NotTo(HaveOccurred())
				trps, err := graph.GetAllTriples()
				Expect(err).NotTo(HaveOccurred())
				Expect(trps).To(ContainElement(testTriples[0]))
				Expect(trps).NotTo(ContainElement(testTriples[1]))
				Expect(trps).To(ContainElement(testTriples[2]))
				Expect(trps).NotTo(ContainElement(testTriples[3]))
				Expect(trps).To(ContainElement(testTriples[4]))
				Expect(trps).To(ContainElement(testTriples[5]))
				Expect(trps).To(ContainElement(testTriples[6]))
			})
		})
		Context("when there is no match", func() {
			It("should leave the store unchanged", func() {
				err := graph.DeleteAllMatches(fmt.Sprintf("<%s#42>", graphUri), "", "")
				Expect(err).NotTo(HaveOccurred())
				trps, err := graph.GetAllTriples()
				Expect(err).NotTo(HaveOccurred())
				Expect(trps).To(ConsistOf(testTriples))
			})
		})
	})

	Describe("Retrieving all triples in the store", func() {
		It("should return all expected test triples", func() {
			trps, err := graph.GetAllTriples()
			Expect(err).NotTo(HaveOccurred())
			Expect(trps).To(ConsistOf(testTriples))
		})
	})

	Describe("Adding a triple", func() {
		Context("when the triple does not exist", func() {
			It("should add the triple to the store", func() {
				// Add triple
				trp, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#c"))
				Expect(err).NotTo(HaveOccurred())
				err = graph.AddTriple(*trp)
				Expect(err).NotTo(HaveOccurred())
				// Check that triple is in store
				trps, err := graph.GetAllTriples()
				Expect(err).NotTo(HaveOccurred())
				Expect(trps).To(ContainElement(Triple{
					Subject:   Term(fmt.Sprintf("<%s#a>", graphUri)),
					Predicate: Term(fmt.Sprintf("<%s#rel-2>", graphUri)),
					Object:    Term(fmt.Sprintf("<%s#c>", graphUri)),
				}))
			})
		})
		Context("when the triple already exists", func() {
			Context("and the adding is checked", func() {
				It("should error with a conflict", func() {
					// Add triple
					trp, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#b"))
					Expect(err).NotTo(HaveOccurred())
					err = graph.AddTriple(*trp)
					// Check error and make sure the store is unchanged
					Expect(err).To(Equal(ErrTripleAlreadyExists))
					trps, err := graph.GetAllTriples()
					Expect(err).NotTo(HaveOccurred())
					Expect(trps).To(ConsistOf(testTriples))
				})
			})
			Context("and the adding is unchecked", func() {
				It("should not error", func() {
					// Add triple
					By("returning ErrTripleAlreadyExists")
					trp, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#b"))
					Expect(err).NotTo(HaveOccurred())
					err = graph.AddTripleUnchecked(*trp)
					// Check error and make sure the store is unchanged
					Expect(err).NotTo(HaveOccurred())
					trps, err := graph.GetAllTriples()
					Expect(err).NotTo(HaveOccurred())
					Expect(trps).To(ConsistOf(testTriples))
				})
			})
		})
	})

	Describe("Adding multiple triples", func() {
		Context("when none of the triples exist", func() {
			It("should add the triples to the store", func() {
				// Add triples
				trp8, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#d"))
				Expect(err).NotTo(HaveOccurred())
				trp9, err := NewTriple(NewResourceTerm(graphUri+"#d"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#e"))
				Expect(err).NotTo(HaveOccurred())
				trp10, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-6"), NewLiteralTerm("lit", "en", ""))
				Expect(err).NotTo(HaveOccurred())
				err = graph.AddTriples([]Triple{*trp8, *trp9, *trp10})
				Expect(err).NotTo(HaveOccurred())
				// Check that all triples are in store
				trps, err := graph.GetAllTriples()
				Expect(err).NotTo(HaveOccurred())
				Expect(trps).To(ContainElements(*trp8, *trp9, *trp10))
			})
		})
		Context("when some of the triples already exist", func() {
			Context("and the adding is checked", func() {
				It("should error with a conflict", func() {
					// Add triples
					trp8, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#d"))
					Expect(err).NotTo(HaveOccurred())
					trp9, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#b"))
					Expect(err).NotTo(HaveOccurred())
					trp10, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-6"), NewLiteralTerm("lit", "en", ""))
					Expect(err).NotTo(HaveOccurred())
					err = graph.AddTriples([]Triple{*trp8, *trp9, *trp10})
					// Check error and make sure the store is unchanged
					Expect(err).To(Equal(ErrTripleAlreadyExists))
					trps, err := graph.GetAllTriples()
					Expect(err).NotTo(HaveOccurred())
					Expect(trps).To(ConsistOf(testTriples))
				})
			})
			Context("and the adding is unchecked", func() {
				It("should not error and add the non-existent triples", func() {
					// Add triples
					trp8, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#d"))
					Expect(err).NotTo(HaveOccurred())
					trp9, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#b"))
					Expect(err).NotTo(HaveOccurred())
					trp10, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-6"), NewLiteralTerm("lit", "en", ""))
					Expect(err).NotTo(HaveOccurred())
					err = graph.AddTriplesUnchecked([]Triple{*trp8, *trp9, *trp10})
					// Check error
					Expect(err).NotTo(HaveOccurred())
					// Check that the non-existent triples are now in the store
					trps, err := graph.GetAllTriples()
					Expect(err).NotTo(HaveOccurred())
					Expect(trps).To(ContainElements(*trp8, *trp10))
					// Check that the existent triples was not added twice
					trps, err = graph.GetAllMatches(trp9.Subject.String(), trp9.Predicate.String(), trp9.Object.String())
					Expect(err).NotTo(HaveOccurred())
					Expect(trps).To(HaveLen(1))
				})
			})
		})
	})

	Describe("Deleting a triple", func() {
		Context("when the triple exists", func() {
			It("should remove the triple from the store", func() {
				// Delete triple
				trp, err := NewTriple(NewResourceTerm(graphUri), NewResourceTerm(graphUri+"#rel-1"), NewResourceTerm(graphUri+"#c"))
				Expect(err).NotTo(HaveOccurred())
				err = graph.DeleteTriple(*trp)
				Expect(err).NotTo(HaveOccurred())
				// Check that triple is not in the store anymore
				trps, err := graph.GetAllTriples()
				Expect(err).NotTo(HaveOccurred())
				Expect(trps).NotTo(ContainElement(Triple{
					Subject:   Term(fmt.Sprintf("<%s>", graphUri)),
					Predicate: Term(fmt.Sprintf("<%s#rel-1>", graphUri)),
					Object:    Term(fmt.Sprintf("<%s#c>", graphUri)),
				}))
			})
		})
		Context("when the triple does not exist", func() {
			Context("and the deletion is checked", func() {
				It("should error with a conflict", func() {
					// Delete triple
					trp, err := NewTriple(NewResourceTerm(graphUri), NewResourceTerm(graphUri+"#rel-42"), NewResourceTerm(graphUri+"#c"))
					Expect(err).NotTo(HaveOccurred())
					err = graph.DeleteTriple(*trp)
					// Check error and make sure the store is unchanged
					Expect(err).To(Equal(ErrTripleDoesNotExist))
					trps, err := graph.GetAllTriples()
					Expect(err).NotTo(HaveOccurred())
					Expect(trps).To(ConsistOf(testTriples))
				})
			})
			Context("and the deletion is unchecked", func() {
				It("should not error", func() {
					// Delete triple
					trp, err := NewTriple(NewResourceTerm(graphUri), NewResourceTerm(graphUri+"#rel-42"), NewResourceTerm(graphUri+"#c"))
					Expect(err).NotTo(HaveOccurred())
					err = graph.DeleteTripleUnchecked(*trp)
					// Check error and make sure the store is unchanged
					Expect(err).NotTo(HaveOccurred())
					trps, err := graph.GetAllTriples()
					Expect(err).NotTo(HaveOccurred())
					Expect(trps).To(ConsistOf(testTriples))
				})
			})
		})
	})

	Describe("Deleting multiple triples", func() {
		Context("when all of the triples exist", func() {
			It("should remove the triples from the store", func() {
				// Delete triples
				trp1, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#b"))
				Expect(err).NotTo(HaveOccurred())
				trp4, err := NewTriple(NewResourceTerm(graphUri), NewResourceTerm(graphUri+"#rel-1"), NewResourceTerm(graphUri+"#a"))
				Expect(err).NotTo(HaveOccurred())
				trp5, err := NewTriple(NewResourceTerm(graphUri+"#c"), NewResourceTerm(graphUri+"#rel-3"), NewLiteralTerm("lit1", "", ""))
				Expect(err).NotTo(HaveOccurred())
				err = graph.DeleteTriples([]Triple{*trp1, *trp4, *trp5})
				Expect(err).NotTo(HaveOccurred())
				// Check that the deleted triples are not in the store anymore
				trps, err := graph.GetAllTriples()
				Expect(err).NotTo(HaveOccurred())
				Expect(trps).NotTo(ContainElement(*trp1))
				Expect(trps).NotTo(ContainElement(*trp4))
				Expect(trps).NotTo(ContainElement(*trp5))
			})
		})
		Context("when some of the triples do not exist", func() {
			Context("and the deletion is checked", func() {
				It("should error with a conflict", func() {
					// Delete triples
					trp1, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#b"))
					Expect(err).NotTo(HaveOccurred())
					trp4, err := NewTriple(NewResourceTerm(graphUri), NewResourceTerm(graphUri+"#rel-1"), NewResourceTerm(graphUri+"#a"))
					Expect(err).NotTo(HaveOccurred())
					trp5, err := NewTriple(NewResourceTerm(graphUri+"#c"), NewResourceTerm(graphUri+"#rel-42"), NewLiteralTerm("lit1", "", ""))
					Expect(err).NotTo(HaveOccurred())
					err = graph.DeleteTriples([]Triple{*trp1, *trp4, *trp5})
					// Check for error and that the store is left unchanged
					Expect(err).To(Equal(ErrTripleDoesNotExist))
					trps, err := graph.GetAllTriples()
					Expect(err).NotTo(HaveOccurred())
					Expect(trps).To(ConsistOf(testTriples))
				})
			})
			Context("and the deletion is unchecked", func() {
				It("should not error and delete the existing triples", func() {
					// Delete triples
					trp1, err := NewTriple(NewResourceTerm(graphUri+"#a"), NewResourceTerm(graphUri+"#rel-2"), NewResourceTerm(graphUri+"#b"))
					Expect(err).NotTo(HaveOccurred())
					trp4, err := NewTriple(NewResourceTerm(graphUri), NewResourceTerm(graphUri+"#rel-1"), NewResourceTerm(graphUri+"#a"))
					Expect(err).NotTo(HaveOccurred())
					trp5, err := NewTriple(NewResourceTerm(graphUri+"#c"), NewResourceTerm(graphUri+"#rel-42"), NewLiteralTerm("lit1", "", ""))
					Expect(err).NotTo(HaveOccurred())
					err = graph.DeleteTriplesUnchecked([]Triple{*trp1, *trp4, *trp5})
					// Check for error and that the store does not contain the existing triples anymore
					Expect(err).NotTo(HaveOccurred())
					trps, err := graph.GetAllTriples()
					Expect(err).NotTo(HaveOccurred())
					Expect(trps).NotTo(ContainElement(trp1))
					Expect(trps).NotTo(ContainElement(trp4))
				})
			})
		})
	})

	Describe("Droping the graph store", func() {
		It("should render the store unusable", func() {
			_ = graph.Drop()
			Expect(graph.GetURI()).To(BeEmpty())
		})
	})

	Describe("Serializing the graph store to and lading from TTL", func() {
		It("should generate and load valid TTL without pretty printed serialization", func() {
			var ttlContent strings.Builder
			// Serialize to turtle
			By("serializing the graph to TTL")
			err := graph.SerializeToTurtle(&ttlContent, false)
			Expect(err).NotTo(HaveOccurred())
			// Load turtle back into graph store and compare to test triples
			By("loading the TTL back to a store")
			loadedGraph, err := ParseFromTurtle(strings.NewReader(ttlContent.String()))
			Expect(err).NotTo(HaveOccurred())
			By("matching the test triples exactly")
			trps, err := loadedGraph.GetAllTriples()
			Expect(err).NotTo(HaveOccurred())
			Expect(trps).To(ConsistOf(testTriples))
		})
		It("should generate and load valid TTL with pretty printed serialization", func() {
			var ttlContent strings.Builder
			// Serialize to turtle
			By("serializing the graph to TTL")
			err := graph.SerializeToTurtle(&ttlContent, true)
			Expect(err).NotTo(HaveOccurred())
			// Load turtle back into graph store and compare to test triples
			By("loading the TTL back to a store")
			loadedGraph, err := ParseFromTurtle(strings.NewReader(ttlContent.String()))
			Expect(err).NotTo(HaveOccurred())
			By("matching the test triples exactly")
			trps, err := loadedGraph.GetAllTriples()
			Expect(err).NotTo(HaveOccurred())
			Expect(trps).To(ConsistOf(testTriples))
		})
	})

	Describe("Retrieving the size of the graph store", func() {
		It("should return the expected number of triples", func() {
			Expect(graph.Size()).To(Equal(len(testTriples)))
		})
	})
})
