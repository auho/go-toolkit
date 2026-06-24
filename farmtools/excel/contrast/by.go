package contrast

type by struct{}

// Convert 0-based index to 1-based numbering
func (b *by) indexToNo(index int) int {
	return index + 1
}
