package storage

import (
	"errors"
	"fmt"
	"log"

	"github.com/EatsLemons/fa_currencies/store"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var collection = "fa_currencies"

type MongoDB struct {
	session *mgo.Session
	dbname  string
}

func NewMongoDB(address, usrname, pwd, dbname string) *MongoDB {
	connectInfo := mgo.DialInfo{
		Addrs:    []string{address},
		Database: dbname,
		Username: usrname,
		Password: pwd,
	}

	mgo, err := mgo.DialWithInfo(&connectInfo)
	if err != nil {
		log.Printf("[WARN] connect to mongo has failed %s", err)
		return nil
	}

	mongoClient := MongoDB{
		session: mgo,
		dbname:  dbname,
	}

	return &mongoClient
}

func (mdb *MongoDB) GetCurrPair(from, to string) (*store.Ratio, error) {
	s := mdb.session.Copy()
	defer s.Close()

	c := s.DB(mdb.dbname).C(collection)
	record := priceRecord{}
	err := c.Find(bson.M{"from": from}).One(&record)
	if err != nil {
		log.Printf("[INFO] failed to find record for %s / %s, error: %s", from, to, err)
		return nil, err
	}

	if _, ok := record.To[to]; !ok {
		log.Printf("[INFO] there is no record %s / %s", from, to)
		return nil, errors.New(fmt.Sprintf("there is no record %s / %s", from, to))
	}

	result := store.Ratio{
		From: from,
		To:   make(map[string]float64),
	}

	result.To[to] = record.To[to]

	return &result, nil
}

func (mdb *MongoDB) Save(currRates []store.Ratio) error {
	if currRates == nil || len(currRates) == 0 {
		return errors.New("nothing to save")
	}

	s := mdb.session.Copy()
	defer s.Close()

	c := s.DB(mdb.dbname).C(collection)
	for _, rate := range currRates {
		row := priceRecord{
			From: rate.From,
			To:   rate.To,
		}

		// TODO: maybe there is more efficiency way to save
		_, err := c.Upsert(bson.M{"from": row.From}, row)
		if err != nil {
			return err
		}
	}

	return nil
}

type priceRecord struct {
	ID   bson.ObjectId `bson:"_id,omitempty"`
	From string
	To   map[string]float64
}
