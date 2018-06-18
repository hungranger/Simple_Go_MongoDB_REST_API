package store

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	mgo "gopkg.in/mgo.v2"
	bson "gopkg.in/mgo.v2/bson"
)

type MongoDBRepository struct {
	session *mgo.Session
}

const AutoRemove = true // auto remove every second
const MinSecond = 1
const MaxSecond = 100
const MinValue = 1000
const MaxValue = 10000
const MaxLen = 10000
const LatestAmount = 20 // 20 latest pair
const PerPage = 10      // 10 oldest pair
const AddingTimeOut = 600

// SERVER the DB server
const SERVER = "mongodb://hung:a123456@ds261570.mlab.com:61570/pairs"

// DBNAME the name of the DB instance
const DBNAME = "pairs"

// COLLECTION is the name of the collection in DB
const COLLECTION = "store"

func (r *MongoDBRepository) connect() bool {
	var err error
	r.session, err = mgo.Dial(SERVER)

	if err != nil {
		panic(err)
	}

	fmt.Println("DB Connected.")
	return true
}

func (r *MongoDBRepository) close() {
	r.session.Close()
	fmt.Println("DB Closed.")
}

func (r *MongoDBRepository) dropDatabase() bool {
	err := r.session.DB(DBNAME).DropDatabase()
	if err != nil {
		panic(err)
	}
	fmt.Println("DB droped.")
	return true
}

// AddProduct adds a Product in the DB
func (r *MongoDBRepository) addPair(pair Pair) int32 {
	// Collection People
	store := r.session.DB(DBNAME).C(COLLECTION)

	result := []Pair{}
	err := store.Find(bson.M{}).All(&result)
	if len(result) == MaxLen {
		p := Pair{}
		store.Find(nil).Sort("timestamp").Limit(1).One(&p)

		err = store.RemoveId(p.ID)
		if err != nil {
			fmt.Printf("remove fail %v: id= %v, s= %v\n", err, p.Timestamp, p.Second)
			os.Exit(1)
		}
		fmt.Printf("removed oldest: %v, s= %v\n", p.Timestamp, p.Second)
	}

	err = store.Insert(pair)
	if err != nil {
		panic(err)
	}

	return pair.Timestamp
}

func (r *MongoDBRepository) deletePairAfter(second int, timestamp int32) {
	after := time.After(time.Duration(second) * time.Second)
	<-after
	r.deletePairByTimestamp(timestamp)
}

func (r *MongoDBRepository) deletePairByTimestamp(timestamp int32) {
	c := r.session.DB(DBNAME).C(COLLECTION)

	result := Pair{}
	err := c.Find(bson.M{"timestamp": timestamp}).One(&result)
	if err != nil {
		return
	}

	if result.IsInUse == false {
		//remove record
		err = c.Remove(bson.M{"timestamp": timestamp})
		if err != nil {
			fmt.Printf("remove fail %v: ts= %v, s= %v\n", err, result.Timestamp, result.Second)
			os.Exit(1)
		}
		fmt.Printf("DB removed Timestamp: %v After %v s \n", result.Timestamp, result.Second)
	} else {
		fmt.Printf("Cannot remove pair is in used: %v \n", result.Timestamp)
	}
}

func (r *MongoDBRepository) addPairEverySecond() {
	r.connect()
	r.dropDatabase()
	tick := time.Tick(time.Second)
	timeout := time.After(AddingTimeOut * time.Second)
	for {
		select {
		case <-tick:
			p := getRndPair()
			pair := Pair{Second: p[0], Value: p[1], IsInUse: false, Timestamp: int32(time.Now().Unix())}

			timestamp := r.addPair(pair)
			fmt.Println("DB added. Timestamp: ", timestamp)

			if AutoRemove {
				go r.deletePairAfter(p[0], timestamp)
			}
		case <-timeout:
			fmt.Println("DB stop adding")
			// r.close()
			return
		}
	}
}

func getRndPair() [2]int {
	var pair [2]int
	pair[0] = getRndInt(MinSecond, MaxSecond)
	pair[1] = getRndInt(MinValue, MaxValue)

	return pair
}

func getRndInt(min, max int) int {
	seed := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(seed)

	rNum := math.Floor(r1.Float64()*(float64(max)-float64(min)+1.0)) + float64(min)

	return int(rNum)
}

func getRndBool() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Float32() < 0.5
}

func (r *MongoDBRepository) get20Latest() []Pair {
	c := r.session.DB(DBNAME).C(COLLECTION)
	pairs := []Pair{}
	err := c.Find(nil).Sort("-timestamp").Limit(LatestAmount).All(&pairs)
	if err != nil {
		fmt.Printf("20 latest pairs not found %v\n", err)
		return nil
	}

	r.updateStatus(pairs, true)

	return pairs
}

func (r *MongoDBRepository) updateStatus(pairs []Pair, status bool) bool {
	c := r.session.DB(DBNAME).C(COLLECTION)
	for _, pair := range pairs {
		colQuerier := bson.M{"timestamp": pair.Timestamp}
		change := bson.M{"$set": bson.M{"isinuse": status}}
		err := c.Update(colQuerier, change)
		if err != nil {
			fmt.Printf("Update not success: %v , %v, %v \n", err, pair.Timestamp, status)
		}
	}
	return true
}

func (r *MongoDBRepository) get10Oldest(pageNo int) []Pair {
	c := r.session.DB(DBNAME).C(COLLECTION)
	pairs := []Pair{}
	err := c.Find(nil).Sort("timestamp").Skip((pageNo - 1) * PerPage).Limit(PerPage).All(&pairs)
	if err != nil {
		fmt.Printf("10 Oldest pairs not found %v\n, ", err)
		return nil
	}

	r.updateStatus(pairs, true)
	return pairs
}
