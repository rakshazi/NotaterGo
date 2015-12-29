package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

var appName = "NotaterGo"
var configDir, dataDir string
var start, end gtk.TextIter //For getting text from textview
var notes []string
var menu *gtk.Menu

func main() {
	if envVar := os.Getenv("XDG_DATA_HOME"); envVar == "" {
		dataDir = os.Getenv("HOME") + "/.local/share/" + appName
	} else {
		dataDir = os.Getenv("XDG_DATA_HOME") + "/" + appName
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		os.MkdirAll(dataDir, 0700)
	}

	gtk.Init(&os.Args)
	glib.SetApplicationName(appName)
	createSystray()

	gtk.Main()
}

func createSystray() {
	menu = gtk.NewMenu()
	updateList()

	statusIcon := gtk.NewStatusIconFromStock(gtk.STOCK_FILE)
	statusIcon.SetTitle(appName)
	statusIcon.SetTooltipMarkup(appName)
	statusIcon.Connect("popup-menu", func(cbx *glib.CallbackContext) {
		menu.Popup(nil, nil, gtk.StatusIconPositionMenu, statusIcon, uint(cbx.Args(0)), uint32(cbx.Args(1)))
	})
}

func createEditorWindow(note string) {
	//Window
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Note " + note)
	window.SetIconName("gtk-about")
	vbox := gtk.NewVBox(false, 1)

	//Menu
	menubar := gtk.NewMenuBar()
	vbox.PackStart(menubar, false, false, 0)
	cascademenu := gtk.NewMenuItemWithMnemonic("_File")
	menubar.Append(cascademenu)
	submenu := gtk.NewMenu()
	cascademenu.SetSubmenu(submenu)

	menuitem := gtk.NewMenuItemWithMnemonic("D_elete note")
	menuitem.Connect("activate", func() {
		deleteNote(note)
		window.Destroy()
	})
	submenu.Append(menuitem)

	//Text view
	textview := gtk.NewTextView()
	textview.SetEditable(true)
	textview.SetCursorVisible(true)
	buffer := textview.GetBuffer()
	buffer.SetText(string(readNote(note)))
	textview.SetSizeRequest(500, 300)
	vbox.PackStart(textview, false, false, 0)

	window.Connect("destroy", func() {
		//Remove (clear) file
		deleteNote(note)
		//Save text to file
		buffer.GetStartIter(&start)
		buffer.GetEndIter(&end)
		writeNote(note, buffer.GetText(&start, &end, true))

		updateList()
	})
	window.Add(vbox)
	window.SetSizeRequest(500, 320)
	window.ShowAll()
}

func updateList() {
	var i uint
	items := menu.GetChildren()
	for i = 0; i < items.Length(); i++ {
		w := gtk.WidgetFromNative(items.NthData(i))
		menu.Remove(w)
	}
	notes = getNotes()
	for i := 0; i < len(notes); i++ {
		note := notes[i]
		note = note[len(dataDir)+1:]
		item := gtk.NewMenuItemWithLabel(note)
		item.Connect("activate", func() {
			createEditorWindow(note)
		})
		menu.Append(item)
	}

	item := gtk.NewMenuItemWithLabel("New note")
	item.Connect("activate", func() {
		createEditorWindow(getCurrentTime() + ".txt")
	})
	menu.Append(item)

	item = gtk.NewMenuItemWithLabel("Exit")
	item.Connect("activate", func() {
		gtk.MainQuit()
	})
	menu.Append(item)
	menu.ShowAll()
}

func getCurrentTime() string {
	currentTime := time.Now()

	return currentTime.Format("02.01.2006 15:04:05")
}

func checkError(e error) {
	if e != nil {
		panic(e)
	}
}

func getNotes() []string {
	files, e := filepath.Glob(filepath.Join(dataDir, "*.txt"))
	checkError(e)
	return files
}

func writeNote(file string, text string) {
	f, e := os.OpenFile(filepath.Join(dataDir, file), os.O_CREATE|os.O_WRONLY, 0666)
	checkError(e)
	f.WriteString(text)
	f.Close()
}

func readNote(file string) string {
	if _, err := os.Stat(filepath.Join(dataDir, file)); err == nil {
		content, e := ioutil.ReadFile(filepath.Join(dataDir, file))
		checkError(e)

		return string(content)
	}
	return ""
}

func deleteNote(file string) {
	os.Remove(dataDir + "/" + file)
}
