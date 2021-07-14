package api_test

import (
	"api"
	"fmt"
	"models"
	"testing"
)

const SAMPLE = "testdata/account_gb.json"
const SAMPLE_ID = "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"

func setUp() {
	ACCOUNT_SAMPLE := models.CreateAccountFromJSONFile(SAMPLE)
	api.Create(&ACCOUNT_SAMPLE)
}

func tearDown() {
	api.Delete(SAMPLE_ID, 0)
}

func setUpFakeUrl() string {
	backup := api.URL
	api.URL = "https://api.form3.tech/"
	return backup
}

func tearDownFakeUrl(backup string) {
	api.URL = backup
}

func TestCheck(t *testing.T) {
	err := api.Check()
	if err != nil {
		t.Errorf("Service unreachable")
	}
	backup := setUpFakeUrl()
	err = api.Check()
	if err.Error() != "API Unavailable" {
		t.Errorf("We should not be authenticated")
	}
	tearDownFakeUrl(backup)
}

func TestCreate(t *testing.T) {
	ACCOUNT_SAMPLE := models.CreateAccountFromJSONFile(SAMPLE)
	id, err := api.Create(&ACCOUNT_SAMPLE)
	// Nominal case
	if id != SAMPLE_ID && err != nil {
		t.Errorf("Invalid id, got: %s, want: %s.", id, SAMPLE_ID)
	}
	// Create twice the ressource will result into an error (duplicate ressource)
	_, err = api.Create(&ACCOUNT_SAMPLE)
	if err.Error() != fmt.Sprintf("Record \"%s\" not created", SAMPLE_ID) {
		t.Errorf("Invalid error, got: %s", err)
	}
	tearDown()
	backup := setUpFakeUrl()
	_, err = api.Create(&ACCOUNT_SAMPLE)
	if err.Error() != "API Unavailable" {
		t.Errorf("We should not be authenticated")
	}
	tearDownFakeUrl(backup)
}

func TestFetch(t *testing.T) {
	setUp()
	// Nominal case
	account, _ := api.Fetch(SAMPLE_ID)
	if account.ID != SAMPLE_ID && *account.Version != 0 {
		t.Errorf("Invalid id, got: %s, want: %s.", account.ID, SAMPLE_ID)
	}
	tearDown()

	tearDown()
	backup := setUpFakeUrl()
	_, err := api.Fetch(SAMPLE_ID)
	if err.Error() != "API Unavailable" {
		t.Errorf("We should not be authenticated")
	}
	tearDownFakeUrl(backup)
}

func TestDelete(t *testing.T) {
	setUp()
	// Delete vith invalid version
	err := api.Delete(SAMPLE_ID, 1)
	if err.Error() != fmt.Sprintf("Invalid version \"1\" for record \"%s\". Hint: Try with version \"0\"", SAMPLE_ID) {
		t.Errorf("Invalid error, got: %s", err)
	}
	// Delete with valid version but invalid record id
	err = api.Delete("id", 0)
	if err.Error() != "Record \"id\" not found" {
		t.Errorf("Invalid error, got: %s", err)
	}
	// Nominal case
	err = api.Delete(SAMPLE_ID, 0)
	if err != nil {
		t.Errorf("Invalid error, got: %s", err)
	}
	// Once deleted, the ressource is no longer accessible
	_, err = api.Fetch(SAMPLE_ID)
	if err.Error() != fmt.Sprintf("Record \"%s\" not found", SAMPLE_ID) {
		t.Errorf("Invalid error, got: %s", err)
	}
	_, err = api.Fetch(SAMPLE_ID)
	if err.Error() != fmt.Sprintf("Record \"%s\" not found", SAMPLE_ID) {
		t.Errorf("Invalid error, got: %s", err)
	}
	backup := setUpFakeUrl()
	err = api.Delete(SAMPLE_ID, 0)
	if err.Error() != "API Unavailable" {
		t.Errorf("We should not be authenticated")
	}
	tearDownFakeUrl(backup)
}

/*func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 1 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}*/
