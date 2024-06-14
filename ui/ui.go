package ui

import "github.com/rivo/tview"

func DisplayText(text string) {
	app := tview.NewApplication()
	textView := tview.NewTextView().
		SetDynamicColors(true).
		SetText(text).
		SetTextAlign(tview.AlignLeft).
		SetScrollable(true).
		SetChangedFunc(func() {
			app.Draw()
		})

	if err := app.SetRoot(textView, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}
