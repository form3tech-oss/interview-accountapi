package main_test

import (
	"testing"

	. "submission"
	. "submission/env"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

func TestFetchInvalidId(t *testing.T) {
	t.Parallel()

	c := NewClient(&NewClientParams{
		BaseUrl: GetServerHost(),
	})

	data, error := c.FetchAccount(FetchAccountParams{
		ID: "1",
	})

	assert.Nil(t, data)
	assert.NotNil(t, error)
	assert.EqualError(t, error, "id is not a valid uuid")
}
