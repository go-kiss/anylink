// +build !prod

package admin

import (
	"net/http"
	"os"

	"github.com/go-kiss/anylink/base"
)

func getUiFS() http.FileSystem {
	return http.FS(os.DirFS(base.Cfg.UiPath))
}
