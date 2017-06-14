[![Stories in Ready](https://badge.waffle.io/as27/govuegui.png?label=ready&title=Ready)](https://waffle.io/as27/govuegui?utm_source=badge)
[![GoDoc](https://godoc.org/github.com/as27/govuegui?status.svg)](https://godoc.org/github.com/as27/govuegui)
[![Go Report Card](https://goreportcard.com/badge/github.com/as27/govuegui)](https://goreportcard.com/report/github.com/as27/govuegui)

# govuegui

A web GUI for Go

## Important

This project is still in development and not production ready. You can try it and you can use it for simple internal apps or prototyping. But there are still a lot of open issues. 

## Target of this project

The idea to use a browser as a gui is not new. For example [Electron](https://electron.atom.io/) uses the same principle. It includes a chrome browser and a node.js server. Because of that the binaries are huge. Why not use the allready installed browser? So govuegui just includes a server which is talking to the browser via a websocket connection.

The api of govuegui should be very simple that it is very easy to use.

# Used projects

This package includes following great project:

* [Bulma](http://bulma.io/) as CSS library
* [Gorillatoolkit](http://www.gorillatoolkit.org/) for the websocket
