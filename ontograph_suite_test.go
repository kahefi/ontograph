package ontograph_test

import (
	"testing"

	. "github.com/kahefi/ontograph"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var endpoint *BlazegraphEndpoint
var testNamespace string

func TestOntograph(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ontograph Suite")
}

var _ = BeforeSuite(func() {
	endpoint = NewBlazegraphEndpoint("http://127.0.0.1:5060")
	testNamespace = "test-ns"
	err := endpoint.CreateNamespace(testNamespace)
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	err := endpoint.DropNamespace(testNamespace)
	Expect(err).NotTo(HaveOccurred())
})
