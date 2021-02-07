package ontograph_test

import (
	"testing"

	//. "github.com/kahefi/ontograph"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestOntograph(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ontograph Suite")
}

// var _ = BeforeSuite(func() {
//     dbRunner = db.NewRunner()
//     err := dbRunner.Start()
//     Expect(err).NotTo(HaveOccurred())

//     dbClient = db.NewClient()
//     err = dbClient.Connect(dbRunner.Address())
//     Expect(err).NotTo(HaveOccurred())
// })
