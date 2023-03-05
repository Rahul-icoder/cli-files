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
	info, err := os.Stat(*navigationPath)
	shared.CheckError(err)
	if info.IsDir() {
		files := shared.ReadDir(*navigationPath)
		list.Rows = []string{}
		for _, file := range files {

			if file.IsDir() && file.Name()[0] != '.' {
				list.Rows = append(list.Rows, file.Name())
			}
		}
	}
}

func forwardNavigation(list *widgets.List, navigationPath *string, selectedFile string) {
	*navigationPath = path.Join(*navigationPath, selectedFile)
	info, err := os.Stat(*navigationPath)
	shared.CheckError(err)
	if info.IsDir() {
		files := shared.ReadDir(*navigationPath)
		list.Rows = []string{}
		for _, file := range files {

			if file.IsDir() && file.Name()[0] != '.' {
				list.Rows = append(list.Rows, file.Name())
			}
		}
	}
}

func keyEvent(l *widgets.List, navigationPath *string) {
	previousKey := ""
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
		case "<C-d>":
			l.ScrollHalfPageDown()
		case "<C-u>":
			l.ScrollHalfPageUp()
		case "<C-f>":
			l.ScrollPageDown()
		case "<C-b>":
			l.ScrollPageUp()
		case "g":
			if previousKey == "g" {
				l.ScrollTop()
			}
		// navigate in backward direction
		case "h":
			backwardNavigation(l, navigationPath)

		// navigate in forward direction
		case "l":
			selectedFile := l.Rows[l.SelectedRow]
			forwardNavigation(l, navigationPath, selectedFile)
		case "<navigationPath>":
			l.ScrollTop()
		case "G", "<End>":
			l.ScrollBottom()
		}

		if previousKey == "g" {
			previousKey = ""
		} else {
			previousKey = e.ID
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
	list.Title = "List"
	list.Rows = []string{}
	list.TextStyle = ui.NewStyle(ui.ColorYellow)
	list.WrapText = false
	list.SetRect(0, 0, 88, 20)

	for _, file := range files {
		if file.IsDir() && file.Name()[0] != '.' {
			list.Rows = append(list.Rows, file.Name())
		}
	}
	ui.Render(list)
	keyEvent(list, navigationPath) // keyboard event
}
func main() {
	path, err := os.UserHomeDir()
	shared.CheckError(err)
	navigationPath := &path
	files := shared.ReadDir(*navigationPath)
	renderFiles(files, navigationPath) // render ui

}
