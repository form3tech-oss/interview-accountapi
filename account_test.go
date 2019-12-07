package account

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/lioda/interview-accountapi/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const path = "/v1/organisation/accounts"

var anID = uuid.MustParse("87654321-abcd-abcd-abcd-210987654321")
var aNewAccount = model.Account{
	ID:      uuid.Nil,
	Country: "country",
	BankID:  "bankid",
}
var anAccountWithID = model.Account{
	ID:      anID,
	Country: "country",
	BankID:  "bankid",
}
var organizationID = uuid.MustParse("12345678-abcd-efab-cdef-123456789012")

type mockRemoteAPI struct {
	mock.Mock
}

func (m *mockRemoteAPI) Get(path string, queryParams string) string {
	args := m.Called(path, queryParams)
	return args.String(0)
}
func (m *mockRemoteAPI) GetArray(path string, queryParams string) ([]model.Account, error) {
	args := m.Called(path, queryParams)
	return args.Get(0).([]model.Account), nil
}
func (m *mockRemoteAPI) Post(path string, organisationID uuid.UUID, account model.Account) (uuid.UUID, error) {
	args := m.Called(path, organisationID, account)
	return organisationID, args.Error(0)
}

var usedRemoteAPI = new(mockRemoteAPI)

type mockIDGenerator struct {
	mock.Mock
}

func (m *mockIDGenerator) Next() uuid.UUID {
	args := m.Called()
	result := args.Get(0).(uuid.UUID)
	return result
}

var usedIDGenerator = new(mockIDGenerator)

func TestWhenRemoteSucceedThenAccountCreated(t *testing.T) {
	usedIDGenerator.On("Next").Return(anID)
	usedRemoteAPI.On("Post", path, organizationID, anAccountWithID).Times(1).Return(nil)

	accounts := New(usedRemoteAPI, usedIDGenerator, organizationID)
	newID, _ := accounts.Create(aNewAccount)

	assert.Equal(t, anID, newID)
	usedRemoteAPI.AssertCalled(t, "Post", path, organizationID, anAccountWithID)
}

func TestWhenRemoteFailsThenAccountIsNotCreated(t *testing.T) {
	usedIDGenerator.On("Next").Return(anID)
	usedRemoteAPI.On("Post", path, organizationID, anAccountWithID).Return(errors.New("Cannot create account"))

	accounts := New(usedRemoteAPI, usedIDGenerator, organizationID)
	newID, err := accounts.Create(aNewAccount)

	assert.Equal(t, uuid.Nil, newID)
	assert.Equal(t, errors.New("Cannot create account"), err)
	usedRemoteAPI.AssertCalled(t, "Post", path, organizationID, anAccountWithID)
}

func TestWhenRemoteSucceedThenReturnAccountList(t *testing.T) {
	usedRemoteAPI.On("GetArray", path, "").Return([]model.Account{anAccountWithID})

	accounts := New(usedRemoteAPI, usedIDGenerator, organizationID)
	accountList, err := accounts.List()

	assert.Equal(t, nil, err)
	assert.ElementsMatch(t, [1]model.Account{anAccountWithID}, accountList)
	usedRemoteAPI.AssertCalled(t, "GetArray", path, "")
}

//TODO : when error,
