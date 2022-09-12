package storage

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
)

type StringStore struct {
	Name string `json:"name"`
	lock sync.Mutex
	DB   map[string]any `json:"db"`
}

func createStringStore() *StringStore {
	return &StringStore{
		DB:   make(map[string]any),
		lock: sync.Mutex{},
	}
}

func (s *StringStore) setName(name string) {
	s.Name = name
}

func (s *StringStore) put(key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Name = "string"
	s.DB[key] = value
}

func (s *StringStore) get(key string) (interface{}, bool) {
	if ele, ok := s.DB[key]; ok {
		return fmt.Sprint(ele), ok
	}
	return nil, false
}

func (s *StringStore) remove(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.DB[key]; ok {
		delete(s.DB, key)
	}
}

func (s *StringStore) len() int {
	return len(s.DB)
}

func (s *StringStore) getKeys() []string {
	keys := make([]string, 0, len(s.DB))
	for k := range s.DB {
		keys = append(keys, k)
	}
	return keys
}

func (s *StringStore) view(amount int) string {
	s.lock.Lock()
	defer s.lock.Unlock()

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Index", "Key", "Value"})
	table.SetRowLine(true)

	keys := s.getKeys()
	length := s.len()

	if length < amount {
		amount = length
	}
	data := make([][]string, amount)
	for index, key := range keys[:amount] {
		if v, ok := s.DB[key]; ok {
			temp := []string{strconv.Itoa(index + 1), key, fmt.Sprint(v)}
			data[index] = temp
		}
	}
	table.SetFooter([]string{"", "", fmt.Sprintf("Table: %s\nSelect Row Amount:%d\nTable Total Row: %d", s.Name, amount, len(s.DB))})
	table.AppendBulk(data)
	table.Render()
	return tableString.String()
}

func (s *StringStore) clear() {
	s.lock.Lock()
	defer s.lock.Unlock()

	s.DB = make(map[string]any)
}

func (s *StringStore) exists(key string) bool {
	if _, ok := s.DB[key]; ok {
		return true
	}
	return false
}
