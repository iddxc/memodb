package storage

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStingStorage(t *testing.T) {
	var stat bool
	var inter interface{}
	store := createStringStore()
	count := 100
	for i := 0; i < count; i++ {
		store.put(strconv.Itoa(i), i)
	}

	// Get exists value
	inter, stat = store.get("66")

	if stat {
		assert.Equal(t, "66", fmt.Sprint(inter))
	}

	// Get not exists value
	inter, stat = store.get("666")
	if stat {
		assert.Equal(t, "66", fmt.Sprint(inter))
	}

	// View
	text := store.view(10)
	fmt.Println(text)

	// Remove
	store.remove("10")
	_, stat = store.get("10")
	assert.Equal(t, stat, false)

	// GetKeys
	keys := store.getKeys()
	fmt.Println(keys)

	// not Exists
	stat = store.exists("10")
	fmt.Println("10 is exist?", stat)
	stat = store.exists("11")
	fmt.Println("11 is exist?", stat)

	// clear
	store.clear()
	text2 := store.view(10)
	fmt.Println(text2)
}
