package describes

import (
	"fmt"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/segmentio/kafka-go"
	"slices"
	"sort"
	"strings"
)

func DescribeAllTopics(network, address string) (table.Writer, error) {
	t, err := NewTopics(network, address)
	if err != nil {
		return nil, err
	}

	_table, err := t.All()
	if err != nil {
		return nil, err
	}

	return _table, nil
}

type TopicsOverview struct {
	connection *kafka.Conn
}

func NewTopics(network, address string) (*TopicsOverview, error) {
	var err error
	t := &TopicsOverview{}
	t.connection, err = kafka.Dial(network, address)
	if err != nil {
		return nil, err
	}

	return t, err
}

func (t *TopicsOverview) All() (table.Writer, error) {
	topics, err := t.GetAllTopics()
	if err != nil {
		return nil, err
	}

	_table := table.NewWriter()
	_table.AppendHeader(table.Row{"topic", "partition id", "in-sync", "replica", "leader", "in-sync ids", "replica ids"})

	for _, topic := range topics {
		_table.AppendRow(table.Row{topic.Topic})

		sort.SliceStable(topic.Partitions, func(i, j int) bool {
			return topic.Partitions[i].Id < topic.Partitions[j].Id
		})

		for _, partition := range topic.Partitions {
			var inSyncNodesId, replicaNodesId []string

			for _, _p := range partition.Isr {
				inSyncNodesId = append(inSyncNodesId, fmt.Sprintf("%d", _p.ID))
			}

			slices.Sort(inSyncNodesId)

			for _, _p := range partition.Replicas {
				replicaNodesId = append(replicaNodesId, fmt.Sprintf("%d", _p.ID))
			}

			slices.Sort(replicaNodesId)

			_row := table.Row{"", partition.Id, len(partition.Isr), len(partition.Replicas), partition.Leader.ID,
				strings.Join(inSyncNodesId, ","), strings.Join(replicaNodesId, ","),
			}

			_table.AppendRow(_row)
		}

		_table.AppendSeparator()
	}

	return _table, nil
}

func (t *TopicsOverview) GetAllTopics() ([]Topic, error) {
	partitions, err := t.connection.ReadPartitions()
	if err != nil {
		return nil, err
	}

	topicFlags := make(map[string]int)
	var topics []Topic

	for _, _p := range partitions {
		if _p.Topic == "__consumer_offsets" {
			continue
		}

		_tp := Partition{
			Id:       _p.ID,
			Leader:   _p.Leader,
			Replicas: _p.Replicas,
			Isr:      _p.Isr,
		}

		if index, ok := topicFlags[_p.Topic]; ok {
			topics[index].Partitions = append(topics[index].Partitions, _tp)
		} else {
			topics = append(topics, Topic{
				Topic:      _p.Topic,
				Partitions: []Partition{_tp},
			})

			topicFlags[_p.Topic] = len(topicFlags)
		}
	}

	return topics, nil
}

type Topic struct {
	Topic      string
	Partitions []Partition
}

type Partition struct {
	Id       int
	Leader   kafka.Broker
	Replicas []kafka.Broker
	Isr      []kafka.Broker
}
