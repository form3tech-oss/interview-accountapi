package form3client_test

import (
	"errors"
	"fmt"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/ahermant/form3client"
)

var _ = Describe("Form3client", func() {

	Describe("Client Error", func() {
		Context("with proper error", func() {
			It("should return a client error", func() {
				err := form3client.ClientError(errors.New("test"))
				Expect(fmt.Sprint(err)).To(Equal("client error: test"))
			})
		})
	})

})
