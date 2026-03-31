public class Reactions: Gtk.Application {
	public Reactions() {
		Object(application_id: "org.readf0x.Reactions");
	}

	public override void activate() {
		var win = new Gtk.ApplicationWindow(this);
		win.set_title("Siffrin Jail");
		win.set_default_size(400, 800);

		var view = new Gtk.ScrolledWindow();

		var flowBox = new Gtk.FlowBox();
		flowBox.selection_mode = Gtk.SelectionMode.NONE;
		flowBox.max_children_per_line = uint.MAX;
		flowBox.row_spacing = 8;
		flowBox.column_spacing = 8;

		var path = GLib.Environment.get_variable("REACTION_PATH");
		// if (path == null) {
		// 	var dialog = new Gtk.FileChooserDialog(
		// 		"Select Reactions Folder",
		// 		win,
		// 		Gtk.FileChooserAction.SELECT_FOLDER,
		// 		"_Cancel", Gtk.ResponseType.CANCEL,
		// 		"_Select", Gtk.ResponseType.ACCEPT
		// 	);

		// 	dialog.response.connect((response) => {
		// 		if (response == Gtk.ResponseType.ACCEPT) {
		// 			path = dialog.get_file().get_path();
		// 		}
		// 		dialog.close();
		// 	});

		// 	dialog.present();
		// }
		var file = File.new_for_path(path);
		try {
			var enumerator = file.enumerate_children("standard::name", 0);
			FileInfo? info;
			File? handle;
			while (enumerator.iterate(out info, out handle)) {
				if (info == null) break;
				Gdk.Texture tex;
				try {
					tex = Gdk.Texture.from_file(handle);
				} catch {
					continue;
				}
				var pic = new Gtk.Picture.for_paintable(tex);
				pic.hexpand = false;
				pic.vexpand = false;
				pic.can_shrink = false;

				var drag = new Gtk.DragSource();
				drag.actions = Gdk.DragAction.COPY;

				pic.add_controller(drag);
				drag.content = new Gdk.ContentProvider.for_bytes(
					"text/uri-list",
					new GLib.Bytes(handle.get_uri().data)
				);
				drag.set_icon(pic.paintable, 0, 0);

				flowBox.append(pic);
			}
		} catch (Error e) {
			GLib.error(e.message);
		}

		view.child = flowBox;

		win.child = view;
		win.present();
	}

	public static int main(string[] args) {
		var app = new Reactions();
		return app.run(args);
	}
}
