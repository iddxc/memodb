package storage

type ExpireNode struct {
	Mode      string `json:"mode"`
	Table     string `json:"table"`
	Key       string `json:"key"`
	Index     int    `json:"index"`
	TimeStamp int64  `json:"expire"`
}

type ExpireStore struct {
	Nodes []*ExpireNode `json:"nodes"`
}

func createExpireStore() *ExpireStore {
	return &ExpireStore{
		Nodes: make([]*ExpireNode, 0),
	}
}
