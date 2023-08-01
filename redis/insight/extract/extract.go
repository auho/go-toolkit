package extract

import (
	"fmt"
	"regexp"
)

type KeyPattern string

type extract struct {
	regexps     map[KeyPattern]*regexp.Regexp // map[key pattern]key pattern representation of a compiled regular expression
	keyPatterns map[KeyPattern]string         // map[key pattern]key pattern regular expression
}

func (e *extract) init() {
	e.regexps = make(map[KeyPattern]*regexp.Regexp)
	e.keyPatterns = make(map[KeyPattern]string)
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
	e.keyPatterns[KeyPattern(kp)] = regex
	e.regexps[KeyPattern(kp)] = regexp.MustCompile(regex)
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

func (e *extract) matchKeyPattern(kp KeyPattern, s string) (string, bool) {
	re := e.regexps[kp]
	ok := re.MatchString(s)
	if ok {
		return e.formatKeyPattern(kp), true
	} else {
		return s, false
	}
}

func (e *extract) formatKeyPattern(kp KeyPattern) string {
	return "<" + string(kp) + ">"
}
