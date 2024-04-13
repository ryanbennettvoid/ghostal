package postgres_db_operator

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

const ValidPostgresURL = "postgresql://admin:pw@localhost/main"

func TestUnit_ParsePostgresURL(t *testing.T) {
	pgURL, err := ParsePostgresURL(ValidPostgresURL)
	assert.NoError(t, err)
	assert.NotNil(t, pgURL)

	assert.Equal(t, "admin", pgURL.Username())
	assert.Equal(t, "main", pgURL.DBName())
}
