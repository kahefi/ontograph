package ontograph_test

import (
    "testing"

    //. "github.com/kahefi/ontograph"
    . "github.com/onsi/ginkgo"
    . "github.com/onsi/gomega"
)

var blazegraphHost string

func TestOntograph(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Ontograph Suite")
}

var _ = BeforeSuite(func() {
    blazegraphHost = "http://127.0.0.1:5060"
})
