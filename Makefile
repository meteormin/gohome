PRJ_NAME=IVDDMS DM Server
PRJ_DESC=$(PRJ_NAME) Deployments Makefile
PRJ_BASE=$(shell pwd)
BUILD_DIR=$(PRJ_BASE)/build

.DEFAULT: help
.SILENT:;

##help: helps (default)
.PHONY: help
help: Makefile
	echo ""
	echo " $(PRJ_DESC)"
	echo ""
	echo " Usage:"
	echo ""
	echo "	make {command}"
	echo ""
	echo " Commands:"
	echo ""
	sed -n 's/^##/	/p' $< | column -t -s ':' |  sed -e 's/^/ /'
	echo ""

##clean: clean build directory
.PHONY: clean
clean:
	rm -rf $(BUILD_DIR)

##run cmd={gohome-cli,gohome-gui,detector} flags={go flags}: run $(cmd)
.PHONY: run
run:
	go run $(PRJ_BASE)/cmd/$(cmd)/main.go $(flags)

##build cmd={gohome-cli,gohome-gui,detector}: build $(cmd)
.PHONY: build
build:
	go build -o $(BUILD_DIR)/$(cmd) $(PRJ_BASE)/cmd/$(cmd)/main.go