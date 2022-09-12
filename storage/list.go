package storage

import (
	"container/list"
	"fmt"
	"strings"
	"sync"

	"github.com/olekukonko/tablewriter"
)

type ListStore struct {
	Name string                `json:"name"`
	DB   map[string]*list.List `json:"db"`
	lock sync.Mutex
}

func createListStore() *ListStore {
	return &ListStore{
		DB:   make(map[string]*list.List),
		lock: sync.Mutex{},
	}
}

func (s *ListStore) put(table string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.DB[table]; !ok {
		s.DB[table] = list.New()
		s.Name = table
	}
	s.DB[table].PushBack(value)
}

func (s *ListStore) len(table string) int {
	if m, ok := s.DB[table]; ok {
		return m.Len()
	}
	return 0
}

func (s *ListStore) getByIndex(table string, index int) any {
	if index > s.len(table) || index < 0 {
		return ""
	}
	if m, ok := s.DB[table]; ok {
		cur := 0
		ele := m.Back().Value
		for e := m.Front(); e != nil; e = e.Next() {
			cur += 1
			ele = e.Value
			if cur == index {
				break
			}
		}
		return ele
	}
	return ""
}

func (s *ListStore) view(table string, amount int) string {
	s.lock.Lock()
	defer s.lock.Unlock()
	if m, ok := s.DB[table]; ok {
		cur := 0
		if amount > m.Len() {
			amount = m.Len()
		}

		tableString := &strings.Builder{}
		table := tablewriter.NewWriter(tableString)
		table.SetHeader([]string{"Index", "Value"})
		table.SetRowLine(true)

		data := make([][]string, amount)
		for e := m.Front(); e != nil; e = e.Next() {
			cur += 1
			data[cur-1] = []string{fmt.Sprint(cur - 1), fmt.Sprint(e.Value)}
			if cur == amount {
				break
			}
		}
		table.SetFooter([]string{"", fmt.Sprintf("Table: %s\nSelect Row Amount:%d\nTable Total Row: %d", s.Name, amount, m.Len())})
		table.AppendBulk(data)
		table.Render()
		return tableString.String()
	}
	return ""
}

func (s *ListStore) rangeByIndex(table string, begin, end int) []any {
	s.lock.Lock()
	defer s.lock.Unlock()

	if begin < 0 || begin > end || end < 0 {
		return []any{}
	}
	result := make([]any, 0, end-begin+1)
	if m, ok := s.DB[table]; ok {
		index := 0
		cur := m.Back()
		for e := m.Front(); e != nil; e = e.Next() {
			if index == begin {
				break
			}
			cur = e
			index += 1
		}
		for e := cur; e != nil; e = e.Next() {
			if index == end {
				break
			}
			result = append(result, e.Value)
			index += 1
		}
		return result
	}

	return []any{}
}

func (s *ListStore) removeByEle(table string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if m, ok := s.DB[table]; ok {
		ele := &list.Element{}
		for e := m.Front(); e != nil; e = e.Next() {
			if fmt.Sprint(e.Value) == fmt.Sprint(value) {
				ele = e
				break
			}
		}
		s.DB[table].Remove(ele)
	}
}

func (s *ListStore) removeByIndex(table string, index int) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.DB[table]; ok {
		ele := s.getByIndex(table, index)
		s.lock.Unlock()
		s.removeByEle(table, ele)
		s.lock.Lock()
	}
}

func (s *ListStore) clear(table string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if m, ok := s.DB[table]; ok {
		m.Init()
	}
}

func (s *ListStore) getTables() []string {
	keys := make([]string, 0, len(s.DB))
	for k := range s.DB {
		keys = append(keys, k)
	}
	return keys
}
