package FxServicesSlide

import (
	"fmt"
	"github.com/rivo/tview"
)

type FxAppManagerPlateContent struct {
	Grid []IdAndName
}

func newFxAppManagerPlateContent(list []IdAndName) *FxAppManagerPlateContent {
	return &FxAppManagerPlateContent{
		Grid: list,
	}
}

func (self *FxAppManagerPlateContent) GetCell(row, column int) *tview.TableCell {
	switch column {
	case 0:
		switch row {
		case 0:
			return tview.NewTableCell("*").
				SetSelectable(false).
				SetAlign(tview.AlignRight)
		default:
			return tview.NewTableCell("")
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
			return tview.NewTableCell(fmt.Sprintf("%08X", self.Grid[row-1].ServiceId)).
				SetAlign(tview.AlignRight)
		}
	case 4:
		switch row {
		case 0:

			return tview.NewTableCell("Depends On").
				SetSelectable(false).
				SetAlign(tview.AlignRight)
		default:
			return tview.NewTableCell(fmt.Sprintf("%08X", self.Grid[row-1].ServiceDependency)).
				SetAlign(tview.AlignRight)
		}
	}
	return tview.NewTableCell("")
}

func (self *FxAppManagerPlateContent) GetRowCount() int {
	return len(self.Grid) + 1
}

func (self *FxAppManagerPlateContent) GetColumnCount() int {
	return 5
}

func (self *FxAppManagerPlateContent) SetCell(_, _ int, _ *tview.TableCell) {
}

func (self *FxAppManagerPlateContent) RemoveRow(_ int) {
}

func (self *FxAppManagerPlateContent) RemoveColumn(_ int) {
}

func (self *FxAppManagerPlateContent) InsertRow(_ int) {
}

func (self *FxAppManagerPlateContent) InsertColumn(_ int) {
}

func (self *FxAppManagerPlateContent) Clear() {
}
