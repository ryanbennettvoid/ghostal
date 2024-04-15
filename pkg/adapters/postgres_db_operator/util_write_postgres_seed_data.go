package postgres_db_operator

func WritePostgresSeedData(dbURL string) {
	PostgresRunQuery(dbURL, `
		CREATE TABLE vehicles (
			id SERIAL PRIMARY KEY,
			make VARCHAR(50),
			model VARCHAR(50),
			year INT,
			color VARCHAR(50)
		);
	`)
	PostgresRunQuery(dbURL, `
		INSERT INTO vehicles (make, model, year, color) VALUES
		('Toyota', 'Camry', 2022, 'Black'),
		('Ford', 'Mustang', 2021, 'Red'),
		('Honda', 'Civic', 2020, 'Blue'),
		('Tesla', 'Model 3', 2023, 'White'),
		('Chevrolet', 'Impala', 2019, 'Silver');
	`)
}
