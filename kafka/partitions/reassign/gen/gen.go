package gen

import (
	"encoding/json"
	"sort"

	"github.com/segmentio/kafka-go"
)

func ReassignPartitionsToJson(network, address string, req Req) (string, error) {
	gen := &Gen{}

	return gen.ReassignToJson(network, address, req)
}

type Gen struct{}

func (g *Gen) ReassignToJson(network, address string, req Req) (string, error) {
	body, err := g.Reassign(network, address, req)
	if err != nil {
		return "", err
	}

	_b, err := json.Marshal(body)
	if err != nil {
		return "", err
	}

	return string(_b), nil
}

func (g *Gen) Reassign(network, address string, req Req) (Body, error) {
	body := Body{}
	if req.Version <= 0 {
		req.Version = 1
	}

	body.Version = req.Version

	conn, err := kafka.Dial(network, address)
	if err != nil {
		return body, err
	}

	bs, err := conn.Brokers()
	if err != nil {
		return body, err
	}

	sort.SliceStable(bs, func(i, j int) bool {
		return bs[i].ID < bs[j].ID
	})

	var ids []int
	for _, b := range bs {
		ids = append(ids, b.ID)
	}

	reassignPartitions := g.combination(ids, req.Partition, req.Replica)

	for _p, _ps := range reassignPartitions {
		body.Partitions = append(body.Partitions, BodyPartition{
			Topic:     req.Topic,
			Partition: _p,
			Replicas:  _ps,
		})
	}

	return body, nil
}

func (g *Gen) combination(ids []int, partition, replica int) map[int][]int {
	usedAmount := make(map[int][]int) // map[id used amount][]id
	usedAmount[0] = ids

	rets := make(map[int][]int)

	currentRound := 0
	for i := 0; i < partition; i++ {
		var roundIds []int
		for j := 0; j < replica; j++ {
			if len(usedAmount[currentRound]) <= 0 {
				currentRound++
			}

			_id := usedAmount[currentRound][0]
			roundIds = append(roundIds, _id)
			usedAmount[currentRound] = usedAmount[currentRound][1:]
			usedAmount[currentRound+1] = append(usedAmount[currentRound+1], _id)
		}

		rets[i] = roundIds
	}

	return rets
}
