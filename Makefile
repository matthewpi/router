SRC_PATH = "go.matthewp.io/router"

PKG_LIST := $(shell go list ${SRC_PATH}/... | grep -v /vendor/)

test:
	@go test -short ${PKG_LIST}

race:
	@go test -race -short ${PKG_LIST}

mem_san:
	@go test -msan -short ${PKG_LIST}
