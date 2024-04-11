package postgres_db_adapter

import (
	"net/url"
)

type PostgresURL struct {
	dbURL *url.URL
}

func (p *PostgresURL) cloneURL() *url.URL {
	clone := *p.dbURL
	clone.User = &(*p.dbURL.User)
	return &clone
}

func (p *PostgresURL) WithoutSSL() *url.URL {
	clone := p.cloneURL()
	params := clone.Query()
	params.Set("sslmode", "disable")
	clone.RawQuery = params.Encode()
	return clone
}

func (p *PostgresURL) Scheme() string {
	return p.dbURL.Scheme
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
