// +build prod

package admin

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed web/ui
var embededFiles embed.FS

func getUiFS() http.FileSystem {
	fsys, err := fs.Sub(embededFiles, "web/ui")
	if err != nil {
		panic(err)
	}

	return http.FS(fsys)
}
