build:
	make clean & make darwin && make cp && cd bin

.PHONY: clean darwin linux windows build run cp
cp:
	cp bin/erproxy-darwin test/erproxy-darwin 
run:
	cd bin && ./erproxy-darwin

darwin:
	go build -o bin/erproxy-darwin erproxy

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/erproxy-linux erproxy

windows:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o bin/erproxy-windows.exe erproxy

all:
	make darwin && make linux && make windows

clean:
	rm -r bin/erproxy-* bin/erproxy.log