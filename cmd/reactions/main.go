package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	showdialog := true
	location := flag.String("path", "", "Location of reactions")
	flag.Parse()
	if env, set := os.LookupEnv("REACTION_PATH"); set && env != "" {
		*location = env
	}
	if *location != "" {
		showdialog = false
	}

	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal(err)
	}
	win.SetTitle("Siffrin Jail")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	if showdialog {
		dialog, _ := gtk.FileChooserDialogNewWith2Buttons(
			"Select a folder",
			win,
			gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER,
			"Cancel", gtk.RESPONSE_CANCEL,
			"Select", gtk.RESPONSE_ACCEPT,
		)

		dialog.SetCurrentFolder(os.Getenv("HOME"))

		response := dialog.Run()
		if response == gtk.RESPONSE_ACCEPT {
			*location = dialog.GetFilename()
		}

		dialog.Destroy()
	}

	flowBox, err := gtk.FlowBoxNew()
	flowBox.SetSelectionMode(gtk.SELECTION_NONE)
	flowBox.SetMaxChildrenPerLine(2)
	flowBox.SetRowSpacing(8)
	flowBox.SetColumnSpacing(8)

	files, err := os.ReadDir(*location)
	if err != nil {
		log.Fatal(err)
	}
	var pics []Draggable
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".png") {
			d, err := NewDraggable(fmt.Sprintf("%s%c%s", *location, os.PathSeparator, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			pics = append(pics, d)
			flowBox.Add(d.EventBox)
		}
	}

	scroll, err := gtk.ScrolledWindowNew(nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	scroll.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scroll.Add(flowBox)

	win.Add(scroll)
	win.ShowAll()
	gtk.Main()
}

type Draggable struct {
	Picture *gtk.Image
	EventBox *gtk.EventBox
}

func NewDraggable(path string) (d Draggable, err error) {
	d.Picture, err = gtk.ImageNewFromFile(path)
	if err != nil {
		return
	}
	d.EventBox, err = gtk.EventBoxNew()
	if err != nil {
		return
	}
	d.EventBox.Add(d.Picture)
	ta, _ := gtk.TargetEntryNew("text/uri-list", 0, 0)
	d.EventBox.DragSourceSet(gdk.BUTTON1_MASK, []gtk.TargetEntry{*ta}, gdk.ACTION_COPY)
	d.EventBox.Connect("drag-begin", func(_ *gtk.EventBox, context *gdk.DragContext) {
		gtk.DragSetIconPixbuf(context, d.Picture.GetPixbuf(), 0, 0)
	})
	d.EventBox.Connect("drag-data-get", func(_ *gtk.EventBox, _ *gdk.DragContext, selectionData *gtk.SelectionData, _ uint, _ uint) {
		selectionData.SetURIs([]string{"file://" + path})
	})
	return
}

