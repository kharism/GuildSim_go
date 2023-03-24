build_text:
	go build -o game.exe main.go
build_gui:
	cd ui && go build -o game.exe *.go
	mkdir dist
	cp ui/game.exe ./dist
	cp -r ui/img/ ./dist
test:
	go test ./internal/cards/cards_test