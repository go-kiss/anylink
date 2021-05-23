package handler

import (
	"github.com/go-kiss/anylink/admin"
	"github.com/go-kiss/anylink/dbdata"
	"github.com/go-kiss/anylink/sessdata"
)

func Start() {
	dbdata.Start()
	sessdata.Start()

	checkTun()
	go admin.StartAdmin()
	go startTls()
	go startDtls()
}

func Stop() {
	_ = dbdata.Stop()
}
