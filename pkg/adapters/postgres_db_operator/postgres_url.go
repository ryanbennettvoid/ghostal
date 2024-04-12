package postgres_db_operator

import (
	"net/url"
)

type PostgresURL struct {
	dbURL *url.URL
}

func (p *PostgresURL) Clone() *url.URL {
	clone := *p.dbURL
	clone.User = &(*p.dbURL.User)
	return &clone
}

func (p *PostgresURL) Username() string {
	return p.dbURL.User.Username()
}

func (p *PostgresURL) DBName() string {
	return p.dbURL.Path[1:]
}

func ParsePostgresURL(dbURL string) (*PostgresURL, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, err
	}
	return &PostgresURL{
		dbURL: u,
	}, nil
}
