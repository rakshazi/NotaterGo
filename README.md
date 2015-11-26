# NotaterGo

> Based on GTK2

Very simple notes app with tray icon.

### Features

* Create note from tray
* Autosave any note
* All notes in plain text format
That's all. Just create/read/delete notes, no more.

# Folders
NotaterGo releases [FreeDesktop specification for data dir](http://standards.freedesktop.org/basedir-spec/latest/ar01s03.html).

Default: `$XDA_DATA_HOME/NotaterGo`

If $XDA_DATA_HOME env not set, will be used `$HOME/.local/share/NotaterGo` dir.

# Downloads

See in ./bin dir

# Compilation

**go build**
```
go build app.go
```

**gccgo** (doesn't work for now, [more info](https://github.com/mattn/go-gtk/issues/253))
```
go build -a -gccgoflags "-march=native -O3" -compiler gccgo app.go
```
