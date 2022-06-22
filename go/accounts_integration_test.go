package form3client_test

import (
	"github.com/ahermant/form3client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"log"
	"net/http"
)

func getBasicAccountData(id, version string) string {
	return `{
		"data":{
			"type": "accounts",
			"id": "` + id + `",
			"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			"version" :"` + version + `",
			"attributes": {
				"name": ["ahtest"],
				"country": "GB"
			}
		}
	}`
}

var _ = Describe("Form3client", func() {
	Describe("Create account", func() {

		Context("with no data", func() {
			It("should trigger an error", func() {
				_, err := form3client.CreateAccount("")
				Expect(err).To(MatchError(form3client.ClientError(form3client.ErrNoAccountData)))
			})
		})

		Context("With data", func() {
			id := "ad27e265-9605-4b4b-a0e5-3004ea9cc8d9"
			version := "0"

			BeforeEach(func() {
				form3client.DeleteAccount(id, version)
			})

			It("should send a correct response", func() {
				accountData := getBasicAccountData(id, version)
				resp, err := form3client.CreateAccount(accountData)
				if err != nil {
					log.Println(err)
				}
				log.Println(resp)
				Expect(resp.StatusCode).To(Equal(http.StatusCreated))

				_, err = form3client.DeleteAccount(id, version)
				if err != nil {
					log.Println(err)
				}
			})
		})

	})

	Describe("Delete account", func() {
		Context("with ID which does not exist", func() {
			id := "ad27e265-9605-4b4b-a0e5-3000ea9cc8d0"
			version := "0"

			BeforeEach(func() {
				form3client.DeleteAccount(id, version)
			})

			It("should trigger an error", func() {
				resp, err := form3client.DeleteAccount(id, version)
				Expect(err).To(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("with version which does not exist", func() {
			id := "ad27e265-9605-4b4b-a0e5-3006ea9cc8d0"
			BeforeEach(func() {
				accountData := getBasicAccountData(id, "0")
				_, err := form3client.CreateAccount(accountData)
				if err != nil {
					log.Println(err)
				}
			})

			It("should trigger an error", func() {
				resp, err := form3client.DeleteAccount(id, "1")
				Expect(err).To(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusConflict))
			})

			AfterEach(func() {
				_, err := form3client.DeleteAccount(id, "0")
				if err != nil {
					log.Println(err)
				}
			})

		})

		Context("with existing ID and version", func() {
			id := "ad27e265-9605-4b4b-a0e5-3007ea9cc8d0"
			version := "0"

			BeforeEach(func() {
				accountData := getBasicAccountData(id, version)
				_, err := form3client.CreateAccount(accountData)
				if err != nil {
					log.Println(err)
				}
			})

			It("should delete the account properly", func() {
				resp, _ := form3client.DeleteAccount(id, version)
				Expect(resp.StatusCode).To(Equal(http.StatusNoContent))
			})
		})
	})

	Describe("Fetch account", func() {
		Context("with ID which does not exist", func() {
			id := "ad27e265-9604-4b4b-a0e5-3000ea9cc8d0"
			version := "0"

			BeforeEach(func() {
				_, err := form3client.DeleteAccount(id, version)
				if err != nil {
					log.Println(err)
				}
			})

			It("should trigger an error", func() {
				resp, err := form3client.FetchAccount(id)
				Expect(err).To(HaveOccurred())
				Expect(resp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})

		Context("with existing ID", func() {
			id := "ad27e265-9605-6b4b-a0e5-3007ea9cc8d0"
			version := "0"

			BeforeEach(func() {
				accountData := getBasicAccountData(id, version)
				_, err := form3client.CreateAccount(accountData)
				if err != nil {
					log.Println(err)
				}
			})

			It("should delete the account properly", func() {
				resp, _ := form3client.FetchAccount(id)
				Expect(resp.StatusCode).To(Equal(http.StatusOK))
			})

			AfterEach(func() {
				_, err := form3client.DeleteAccount(id, version)
				if err != nil {
					log.Println(err)
				}
			})

		})
	})

})
