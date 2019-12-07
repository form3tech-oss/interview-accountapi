package account

import (
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	interviewaccountapi "github.com/lioda/interview-accountapi"
	"github.com/lioda/interview-accountapi/model"
)

const bankID = "BANKID"
const bic = "BICCOD00"
const bankIDCode = "GBDSC"

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

type E2ETestSuite struct {
	suite.Suite
	organizationID uuid.UUID
	accounts       interviewaccountapi.Accounts
}

// SetupTest prepares all tests
func (suite *E2ETestSuite) SetupTest() {
	apiURL := GetEnvironnmentURL()
	suite.organizationID = uuid.MustParse("110e8400-e29b-11d4-a716-446655440000")
	suite.accounts = interviewaccountapi.NewDefault(apiURL, suite.organizationID)
}

// TestCreate an account
func (suite *E2ETestSuite) TestWhenCreateGBAccountThenItAppearsInList() {
	list, err := suite.accounts.List()
	assert.Equal(suite.T(), 0, len(list))
	assert.Equal(suite.T(), nil, err)
	suite.accounts.Create(model.NewGbAccount(bankID, bic, bankIDCode))
	list, err = suite.accounts.List()
	assert.Equal(suite.T(), 1, len(list))
	assert.Equal(suite.T(), nil, err)
}

// TestE2ETestSuite runs the suite
func TestE2ETestSuite(t *testing.T) {
	// suite.Run(t, new(E2ETestSuite))
}
