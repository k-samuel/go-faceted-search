package facet

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"testing"
	"time"
)

var index *Index
var datasetFilePrefix = ".test.dataset."
var indexSize uint64
var indexLoad time.Duration
var results = 100000
var datasetFile string

func init() {
	datasetFile = datasetFilePrefix + strconv.Itoa(results)
	if _, err := os.Stat(datasetFile); errors.Is(err, os.ErrNotExist) {
		CreateDataset()
	}
	index = CreateIndex()
}

func CreateDataset() {
	start := time.Now()
	colors := []string{"red", "green", "blue", "yellow", "black", "white"}
	brands := []string{
		"Nike",
		"H&M",
		"Zara",
		"Adidas",
		"Louis Vuitton",
		"Cartier",
		"Hermes",
		"Gucci",
		"Uniqlo",
		"Rolex",
		"Coach",
		"Victoria\"s Secret",
		"Chow Tai Fook",
		"Tiffany & Co.",
		"Burberry",
		"Christian Dior",
		"Polo Ralph Lauren",
		"Prada",
		"Under Armour",
		"Armani",
		"Puma",
		"Ray-Ban"}

	warehouses := []int{1, 10, 23, 345, 43, 5476, 34, 675, 34, 24, 789, 45, 65, 34, 54, 511, 512, 520}
	itemType := []string{"normal", "middle", "good"}

	f, err := os.Create(datasetFile)
	check(err)
	defer f.Close()

	for i := 1; i < results+1; i++ {

		countWh := rand.Int63n(int64(len(warehouses)))
		wh := make([]int64, 0)
		for j := 0; j < int(countWh); j++ {
			wh = append(wh, rand.Int63n(int64(len(warehouses))-1))
		}

		randType := rand.Int63n(int64(len(itemType) - 1))

		record := map[string]interface{}{
			"id":         i,
			"color":      colors[rand.Int31n(5)],
			"back_color": colors[rand.Int31n(5)],
			"size":       randNum(34, 50),
			"brand":      brands[rand.Int63n(int64(len(brands))-1)],
			"price":      randNum(10000, 100000),
			"discount":   rand.Int31n(10),
			"combined":   rand.Int31n(1),
			"quantity":   rand.Int31n(100),
			"warehouse":  unique(wh),
			"type":       itemType[randType],
		}

		s, e := json.Marshal(record)
		check(e)
		f.Write(s)
		f.WriteString("\n")
	}
	fmt.Println("Dataset: ", time.Since(start))
}

func CreateIndex() *Index {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	startM := m.Alloc
	start := time.Now()
	var result map[string]interface{}

	var index *Index
	index = NewIndex()

	file, err := os.Open(datasetFile)
	check(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		json.Unmarshal([]byte(scanner.Text()), &result)
		id := int64(result["id"].(float64))
		delete(result, "id")
		index.Add(id, result)
	}
	indexLoad = time.Since(start)
	runtime.GC()
	runtime.ReadMemStats(&m)
	indexSize = m.Alloc - startM
	return index
}

func randNum(min, max int64) int64 {
	return rand.Int63n(max-min+1) + min
}

func unique(intSlice []int64) []int64 {
	keys := make(map[int64]bool)
	list := []int64{}
	for _, entry := range intSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// -----
// go test facet -bench . -benchmem
// go test facet -bench . -benchmem -cpuprofile=cpu.out -memprofile=mem.out -memprofilerate=1 performance_test.go
// go tool pprof -callgrind -output callgrind.c.out cpu.out
// go tool pprof -callgrind -output callgrind.m.out mem.out

func BenchmarkFind(b *testing.B) {
	var recordFilter []int64
	search := NewSearch(index)
	filters := make([]FilterInterface, 0, 3)
	filters = append(filters, &ValueFilter{FieldName: "color", Values: []string{"black"}})
	filters = append(filters, &ValueFilter{FieldName: "warehouse", Values: []string{"789", "45", "65", "1", "10"}})
	filters = append(filters, &ValueFilter{FieldName: "type", Values: []string{"normal", "middle"}})
	search.Find(filters, recordFilter)
}

func BenchmarkAggregateFilters(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	var recordFilter []int64
	search := NewSearch(index)
	filters := make([]FilterInterface, 0, 3)
	filters = append(filters, &ValueFilter{FieldName: "color", Values: []string{"black"}})
	filters = append(filters, &ValueFilter{FieldName: "warehouse", Values: []string{"789", "45", "65", "1", "10"}})
	filters = append(filters, &ValueFilter{FieldName: "type", Values: []string{"normal", "middle"}})
	search.AggregateFilters(filters, recordFilter)
}

func BenchmarkSearch(b *testing.B) {

	runtime.GOMAXPROCS(runtime.NumCPU())

	start := time.Now()
	fmt.Printf("Alloc: %v MiB ", bToMb(indexSize))
	fmt.Print("Load: ", indexLoad)

	search := NewSearch(index)
	filters := make([]FilterInterface, 0, 3)
	filters = append(filters, &ValueFilter{FieldName: "color", Values: []string{"black"}})
	filters = append(filters, &ValueFilter{FieldName: "warehouse", Values: []string{"789", "45", "65", "1", "10"}})
	filters = append(filters, &ValueFilter{FieldName: "type", Values: []string{"normal", "middle"}})

	var recordFilter []int64
	start = time.Now()
	res, _ := search.Find(filters, recordFilter)
	duration := time.Since(start)
	fmt.Print(" Find: ", duration)
	fmt.Printf(" Results: %d ", len(res))
	fmt.Print(" Items: ", index.GetItemsCount())

	var sorter = NewFieldSorter(index)
	start = time.Now()
	sortedRecords, err := sorter.Sort(res, "quantity", SORT_DESC)
	if err != nil {
		panic(err)
	}
	duration = time.Since(start)
	fmt.Print(" Sort by field: ", duration, " sorted: ", len(sortedRecords))

	start = time.Now()
	filterRes, _ := search.AggregateFilters(filters, recordFilter)
	duration = time.Since(start)
	fmt.Println(" Aggregate filters: ", duration, " filters: ", len(filterRes))
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
