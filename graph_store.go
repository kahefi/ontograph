package ontograph

import (
	"errors"
	"io"
)

// GraphStore provides methods to create, read, update and delete RDF triples for graphs.
type GraphStore interface {
	// GetUri should return the named graph URI.
	GetUri() string

	// GetFirstMatch should retrieve the first triple that matches the pattern. Empty strings in subject, predicate or object should be treated as wildcards.
	GetFirstMatch(subj, pred, obj string) (Triple, error)
	// GetAllMatches should retrieve all triples that match the pattern. Empty strings in subject, predicate or object should be treated as wildcards.
	GetAllMatches(subj, pred, obj string) ([]Triple, error)

	// DeleteAllMatches should remove all triples that match the pattern. Empty strings in subject, predicate or object should be treated as wildcards.
	DeleteAllMatches(subj, pred, obj string) error

	// GetAllTriples should return all triples in the store. The operation should be equivalent to GetAllMatches("", "", "").
	GetAllTriples() ([]Triple, error)

	// AddTriple should add the given triple to the store. It should error if the triple already exists.
	AddTriple(trp Triple) error
	// AddTriples should add all the given triples to the store. It should error if one of the triples already exists.
	AddTriples(trps []Triple) error
	// AddTripleUnchecked should add the given triple to the store. It should not error if the triple already exists.
	AddTripleUnchecked(trp Triple) error
	// AddTriplesUnchecked should add all the given triples to the store. It should not error if any of the triples already exists.
	AddTriplesUnchecked(trps []Triple) error

	// DeleteTriple should remove the given triple from the store.
	DeleteTriple(trp Triple) error
	// DeleteTriples should remove all the given triples from the store.
	DeleteTriples(trps []Triple) error
	// DeleteTripleUnchecked should remove the given triple from the store. It should not error if the triple does not exist.
	DeleteTripleUnchecked(trp Triple) error
	// DeleteTriplesUnchecked should remove all the given triples from the store. It should not error if any of the triples does not exist.
	DeleteTriplesUnchecked(trps []Triple) error

	// Drop should remove all triples and clear the store completely.
	Drop() error

	// SerializeToTurtle should write the entire store into the writer in Turtle (TTL) format. If pretty is set to true, the method should pretty print the turtle data.
	SerializeToTurtle(w io.Writer, pretty bool) error

	// Size should return the total number of triples in the store.
	Size() (int, error)
}

// *****************
// * Shared Errors *
// *****************

// ErrTripleAlreadyExists is raised on conflict errors when a triple already exists (i.e. adding triples).
var ErrTripleAlreadyExists error = errors.New("Triple already exists")

// ErrTripleDoesNotExist is raised on conflict errors when a triple does not yet exist (i.e. deleting triples).
var ErrTripleDoesNotExist error = errors.New("Triple does not exist")
