package ontograph_test

import (
	. "github.com/kahefi/ontograph"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Triple", func() {

	Describe("Creating a new resource term", func() {
		It("should return the expected representation", func() {
			Expect(NewResourceTerm("https://www.ontograph.com/test").String()).To(Equal("<https://www.ontograph.com/test>"))
			Expect(NewResourceTerm("https://www.ontograph.com/test#a").String()).To(Equal("<https://www.ontograph.com/test#a>"))
		})
	})

	Describe("Creating a new literal term", func() {
		It("should return the expected representation", func() {
			Expect(NewLiteralTerm("Lorem ipsum", "", "").String()).To(Equal("\"Lorem ipsum\""))
			Expect(NewLiteralTerm("Lorem ipsum", "en", "").String()).To(Equal("\"Lorem ipsum\"@en"))
			Expect(NewLiteralTerm("Lorem ipsum", "", "http://www.w3.org/2001/XMLSchema#int").String()).To(Equal("\"Lorem ipsum\"^^<http://www.w3.org/2001/XMLSchema#int>"))
		})
	})

	Describe("Checking if a term is a resource", func() {
		Context("when the term has a valid NTriple resource representation", func() {
			It("should confirm the term", func() {
				Expect(Term("<https://www.ontograph.com/test>").IsResource()).To(BeTrue())
				Expect(Term("<https://www.ontograph.com/test#a>").IsResource()).To(BeTrue())
			})
		})
		Context("when the term is a valid NTriple literal", func() {
			It("should reject the term", func() {
				Expect(Term(`"some literal"`).IsResource()).To(BeFalse())
				Expect(Term(`"some literal"@de`).IsResource()).To(BeFalse())
				Expect(Term(`"some literal"^^<https://www.ontograph.com/test#literal>`).IsResource()).To(BeFalse())
			})
		})
		Context("when the term is not formatted as NTriple", func() {
			It("should reject the term", func() {
				Expect(Term("https://www.ontograph.com/test").IsResource()).To(BeFalse())
				Expect(Term("<https://www.ontograph.com/test").IsResource()).To(BeFalse())
				Expect(Term("https://www.ontograph.com/test>").IsResource()).To(BeFalse())
				Expect(Term("<>").IsResource()).To(BeFalse())
				Expect(Term("").IsResource()).To(BeFalse())
			})
		})
	})

	Describe("Checking if a term is a literal", func() {
		Context("when the term has a valid NTriple literal representation", func() {
			It("should confirm the term", func() {
				Expect(Term(`"some literal"`).IsLiteral()).To(BeTrue())
				Expect(Term(`"some literal"@de`).IsLiteral()).To(BeTrue())
				Expect(Term(`"some literal"^^<https://www.ontograph.com/test#literal>`).IsLiteral()).To(BeTrue())
			})
		})
		Context("when the term is a valid NTriple resource", func() {
			It("should reject the term", func() {
				Expect(Term("<https://www.ontograph.com/test>").IsLiteral()).To(BeFalse())
				Expect(Term("<https://www.ontograph.com/test#a>").IsLiteral()).To(BeFalse())
			})
		})
		Context("when the term is not formatted as NTriple", func() {
			It("should reject the term", func() {
				Expect(Term(`some literal`).IsLiteral()).To(BeFalse())
				Expect(Term(`"some literal`).IsLiteral()).To(BeFalse())
				Expect(Term(`some literal"`).IsLiteral()).To(BeFalse())
				Expect(Term(`"some literal"de`).IsLiteral()).To(BeFalse())
				Expect(Term(`"some literal"^<https://www.ontograph.com/test#literal>`).IsLiteral()).To(BeFalse())
				Expect(Term("\"\"").IsLiteral()).To(BeFalse())
				Expect(Term("").IsLiteral()).To(BeFalse())
			})
		})
	})

	Describe("Parsing the value from a term", func() {
		Context("when the term is a resource", func() {
			It("should return the expected URI", func() {
				Expect(Term("<https://www.ontograph.com/test>").Value()).To(Equal("https://www.ontograph.com/test"))
				Expect(Term("<https://www.ontograph.com/test#a>").Value()).To(Equal("https://www.ontograph.com/test#a"))
			})
		})
		Context("when the term is a literal", func() {
			It("should return the expected literal", func() {
				Expect(Term(`"some literal"`).Value()).To(Equal("some literal"))
				Expect(Term(`"some literal"@de`).Value()).To(Equal("some literal"))
				Expect(Term(`"some literal"^^<https://www.ontograph.com/test#literal>`).Value()).To(Equal("some literal"))
			})
		})
		Context("when the term is invalid", func() {
			It("should return an empty string", func() {
				Expect(Term(`some literal`).Value()).To(Equal(""))
				Expect(Term(`"some literal`).Value()).To(Equal(""))
				Expect(Term(`some literal"`).Value()).To(Equal(""))
				Expect(Term(`"some literal"de`).Value()).To(Equal(""))
				Expect(Term(`"some literal"^<https://www.ontograph.com/test#literal>`).Value()).To(Equal(""))
				Expect(Term("https://www.ontograph.com/test").Value()).To(Equal(""))
				Expect(Term("<https://www.ontograph.com/test").Value()).To(Equal(""))
				Expect(Term("https://www.ontograph.com/test>").Value()).To(Equal(""))
				Expect(Term("<>").Value()).To(Equal(""))
				Expect(Term("\"\"").Value()).To(Equal(""))
				Expect(Term("").Value()).To(Equal(""))
			})
		})
	})

	Describe("Parsing the language tag from a term", func() {
		Context("when the term is a literal", func() {
			It("should return the expected language", func() {
				Expect(Term(`"some literal"`).Language()).To(Equal(""))
				Expect(Term(`"some literal"@de`).Language()).To(Equal("de"))
				Expect(Term(`"some literal"^^<https://www.ontograph.com/test#literal>`).Language()).To(Equal(""))
			})
		})
		Context("when the term is invalid", func() {
			It("should return an empty string", func() {
				Expect(Term(`@de"some literal"`).Language()).To(Equal(""))
				Expect(Term(`de@"some literal"`).Language()).To(Equal(""))
				Expect(Term("<>@de").Language()).To(Equal(""))
				Expect(Term("\"\"").Language()).To(Equal(""))
				Expect(Term("").Language()).To(Equal(""))
			})
		})
	})

	Describe("Parsing the data type tag from a term", func() {
		Context("when the term is a literal", func() {
			It("should return the expected data type", func() {
				Expect(Term(`"some literal"`).Datatype()).To(Equal(""))
				Expect(Term(`"some literal"@de`).Datatype()).To(Equal(""))
				Expect(Term(`"some literal"^^<https://www.ontograph.com/test#literal>`).Datatype()).To(Equal("https://www.ontograph.com/test#literal"))
			})
		})
		Context("when the term is invalid", func() {
			It("should return an empty string", func() {
				Expect(Term(`^^de"some literal"`).Datatype()).To(Equal(""))
				Expect(Term(`de^^"some literal"`).Datatype()).To(Equal(""))
				Expect(Term(`"some literal"^^https://www.ontograph.com/test#literal`).Datatype()).To(Equal(""))
				Expect(Term("<>^^de").Datatype()).To(Equal(""))
				Expect(Term("\"\"").Datatype()).To(Equal(""))
				Expect(Term("").Datatype()).To(Equal(""))
			})
		})
	})

	Describe("Creating a new triple", func() {
		Context("when all terms are valid NTriple resources", func() {
			It("should return a valid triple", func() {
				trp, err := NewTriple("<https://www.ontograph.com/test>", "<https://www.ontograph.com/test#rel>", "<https://www.ontograph.com/test#a>")
				Expect(err).NotTo(HaveOccurred())
				Expect(trp.Subject.Value()).To(Equal("https://www.ontograph.com/test"))
				Expect(trp.Predicate.Value()).To(Equal("https://www.ontograph.com/test#rel"))
				Expect(trp.Object.Value()).To(Equal("https://www.ontograph.com/test#a"))

			})
		})
		Context("when the object term is a NTriple literal", func() {
			It("should return a valid triple", func() {
				trp, err := NewTriple("<https://www.ontograph.com/test>", "<https://www.ontograph.com/test#rel>", "\"some literal\"")
				Expect(err).NotTo(HaveOccurred())
				Expect(trp.Subject.Value()).To(Equal("https://www.ontograph.com/test"))
				Expect(trp.Predicate.Value()).To(Equal("https://www.ontograph.com/test#rel"))
				Expect(trp.Object.Value()).To(Equal("some literal"))
			})
		})
		Context("when the object term is a NTriple literal with language", func() {
			It("should return a valid triple", func() {
				trp, err := NewTriple("<https://www.ontograph.com/test>", "<https://www.ontograph.com/test#rel>", "\"some literal\"@en")
				Expect(err).NotTo(HaveOccurred())
				Expect(trp.Subject.Value()).To(Equal("https://www.ontograph.com/test"))
				Expect(trp.Predicate.Value()).To(Equal("https://www.ontograph.com/test#rel"))
				Expect(trp.Object.Value()).To(Equal("some literal"))
				Expect(trp.Object.Language()).To(Equal("en"))
			})
		})
		Context("when the object term is a NTriple literal with datatype", func() {
			It("should return a valid triple", func() {
				trp, err := NewTriple("<https://www.ontograph.com/test>", "<https://www.ontograph.com/test#rel>", "\"some literal\"^^<https://www.ontograph.com/test#literal>")
				Expect(err).NotTo(HaveOccurred())
				Expect(trp.Subject.Value()).To(Equal("https://www.ontograph.com/test"))
				Expect(trp.Predicate.Value()).To(Equal("https://www.ontograph.com/test#rel"))
				Expect(trp.Object.Value()).To(Equal("some literal"))
				Expect(trp.Object.Datatype()).To(Equal("https://www.ontograph.com/test#literal"))
			})
		})
		Context("when the subject is a valid NTriple literal", func() {
			It("should error", func() {
				_, err := NewTriple("\"some literal\"", "<https://www.ontograph.com/test#rel>", "<https://www.ontograph.com/test#a>")
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when the predicate is a valid NTriple literal", func() {
			It("should error", func() {
				_, err := NewTriple("<https://www.ontograph.com/test>", "\"some literal\"", "<https://www.ontograph.com/test#a>")
				Expect(err).To(HaveOccurred())
			})
		})
		Context("when one of the terms is invalid", func() {
			It("should error", func() {
				_, err := NewTriple("https://www.ontograph.com/test>", "<https://www.ontograph.com/test#rel>", "<https://www.ontograph.com/test#a>")
				Expect(err).To(HaveOccurred())
				_, err = NewTriple("<https://www.ontograph.com/test>", "<https://www.ontograph.com/test#rel", "<https://www.ontograph.com/test#a>")
				Expect(err).To(HaveOccurred())
				_, err = NewTriple("<https://www.ontograph.com/test>", "<https://www.ontograph.com/test#rel>", "https://www.ontograph.com/test#a>")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
