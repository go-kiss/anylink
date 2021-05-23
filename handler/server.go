package handler

import (
	"crypto/tls"
	"log"
	"net"
	"net/http"

	"github.com/go-kiss/anylink/base"
	"github.com/gorilla/mux"
)

func startTls() {
	addr := base.Cfg.ServerAddr
	certFile := base.Cfg.CertFile
	keyFile := base.Cfg.CertKey

	tlsConfig := &tls.Config{
		NextProtos:         []string{"http/1.1"},
		MinVersion:         tls.VersionTLS12,
		InsecureSkipVerify: true,
	}
	srv := &http.Server{
		Addr:      addr,
		Handler:   initRoute(),
		TLSConfig: tlsConfig,
		ErrorLog:  base.GetBaseLog(),
	}

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	base.Info("listen server", addr)
	err = srv.ServeTLS(ln, certFile, keyFile)
	if err != nil {
		base.Fatal(err)
	}
}

func initRoute() http.Handler {
	r := mux.NewRouter()
	r.HandleFunc("/", LinkHome).Methods(http.MethodGet)
	r.HandleFunc("/", LinkAuth).Methods(http.MethodPost)
	r.HandleFunc("/CSCOSSLC/tunnel", LinkTunnel).Methods(http.MethodConnect)
	r.HandleFunc("/otp_qr", LinkOtpQr).Methods(http.MethodGet)
	r.PathPrefix("/files/").Handler(
		http.StripPrefix("/files/",
			http.FileServer(http.Dir(base.Cfg.FilesPath)),
		),
	)

	return r
}
