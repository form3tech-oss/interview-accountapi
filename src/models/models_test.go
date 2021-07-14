package models_test

import (
	"fmt"
	"models"
	"os"
	"testing"
)

func TestAccountCreationFromFile(t *testing.T) {
	data_file := "testdata/account_gb.json"
	data := models.CreateAccountFromJSONFile(data_file)
	expected_id := "ad27e265-9605-4b4b-a0e5-3003ea9cc4dc"
	if data.ID != expected_id {
		t.Errorf("Invalid id, got: %s, want: %s.", data.ID, expected_id)
	}
}

func TestMain(m *testing.M) {
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
}
