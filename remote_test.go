package account

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/lioda/interview-accountapi/model"

	"github.com/stretchr/testify/mock"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

// RemoteTestSuite runs the unit tests on remote HTTP implementation
type RemoteTestSuite struct {
	suite.Suite
	server     *httptest.Server
	httpRemote HTTPRemoteAPI
}

type mockHandler struct {
	mock.Mock
}

func (mock *mockHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadAll(r.Body)
	header := r.Header["Content-Type"]
	contentType := ""
	if len(header) > 0 {
		contentType = header[0]
	}
	args := mock.MethodCalled("ServeHTTP", r.URL.String(), r.Method, contentType, string(body))
	w.WriteHeader(args.Int(0))
	fmt.Fprint(w, args.String(1))
}

var handler *mockHandler

var anAccount = model.Account{
	ID:                          uuid.MustParse("12345678-abcd-abcd-abcd-123456789012"),
	Country:                     "GB",
	BaseCurrency:                "GBP",
	AccountNumber:               "41426819",
	BankID:                      "400300",
	BankIDCode:                  "GBDSC",
	Bic:                         "NWBKGB22",
	Iban:                        "GB11NWBK40030041426819",
	Title:                       "Ms",
	FirstName:                   "Samantha",
	BankAccountName:             "AccountOfSamanthaHolder",
	AlternativeBankAccountNames: [3]string{"SamHolder"},
	AccountClassification:       "Personal",
	JointAccount:                false,
	AccountMatchingOptOut:       false,
	SecondaryIdentification:     "A1B2C3D4",
}

var jsonAccount = fmt.Sprintf(`{
	"data": {
		"type": "accounts",
		"id": "%v",
		"organisation_id": "%v",
		"version": "0",
		"attributes": {
			"country": "%v",
			"base_currency": "%v",
			"bank_id": "%v",
			"bank_id_code": "%v",
			"account_number": "%v",
			"bic": "%v",
			"iban": "%v",
			"customer_id": "%v",
			"title": "%v",
			"first_name": "%v",
			"bank_account_name": "%v",
			"alternative_bank_account_names": [
				"%v", "", ""
			],
			"account_classification": "%v",
			"joint_account": %v,
			"account_matching_opt_out": %v,
			"secondary_identification": "%v"
		}
	}
}`,
	anAccount.ID,
	organizationID,
	anAccount.Country,
	anAccount.BaseCurrency,
	anAccount.BankID,
	anAccount.BankIDCode,
	anAccount.AccountNumber,
	anAccount.Bic,
	anAccount.Iban,
	anAccount.CustomerID,
	anAccount.Title,
	anAccount.FirstName,
	anAccount.BankAccountName,
	anAccount.AlternativeBankAccountNames[0],
	anAccount.AccountClassification,
	anAccount.JointAccount,
	anAccount.AccountMatchingOptOut,
	anAccount.SecondaryIdentification)

// SetupTest assign variables
func (suite *RemoteTestSuite) SetupTest() {
	handler = new(mockHandler)
	suite.server = httptest.NewServer(handler)
	suite.httpRemote = NewHTTPRemote(suite.server.URL)
}

func (suite *RemoteTestSuite) TestWhenGetArrayThenSendConcatPathAndParamsBeforeUnmarshallingResponse() {
	handler.On("ServeHTTP", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(200, `{
		"data": [
			{
				"type": "accounts",
				"id": "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc",
				"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
				"version": 0,
				"attributes": {
					"account_number": "41426819",
					"bank_id": "400300",
					"bank_id_code": "GBDSC"
				}
			},
			{
				"type": "accounts",
				"id": "ea6239c1-99e9-42b3-bca1-92f5c068da6b",
				"organisation_id": "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
				"version": 0,
				"attributes": {
					"country": "GB",
					"base_currency": "GBP",
					"account_number": "41426819",
					"bank_id": "400300",
					"bank_id_code": "GBDSC",
					"bic": "NWBKGB22",
					"iban": "GB11NWBK40030041426819",
					"title": "Ms",
					"first_name": "Samantha",
					"bank_account_name": "Samantha Holder",
					"alternative_bank_account_names": [
						"Sam Holder"
					],
					"account_classification": "Personal",
					"joint_account": false,
					"account_matching_opt_out": false,
					"secondary_identification": "A1B2C3D4"
				}
			}
		]
	}`)
	arr, err := suite.httpRemote.GetArray("/path/to/resource", "queryParam1=value1&queryParam2=value2")
	assert.ElementsMatch(suite.T(), [2]model.Account{
		model.Account{
			ID:            uuid.MustParse("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"),
			AccountNumber: "41426819",
			BankID:        "400300",
			BankIDCode:    "GBDSC",
		},
		model.Account{
			ID:                    uuid.MustParse("ea6239c1-99e9-42b3-bca1-92f5c068da6b"),
			AccountClassification: "Personal",
			AccountMatchingOptOut: false,
			AccountNumber:         "41426819",
			BankAccountName:       "Samantha Holder",
			BankID:                "400300",
			BankIDCode:            "GBDSC",
			BaseCurrency:          "GBP",
			Bic:                   "NWBKGB22",
			Country:               "GB",
			Iban:                  "GB11NWBK40030041426819",
			FirstName:             "Samantha",
			AlternativeBankAccountNames: [3]string{
				"Sam Holder",
			},
			JointAccount:            false,
			SecondaryIdentification: "A1B2C3D4",
			Title:                   "Ms",
		},
	}, arr)
	assert.Equal(suite.T(), nil, err)

	handler.AssertCalled(suite.T(), "ServeHTTP", "/path/to/resource?queryParam1=value1&queryParam2=value2", "GET", mock.Anything, mock.Anything)
}
func (suite *RemoteTestSuite) TestWhenGetArrayFailsThenReturnsError() {
	handler.On("ServeHTTP", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(500, `reason`)
	arr, err := suite.httpRemote.GetArray("/path/to/resource", "queryParam1=value1&queryParam2=value2")
	assert.Equal(suite.T(), errors.New("reason"), err)
	assert.Equal(suite.T(), []model.Account{}, arr)

	handler.AssertCalled(suite.T(), "ServeHTTP", "/path/to/resource?queryParam1=value1&queryParam2=value2", "GET", mock.Anything, mock.Anything)
}
func (suite *RemoteTestSuite) TestWhenPostThenMarshallData() {
	handler.On("ServeHTTP", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(201, jsonAccount)

	id, err := suite.httpRemote.Post("/path/to/resource", organizationID, anAccount)
	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), anAccount.ID, id)
	handler.AssertCalled(suite.T(), "ServeHTTP", "/path/to/resource", "POST", "application/json", rawJSON(jsonAccount))
}

func (suite *RemoteTestSuite) TestWhenPostNotReturn200ThenError() {
	handler.On("ServeHTTP", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(500, "reason")

	id, err := suite.httpRemote.Post("/path/to/resource", organizationID, anAccount)
	assert.Equal(suite.T(), errors.New("reason"), err)
	assert.Equal(suite.T(), uuid.Nil, id)
	handler.AssertCalled(suite.T(), "ServeHTTP", "/path/to/resource", "POST", "application/json", rawJSON(jsonAccount))
}

func rawJSON(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(s, "\n", ""), "\t", ""), " ", "")
}

// TestRemoteTestSuite run all tests in suite
func TestRemoteTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteTestSuite))
}
