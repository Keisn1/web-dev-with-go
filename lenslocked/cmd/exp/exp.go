package main

import (
	"fmt"
	"github.com/keisn1/lenslocked/models"
)

func main() {
	db, err := models.Open(models.DefaultPostgresConfig())
	if err != nil {
		panic(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected!")

	_, err = db.Exec(`
CREATE TABLE IF NOT EXISTS users (
  id SERIAL PRIMARY KEY,
  name TEXT,
  email TEXT UNIQUE NOT NULL
);


CREATE TABLE IF NOT EXISTS orders (
  id SERIAL PRIMARY KEY,
  user_id INT NOT NULL,
  amount INT,
  description TEXT);`)

	if err != nil {
		panic(err)
	}
	fmt.Println("Tables created")

	// name := "kay"
	// email := "kay@email.com"
	// 	_, err = db.Exec(`
	// INSERT INTO users (name, email) VALUES
	//        ($1, $2);
	// `,
	// 		name, email)

	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("User created")

	// 	name = "celine"
	// 	email = "celine@email"
	// 	row := db.QueryRow(`
	// INSERT INTO users  (name, email) VALUES
	// ($1, $2)
	// RETURNING id;`,
	// 		name, email,
	// 	)
	// 	var id int
	// 	err = row.Scan(&id)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println("User created. id = ", id)

	// 	id := 4
	// 	row, _ := db.Query(`
	// SELECT name, email
	//   FROM users
	//   WHERE id=$1;`, id)
	// 	var name string
	// 	var email string
	// 	row.Next()
	// 	err = row.Scan(&name, &email)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	fmt.Println(name, email)

	// 	userID := 4
	// 	for i := 1; i <= 5; i++ {
	// 		amount := i * 100
	// 		desc := fmt.Sprintf("Fake order @%d", i)
	// 		_, err := db.Exec(`
	// INSERT INTO orders(user_id, amount, description)
	// VALUES($1, $2, $3)`, userID, amount, desc)

	//		if err != nil {
	//			panic(err)
	//		}
	//	}
	//
	// fmt.Println("Created fake orders")

	var orders []Order
	userID := 4
	rows, err := db.Query(`
SELECT id, amount, description
FROM orders
WHERE user_id=$1;`,
		userID,
	)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		o := Order{UserID: 1}
		rows.Scan(&o.ID, &o.Amount, &o.Description)
		orders = append(orders, o)
	}
	err = rows.Err()
	if err != nil {
		panic(err)
	}
	fmt.Println(orders)
}

type Order struct {
	ID          int
	UserID      int
	Amount      int
	Description string
}
