// +build !prod

package admin

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-kiss/anylink/base"
)

func getUiFS() http.FileSystem {
	fmt.Println(base.Cfg.UiPath)
	return http.FS(os.DirFS(base.Cfg.UiPath))
}
