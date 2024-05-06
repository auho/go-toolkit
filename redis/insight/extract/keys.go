package extract

import "strings"

type Keys struct {
	NumberExtract    *Number
	StringExtract    *String
	FullMatchExtract *FullMatch
}

func NewKeys() *Keys {
	ks := &Keys{}
	ks.NumberExtract = NewNumber()
	ks.StringExtract = NewString()
	ks.FullMatchExtract = NewFullMatch()

	return ks
}

func (ks *Keys) ExtractKey(key string, sep string) string {
	if sep == "" {
		sep = ":"
	}

	// full match
	if _key, ok := ks.FullMatchExtract.Match(key); ok {
		return _key
	} else {
		keyNodes := make([]string, 0, 10)

		nodes := strings.Split(key, sep)
		for _, node := range nodes {
			// number match
			// string match
			keyNodes = append(keyNodes, ks.extractKeyNode(node))
		}

		return strings.Join(keyNodes, ":")
	}
}

func (ks *Keys) extractKeyNode(s string) string {
	ss, ok := ks.NumberExtract.Match(s)
	if ok {
		return ss
	}

	ss, _ = ks.StringExtract.Match(s)
	return ss
}
