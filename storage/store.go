package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/iddxc/memodb/utils"
)

var wg sync.WaitGroup

type Store struct {
	TimePeriod  int `json:"timePeriod"`
	location    string
	String      *StringStore `json:"string"`
	List        *ListStore   `json:"list"`
	Dict        *DictStore   `json:"dict"`
	ExpireStore *ExpireStore `json:"expire"`
	FlushTime   int64        `json:"flush_time"`
	storeChan   chan *Store
}

func isExist(filename string) bool {
	_, err := os.Stat(filename)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func InitStore(filename string, timePeriod int) *Store {
	stat := isExist(filename)
	fmt.Println(filename, "isExist?", stat)
	store := &Store{
		location:   filename,
		TimePeriod: timePeriod,
		storeChan:  make(chan *Store, 10),
	}
	if stat {
		store.loads()
	} else {
		store.Dict = createDictStore()
		store.List = createListStore()
		store.String = createStringStore()
		store.ExpireStore = createExpireStore()
	}
	store.FlushTime = time.Now().Unix() + int64(timePeriod)
	return store
}

func (s *Store) dumps() {
	data, err := json.Marshal(&s)
	if err != nil {
		fmt.Println("数据json序列化失败", err)
		return
	}
	utils.ZipString(s.location, data)
	fmt.Println("数据写入成功")
}

func (s *Store) loads() error {
	content := utils.UnzipString(s.location)
	err := json.Unmarshal(content, &s)
	if err != nil {
		fmt.Println("解码失败", err.Error())
	} else {
		fmt.Println("解码成功")
	}
	return err
}

func (s *Store) addExpire(mode, table, key string, index, timeout int) {
	s.ExpireStore.Nodes = append(s.ExpireStore.Nodes, &ExpireNode{
		Mode:      mode,
		Table:     table,
		Key:       key,
		Index:     index,
		TimeStamp: int64(timeout) + time.Now().Unix(),
	})
}

func (s *Store) watch() {
	for _, node := range s.ExpireStore.Nodes {
		if node.TimeStamp <= time.Now().Unix() {
			switch node.Mode {
			case "dict":
				s.DRemove(node.Table, node.Key)
			case "string":
				s.Remove(node.Key)
			}
		}
	}
	wg.Done()
}

func (s *Store) flush() {
	if s.TimePeriod > 0 {
		s.storeChan <- s
		sc := <-s.storeChan
		if s.FlushTime <= time.Now().Unix() {
			s.dumps()
			s.FlushTime = time.Now().Unix() + int64(sc.TimePeriod)
		}
	}
	wg.Done()
}

func (s *Store) Run() {
	for {
		wg.Add(2)
		go s.watch()
		go s.flush()
		wg.Wait()
		time.Sleep(time.Second)
		runtime.Gosched()
	}
}

func (s *Store) Put(key string, value interface{}) { s.String.put(key, value) }

func (s *Store) Get(key string) (interface{}, bool) { return s.String.get(key) }

func (s *Store) Remove(key string) { s.String.remove(key) }

func (s *Store) Len() int { return s.String.len() }

func (s *Store) GetKeys() []string { return s.String.getKeys() }

func (s *Store) View(amount int) string { return s.String.view(amount) }

func (s *Store) Clear() { s.String.clear() }

func (s *Store) Exists(key string) bool { return s.String.exists(key) }

func (s *Store) Expire(key string, timeout int) { s.addExpire("string", "", key, 0, timeout) }

func (s *Store) DPut(table, key string, value interface{}) { s.Dict.put(table, key, value) }

func (s *Store) DGet(table, key string) (interface{}, bool) { return s.Dict.get(table, key) }

func (s *Store) DRemove(table, key string) { s.Dict.remove(table, key) }

func (s *Store) DClear(table string) { s.Dict.clear(table) }

func (s *Store) DLen(table string) int { return s.Dict.len(table) }

func (s *Store) DGetKeys(table string) []string { return s.Dict.getKeys(table) }

func (s *Store) DView(table string, amount int) string { return s.Dict.view(table, amount) }

func (s *Store) HasTable(table string) bool { return s.Dict.hasTable(table) }

func (s *Store) DExists(table string, key string) bool { return s.Dict.exists(table, key) }

func (s *Store) DGetTables() []string { return s.Dict.getTables() }

func (s *Store) DExpire(table, key string, timeout int) { s.addExpire("dict", table, key, 0, timeout) }

func (s *Store) LPut(table string, value interface{}) { s.List.put(table, value) }

func (s *Store) LGetByIndex(table string, index int) any { return s.List.getByIndex(table, index) }

func (s *Store) LView(table string, amount int) string { return s.List.view(table, amount) }

func (s *Store) LRange(table string, begin, end int) []any {
	return s.List.rangeByIndex(table, begin, end)
}

func (s *Store) LRemoveByEle(table string, value interface{}) { s.List.removeByEle(table, value) }

func (s *Store) LRemoveByIndex(table string, index int) { s.List.removeByIndex(table, index) }

func (s *Store) LClear(table string) { s.List.clear(table) }

func (s *Store) LGetTables() []string { return s.List.getTables() }

func (s *Store) Dumps() { s.dumps() }
