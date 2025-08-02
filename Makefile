.PHONY: run
run:
	go run ./cmd/client login 123
	
.PHONY: login
login:
	go run ./cmd/client login -u test -p password

.PHONY: register
register:
	go run ./cmd/client register -u test -p password -e e@mail.ee

.PHONY: encrypt
encrypt:
	go run ./cmd/client encrypt -u test -p password -i /home/aube/Videos/js-middle.mkv -o js-middle.bin

.PHONY: decrypt
decrypt:
	go run ./cmd/client decrypt -u test -p password -o /home/aube/Videos/js-middle2.mkv -i js-middle.bin

.PHONY: download
download:
	go run ./cmd/client download -u test -i js-middle.bin'

.PHONY: sync
sync:
	go run ./cmd/client sync -u test

.PHONY: build
build:
	go build -ldflags "                                     \
	-X main.buildVersion=v1.0.1                             \
	-X 'main.buildTime=$$(date +'%Y/%m/%d %H:%M:%S')'       \
	-X 'main.buildCommit=$$(git rev-parse --short HEAD)'    \
	"  -o ./cmd/client/client ./cmd/client/main.go \

.PHONY: buildlint
buildlint:
	go build -o ./cmd/staticlint/staticlint ./cmd/staticlint/main.go

.PHONY: profbase
profbase:
	curl -sK -v http://localhost:8080/debug/pprof/heap?seconds=10 > ./profiles/base.pprof

.PHONY: profres
profres:
	curl -sK -v http://localhost:8080/debug/pprof/heap?seconds=10 > ./profiles/result.pprof

.PHONY: profdiff
profdiff:
	go tool pprof -top -diff_base=profiles/base.pprof profiles/result.pprof

.PHONY: mocks
mocks:
	mockgen -source=./internal/app/store/interfaces.go -destination=./mocks/mocks.go

.PHONY: test
test:
	go test -v -timeout 30s ./...

.PHONY: race
race:
	go test -v -race -timeout 30s ./...

.PHONY: cover
cover:
	go test -v -short -count=1 -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o test-coverage.html
	rm coverage.out

.PHONY: staticcheck
staticcheck:
	staticcheck ./...

.PHONY: protoc
protoc:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/api/grpc/proto/urlclient.proto

.DEFAULT_GOAL := run
