package main

import (
	"io/ioutil"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

type config struct {
	EditWidget *widget.Entry
	PreviewWidget *widget.RichText
	CurrentFile fyne.URI
	SaveMenuItem *fyne.MenuItem
}

var cnf config

func main(){
	// create a fyne app
	a := app.New()

	//create a window for the app
	win := a.NewWindow("Markdown")

	//get the user interface
	edit, preview := cnf.makeUI()
	cnf.createMenuItems(win)

	// set the content of the window
	win.SetContent(container.NewHSplit(edit, preview))

	// show window and run app
	win.Resize(fyne.Size{Width: 800, Height: 600})
	win.CenterOnScreen()
	win.ShowAndRun()
}

func(app *config) makeUI() (*widget.Entry, *widget.RichText){
	edit := widget.NewMultiLineEntry()
	preview := widget.NewRichTextFromMarkdown("")
	app.EditWidget = edit
	app.PreviewWidget = preview

	edit.OnChanged = preview.ParseMarkdown

	return edit, preview
	
}

func (app *config) createMenuItems(win fyne.Window){
	openMenuItem := fyne.NewMenuItem("Open...", app.openFunc(win))

	saveMenuItem := fyne.NewMenuItem("Save", app.saveFunc(win))
	app.SaveMenuItem = saveMenuItem
	app.SaveMenuItem.Disabled = true
	saveAsMenuItem := fyne.NewMenuItem("Save As...", app.saveAsFunc(win))

	fileMenu := fyne.NewMenu("File", openMenuItem, saveMenuItem, saveAsMenuItem)

	menu := fyne.NewMainMenu(fileMenu)

	win.SetMainMenu(menu)
}


var filter = storage.NewExtensionFileFilter([]string{".md", ".MD"})
func (app *config) saveFunc(win fyne.Window) func(){
	return func(){
		if app.CurrentFile != nil{
			write, err := storage.Writer(app.CurrentFile)
			if err != nil{
				dialog.ShowError(err, win)
				return
			}
			write.Write([]byte(app.EditWidget.Text))
			defer write.Close()
		}
	}
}

func (app *config) openFunc(win fyne.Window) func(){
	return func(){
		opednDialog := dialog.NewFileOpen(func(read fyne.URIReadCloser, err error){
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			
			if read == nil{
				return
			}
			defer read.Close()
			
			data, err := ioutil.ReadAll(read)
			if err != nil {
				dialog.ShowError(err, win)
				return
			}
			
			app.EditWidget.SetText(string(data))
			app.CurrentFile = read.URI()
			win.SetTitle(win.Title() + " - " + read.URI().Name())
			app.SaveMenuItem.Disabled = false
			
			}, win)
			opednDialog.SetFilter(filter)
		opednDialog.Show()
	}
}


func (app *config) saveAsFunc(win fyne.Window) func(){
	return func(){
		saveDialog := dialog.NewFileSave(func(write fyne.URIWriteCloser, err error){
			if err != nil{
				dialog.ShowError(err, win)
				return
			}
			
			if write == nil {
				// user cancelled
				return
			}

			if !strings.HasSuffix(strings.ToLower(write.URI().String()), ".md"){
				dialog.ShowInformation("Error", "Pllease name a file with .md extension", win)
				return
			}

			// save file
			write.Write([]byte(app.EditWidget.Text))
			app.CurrentFile = write.URI()
			defer write.Close()

			win.SetTitle(win.Title() + " - " + write.URI().Name())
			app.SaveMenuItem.Disabled = false
		}, win)
		saveDialog.SetFileName("untitled.md")
		saveDialog.SetFilter(filter)
		saveDialog.Show()
	}
}