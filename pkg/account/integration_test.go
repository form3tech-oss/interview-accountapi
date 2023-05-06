package account

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strconv"
	"testing"

	_ "github.com/lib/pq"
)

func TestFetchAccountIntegration(t *testing.T) {

	account := createAccount()

	db, err := openDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	err = initDB(db, []AccountData{*account})
	if err != nil {
		t.Fatal(err)
	}
}

func TestCreateAccountIntegration(t *testing.T) {
	db, err := openDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	initDB(db, []AccountData{})
}

func TestDeleteAccountIntegration(t *testing.T) {
	db, err := openDB()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	initDB(db, []AccountData{})
}

func openDB() (*sql.DB, error) {
	const (
		host     = "localhost"
		port     = 5432
		user     = "interview_accountapi_user"
		password = "123"
		dbname   = "interview_accountapi"
	)

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

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
			('` + id + `', '` + organisationID + `', ` + strconv.FormatInt(version, 10) + `, false, false, current_timestamp, current_timestamp,'` + string(record) + `'::jsonb , DEFAULT)`)
		if err != nil {
			return err
		}
	}
	fmt.Println("[DB] successfully initialized")

	return nil
}

// INSERT INTO "films" (id, organisation_id, version, is_deleted, is_locked, created_on, modified_on, record, pagination_id)
// VALUES ("ad27e265-9605-4b4b-a0e5-3003ea9cc4dc", "eb0bd6f5-c3f5-44b2-b677-acd23cdde73c", 0, false, false, current_timestamp, current_timestamp, record, DEFAULT);
