package storage

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDictStore(t *testing.T) {
	var stat bool
	var inter interface{}
	store := createDictStore()
	count := 100
	for i := 0; i < count; i++ {
		store.put("test_dict", strconv.Itoa(i), i)
	}

	// Get exists value
	inter, stat = store.get("test_dict", "66")

	if stat {
		assert.Equal(t, "66", fmt.Sprint(inter))
	}

	// Get not exists value
	inter, stat = store.get("test_dict", "666")
	if stat {
		assert.Equal(t, "66", fmt.Sprint(inter))
	}

	// View
	text := store.view("test_dict", 10)
	fmt.Println(text)

	// Remove
	store.remove("test_dict", "10")
	_, stat = store.get("test_dict", "10")
	assert.Equal(t, stat, false)

	// GetKeys
	keys := store.getKeys("test_dict")
	fmt.Println(keys)

	// not Exists
	stat = store.exists("test_dict", "10")
	fmt.Println("10 is exist?", stat)
	stat = store.exists("test_dict", "11")
	fmt.Println("11 is exist?", stat)

	// clear
	store.clear("test_dict")
	text2 := store.view("test_dict", 10)
	fmt.Println(text2)

	// tables
	tables := store.getTables()
	fmt.Println("tables:", tables)
}
