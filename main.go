package main

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"time"
)

type Thermometer struct {
	Id          int64 `id`
	Temperature int   `temperature`
	Created     int64 `created_at`
	Updated     int64 `updated_at`
}

func initDb() *gorp.DbMap {
	// connect to db using standard Go database/sql API
	// use whatever database/sql driver you wish
	db, err := sql.Open("sqlite3", "thermistor_db.bin")
	if err != nil {
		log.Fatalln("sql.Open failed", err)
	}

	// construct a gorp DbMap
	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}

	// add a table, setting the table name to 'posts' and
	// specifying that the Id property is an auto incrementing PK
	dbmap.AddTableWithName(Thermometer{}, "thermometer").SetKeys(true, "Id")

	// create the table. in a production system you'd generally
	// use a migration tool, or create the tables via scripts
	if err = dbmap.CreateTablesIfNotExists(); err != nil {
		log.Fatalln("Create tables failed", err)
	}
	return dbmap
}

func NewThermometer(temperature int) Thermometer {
	return Thermometer{
		Temperature: temperature,
		Created:     time.Now().UnixNano(),
		Updated:     time.Now().UnixNano(),
	}
}

func main() {
	// initialize the DbMap
	dbmap := initDb()
	defer dbmap.Db.Close()

	// delete any existing rows
	if err := dbmap.TruncateTables(); err != nil {
		log.Printf("TruncateTables failed: %s\n", err.Error())
	}

	t1 := NewThermometer(25)
	t2 := NewThermometer(29)

	if err := dbmap.Insert(&t1, &t2); err != nil {
		log.Fatalln("Insert failed", err)
	}

	//select one
	if err := dbmap.SelectOne(&t2, "select * from thermometer where id=?", t2.Id); err != nil {
		log.Fatalln("Select One failed")
	}
	log.Println("t2 row:", t2)

	//delete
	count, err := dbmap.Delete(&t1)
	if err != nil {
		log.Fatalln("Delete failed", err)
	}
	log.Println("Rows deleted:", count)

	count, err = dbmap.Delete(&t2)
	if err != nil {
		log.Fatalln("Delete failed", err)
	}
	log.Println("Rows deleted:", count)

	log.Println("Done!")
}
