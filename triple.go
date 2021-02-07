package ontograph

import (
	"fmt"
	"strings"
)

// ********************
// * Term Definitions *
// ********************

// Term encodes a subject, predicate or objcect in NTriple format.
type Term string

// NewResourceTerm creates a new resource term in NTriple format.
func NewResourceTerm(uri string) Term {
	return Term(fmt.Sprintf("<%s>", uri))
}

// NewLiteralTerm creates a new literal term in NTriple format.
func NewLiteralTerm(literal, language, datatype string) Term {
	t := fmt.Sprintf("\"%s\"", literal)
	if language != "" {
		t += fmt.Sprintf("@%s", language)
	}
	if datatype != "" {
		t += fmt.Sprintf("^^<%s>", datatype)
	}
	return Term(t)
}

// String converts the term into a string. Equivalent to direct casting with string(t).
func (t Term) String() string {
	return string(t)
}

// IsResource returns true if the term is a resource.
func (t Term) IsResource() bool {
	s := string(t)
	return len(s) > 2 && string(s[0]) == "<" && string(s[len(s)-1]) == ">"
}

// IsLiteral returns true if the term is a literal.
func (t Term) IsLiteral() bool {
	s := string(t)
	return len(s) > 2 && string(s[0]) == "\"" && (string(s[len(s)-1]) == "\"" || strings.Contains(s, "\"@") || strings.Contains(s, "\"^^"))
}

// Value returns the value of the term (i.e. the URI or literal).
func (t Term) Value() string {
	s := string(t)
	if len(s) > 2 {
		if string(s[0]) == "<" && string(s[len(s)-1]) == ">" {
			return s[1 : len(s)-1]
		} else if string(s[0]) == "\"" && string(s[len(s)-1]) == "\"" {
			return s[1 : len(s)-1]
		} else if string(s[0]) == "\"" && strings.Contains(s, "\"@") {
			atPos := strings.LastIndex(s, "@")
			return s[1 : atPos-1]
		} else if string(s[0]) == "\"" && strings.Contains(s, "\"^^") {
			atPos := strings.LastIndex(s, "^^")
			return s[1 : atPos-1]
		} else {
			return ""
		}
	}
	return ""
}

// Language returns the language of the term. Will be the empty string if the term is a not a literal or does not contain a language.
func (t Term) Language() string {
	s := string(t)
	if len(s) > 2 && string(s[0]) == "\"" && strings.Contains(s, "\"@") {
		atPos := strings.LastIndex(s, "@")
		return s[atPos+1:]
	}
	return ""
}

// Datatype returns the data type of the term. Will be the empty string if the term is a not a literal or does not contain a data type.
func (t Term) Datatype() string {
	s := string(t)
	if len(s) > 2 && string(s[0]) == "\"" && strings.Contains(s, "\"^^") {
		atPos := strings.LastIndex(s, "^^")
		return Term(s[atPos+2:]).Value()
	}
	return ""
}

// **********************
// * Triple Definitions *
// **********************

// Triple represents a subject-predicate-object triple term from an ontology.
type Triple struct {
	Subject   Term
	Predicate Term
	Object    Term
}

// NewTriple creates a new triple from the given string terms. The terms are checked and parsed. If you are sure that the terms are valid NTriples, initialize directly with the Triple structure.
func NewTriple(subj, pred, obj Term) (*Triple, error) {
	// Sanity check terms
	if !subj.IsResource() {
		return nil, fmt.Errorf("Subject '%s' is not a resource", subj)
	}
	if !pred.IsResource() {
		return nil, fmt.Errorf("Predicate '%s' is not a resource", pred)
	}
	if !obj.IsResource() && !obj.IsLiteral() {
		return nil, fmt.Errorf("Object '%s' is not a resource or literal", obj)
	}
	// All fine, return triple
	trp := Triple{
		Subject:   subj,
		Predicate: pred,
		Object:    obj,
	}
	return &trp, nil
}
