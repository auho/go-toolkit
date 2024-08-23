package gen

type Req struct {
	Version   int    `json:"version"`
	Topic     string `json:"topic"`
	Partition int    `json:"partition"`
	Replica   int    `json:"replica"`
}

type Body struct {
	Version    int             `json:"version"`
	Partitions []BodyPartition `json:"partitions"`
}

type BodyPartition struct {
	Topic     string `json:"topic"`
	Partition int    `json:"partition"`
	Replicas  []int  `json:"replicas"`
}
