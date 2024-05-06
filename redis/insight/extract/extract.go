package extract

import (
	"fmt"
	"regexp"
)

type extract struct {
	regexps     map[string]*regexp.Regexp // map[key pattern]key pattern representation of a compiled regular expression
	keyPatterns map[string]string         // map[key pattern]key pattern regular expression
}

func (e *extract) init() {
	e.regexps = make(map[string]*regexp.Regexp)
	e.keyPatterns = make(map[string]string)
}

func (e *extract) ShowKeyPattern() {
	for k, v := range e.keyPatterns {
		fmt.Println(k, v)
	}
}

func (e *extract) RegisterKeyPatterns(m map[string]string) {
	for k, v := range m {
		e.RegisterKeyPattern(k, v)
	}
}

func (e *extract) RegisterKeyPattern(kp, regex string) {
	e.keyPatterns[kp] = regex
	e.regexps[kp] = regexp.MustCompile(regex)
}

func (e *extract) match(s string) (string, bool) {
	for k := range e.keyPatterns {
		kp, ok := e.matchKeyPattern(k, s)
		if ok {
			return kp, true
		}
	}

	return s, false
}

func (e *extract) matchKeyPattern(kp string, s string) (string, bool) {
	re := e.regexps[kp]
	ok := re.MatchString(s)
	if ok {
		return e.formatKeyPattern(kp), true
	} else {
		return s, false
	}
}

func (e *extract) formatKeyPattern(kp string) string {
	return "<" + string(kp) + ">"
}
