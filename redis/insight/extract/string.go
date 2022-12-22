package extract

type String struct {
	extract
}

func NewString() *String {
	s := &String{}
	s.init()

	return s
}

func (s *String) Match(t string) (string, bool) {
	return s.matchKeyPatterns(t)
}
