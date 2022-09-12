package storage

import "sync"

type DictStore struct {
	lock sync.Mutex
	DB   map[string]*StringStore `json:"db"`
}

func createDictStore() *DictStore {
	return &DictStore{
		lock: sync.Mutex{},
		DB:   make(map[string]*StringStore),
	}
}

func (s *DictStore) put(table, key string, value interface{}) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if _, ok := s.DB[table]; !ok {
		s.DB[table] = createStringStore()
		s.DB[table].setName(table)
	}

	s.DB[table].put(key, value)
}

func (s *DictStore) get(table, key string) (interface{}, bool) {
	if m, ok := s.DB[table]; ok {
		return m.get(key)
	}
	return nil, false
}

func (s *DictStore) remove(table, key string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if m, ok := s.DB[table]; ok {
		m.remove(key)
	}
}

func (s *DictStore) clear(table string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if m, ok := s.DB[table]; ok {
		m.clear()
	}
}

func (s *DictStore) len(table string) int {
	if m, ok := s.DB[table]; ok {
		return m.len()
	}
	return 0
}

func (s *DictStore) getKeys(table string) []string {
	if m, ok := s.DB[table]; ok {
		return m.getKeys()
	}
	return []string{}
}

func (s *DictStore) view(table string, amount int) string {
	if m, ok := s.DB[table]; ok {
		return m.view(amount)
	}
	return ""
}

func (s *DictStore) hasTable(table string) bool {
	_, ok := s.DB[table]
	return ok
}

func (s *DictStore) exists(table string, key string) bool {
	if m, ok := s.DB[table]; ok {
		return m.exists(key)
	}
	return false
}

func (s *DictStore) getTables() []string {
	tables := make([]string, 0, len(s.DB))
	for table := range s.DB {
		tables = append(tables, table)
	}
	return tables
}
