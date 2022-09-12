# MemoDB
[MemoDB](https://github.com/iddxc/MemoDB)是一个轻量级且简单的键值存储应用,启发于Redis，使用Go语言进行开发，能够做到开包即用的效果。

特点:
- 支持任意类型
- 支持并发操作
- 支持视图字符串输出
- 支持过期删除，单位为秒
- 支持周期快照存储，json字符串压缩文件
- 支持快照文件重载

## 导入与使用

### 导入的两种方式
- 使用go命令`go get "github.com/iddxc/memodb"` 
- 使用在源代码中使用`import "github.com/iddxc/memodb"`进行导入，并运行`go mod tidy`进行环境管理

### 启用快照和过期删除功能
```go
package main

import (
	"github.com/iddxc/memodb"
)

func main() {
	store := memodb.New("test.db", 20)
	count := 1000

	rand.Seed(time.Now().Unix())
	for i := 0; i < count; i++ {
		r := rand.Intn(100)
		store.Put(strconv.Itoa(i), r)
		store.Expire(strconv.Itoa(i), r)
	}
	var wait sync.WaitGroup
	wait.Add(1)
	go store.Run()
	defer wait.Wait()

	fmt.Println(store.View(10))
}
```
该功能做到每20秒进行快照存储1次,并且对存在过期时间限制元素进行判断

### String
String: key-value key为字符串类型且为唯一，而value可以为任意类型，支持过期删除 
```go
package main

import (
	"fmt"
	"strconv"

	"github.com/iddxc/memodb"
)

func main() {
	var stat bool
	var inter interface{}

	memo := memodb.New("test.db", 0)

	count := 100
    // 设置键值对
	for i := 0; i < count; i++ {
		memo.Put(strconv.Itoa(i), i)
	}

	// 获取存在的值
	inter, stat = memo.Get("66")

	if stat {
		assert.Equal(t, "66", fmt.Sprint(inter))
	}

	// 获取不存在的值
	inter, stat = memo.Get("666")
	if stat {
		assert.Equal(t, "66", fmt.Sprint(inter))
	}

	// 获取String键值对的前10行视图字符串
	text := memo.View(10)
	fmt.Println(text)

	// 删除指定建
	memo.Remove("10")
	_, stat = memo.Get("10")
	assert.Equal(t, stat, false)

	// 获取所有键值对的键
	keys := memo.GetKeys()
	fmt.Println(keys)

	// 判断元素是否存在，key="10"的键已被删除
	stat = memo.Exists("10")
	fmt.Println("10 is exist?", stat)
    // 判断元素是否存在，key="11"的键仍存在
	stat = memo.Exists("11")
	fmt.Println("11 is exist?", stat)

	// 清除所有记录
	memo.Clear()
	text2 := memo.View(10)
	fmt.Println(text2)

    // 设置过期时间 **此处存在缺陷，需要在memo.Run()进行监控情况下才能进行
    memo.Expire("10", 10)
    // 10秒中之后进行删除
}
```

#### 运行结果
```shell
test.db isExist? false
+-------+-----+----------------------+
| INDEX | KEY |        VALUE         |
+-------+-----+----------------------+
|     1 |  46 |                   46 |
+-------+-----+----------------------+
|     2 |  47 |                   47 |
+-------+-----+----------------------+
|     3 |  48 |                   48 |
+-------+-----+----------------------+
|     4 |  68 |                   68 |
+-------+-----+----------------------+
|     5 |  75 |                   75 |
+-------+-----+----------------------+
|     6 |  26 |                   26 |
+-------+-----+----------------------+
|     7 |  31 |                   31 |
+-------+-----+----------------------+
|     8 |  43 |                   43 |
+-------+-----+----------------------+
|     9 |  96 |                   96 |
+-------+-----+----------------------+
|    10 |  88 |                   88 |
+-------+-----+----------------------+
|               TABLE: STRING SELECT |
|               ROW AMOUNT:10 TABLE  |
|                  TOTAL ROW: 100    |
+-------+-----+----------------------+

[73 89 91 98 20 28 69 76 6 32 36 52 70 95 13 34 37 56 78 90 92 2 7 55 80 81 83 12 29 50 49 84 86 0 16 39 17 62 71 82 94 3 8 14 57 63 66 25 42 44 40 61 64 93 9 24 72 1 35 51 67 18 38 65 27 30 60 74 85 4 19 22 59 79 87 21 33 54 46 47 48 68 75 26 31 43 96 88 99 23 41 58 45 53 77 97 5 11 15]
10 is exist? false
11 is exist? true
+-------+-----+---------------------+
| INDEX | KEY |        VALUE        |
+-------+-----+---------------------+
|                  TABLE: STRING    |
|               SELECT ROW AMOUNT:0 |
|               TABLE TOTAL ROW: 0  |
+-------+-----+---------------------+
```

### List
List: table:key-value 寻找并操作table下名为key的列表，并添加操作值，其中table和key必为字符串类型，value为任意类型，不支持过期删除

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	var inter interface{}
	memo := New("test.db", 10)
	count := 100
	for i := 0; i < count; i++ {
		memo.LPut("test_list", strconv.Itoa(i))
	}

	// Get exists value
	inter = memo.LGetByIndex("test_list", 10)
	fmt.Println("Get exists value on test_list index=10:", inter)
	// index > len or index < 0
	inter = memo.LGetByIndex("test_list", 101)
	fmt.Println("Get value on test_list index=101:", inter)

	inter = memo.LGetByIndex("test_list", -1)
	fmt.Println("Get value on test_list index=101:", inter)

	// View
	text := memo.LView("test_list", 10)
	fmt.Println(text)

	// Remove index = 11, value = 10
	memo.LRemoveByIndex("test_list", 11)
	// Range 10-20
	items := memo.LRange("test_list", 10, 20)
	fmt.Println(items)

	// tables
	tables := memo.LGetTables()
	fmt.Println("tables:", tables)

	// clear
	memo.LClear("test_list")
	text2 := memo.LView("test_list", 10)
	fmt.Println(text2)
}

```
#### 输出结果
``` shell
Get exists value on test_list index=10: 9
Get value on test_list index=101:
Get value on test_list index=101:
+-------+----------------------+
| INDEX |        VALUE         |
+-------+----------------------+
|     0 |                    0 |
+-------+----------------------+
|     1 |                    1 |
+-------+----------------------+
|     2 |                    2 |
+-------+----------------------+
|     3 |                    3 |
+-------+----------------------+
|     4 |                    4 |
+-------+----------------------+
|     5 |                    5 |
+-------+----------------------+
|     6 |                    6 |
+-------+----------------------+
|     7 |                    7 |
+-------+----------------------+
|     8 |                    8 |
+-------+----------------------+
|     9 |                    9 |
+-------+----------------------+
|           TABLE: TEST LIST   |
|         SELECT ROW AMOUNT:10 |
|         TABLE TOTAL ROW: 100 |
+-------+----------------------+

[9 11 12 13 14 15 16 17 18 19]
+-------+---------------------+
| INDEX |        VALUE        |
+-------+---------------------+
|          TABLE: TEST LIST   |
|         SELECT ROW AMOUNT:0 |
|         TABLE TOTAL ROW: 0  |
+-------+---------------------+

tables: [test_list]

```

### Dict
Dict table:key-value 寻找并操作table名下key-value键值对，基础结构为String结构体的封装，支持过期删除

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {

	var stat bool
	var inter interface{}
	memo := New("test.db", 0)
	count := 100
	for i := 0; i < count; i++ {
		memo.DPut("test_dict", strconv.Itoa(i), i)
	}

	// Get exists value
	inter, stat = memo.DGet("test_dict", "66")
	fmt.Println("test_dict key=66:", inter)

	// Get not exists value
	inter, stat = memo.DGet("test_dict", "666")
	fmt.Println("test_dict key=666:", inter)

	// View
	text := memo.DView("test_dict", 10)
	fmt.Println(text)

	// Remove
	memo.DRemove("test_dict", "10")
	_, stat = memo.DGet("test_dict", "10")

	// GetKeys
	keys := memo.DGetKeys("test_dict")
	fmt.Println(keys)

	// not Exists
	stat = memo.DExists("test_dict", "10")
	fmt.Println("10 is exist?", stat)
	stat = memo.DExists("test_dict", "11")
	fmt.Println("11 is exist?", stat)

	// clear
	memo.DClear("test_dict")
	text2 := memo.DView("test_dict", 10)
	fmt.Println(text2)

	// tables
	tables := memo.DGetTables()
	fmt.Println("tables:", tables)
}
```
#### 输出结果
```shell
test.db isExist? false
test_dict key=66: 66
test_dict key=666: <nil>
+-------+-----+----------------------+
| INDEX | KEY |        VALUE         |
+-------+-----+----------------------+
|     1 |  28 |                   28 |
+-------+-----+----------------------+
|     2 |  56 |                   56 |
+-------+-----+----------------------+
|     3 |  64 |                   64 |
+-------+-----+----------------------+
|     4 |  68 |                   68 |
+-------+-----+----------------------+
|     5 |   1 |                    1 |
+-------+-----+----------------------+
|     6 |   5 |                    5 |
+-------+-----+----------------------+
|     7 |  10 |                   10 |
+-------+-----+----------------------+
|     8 |  14 |                   14 |
+-------+-----+----------------------+
|     9 |  71 |                   71 |
+-------+-----+----------------------+
|    10 |  73 |                   73 |
+-------+-----+----------------------+
|               TABLE: STRING SELECT |
|               ROW AMOUNT:10 TABLE  |
|                  TOTAL ROW: 100    |
+-------+-----+----------------------+
[28 56 64 68 1 5 14 71 73 77 43 55 66 69 3 7 26 40 22 35 54 59 30 46 84 75 83 24 50 65 70 20 23 72 79 11 12 16 19 78 4 29 31 48 17 51 74 94 2 8 44 88 47 63 92 6 18 36 41 62 81 89 90 0 32 49 53 91 60 67 76 86 13 33 42 45 87 96 38 58 80 99 27 39 97 98 82 85 93 15 21 25 57 61 95 9 34 37 52]
10 is exist? false
11 is exist? true
+-------+-----+---------------------+
| INDEX | KEY |        VALUE        |
+-------+-----+---------------------+
|                  TABLE: STRING    |
|               SELECT ROW AMOUNT:0 |
|               TABLE TOTAL ROW: 0  |
+-------+-----+---------------------+

tables: [test_dict]

```
