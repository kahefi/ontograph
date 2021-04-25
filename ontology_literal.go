package ontograph

import (
    "errors"
    "strconv"
)

// GenericLiteral represents a generic literal term (i.e. containing a value and a corresponding datatype).
// Generic literals can be parsed into specific literals using the corresponding methods.
type GenericLiteral struct {
    value    Term
    datatype OntologyDatatype
}

// NewGenericLiteral creates a new generic literal from the given term.
func NewGenericLiteral(t Term) *GenericLiteral {
    return &GenericLiteral{
        value: t,
        datatype: OntologyDatatype{
            URI: t.Datatype(),
        },
    }
}

// Term returns the term representation of the literal.
func (l *GenericLiteral) Term() Term {
    return l.value
}

// Type returns the ontological datatype of the literal.
func (l *GenericLiteral) Type() OntologyDatatype {
    return l.datatype
}

// Value returns a string representation of the value of the literal.
func (l *GenericLiteral) Value() string {
    return l.value.Value()
}

// String returns a string representation of the whole literal in NTriple format.
// This method is equivalent to `l.Term().String()`.
func (l *GenericLiteral) String() string {
    return l.value.String()
}

// ToXSDString parses the literal into a xsd:string literal. If the literal is not of type xsd:string, an `ErrLiteralTypeMismatch` is returned.
func (l *GenericLiteral) ToXSDString() (XSDStringLiteral, error) {
    // Check for type mismatch
    if l.Type().URI != XSDString {
        return "", ErrLiteralTypeMismatch
    }
    // Parse literal
    return XSDStringLiteral(l.Value()), nil
}

// ToXSDInteger parses the literal into a xsd:integer literal. If the literal is not of type xsd:integer, an `ErrLiteralTypeMismatch` is returned.
func (l *GenericLiteral) ToXSDInteger() (XSDIntegerLiteral, error) {
    // Check for type mismatch
    if l.Type().URI != XSDInteger {
        return 0, ErrLiteralTypeMismatch
    }
    // Parse literal
    val, err := strconv.Atoi(l.Value())
    if err != nil {
        return 0, err
    }
    return XSDIntegerLiteral(val), nil
}

// ErrLiteralTypeMismatch is raised when a generic literal is attempted to be converted into a specific literal of a certain datatype, but the datatype does not match.
var ErrLiteralTypeMismatch error = errors.New("The literal is not of the expected type")

// **************
// * xsd:string *
// **************

type XSDStringLiteral string

func (l XSDStringLiteral) Generic() GenericLiteral {
    t := NewLiteralTerm(string(l), "", XSDString)
    return *NewGenericLiteral(t)
}

// ***************
// * xsd:integer *
// ***************

type XSDIntegerLiteral int

func (l XSDIntegerLiteral) Generic() GenericLiteral {
    t := NewLiteralTerm(strconv.Itoa(int(l)), "", XSDInteger)
    return *NewGenericLiteral(t)
}
