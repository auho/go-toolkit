package extract

import "strings"

type Key struct {
	NumberExtract *Number
	StringExtract *String
}

func NewKey() *Key {
	k := &Key{}
	k.NumberExtract = NewNumber()
	k.StringExtract = NewString()

	return k
}

func (k *Key) ExtractKey(key string) []string {
	keyNodes := make([]string, 0, 10)

	nodes := strings.Split(key, ":")
	for _, node := range nodes {
		keyNodes = append(keyNodes, k.extractKeyNode(node))
	}

	return keyNodes
}

func (k *Key) extractKeyNode(s string) string {
	ss, ok := k.NumberExtract.Match(s)
	if ok {
		return ss
	}

	ss, ok = k.StringExtract.Match(s)
	return ss
}
