package cache

import (
	"fmt"
	"testing"
	"time"

	"github.com/jadogg22/go-sharpGraphs/pkg/getData"
	"github.com/jadogg22/go-sharpGraphs/pkg/models"
)

type Person struct {
	Name string
	Age  int
}

type Product struct {
	ID    int
	Name  string
	Price float64
}

func TestCache(t *testing.T) {
	persons := []Person{{Name: "Alice", Age: 30}, {Name: "Bob", Age: 25}}
	products := []Product{{ID: 1, Name: "Laptop", Price: 999.99}, {ID: 2, Name: "Phone", Price: 499.99}}

	var tests = []struct {
		name   string
		key    string
		value  interface{}
		typeID string
		ttl    time.Duration
		pass   bool
	}{
		{"test1", "person", persons, "[]Person", 1 * time.Second, true},
		{"test2", "product", products, "[]Product", 1 * time.Second, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MyCache.Set(tt.key, tt.value, tt.typeID, tt.ttl)
			val, typeID, found := MyCache.Get(tt.key)
			if !tt.pass && found {
				if typeID != tt.typeID {
					t.Errorf("expected typeID %s, got %s", tt.typeID, typeID)
					return
				}
				t.Errorf("expected not to find key %s", tt.key)
			}

			if tt.pass && !found {
				t.Errorf("expected to find key %s", tt.key)
			}
			if tt.pass && typeID != tt.typeID {
				t.Errorf("expected typeID %s, got %s", tt.typeID, typeID)
			}
			if tt.pass && fmt.Sprintf("%v", val) != fmt.Sprintf("%v", tt.value) {
				t.Errorf("expected value %v, got %v", tt.value, val)
			}
		})
	}

}

func Daily_Ops_Test() ([]*models.DailyOpsData, error, bool) {
	cacheKey := "dailyOpsData"
	usedCache := false

	cachedData, typeID, found := MyCache.Get(cacheKey)
	if found {
		if typeID == "[]*models.DailyOpsData" {
			if cachedData, ok := cachedData.([]*models.DailyOpsData); ok {
				usedCache = true
				return cachedData, nil, usedCache
			}
		} else {
			return nil, fmt.Errorf("Error casting the data"), usedCache
		}
	}
	// cache miss, get the data from the database.

	today := time.Now()
	startDate := time.Date(today.Year(), today.Month(), today.Day(), 0, 0, 0, 0, today.Location())

	data, err := getdata.GetTransportationDailyOps(startDate, today)
	if err != nil {
		return nil, err, usedCache
	}

	// Set the cache
	MyCache.Set(cacheKey, data, "[]*models.DailyOpsData", time.Second*2)

	//Finally update the Response with the json data
	return data, nil, usedCache
}

func Test_Daily_Ops_Test(t *testing.T) {
	// test the function
	var key = "dailyOpsData"

	// ------- test the cache hit with nothing in it ------
	_, _, found := MyCache.Get(key)
	if found {
		t.Errorf("expected not to find key %s", key)
	}

	// ------- test cache miss ------
	data, err, usedCache := Daily_Ops_Test()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if usedCache {
		t.Errorf("expected not to use the cache")
	}
	myTime := time.Now() // save the time for an estimate when the cache will expire

	if len(data) == 0 {
		t.Errorf("expected data to have length > 0")
	}

	// ------- test the cache hit ------
	data2, err, usedCache := Daily_Ops_Test()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if !usedCache {
		t.Errorf("expected to use the cache")
		if time.Now().After(myTime.Add(time.Second * 2)) {
			fmt.Println("Cache expired")
		}
	}

	if len(data2) == 0 {
		t.Errorf("expected data to have length > 0")
	}

	fmt.Println("waiting for cache to expire")
	// wait for the cache to expire
	time.Sleep(time.Second * 2)

	// ------- test cache miss ------
	data3, err, usedCache := Daily_Ops_Test()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if usedCache {
		t.Errorf("expected not to use the cache")
	}

	if len(data3) == 0 {
		t.Errorf("expected data to have length > 0")
	}

	fmt.Println(data3)
}
