package main

import (
	"log"
	"os"
	"strings"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

const path = "/home/readf0x/Pictures/Reactions"

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal(err)
	}
	win.SetTitle("Test")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	flowBox, err := gtk.FlowBoxNew()
	flowBox.SetSelectionMode(gtk.SELECTION_NONE)
	flowBox.SetMaxChildrenPerLine(2)
	flowBox.SetRowSpacing(8)
	flowBox.SetColumnSpacing(8)

	files, err := os.ReadDir(path)
	if err != nil {
		log.Fatal(err)
	}
	var pics []Draggable
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".png") {
			d, err := NewDraggable(path + "/" + file.Name())
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

