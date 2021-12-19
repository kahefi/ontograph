package ontograph

import (
    "errors"
    "fmt"
    "strconv"
    "time"
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

// ToXSDString parses the literal into a xsd:string literal. If the literal is not of type xsd:string, an `ErrLiteralTypeMismatch` is returned.
func (l *GenericLiteral) ToXSDString() (XSDStringLiteral, error) {
    // Check for type mismatch
    if l.Type().URI != XSDString {
        return "", ErrLiteralTypeMismatch
    }
    // Parse literal
    return XSDStringLiteral(l.Value()), nil
}

// ***************
// * xsd:integer *
// ***************

type XSDIntegerLiteral int

func (l XSDIntegerLiteral) Generic() GenericLiteral {
    t := NewLiteralTerm(strconv.Itoa(int(l)), "", XSDInteger)
    return *NewGenericLiteral(t)
}

// ***************
// * xsd:decimal *
// ***************

type XSDDecimalLiteral float64

func (l XSDDecimalLiteral) Generic() GenericLiteral {
    t := NewLiteralTerm(fmt.Sprintf("%f", float64(l)), "", XSDDecimal)
    return *NewGenericLiteral(t)
}

// ToXSDDecimalLiteral parses the literal into a xsd:decimal literal. If the literal is not a number, an `ErrLiteralTypeMismatch` is returned.
func (l *GenericLiteral) ToXSDDecimal() (XSDDecimalLiteral, error) {
    // Check for type mismatch
    if l.Type().URI != XSDDecimal {
        return 0, ErrLiteralTypeMismatch
    }
    // Parse literal
    val, err := strconv.ParseFloat(l.Value(), 64)
    if err != nil {
        return 0, err
    }
    return XSDDecimalLiteral(val), nil
}

// ***************
// * xsd:boolean *
// ***************

type XSDBooleanLiteral bool

func (l XSDBooleanLiteral) Generic() GenericLiteral {
    t := NewLiteralTerm(strconv.FormatBool(bool(l)), "", XSDBoolean)
    return *NewGenericLiteral(t)
}

// ToXSDBoolean parses the literal into a xsd:boolean literal. If the literal is not of type xsd:boolean, an `ErrLiteralTypeMismatch` is returned.
func (l *GenericLiteral) ToXSDBoolean() (XSDBooleanLiteral, error) {
    // Check for type mismatch
    if l.Type().URI != XSDBoolean {
        return false, ErrLiteralTypeMismatch
    }
    // Parse literal
    val, err := strconv.ParseBool(l.Value())
    if err != nil {
        return false, err
    }
    return XSDBooleanLiteral(val), nil
}

// ***************
// * xsd:anyURI *
// ***************

type XSDAnyURILiteral string

func (l XSDAnyURILiteral) Generic() GenericLiteral {
    t := NewLiteralTerm(string(l), "", XSDAnyURI)
    return *NewGenericLiteral(t)
}

// ToXSDAnyURI parses the literal into a xsd:anyURI literal. If the literal is not of type xsd:anyURI, an `ErrLiteralTypeMismatch` is returned.
func (l *GenericLiteral) ToXSDAnyURI() (XSDAnyURILiteral, error) {
    // Check for type mismatch
    if l.Type().URI != XSDAnyURI {
        return "", ErrLiteralTypeMismatch
    }
    // Parse literal
    return XSDAnyURILiteral(l.Value()), nil
}

// ***************
// * xsd:dateTime *
// ***************

type XSDDateTimeLiteral time.Time

func (l XSDDateTimeLiteral) Generic() GenericLiteral {
    t := NewLiteralTerm(l.Format(time.RFC3339), "", XSDDateTime)
    return *NewGenericLiteral(t)
}

// ToXSDDateTime parses the literal into a xsd:dateTime literal. If the literal is not of type xsd:dateTime, an `ErrLiteralTypeMismatch` is returned. The value must be formatted according to the RFC3339 standard.
func (l *GenericLiteral) ToXSDDateTime() (XSDDateTimeLiteral, error) {
    var t time.Time
    // Check for type mismatch
    if l.Type().URI != XSDDateTime {
        return t, ErrLiteralTypeMismatch
    }
    // Parse literal
    t, err := time.Parse(time.RFC3339, l.Value())
    if err != nil {
        return t, err
    }
    return XSDDateTimeLiteral(t), nil
}
