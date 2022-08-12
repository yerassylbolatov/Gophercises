package main

import (
	"fmt"
	_ "github.com/lib/pq"
	"log"
	phoneDb "main/db"
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
	must(phoneDb.Reset("postgres", psqlInfo, dbname))

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	must(phoneDb.Migrate("postgres", psqlInfo))

	db, err := phoneDb.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()

	err = db.Seed()
	must(err)

	phones, err := db.AllPhones()
	must(err)
	for _, p := range phones {
		fmt.Printf("Working on... %+v\n", p)
		number := normalize(p.Number)
		if p.Number != number {
			fmt.Println("Updating or removing...", p.Number)
			existing, err := db.FindPhone(number)
			must(err)
			if existing != nil {
				must(db.DeletePhone(p.Id))
			} else {
				p.Number = number
				must(db.UpdatePhone(&p))
			}
		} else {
			fmt.Println("No changes required")
		}
	}
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
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
