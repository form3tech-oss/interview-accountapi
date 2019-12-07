package account

import (
	"github.com/google/uuid"
	"github.com/lioda/interview-accountapi/model"
)

type remoteAPI interface {
	Get(path string, queryParams string) string
	GetArray(path string, queryParams string) ([]model.Account, error)
	Post(path string, organisationID uuid.UUID, account model.Account) (uuid.UUID, error)
}

type idGenerator interface {
	Next() uuid.UUID
}

// Abs returns absolute value
func Abs(i int) int {
	if i < 0 {
		return i * -1
	}
	return i
}

// Accounts manage the API
type Accounts struct {
	organisationID uuid.UUID
	remote         remoteAPI
	idGen          idGenerator
}

// NewDefault creates a new Accounts with all default values
func NewDefault(baseURL string, organisationID uuid.UUID) Accounts { // TODO move in another file
	return Accounts{remote: NewHTTPRemote(baseURL), idGen: NewRandomIDGenerator(), organisationID: organisationID}
}

// New create a new Accounts
func New(remote remoteAPI, idGen idGenerator, organisationID uuid.UUID) Accounts {
	return Accounts{remote: remote, idGen: idGen, organisationID: organisationID}
}

// List return the account list
func (a Accounts) List() ([]model.Account, error) {
	accounts, _ := a.remote.GetArray("/v1/organisation/accounts", "")
	return accounts, nil
}

// Create an account with a generated ID
func (a Accounts) Create(acc model.Account) (uuid.UUID, error) {
	id := a.idGen.Next()
	acc.ID = id
	_, err := a.remote.Post("/v1/organisation/accounts", a.organisationID, acc)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
