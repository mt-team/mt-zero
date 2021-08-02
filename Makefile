PROJECT = $(shell ls src | grep -v gateway | grep -v util )

.PHONY: build
build:
	cd src/$(PROJECT); CGO_ENABLED=0 GO111MODULE=on go build -v -o ../../bin/$(PROJECT) $(PROJECT).go

.PHONY: newrpc
newrpc:
	mkdir -p cd src/$(PROJECT)/rpc; cd src/$(PROJECT)/rpc/; echo 'syntax = "proto3";' >> $(PROJECT).proto; echo '' >> $(PROJECT).proto; echo 'package $(PROJECT);' >> $(PROJECT).proto;

.PHONY: rpc
rpc:
	goctl rpc proto -src src/$(PROJECT)/rpc/$(PROJECT).proto -dir src/$(PROJECT)

.PHONY: gateway
gateway:
	goctl api go -api src/gateway/api/gateway.api -dir src/gateway;
	cd src/gateway; CGO_ENABLED=0 GO111MODULE=on go build -v -o ../../bin/gateway gateway.go

.PNONY: $(PROJECT)
$(PROJECT):
	goctl rpc proto -src src/$(@)/rpc/$(@).proto -dir src/$(@);
	cd src/$(@); CGO_ENABLED=0 GO111MODULE=on go build -v -o ../../bin/$(@) $(@).go

.PHONY: model
model:
	cd src/$(PROJECT); goctl model mysql ddl -s ./sql/*.sql -dir ./internal/model --cache

#.PHONY: swagger
## build
#swagger:
#	cd src/$(PROJECTS); swag init