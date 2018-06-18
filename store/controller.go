package store

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
)

//Controller ...
type Controller struct {
	Repository MongoDBRepository
}

var sumList []int

// Index GET /
func (c *Controller) Index(w http.ResponseWriter, r *http.Request) {
	sumList = nil
	go c.Repository.addPairEverySecond()
	// defer c.Repository.close()

	tmplIndex := template.Must(template.ParseFiles("./static/index.html"))
	tmplIndex.Execute(w, nil)
}

// get20Latest GET /get20Latest
func (c *Controller) Get20Latest(w http.ResponseWriter, r *http.Request) {
	type LatestData struct {
		Value int
	}

	pairs := c.Repository.get20Latest()
	var data []LatestData
	for _, pair := range pairs {
		data = append(data, LatestData{pair.Value})
	}
	jsonData, _ := json.Marshal(data)
	c.writeResponse(w, pairs, jsonData)
	return
}

// get10Oldest GET /get10Oldest/{index}
func (c *Controller) Get10Oldest(w http.ResponseWriter, r *http.Request) {
	type OldestData struct {
		Value int
	}

	vars := mux.Vars(r)
	index, _ := strconv.Atoi(vars["index"])

	fmt.Printf("index: %v \n", index)

	pairs := c.Repository.get10Oldest(index)
	fmt.Printf("10 oldest: %v \n", pairs)
	sum := 0
	for _, pair := range pairs {
		sum += pair.Value
		fmt.Printf("t: %v \n", pair.Timestamp)
	}
	fmt.Printf("sum: %v \n", sum)

	sumList = append(sumList, sum)

	jsonData, _ := json.Marshal(OldestData{sum})
	c.writeResponse(w, pairs, jsonData)
	return
}

// getMedian GET /getMedian
func (c *Controller) GetMedian(w http.ResponseWriter, r *http.Request) {
	type MedianData struct {
		SumList []int
		Value   float64
	}

	median := median(sumList)

	medianData := MedianData{SumList: sumList, Value: median}

	jsonData, _ := json.Marshal(medianData)
	c.writeResponse(w, nil, jsonData)
	return
}

func median(numbers []int) float64 {
	// median of [3, 5, 4, 4, 1, 1, 2, 3] = 3
	median := 0.0
	numsLen := len(numbers)
	sort.Ints(numbers)

	if numsLen%2 == 0 { // is even
		// average of two middle numbers
		median = (float64(numbers[numsLen/2-1]) + float64(numbers[numsLen/2])) / 2.0
	} else { // is odd
		// middle number only
		median = float64(numbers[(numsLen-1)/2])
	}

	return median
}

func (c *Controller) writeResponse(w http.ResponseWriter, pairs []Pair, data []byte) {

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

	if pairs != nil {
		c.Repository.updateStatus(pairs, false)
	}
}
