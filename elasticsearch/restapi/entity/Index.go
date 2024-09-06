package entity

type Index struct {
	DocsCount    string `json:"docs.count"`
	DocsDeleted  string `json:"docs.deleted"`
	Health       string `json:"health"`
	Index        string `json:"index"`
	Pri          string `json:"pri"`
	PriStoreSize string `json:"pri.store.size"`
	Rep          string `json:"rep"`
	Status       string `json:"status"`
	StoreSize    string `json:"store.size"`
	Uuid         string `json:"uuid"`
	Shards       []Shard
}

type Shard struct {
	Index  string `json:"index"`
	Shard  string `json:"shard"`
	Prirep string `json:"prirep"`
	State  string `json:"state"`
	Docs   string `json:"docs"`
	Store  string `json:"store"`
	Ip     string `json:"ip"`
	Node   string `json:"node"`
}
