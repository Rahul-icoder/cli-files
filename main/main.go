package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"

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
			if file.IsDir() && file.Name()[0] == '.' {
				continue
			}
			list.Rows = append(list.Rows, file.Name())
		}
	}
}

func forwardNavigation(list *widgets.List, navigationPath *string, selectedFile string) {
	*navigationPath = path.Join(*navigationPath, selectedFile)
	info, err := os.Stat(*navigationPath)
	shared.CheckError(err)
	if info.IsDir() {
		// changing title
		files := shared.ReadDir(*navigationPath)
		list.Rows = []string{}
		for _, file := range files {
			if file.IsDir() && file.Name()[0] == '.' {
				continue
			}
			list.Rows = append(list.Rows, file.Name())
		}
	} else {
		*navigationPath = path.Dir(*navigationPath)
		// open the file
	}
	list.Title = *navigationPath
}

func setFileDetails(fileDetails *widgets.Paragraph, navigationPath *string, selectedFile string) {
	filePath := path.Join(*navigationPath, selectedFile)
	info, err := os.Stat(filePath)
	shared.CheckError(err)
	modTime := info.ModTime().Format("2006-01-02 15:04")
	if info.IsDir() {
		var fileSize int64 = 0
		filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				fileSize += info.Size()
			}
			return nil
		})
		formatedSize := shared.ReadablefileSize(fileSize)
		result := fmt.Sprintf("%s (S)          %s (M)        DIR (T)", formatedSize, modTime)
		fileDetails.Text = result
	} else {
		fileSize := info.Size()
		formatedSize := shared.ReadablefileSize(fileSize)
		result := fmt.Sprintf("%s (S)          %s (M)        DIR (T)", formatedSize, modTime)
		fileDetails.Text = result
	}
}

func keyEvent(l *widgets.List, fileDetails *widgets.Paragraph, navigationPath *string, grid *ui.Grid) {
	uiEvents := ui.PollEvents()
	for {

		e := <-uiEvents
		switch e.ID {
		case "q", "<C-c>":
			return
		case "j", "<Down>":
			l.ScrollDown()
		case "k", "<Up>":
			l.ScrollUp()
		// navigate in backward direction
		case "h":
			backwardNavigation(l, navigationPath)

		// navigate in forward direction
		case "l":
			// checking condition because of index out of bound error
			if l.SelectedRow < len(l.Rows) {
				selectedFile := l.Rows[l.SelectedRow]
				forwardNavigation(l, navigationPath, selectedFile)
			}
		case "<navigationPath>":
			l.ScrollTop()
		}
		// set file details in footer
		if l.SelectedRow < len(l.Rows) {
			selectedFile := l.Rows[l.SelectedRow]
			setFileDetails(fileDetails, navigationPath, selectedFile)
		}
		ui.Render(grid)
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
	fileDetails.TextStyle.Fg = ui.ColorGreen
	for _, file := range files {
		if file.IsDir() && file.Name()[0] == '.' {
			continue
		}
		list.Rows = append(list.Rows, file.Name())

	}
	ui.Render(grid)
	keyEvent(list, fileDetails, navigationPath, grid) // keyboard event
}
func main() {
	path, err := os.UserHomeDir()
	shared.CheckError(err)
	navigationPath := &path
	files := shared.ReadDir(*navigationPath)
	renderFiles(files, navigationPath) // render ui

}
