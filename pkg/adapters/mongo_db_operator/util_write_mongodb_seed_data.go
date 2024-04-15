package mongo_db_operator

import (
	"context"
	"log"
)

type MockCar struct {
	Make  string `bson:"make"`
	Model string `bson:"model"`
	Year  int    `bson:"year"`
	Color string `bson:"color"`
}

func WriteMongoDBSeedData(dbURL string) {

	vehicles := []interface{}{
		MockCar{"Toyota", "Camry", 2022, "Black"},
		MockCar{"Ford", "Mustang", 2021, "Red"},
		MockCar{"Honda", "Civic", 2020, "Blue"},
		MockCar{"Tesla", "Model 3", 2023, "White"},
		MockCar{"Chevrolet", "Impala", 2019, "Silver"},
	}

	collection, cleanup := GetMongoDBCollection(dbURL)
	defer cleanup()

	_, err := collection.InsertMany(context.TODO(), vehicles)
	if err != nil {
		log.Fatal(err)
	}
}
