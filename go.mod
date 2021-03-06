module github.com/go-kiss/anylink

go 1.16

require (
	github.com/asdine/storm/v3 v3.2.1
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/pion/dtls/v2 v2.0.0-00010101000000-000000000000
	github.com/pion/logging v0.2.2
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	github.com/songgao/water v0.0.0-20200317203138-2b4b6d7c09d8
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	github.com/xhit/go-simple-mail/v2 v2.9.0
	github.com/xlzd/gotp v0.0.0-20181030022105-c8557ba2c119
	go.etcd.io/bbolt v1.3.5
	golang.org/x/crypto v0.0.0-20210513164829-c07d793c2f9a
	golang.org/x/net v0.0.0-20210520170846-37e1c6afe023 // indirect
	golang.org/x/time v0.0.0-20210220033141-f8bda1e9f3ba
)

replace github.com/pion/dtls/v2 => ./dtls-2.0.9
