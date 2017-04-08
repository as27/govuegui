[![GoDoc](https://godoc.org/github.com/as27/govuegui?status.svg)](https://godoc.org/github.com/as27/govuegui)

# govuegui
A web GUI for Go

# Development settings

When developing on the html files or the vue app you need to set the unexported variable `useRice = false`.

That the changes takes effect to the package, after every change inside of the html folder you need to run:

```
rice embed-go
```

# Used projects

This package includes following great project:

* [Pure.css](https://purecss.io/)
* [Gorillatoolkit](http://www.gorillatoolkit.org/)