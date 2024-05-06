package extract

const keyPatternNumber = "number"
const keyPatternRegexNumber = `^\d+$`

type Number struct {
	extract
}

func NewNumber() *Number {
	n := &Number{}
	n.init()

	n.RegisterKeyPattern(string(keyPatternNumber), keyPatternRegexNumber)

	return n
}

func (n *Number) Match(s string) (string, bool) {
	ss, ok := n.matchKeyPattern(keyPatternNumber, s)
	if !ok {
		return ss, false
	}

	sss, ok := n.match(s)
	if ok {
		return sss, true
	} else {
		return ss, true
	}
}
