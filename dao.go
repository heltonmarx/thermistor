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

type DAO struct {
	dbmap *gorp.DbMap
}

func NewDAO() *DAO {

	db, err := sql.Open("sqlite3", "thermistor_db.bin")
	if err != nil {
		log.Println("sql.Open failed", err)
		return nil
	}

	dao := DAO{}
	dao.dbmap = &gorp.DbMap{Db: db, Dialect: gorp.SqliteDialect{}}
	dao.dbmap.AddTableWithName(Thermometer{}, "thermometer").SetKeys(true, "Id")

	if err = dao.dbmap.CreateTablesIfNotExists(); err != nil {
		log.Println("Create tables failed", err)
		return nil
	}
	return &dao
}

func (self *DAO) Insert(temperature int) error {
	//set created and updated with datetime
	t := Thermometer{
		Temperature: temperature,
		Created:     time.Now().UnixNano(),
		Updated:     time.Now().UnixNano(),
	}
	if err := self.dbmap.Insert(&t); err != nil {
		log.Println("Insert failed", err)
		return err
	}
	return nil
}

func (self *DAO) GetById(id int64) (*Thermometer, error) {
	obj, err := self.dbmap.Get(Thermometer{}, id)
	if err != nil {
		log.Println("Select One failed", err)
		return nil, err
	}
	t := obj.(*Thermometer)
	return t, nil
}

func (self *DAO) GetAll() ([]Thermometer, error) {
	var l []Thermometer
	_, err := self.dbmap.Select(&l, "select * from thermometer order by id")
	if err != nil {
		log.Println("Select failed", err)
		return l, err
	}
	return l, nil
}

func (self *DAO) Update(t *Thermometer) (int64, error) {
	//only set updated with datetime
	t.Updated = time.Now().UnixNano()

	count, err := self.dbmap.Update(t)
	if err != nil {
		log.Println("Update failed", err)
		return -1, err
	}
	return count, nil
}

func (self *DAO) Delete(t *Thermometer) (int64, error) {
	count, err := self.dbmap.Delete(t)
	if err != nil {
		log.Println("Delete failed", err)
		return -1, err
	}
	return count, nil
}

func (self *DAO) DeleteAll() error {
	// delete any existing rows
	if err := self.dbmap.TruncateTables(); err != nil {
		log.Println("TruncateTables failed", err)
		return err
	}
	return nil
}

func (self *DAO) Close() {
	self.dbmap.Db.Close()
}
