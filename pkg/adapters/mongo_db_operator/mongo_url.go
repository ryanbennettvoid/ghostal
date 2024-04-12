package mongo_db_operator

import "net/url"

type MongoURL struct {
	dbURL *url.URL
}

func (p *MongoURL) Clone() *url.URL {
	clone := *p.dbURL
	clone.User = &(*p.dbURL.User)
	return &clone
}

func (p *MongoURL) Username() string {
	return p.dbURL.User.Username()
}

func (p *MongoURL) DBName() string {
	return p.dbURL.Path[1:]
}

func ParseMongoURL(dbURL string) (*MongoURL, error) {
	u, err := url.Parse(dbURL)
	if err != nil {
		return nil, err
	}
	return &MongoURL{
		dbURL: u,
	}, nil
}
