package storage

import (
	"fmt"
	"strconv"
	"testing"
)

func TestListStore(t *testing.T) {
	var inter interface{}
	store := createListStore()
	count := 100
	for i := 0; i < count; i++ {
		store.put("test_list", strconv.Itoa(i))
	}

	// Get exists value
	inter = store.getByIndex("test_list", 10)
	fmt.Println("Get exists value on test_list index=10:", inter)
	// index > len or index < 0
	inter = store.getByIndex("test_list", 101)
	fmt.Println("Get value on test_list index=101:", inter)

	inter = store.getByIndex("test_list", -1)
	fmt.Println("Get value on test_list index=101:", inter)

	// View
	text := store.view("test_list", 10)
	fmt.Println(text)

	// Remove index = 11, value = 10
	store.removeByIndex("test_list", 11)
	// Range 10-20
	items := store.rangeByIndex("test_list", 10, 20)
	fmt.Println(items)

	// // Get not exists value
	// inter, stat = store.Get("test_dict", "666")
	// if stat {
	// 	assert.Equal(t, "66", fmt.Sprint(inter))
	// }

	// // View
	// text := store.View("test_dict", 10)
	// fmt.Println(text)

	// // Remove
	// store.Remove("test_dict", "10")
	// _, stat = store.Get("test_dict", "10")
	// assert.Equal(t, stat, false)

	// // GetKeys
	// keys := store.GetKeys("test_dict")
	// fmt.Println(keys)

	// // not Exists
	// stat = store.Exists("test_dict", "10")
	// fmt.Println("10 is exist?", stat)
	// stat = store.Exists("test_dict", "11")
	// fmt.Println("11 is exist?", stat)

	// // clear
	store.clear("test_list")
	text2 := store.view("test_list", 10)
	fmt.Println(text2)

	// tables
	tables := store.getTables()
	fmt.Println("tables:", tables)
}
