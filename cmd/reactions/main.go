package main

import (
	"context"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/diamondburned/gotk4/pkg/gdk/v4"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
)

var showdialog = true
var location string

var imageExts = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
}

func main() {
	flag.StringVar(&location, "path", "", "Location of reactions")
	flag.Parse()
	if env, set := os.LookupEnv("REACTION_PATH"); set && env != "" {
		location = env
	}
	if location != "" {
		showdialog = false
	}

	app := gtk.NewApplication("com.github.readf0x.reactions", gio.ApplicationFlagsNone)
	app.ConnectActivate(func() { activate(app) })

	if code := app.Run(nil); code > 0 {
		os.Exit(code)
	}
}

func activate(app *gtk.Application) {
	win := gtk.NewApplicationWindow(app)
	win.SetTitle("Siffrin Jail")
	win.SetDefaultSize(400, 800)

	if showdialog {
		dialog := gtk.NewFileDialog()

		dialog.SelectFolder(context.Background(), nil, func(res gio.AsyncResulter) {
			file, err := dialog.SelectFolderFinish(res)
			if err != nil {
				log.Fatal(err)
			}
			location = file.ParseName()
			createWindow(win)
		})
	} else {
		createWindow(win)
	}
}

func createWindow(win *gtk.ApplicationWindow) {
	flowBox := gtk.NewFlowBox()
	flowBox.SetSelectionMode(gtk.SelectionNone)
	flowBox.SetHomogeneous(false)
	flowBox.SetMaxChildrenPerLine(2)
	flowBox.SetRowSpacing(8)
	flowBox.SetColumnSpacing(8)

	files, err := os.ReadDir(location)
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if _, ok := imageExts[ext]; ok {
			pic := DraggableImage(filepath.Join(location, file.Name()))
			if err != nil {
				log.Fatal(err)
			}
			flowBox.Append(pic)
		}
	}

	scroll := gtk.NewScrolledWindow()
	scroll.SetPolicy(gtk.PolicyAutomatic, gtk.PolicyAutomatic)
	scroll.SetChild(flowBox)

	win.SetChild(scroll)
	win.SetVisible(true)
}

func DraggableImage(path string) (pic *gtk.Picture) {
	tex, err := gdk.NewTextureFromFilename(path)
	if err != nil {
		log.Fatal(err)
	}
	pic = gtk.NewPictureForPaintable(tex)
	pic.SetHExpand(false)
	pic.SetVExpand(false)
	pic.SetCanShrink(false)

	source := gtk.NewDragSource()
	source.SetActions(gdk.ActionCopy)
	pic.AddController(source)
	source.ConnectPrepare(generatePrepare(path))
	source.SetIcon(pic.Paintable(), 0, 0)

	return
}

func generatePrepare(path string) func(_, _ float64) (contentProvider *gdk.ContentProvider) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}
	mime := imageExts[filepath.Ext(path)]
	gbytes := glib.NewBytes(bytes)
	return func(_, _ float64) (contentProvider *gdk.ContentProvider) {
		uri := gdk.NewContentProviderForBytes("text/uri-list", glib.NewBytes([]byte("file://" + path)))
		raw := gdk.NewContentProviderForBytes(mime, gbytes)
		contentProvider = gdk.NewContentProviderUnion([]*gdk.ContentProvider{uri,raw})
		return
	}
}

