package store

import (
	"gopkg.in/mgo.v2/bson"
)

// Pair represents an db item
type Pair struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Second    int           `json:"second"`
	Value     int           `json:"value"`
	IsInUse   bool          `json:"isinuse"`
	Timestamp int32         `json:"timestamp"`
}

// Pairs is an array of Product objects
type Pairs []Pair
