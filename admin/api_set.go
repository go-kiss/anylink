package admin

import (
	"net/http"

	"github.com/go-kiss/anylink/base"
	"github.com/go-kiss/anylink/dbdata"
	"github.com/go-kiss/anylink/sessdata"
)

func SetHome(w http.ResponseWriter, r *http.Request) {
	data := make(map[string]interface{})

	sess := sessdata.OnlineSess()

	data["counts"] = map[string]int{
		"online": len(sess),
		"user":   dbdata.CountAll(&dbdata.User{}),
		"group":  dbdata.CountAll(&dbdata.Group{}),
		"ip_map": dbdata.CountAll(&dbdata.IpMap{}),
	}

	RespSucess(w, data)
}

func SetSoft(w http.ResponseWriter, r *http.Request) {
	data := base.ServerCfg2Slice()
	RespSucess(w, data)
}

func decimal(f float64) float64 {
	i := int(f * 100)
	return float64(i) / 100
}
