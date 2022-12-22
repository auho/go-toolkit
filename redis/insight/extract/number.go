package extract

const KeyPatternNumber KeyPattern = "number"

const keyPatternRegexNumber = `^\d+$`

type Number struct {
	extract
}

func NewNumber() *Number {
	n := &Number{}
	n.init()

	return n
}

func (n *Number) matchNumber(s string) (string, bool) {
	return n.matchKeyPattern(KeyPatternNumber, keyPatternRegexNumber, s)
}

func (n *Number) Match(s string) (string, bool) {
	ss, ok := n.matchNumber(s)
	if !ok {
		return ss, false
	}

	sss, ok := n.matchKeyPatterns(s)
	if ok {
		return sss, true
	} else {
		return ss, true
	}
}
