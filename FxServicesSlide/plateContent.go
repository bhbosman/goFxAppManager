package FxServicesSlide

import (
	"fmt"
	"github.com/bhbosman/goFxAppManager/FxServicesSlide/internal"
	"github.com/rivo/tview"
	"strconv"
)

type PlateContent struct {
	Grid []internal.IdAndName
}

func newFxAppManagerPlateContent(list []internal.IdAndName) *PlateContent {
	return &PlateContent{
		Grid: list,
	}
}

func (self *PlateContent) GetCell(row, column int) *tview.TableCell {
	if row == -1 || column == -1 {
		return tview.NewTableCell("")
	}

	switch column {
	case 0:
		switch row {
		case 0:
			return tview.NewTableCell("*").
				SetSelectable(false).
				SetAlign(tview.AlignRight)
		default:
			return tview.NewTableCell(strconv.Itoa(row)).SetSelectable(false).SetAlign(tview.AlignRight)
		}
	case 1:
		switch row {
		case 0:
			return tview.NewTableCell("Name").
				SetSelectable(false).
				SetAlign(tview.AlignRight)
		default:
			return tview.NewTableCell(self.Grid[row-1].Name)
		}
	case 2:
		switch row {
		case 0:
			return tview.NewTableCell("Active").
				SetSelectable(false).
				SetAlign(tview.AlignRight)
		default:
			return tview.NewTableCell(fmt.Sprintf("%v", self.Grid[row-1].Active)).
				SetAlign(tview.AlignRight)
		}
	case 3:
		switch row {
		case 0:
			return tview.NewTableCell("Id").
				SetSelectable(false).
				SetAlign(tview.AlignRight)
		default:
			return tview.NewTableCell(fmt.Sprintf("%08X", 0)).
				SetAlign(tview.AlignRight)
		}
	case 4:
		switch row {
		case 0:

			return tview.NewTableCell("Depends On").
				SetSelectable(false).
				SetAlign(tview.AlignRight)
		default:
			return tview.NewTableCell(fmt.Sprintf("%08X", 0)).
				SetAlign(tview.AlignRight)
		}
	}
	return tview.NewTableCell("")
}

func (self *PlateContent) GetRowCount() int {
	return len(self.Grid) + 1
}

func (self *PlateContent) GetColumnCount() int {
	return 5
}

func (self *PlateContent) SetCell(_, _ int, _ *tview.TableCell) {
}

func (self *PlateContent) RemoveRow(_ int) {
}

func (self *PlateContent) RemoveColumn(_ int) {
}

func (self *PlateContent) InsertRow(_ int) {
}

func (self *PlateContent) InsertColumn(_ int) {
}

func (self *PlateContent) Clear() {
}
