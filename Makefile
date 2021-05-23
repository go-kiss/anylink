COMMIT_ID= $(shell git rev-parse HEAD)

default: web
	go build -trimpath -tags prod -ldflags "-X main.COMMIT_ID=$(COMMIT_ID)"

web:
	npm --prefix admin/web install
	npm --prefix admin/web run build

clean:
	rm -rf admin/web/ui
	rm -rf anylink
