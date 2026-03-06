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

const (
	colorBlue = "[#3e5db9]"
	colorRed  = "[#d54646]"
	banner    = "\r\n[::b]" +
		colorBlue + "██████╗ ██╗ ██████╗ ██████╗ " + colorRed + " ██████╗██╗      █████╗ ██╗    ██╗\n" +
		colorBlue + "██╔══██╗██║██╔════╝██╔═══██╗" + colorRed + "██╔════╝██║     ██╔══██╗██║    ██║\n" +
		colorBlue + "██████╔╝██║██║     ██║   ██║" + colorRed + "██║     ██║     ███████║██║ █╗ ██║\n" +
		colorBlue + "██╔═══╝ ██║██║     ██║   ██║" + colorRed + "██║     ██║     ██╔══██║██║███╗██║\n" +
		colorBlue + "██║     ██║╚██████╗╚██████╔╝" + colorRed + "╚██████╗███████╗██║  ██║╚███╔███╔╝\n" +
		colorBlue + "╚═╝     ╚═╝ ╚═════╝ ╚═════╝ " + colorRed + " ╚═════╝╚══════╝╚═╝  ╚═╝ ╚══╝╚══╝\n " +
		"[:]"
)

func applyStyles() {
	tview.Styles.PrimitiveBackgroundColor = tcell.NewRGBColor(12, 13, 22)
	tview.Styles.ContrastBackgroundColor = tcell.NewRGBColor(34, 19, 53)
	tview.Styles.MoreContrastBackgroundColor = tcell.NewRGBColor(18, 18, 32)
	tview.Styles.BorderColor = tcell.NewRGBColor(112, 102, 255)
	tview.Styles.TitleColor = tcell.NewRGBColor(255, 121, 198)
	tview.Styles.GraphicsColor = tcell.NewRGBColor(139, 233, 253)
	tview.Styles.PrimaryTextColor = tcell.NewRGBColor(241, 250, 255)
	tview.Styles.SecondaryTextColor = tcell.NewRGBColor(80, 250, 123)
	tview.Styles.TertiaryTextColor = tcell.NewRGBColor(139, 233, 253)
	tview.Styles.InverseTextColor = tcell.NewRGBColor(12, 13, 22)
	tview.Styles.ContrastSecondaryTextColor = tcell.NewRGBColor(189, 147, 249)
}

func bannerView() *tview.TextView {
	text := tview.NewTextView()
	text.SetDynamicColors(true)
	text.SetTextAlign(tview.AlignCenter)
	text.SetBackgroundColor(tview.Styles.PrimitiveBackgroundColor)
	text.SetText(banner)
	text.SetBorder(false)
	return text
}
