package form3client_test

import (
	"github.com/ahermant/form3client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
)

var _ = Describe("Form3client", func() {
	Describe("Create account", func() {

		Context("with no data", func() {
			It("should trigger an error", func() {
				_, err := form3client.CreateAccount("")
				Expect(err).To(MatchError(form3client.ClientError(form3client.ErrNoAccountData)))
			})
		})

		Context("With data", func() {
			It("should send a correct response", func() {
				id := "ad27e265-9605-4b4b-a0e5-3004ea9cc8d9"
				accountData := `{
					"data":{
						"type": "accounts",
						"id": "` + id + `",
						"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
						"attributes": {
							"name": ["ahtest"],
							"country": "GB"
						}
					}
				}`
				resp, err := form3client.CreateAccount(accountData)
				if err != nil {
					log.Println(err)
				}
				log.Println(resp)
				Expect(resp).To(ContainSubstring(id))

				_, err = form3client.DeleteAccount(id, "0")
				if err != nil {
					log.Println(err)
				}
			})
		})

	})

	Describe("Delete account", func() {
		Context("with ID which does not exist", func() {
			It("should trigger an error", func() {
			})
		})

		Context("with version which does not exist", func() {
			It("should trigger an error", func() {
			})
		})

		Context("with existing ID and version", func() {
			It("should delete the account properly", func() {
			})
		})
	})

	// Test creation for each country

	// Test creation with missing data

	// Test deletion

	// Test deletion when ID does not exist

	// Test fetching

	// Test fetching where ID does not exist

})
