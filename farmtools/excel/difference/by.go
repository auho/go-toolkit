package difference

type by struct{}

func (b *by) indexToNo(index int) int {
	return index + 1
}
