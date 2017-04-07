package govuegui

import (
	"github.com/GeertJohan/go.rice/embedded"
	"time"
)

func init() {

	// define files
	file2 := &embedded.EmbeddedFile{
		Filename:    `app.js`,
		FileModTime: time.Unix(1491552238, 0),
		Content:     string("let a = 123;\r\nlet b = \"mystring\";"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   ``,
		DirModTime: time.Unix(1491552238, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file2, // app.js

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`html`, &embedded.EmbeddedBox{
		Name: `html`,
		Time: time.Unix(1491552238, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"": dir1,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"app.js": file2,
		},
	})
}
