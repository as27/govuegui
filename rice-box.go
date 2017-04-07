package govuegui

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    `app.js`,
		FileModTime: time.Unix(1491572231, 0),
		Content:     string("let a = 1234;\r\nlet b = \"mystring\";\r\na = 4;"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   ``,
		DirModTime: time.Unix(1491564304, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // app.js

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`lib`, &embedded.EmbeddedBox{
		Name: `lib`,
		Time: time.Unix(1491564304, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"app.js": file2,
		},
	})
}

func init() {

	// define files
	file4 := &embedded.EmbeddedFile{
		Filename:    `index.html`,
		FileModTime: time.Unix(1491571532, 0),
		Content:     string("<!doctype html>\r\n<html>\r\n<head>\r\n    <meta charset=\"utf-8\">\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\r\n    <link rel=\"stylesheet\" href=\"{{ .PathPrefix }}/pure.min.css\" >\r\n    </head>\r\n    <body>\r\n        \r\n    </body>\r\n</html>"),
	}

	// define dirs
	dir3 := &embedded.EmbeddedDir{
		Filename:   ``,
		DirModTime: time.Unix(1491564304, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file4, // index.html

		},
	}

	// link ChildDirs
	dir3.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`html`, &embedded.EmbeddedBox{
		Name: `html`,
		Time: time.Unix(1491564304, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir3,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"index.html": file4,
		},
	})
}
