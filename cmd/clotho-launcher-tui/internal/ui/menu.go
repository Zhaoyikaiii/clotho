// Copyright (c) 2026 Clotho contributors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type MenuAction func()

type MenuItem struct {
	Label       string
	Description string
	Action      MenuAction
	Disabled    bool
	MainColor   *tcell.Color
	DescColor   *tcell.Color
}

type Menu struct {
	*tview.Table
	items []MenuItem
}

func NewMenu(title string, items []MenuItem) *Menu {
	table := tview.NewTable().SetSelectable(true, false)
	table.SetBorder(true).SetTitle(title)
	table.SetBorders(false)
	menu := &Menu{Table: table, items: items}
	menu.applyItems(items)
	menu.SetSelectedFunc(func(row, _ int) {
		if row < 0 || row >= len(menu.items) {
			return
		}
		item := menu.items[row]
		if item.Disabled || item.Action == nil {
			return
		}
		item.Action()
	})
	menu.SetSelectedStyle(
		tcell.StyleDefault.Foreground(tview.Styles.InverseTextColor).
			Background(tcell.NewRGBColor(189, 147, 249)),
	)
	return menu
}

func (m *Menu) applyItems(items []MenuItem) {
	m.items = items
	m.Clear()
	for row, item := range items {
		label := item.Label
		if item.Disabled && label != "" {
			label = label + " (disabled)"
		}
		left := tview.NewTableCell(label)
		right := tview.NewTableCell(item.Description).SetAlign(tview.AlignRight)
		if item.MainColor != nil {
			left.SetTextColor(*item.MainColor)
		}
		if item.DescColor != nil {
			right.SetTextColor(*item.DescColor)
		} else {
			right.SetTextColor(tview.Styles.TertiaryTextColor)
		}
		if item.Disabled {
			left.SetTextColor(tcell.ColorGray)
			right.SetTextColor(tcell.ColorGray)
		}
		m.SetCell(row, 0, left)
		m.SetCell(row, 1, right)
	}
}
