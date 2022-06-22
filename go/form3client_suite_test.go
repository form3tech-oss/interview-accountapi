package form3client_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestForm3client(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Form3client Suite")
}
