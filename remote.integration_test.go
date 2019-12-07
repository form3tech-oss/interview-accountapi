package account

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"

	"github.com/lioda/interview-accountapi/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func GetEnvironnmentURL() string {
	host, present := os.LookupEnv("API_HOST")
	if !present {
		return "http://localhost:8080"
	}
	port, present := os.LookupEnv("API_PORT")
	if !present {
		return "http://localhost:8080"
	}
	return fmt.Sprintf("http://%v:%v", host, port)
}

// RemoteIntegrationTestSuite runs the integration tests on remote HTTP implementation (needs local API started)
type RemoteIntegrationTestSuite struct {
	suite.Suite
	httpRemote HTTPRemoteAPI
}

// SetupTest assign variables
func (suite *RemoteIntegrationTestSuite) SetupTest() {
	suite.httpRemote = NewHTTPRemote(GetEnvironnmentURL())
}

func (suite *RemoteIntegrationTestSuite) TestWhenGetThenReturnData() {
	arr, err := suite.httpRemote.GetArray("/v1/organisation/accounts", "")
	assert.ElementsMatch(suite.T(), []model.Account{}, arr)
	assert.ElementsMatch(suite.T(), nil, err)

	account := model.Account{ // TODO create a new GBACCOUNT
		ID:                          uuid.New(),
		Country:                     "GB",
		AccountClassification:       "Personal",
		AlternativeBankAccountNames: [3]string{"ABC", "CDE", "FGH"},
	}
	id, err := suite.httpRemote.Post("/v1/organisation/accounts", uuid.New(), account)

	assert.Equal(suite.T(), nil, err)
	assert.Equal(suite.T(), account.ID, id)

	arr, err = suite.httpRemote.GetArray("/v1/organisation/accounts", "")
	assert.ElementsMatch(suite.T(), []model.Account{account}, arr)
	//TODO DELETE
}

// TestRemoteIntegrationTestSuite run all tests in suite
func TestRemoteIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(RemoteIntegrationTestSuite))
}
