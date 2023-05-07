package account

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

var testConfig = ReadConfig(os.Getenv("TEST_CONTEXT"))

var baseUrl = testConfig["api"]["baseUrl"]
var host = testConfig["database"]["host"]
var port = testConfig["database"]["port"]
var user = testConfig["database"]["user"]
var password = testConfig["database"]["password"]
var dbname = testConfig["database"]["dbname"]

func TestFetchAccountIntegration(t *testing.T) {

	account := createAccount()

	type testCase struct {
		name     string
		id       string
		expected *AccountData
		err      error
	}

	testCases := []testCase{
		{
			"fetch existent account",
			account.ID,
			account,
			nil,
		},
		{
			"invalid uuid",
			"1234",
			nil,
			&ErrorResponse{Code: 400, Message: "id is not a valid uuid"},
		},
		{
			"fetch non existent account",
			"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			nil,
			&ErrorResponse{Code: 404, Message: "record eb0bd6f5-c3f5-44b2-b677-acd23cdde73c does not exist"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db, err := openDB()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			err = initDB(db, []AccountData{*account})
			if err != nil {
				t.Fatal(err)
			}

			accountClient := NewAccountClient(&Config{BaseUrl: baseUrl, Version: "v1"})
			accountData, err := accountClient.FetchAccount(context.Background(), tc.id)
			assert.Equal(t, tc.err, err)
			assert.Equal(t, tc.expected, accountData)
		})
	}
}

func TestCreateAccountIntegration(t *testing.T) {
	existentAccount := createAccount()

	newAccount := createAccount()
	newAccount.ID = "8ceac1ce-ec44-11ed-a05b-0242ac120003"

	invalidUUIDAccount := createAccount()
	invalidUUIDAccount.ID = "1234"

	invalidCode := "XYZ"
	invalidCodeAccount := createAccount()
	invalidCodeAccount.ID = "8ceac610-ec44-11ed-a05b-0242ac120003"
	invalidCodeAccount.Attributes.Country = &invalidCode

	type testCase struct {
		name     string
		account  *AccountData
		expected *AccountData
		err      error
	}

	testCases := []testCase{
		{
			"create account",
			newAccount,
			newAccount,
			nil,
		},
		{
			"create duplicated account",
			account,
			nil,
			&ErrorResponse{Code: 409, Message: "Account cannot be created as it violates a duplicate constraint"},
		},
		{
			"create nil account",
			nil,
			nil,
			&ErrorResponse{Code: 400, Message: "invalid account data"},
		},
		{
			"create empty account",
			&AccountData{},
			nil,
			&ErrorResponse{Code: 400, Message: "validation failure list:\nvalidation failure list:\nattributes in body is required\nid in body is required\norganisation_id in body is required\ntype in body is required"},
		},
		{
			"invalid uuid",
			invalidUUIDAccount,
			nil,
			&ErrorResponse{Code: 400, Message: "validation failure list:\nvalidation failure list:\nid in body must be of type uuid: \"1234\""},
		},
		{
			"invalid code",
			invalidCodeAccount,
			nil,
			&ErrorResponse{Code: 400, Message: "validation failure list:\nvalidation failure list:\nvalidation failure list:\ncountry in body should match '^[A-Z]{2}$'"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db, err := openDB()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			err = initDB(db, []AccountData{*existentAccount})
			if err != nil {
				t.Fatal(err)
			}

			accountClient := NewAccountClient(&Config{BaseUrl: baseUrl, Version: "v1"})
			accountData, err := accountClient.CreateAccount(context.Background(), tc.account)
			assert.Equal(t, tc.err, err)

			if err == nil {
				accountData.CreatedOn = tc.expected.CreatedOn
				assert.Equal(t, tc.expected, accountData)
			}
		})
	}
}

func TestDeleteAccountIntegration(t *testing.T) {
	account := createAccount()

	type testCase struct {
		name    string
		id      string
		version int64
		err     error
	}

	testCases := []testCase{
		{
			"delete existent account",
			account.ID,
			int64(0),
			nil,
		},
		{
			"delete existent account with incorrect version",
			account.ID,
			int64(1),
			&ErrorResponse{Code: 409, Message: "invalid version"},
		},
		{
			"invalid uuid",
			"1234",
			int64(0),
			&ErrorResponse{Code: 400, Message: "id is not a valid uuid"},
		},
		{
			"delete non existent account",
			"eb0bd6f5-c3f5-44b2-b677-acd23cdde73c",
			int64(0),
			&ErrorResponse{Code: 404, Message: ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			db, err := openDB()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			err = initDB(db, []AccountData{*account})
			if err != nil {
				t.Fatal(err)
			}

			accountClient := NewAccountClient(&Config{BaseUrl: baseUrl, Version: "v1"})
			err = accountClient.DeleteAccount(context.Background(), tc.id, tc.version)
			assert.Equal(t, tc.err, err)

			if err == nil {
				_, err = accountClient.FetchAccount(context.Background(), tc.id)
				msg := fmt.Sprintf("record %s does not exist", tc.id)
				assert.Equal(t, &ErrorResponse{Code: 404, Message: msg}, err)
			}

		})
	}
}

func openDB() (*sql.DB, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

	fmt.Printf("[DB] connecting to %s\n", psqlInfo)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	fmt.Println("[DB] successfully connected")
	return db, nil

}

func initDB(db *sql.DB, accounts []AccountData) error {

	_, err := db.Exec(`DELETE FROM "Account"`)
	if err != nil {
		return err
	}

	fmt.Println("[DB] successfully cleaned up")

	for _, a := range accounts {
		id := a.ID
		organisationID := a.OrganisationID
		var version int64 = 0
		if a.Version != nil {
			version = *a.Version
		}

		record, err := json.MarshalIndent(a.Attributes, "", "  ")
		if err != nil {
			return err
		}

		_, err = db.Exec(`
		INSERT INTO "Account"
			(id, organisation_id, version, is_deleted, is_locked, created_on, modified_on, record, pagination_id)
		VALUES
			('` + id + `', '` + organisationID + `', ` + strconv.FormatInt(version, 10) + `, false, false, '` + a.CreatedOn + `', current_timestamp,'` + string(record) + `'::jsonb , DEFAULT)`)
		if err != nil {
			return err
		}
	}
	fmt.Println("[DB] successfully initialized")

	return nil
}

func TestReadTestConfig(t *testing.T) {

	m := ReadConfig("default")
	assert.Equal(t, m["api"]["baseUrl"], "http://localhost:8080")
	assert.Equal(t, m["database"]["host"], "localhost")
	assert.Equal(t, m["database"]["port"], "5432")
	assert.Equal(t, m["database"]["user"], "interview_accountapi_user")
	assert.Equal(t, m["database"]["password"], "123")
	assert.Equal(t, m["database"]["dbname"], "interview_accountapi")

	m = ReadConfig("container")
	assert.Equal(t, m["api"]["baseUrl"], "http://accountapi:8080")
	assert.Equal(t, m["database"]["host"], "postgresql")
	assert.Equal(t, m["database"]["port"], "5432")
	assert.Equal(t, m["database"]["user"], "interview_accountapi_user")
	assert.Equal(t, m["database"]["password"], "123")
	assert.Equal(t, m["database"]["dbname"], "interview_accountapi")
}

func ReadConfig(env string) map[string]map[string]string {

	b, err := os.ReadFile("test_config.json")
	if err != nil {
		panic(err)
	}
	x := map[string]map[string]map[string]string{}
	err = json.Unmarshal(b, &x)
	if err != nil {
		panic(err)
	}

	if config, ok := x[env]; !ok {
		return x["local"]
	} else {
		return config
	}
}
