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
	for kp, kpReg := range fm.keyPatterns {
		ok := fm.matchString(kp, kpReg, s)
		if ok {
			return string(kp), true
		}
	}

	return s, false
}
