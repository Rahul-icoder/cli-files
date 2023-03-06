package main

import (
	"io/fs"
	"log"
	"os"
	"path"

	"github.com/Rahul-icoder/cli-files/shared"
	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

func backwardNavigation(list *widgets.List, navigationPath *string) {
	*navigationPath = path.Join(*navigationPath, "..")
	// changing title
	list.Title = *navigationPath
	info, err := os.Stat(*navigationPath)
	shared.CheckError(err)
	if info.IsDir() {
		files := shared.ReadDir(*navigationPath)
		list.Rows = []string{}
		for _, file := range files {

			if len(file.Name()) > 0 {
				if file.IsDir() && file.Name()[0] != '.' {
					list.Rows = append(list.Rows, file.Name())
				}
			}
		}
	}
}

func forwardNavigation(list *widgets.List, navigationPath *string, selectedFile string) {
	*navigationPath = path.Join(*navigationPath, selectedFile)
	// changing title
	list.Title = *navigationPath
	info, err := os.Stat(*navigationPath)
	shared.CheckError(err)

	if info.IsDir() {
		files := shared.ReadDir(*navigationPath)
		list.Rows = []string{}
		for _, file := range files {
			if len(file.Name()) > 0 {
				if file.IsDir() && file.Name()[0] != '.' {
					list.Rows = append(list.Rows, file.Name())
				}
			}
		}
	}
}

func setFileDetails() {

}

func keyEvent(l *widgets.List, navigationPath *string) {
	uiEvents := ui.PollEvents()
	for {

		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			// selectedFile := l.Rows[l.SelectedRow]
			l.ScrollUp()
		case "<C-u>":
			l.ScrollHalfPageUp()
		// navigate in backward direction
		case "h":
			backwardNavigation(l, navigationPath)

		// navigate in forward direction
		case "l":
			// checking condition because of index out of bound error
			// if l.SelectedRow < len(l.Rows) {
			selectedFile := l.Rows[l.SelectedRow]
			forwardNavigation(l, navigationPath, selectedFile)
			// }
		case "<navigationPath>":
			l.ScrollTop()
		}
		if l.SelectedRow > len(l.Rows) {
			l.SelectedRow = 1
		}
		ui.Render(l)
	}
}

func renderFiles(files []fs.DirEntry, navigationPath *string) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize ui: %v", err)
	}
	defer ui.Close()
	list := widgets.NewList()
	list.Title = *navigationPath
	list.Rows = []string{}
	termWidth, termHeight := ui.TerminalDimensions()
	// Create widgets
	fileDetails := widgets.NewParagraph()

	// Create a grid layout
	grid := ui.NewGrid()
	grid.SetRect(0, 0, termWidth, termHeight)
	// Set the grid layout cells
	grid.Set(
		ui.NewRow(
			0.93,
			list,
		),
		ui.NewRow(0.07, fileDetails),
	)
	// styling list
	list.WrapText = false
	list.Border = false
	list.PaddingLeft = 1
	list.PaddingTop = 1
	list.TextStyle.Fg = ui.ColorGreen
	// styling fileDetails
	fileDetails.Border = false
	fileDetails.PaddingLeft = 1
	fileDetails.Text = "rxrx      1mb    2:02   "
	fileDetails.TextStyle.Fg = ui.ColorGreen
	for _, file := range files {
		if file.IsDir() && file.Name()[0] != '.' {
			list.Rows = append(list.Rows, file.Name())
		}
	}
	ui.Render(grid)
	keyEvent(list, navigationPath) // keyboard event
}
func main() {
	path, err := os.UserHomeDir()
	shared.CheckError(err)
	navigationPath := &path
	files := shared.ReadDir(*navigationPath)
	renderFiles(files, navigationPath) // render ui

}
