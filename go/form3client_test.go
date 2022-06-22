package form3client_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/ahermant/form3client"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var _ = Describe("Form3client", func() {

	Describe("Get base URL", func() {
		Context("No env set", func() {
			initialEnv := os.Getenv("ACCOUNT_API_BASE_URL")
			BeforeEach(func() {
				os.Unsetenv("ACCOUNT_API_BASE_URL")
			})

			It("should provide a proper url", func() {
				url := form3client.GetBaseUrl()
				Expect(url).To(Equal(form3client.DefaultBaseUrl))
			})
			AfterEach(func() {
				os.Setenv("ACCOUNT_API_BASE_URL", initialEnv)
			})
		})

		Context("Env set", func() {
			initialEnv := os.Getenv("ACCOUNT_API_BASE_URL")
			testUrl := "http://test"
			BeforeEach(func() {
				os.Setenv("ACCOUNT_API_BASE_URL", testUrl)
			})

			It("should provide a proper url", func() {
				url := form3client.GetBaseUrl()
				Expect(url).To(Equal(testUrl))
			})

			AfterEach(func() {
				os.Setenv("ACCOUNT_API_BASE_URL", initialEnv)
			})

		})
	})

	Describe("Client Error", func() {
		Context("with proper error", func() {
			It("should return a client error", func() {
				err := form3client.ClientError(errors.New("test"))
				Expect(fmt.Sprint(err)).To(Equal("client error: test"))
			})
		})
	})

	Describe("HTTP responses", func() {
		Context("with proper response and code", func() {
			It("should return a ClientResponse", func() {
				body := "Hello world"
				httpResponse := &http.Response{
					Status:        "200 OK",
					StatusCode:    200,
					Proto:         "HTTP/1.1",
					ProtoMajor:    1,
					ProtoMinor:    1,
					Body:          ioutil.NopCloser(bytes.NewBufferString(body)),
					ContentLength: int64(len(body)),
					Header:        make(http.Header, 0),
				}

				resp, _ := form3client.ParseResponse(httpResponse, 200)
				httpResponse.Body.Close()
				Expect(resp).To(Equal(form3client.ClientResponse{200, body}))
			})
		})
		Context("with bad body format", func() {
			It("should return an empty ClientResponse + error", func() {
				// We need to mock the io.Reader interface to test ioutil.ReadAll
			})
		})
	})

	Describe("Deletion request", func() {
		var server *ghttp.Server
		var url string
		parameters := "version=0"
		accountID := "ad27e265-9605-4b4b-a0e5-3004ea9cc8d9"
		endpoint := form3client.AccountsEndpoint + "/" + accountID
		BeforeEach(func() {
			server = ghttp.NewServer()
			url = server.URL() + endpoint
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", endpoint, parameters),
					ghttp.RespondWith(http.StatusNoContent, ``),
				),
			)
		})

		AfterEach(func() {
			server.Close()
		})

		Context("with proper url", func() {
			It("should return a ClientResponse", func() {
				resp, err := form3client.DeletionRequest(url + "?" + parameters)
				Expect(err).To(BeNil())
				Expect(resp).To(Equal(form3client.ClientResponse{204, ""}))
			})
		})
		Context("with wrong url", func() {
			It("should return an error", func() {
				_, err := form3client.DeletionRequest("")
				log.Println(err)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Fetch request", func() {
		var server *ghttp.Server
		var url string
		accountID := "ad27e265-9605-4b4b-a0e5-3004ea9cc8d9"
		endpoint := form3client.AccountsEndpoint + "/" + accountID
		body := `{
			"data":{
				"type": "accounts",
				"id": "` + accountID + `",
				"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
				"version" :"0",
				"attributes": {
					"name": ["ahtest"],
					"country": "GB"
				}
			}
		}`

		BeforeEach(func() {
			server = ghttp.NewServer()
			url = server.URL() + endpoint
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", endpoint),
					ghttp.RespondWith(http.StatusOK, body),
				),
			)
		})

		AfterEach(func() {
			server.Close()
		})

		Context("with proper url", func() {
			It("should return a ClientResponse", func() {
				resp, err := form3client.FetchRequest(url)
				Expect(err).To(BeNil())
				Expect(resp).To(Equal(form3client.ClientResponse{200, body}))
			})
		})
		Context("with wrong url", func() {
			It("should return an error", func() {
				_, err := form3client.FetchRequest("")
				log.Println(err)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("Post request", func() {
		var server *ghttp.Server
		var url string
		accountID := "ad27e265-9605-4b4b-a0e5-3004ea9cc8d9"
		endpoint := form3client.AccountsEndpoint
		body := `{
			"data":{
				"type": "accounts",
				"id": "` + accountID + `",
				"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
				"version" :"0",
				"attributes": {
					"name": ["ahtest"],
					"country": "GB"
				}
			}
		}`

		BeforeEach(func() {
			server = ghttp.NewServer()
			url = server.URL() + endpoint
			server.AppendHandlers(
				ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", endpoint),
					ghttp.RespondWith(http.StatusCreated, body),
				),
			)
		})

		AfterEach(func() {
			server.Close()
		})

		Context("with proper url", func() {
			It("should return a ClientResponse", func() {
				resp, err := form3client.PostRequest(url, body)
				Expect(err).To(BeNil())
				Expect(resp).To(Equal(form3client.ClientResponse{201, body}))
			})
		})
		Context("with wrong url", func() {
			It("should return an error", func() {
				_, err := form3client.PostRequest("", body)
				log.Println(err)
				Expect(err).To(HaveOccurred())
			})
		})
		Context("with empty body", func() {
			It("should return an error", func() {
				_, err := form3client.PostRequest(url, "")
				log.Println(err)
				Expect(err).To(HaveOccurred())
			})
		})
	})

})
