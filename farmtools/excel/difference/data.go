package difference

type Data struct {
	baseData [][]string
	diffData [][]string
}

func (d *Data) diff() (SheetDiff, error) {
	addedRows := make(map[int][]string)
	deletedRows := make(map[int][]string)
	var modifiedCells []CellDiff

	baseRowNoMax := len(d.baseData)
	diffRowNoMax := len(d.diffData)

	if baseRowNoMax < diffRowNoMax { // length of rows of base less than diff. new cells (row)
		for i := baseRowNoMax + 1; i < diffRowNoMax; i++ {
			addedRows[i] = d.diffData[i]
		}
	} else { // length of rows of base greater than diff. new cells (row)
		for baseRowNo, baseRowValues := range d.baseData { // iteration row
			// rows
			if baseRowNo > diffRowNoMax { // delete cells (row)
				deletedRows[baseRowNo] = baseRowValues
			} else { // change cells (row)

				// cols
				diffColNoMax := len(d.diffData[baseRowNo])
				baseColNoMax := len(baseRowValues)

				if baseColNoMax >= diffColNoMax {
					for baseColNo, baseColValue := range baseRowValues { // iteration col

						if baseColNo <= diffColNoMax { // change cells (col)
							diffColValue := d.diffData[baseRowNo][baseColNo]
							if baseColValue != diffColValue {
								modifiedCells = append(modifiedCells, CellDiff{
									row:    baseRowNo,
									col:    baseColNo,
									action: actionModify,
									value:  diffColValue,
								})
							}
						} else { // delete cells (col)
							modifiedCells = append(modifiedCells, CellDiff{
								row:    baseRowNo,
								col:    baseColNo,
								action: actionDelete,
								value:  baseColValue,
							})
						}
					}
				} else { // new cells (col)
					for i := baseColNoMax + 1; i <= diffColNoMax; i++ {
						modifiedCells = append(modifiedCells, CellDiff{
							row:    baseRowNo,
							col:    i,
							action: actionAdd,
							value:  d.diffData[baseRowNo][i],
						})
					}
				}
			}
		}
	}

	return SheetDiff{
		addedRows:     addedRows,
		deletedRows:   deletedRows,
		modifiedCells: modifiedCells,
	}, nil
}
