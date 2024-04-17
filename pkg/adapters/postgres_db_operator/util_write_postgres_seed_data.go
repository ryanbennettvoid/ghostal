package postgres_db_operator

import "fmt"

func WritePostgresSeedData(dbURL string, tableName string) {
	PostgresRunQuery(dbURL, fmt.Sprintf(`
		CREATE TABLE %s (
			id SERIAL PRIMARY KEY,
			make VARCHAR(50),
			model VARCHAR(50),
			year INT,
			color VARCHAR(50)
		);
	`, tableName))
	PostgresRunQuery(dbURL, fmt.Sprintf(`
		INSERT INTO %s (make, model, year, color) VALUES
		('Toyota', 'Camry', 2022, 'Black'),
		('Ford', 'Mustang', 2021, 'Red'),
		('Honda', 'Civic', 2020, 'Blue'),
		('Tesla', 'Model 3', 2023, 'White'),
		('Chevrolet', 'Impala', 2019, 'Silver');
	`, "vehicles"))
}
