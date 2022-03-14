GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
COUNT=$$(docker ps -a -q)

fmt:
	gofmt -w $(GOFMT_FILES)

vet:
	$Qgo vet ./...

test: vet
	$Qgo test -count 1 -p 1 ./...

generate: deps
	cd api ;buf generate
	# gopatch -p test.patch pkg/proto/flago.pb.gw.go

deps:
	$Qgo install github.com/uber-go/gopatch@latest
run:
	cd cmd/server ;$Qgo run .