package extract

func ExampleNewKeys() {
	_extract := NewKeys()

	_extract.FullMatchExtract.RegisterKeyPatterns(map[string]string{
		"user":    `^cache:user:\d+$`,
		"article": `^cache:article:\d+:list$`,
	})

	_extract.NumberExtract.RegisterKeyPatterns(map[string]string{
		"uid":      `^\d{11}$`,
		"yyyymmdd": `^20\d{6}$`,  // 20200101
		"yyyyww":   `^202\d{3}$`, // 202020
	})

	_extract.StringExtract.RegisterKeyPatterns(map[string]string{
		"uid-uid":    `^\d{11}-\d{11}$`,
		"yyyy-mm-dd": `\d{4}-\d{2}-\d{2}`,
	})

	_ = _extract.ExtractKey("cache:key:name", ":")
}

func ExampleNewFullMatch() {
	_fullMatch := NewFullMatch()
	_fullMatch.RegisterKeyPatterns(map[string]string{
		"user":    `^cache:user:\d+$`,
		"article": `^cache:article:\d+$`,
	})

	_, _ = _fullMatch.Match("cache:key:name")
}

func ExampleNewNumber() {
	_number := NewNumber()
	_number.RegisterKeyPatterns(map[string]string{
		"uid":      `^\d{11}$`,
		"yyyymmdd": `^20\d{6}$`,  // 20200101
		"yyyyww":   `^202\d{3}$`, // 202020
	})

	_, _ = _number.Match("cache")
}

func ExampleNewString() {
	_string := NewString()
	_string.RegisterKeyPatterns(map[string]string{
		"uid-uid":    `^\d{11}-\d{11}$`,
		"yyyy-mm-dd": `\d{4}-\d{2}-\d{2}`,
	})

	_, _ = _string.Match("cache")
}
