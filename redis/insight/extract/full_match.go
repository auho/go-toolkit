package extract

type FullMatch struct {
	extract
}

func NewFullMatch() *FullMatch {
	fm := &FullMatch{}
	fm.init()

	return fm
}

func (fm *FullMatch) Match(s string) (string, bool) {
	return fm.match(s)
}
