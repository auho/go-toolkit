package output

type MultilineText struct {
	content []string
}

func NewMultilineText() *MultilineText {
	return &MultilineText{}
}

// CoverAll 覆盖全部内容
func (m *MultilineText) CoverAll(sss ...[]string) {
	m.content = []string{}
	for _, ss := range sss {
		m.content = append(m.content, ss...)
	}
}

// Print 在第几行打印内容
// line 从 1 开始
func (m *MultilineText) Print(line int, s string) {
	if line < 0 {
		line = 1
	}

	contentLen := len(m.content)
	if line > contentLen {
		for i := 0; i < line-contentLen; i++ {
			m.content = append(m.content, "")
		}
	}

	m.content[line-1] = s
}

func (m *MultilineText) PrintNext(s string) {
	m.content = append(m.content, s)
}

func (m *MultilineText) Content() []string {
	return m.content
}
