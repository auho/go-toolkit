package extract

import (
	"fmt"
	"regexp"
)

type KeyPattern string

type extract struct {
	regexps     map[KeyPattern]*regexp.Regexp
	keyPatterns map[KeyPattern]string // map[key pattern]key pattern regex
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
}

func (e *extract) matchKeyPatterns(s string) (string, bool) {
	for k, v := range e.keyPatterns {
		kp, ok := e.matchKeyPattern(k, v, s)
		if ok {
			return kp, true
		}
	}

	return s, false
}

func (e *extract) matchKeyPattern(kp KeyPattern, regex string, s string) (string, bool) {
	ok := e.matchString(kp, regex, s)
	if ok {
		return e.formatKeyPattern(kp), true
	} else {
		return s, false
	}
}

func (e *extract) matchString(kp KeyPattern, regex string, s string) bool {
	re := e.generateKeyPatternRe(kp, regex)
	return re.MatchString(s)
}

func (e *extract) generateKeyPatternRe(kp KeyPattern, regex string) *regexp.Regexp {
	if re, ok := e.regexps[kp]; ok {
		return re
	} else {
		re = regexp.MustCompile(regex)
		e.regexps[kp] = re

		return re
	}
}

func (e *extract) formatKeyPattern(kp KeyPattern) string {
	return "<" + string(kp) + ">"
}
