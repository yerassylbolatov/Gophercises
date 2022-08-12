package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"regexp"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "yerassyl"
	password = "password"
	dbname   = "gophercises_phone"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	must(err)
	err = resetDB(db, dbname)
	must(err)
	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	err = db.Ping()
	must(err)

	err = createTable(db)
	must(err)

	var id int

	_, err = insertPhone(db, "1234567890")
	must(err)
	_, err = insertPhone(db, "123 456 7891")
	must(err)
	id, err = insertPhone(db, "(123) 456 7892")
	must(err)
	_, err = insertPhone(db, "(123) 456-7893")
	must(err)
	_, err = insertPhone(db, "123-456-7894")
	must(err)
	_, err = insertPhone(db, "123-456-7890")
	must(err)
	_, err = insertPhone(db, "1234567892")
	must(err)
	_, err = insertPhone(db, "(123)456-7892")
	must(err)

	var number string
	number, err = getPhone(db, id)
	must(err)
	fmt.Println("The number is -", number)

	var phones []phoneNumber
	phones, err = allPhones(db)
	must(err)

	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.phone)
		if p.phone != number {
			fmt.Println("Updating or removing...", p.phone)
			existing, err := findPhone(db, number)
			must(err)
			if existing != nil {
				must(deletePhone(db, p.id))
			} else {
				p.phone = number
				must(updatePhone(db, p))
			}
		} else {
			fmt.Println("No changes required")
		}
	}
}

func getPhone(db *sql.DB, id int) (string, error) {
	var number string
	err := db.QueryRow("SELECT * FROM phone_numbers WHERE id=$1", id).Scan(&id, &number)
	if err != nil {
		return "", err
	}
	return number, nil
}

func findPhone(db *sql.DB, number string) (*phoneNumber, error) {
	var p phoneNumber
	err := db.QueryRow("SELECT * FROM phone_numbers WHERE value=$1", number).Scan(&p.id, &p.phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		} else {
			return nil, err
		}
	}
	return &p, nil
}

func updatePhone(db *sql.DB, p phoneNumber) error {
	statement := `UPDATE phone_numbers SET value=$2 WHERE id=$1`
	_, err := db.Exec(statement, p.id, p.phone)
	return err
}

func deletePhone(db *sql.DB, id int) error {
	statement := `DELETE FROM phone_numbers WHERE id=$1`
	_, err := db.Exec(statement, id)
	return err
}

type phoneNumber struct {
	id    int
	phone string
}

func allPhones(db *sql.DB) ([]phoneNumber, error) {
	var ret []phoneNumber
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var p phoneNumber
		err = rows.Scan(&p.id, &p.phone)
		if err != nil {
			return nil, err
		}
		ret = append(ret, p)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

func insertPhone(db *sql.DB, phone string) (int, error) {
	statement := `INSERT INTO phone_numbers (value) VALUES ($1) RETURNING id`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, nil
}

func createTable(db *sql.DB) error {
	statement := `
		CREATE TABLE IF NOT EXISTS phone_numbers (
		    id SERIAL,
		    value VARCHAR(255)
		)`
	_, err := db.Exec(statement)
	return err
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	if err != nil {
		return err
	}
	return nil
}

func normalize(phone string) string {
	re := regexp.MustCompile("\\D")
	return re.ReplaceAllString(phone, "")
}

//func normalize(phone string) string {
//	ret := ""
//	for _, ch := range phone {
//		if ch >= '0' && ch <= '9' {
//			ret += string(ch)
//		}
//	}
//	return ret
//}
