build_text:
	go build -o game.exe main.go
build_gui:
	cd ui && go build -o game.exe *.go
	mkdir dist
	cp ui/game.exe ./dist
	cp -r ui/img/ ./dist