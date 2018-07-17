package storage

import (
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

func (mdb *MongoDB) Update(currRates []store.Ratio) error {
	s := mdb.session.Copy()
	defer s.Close()

	c := s.DB(mdb.dbname).C(collection)
	for _, rate := range currRates {
		row := priceRecord{
			From: rate.From,
			To:   rate.To,
		}

		// TODO: maybe there is more efficiency way
		_, err := c.Upsert(bson.M{"From": row.From}, row)
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
