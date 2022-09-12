package storage

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestStore(t *testing.T) {
	s := InitStore("test.db", 10)
	fmt.Println(s)

	rand.Seed(time.Now().Unix())
	count := 10
	for i := 0; i < count; i++ {
		s.LPut("test_list", strconv.Itoa(i))

	}

	for i := 0; i < count; i++ {
		s.DPut("test_dict", strconv.Itoa(i), i)
		r := rand.Intn(count)
		if r < int(count/2) {
			s.DExpire("test_dict", strconv.Itoa(i), r)
		}
	}

	for i := 0; i < count; i++ {
		s.Put(strconv.Itoa(i), i)
		r := rand.Intn(count)
		if r < int(count/2) {
			s.Expire(strconv.Itoa(i), r)
		}
	}
	s.dumps()
}
