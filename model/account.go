package model

import (
	"github.com/google/uuid"
)

// Account is the model of an account
type Account struct {
	ID                          uuid.UUID
	Country                     string
	BaseCurrency                string
	BankID                      string
	BankIDCode                  string
	AccountNumber               string
	Bic                         string
	Iban                        string
	CustomerID                  string
	Title                       string    //[40]
	FirstName                   string    //[40]
	BankAccountName             string    // [140]
	AlternativeBankAccountNames [3]string //[140]
	AccountClassification       string
	JointAccount                bool
	AccountMatchingOptOut       bool
	SecondaryIdentification     string // [140]
}
