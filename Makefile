BINARY=terraform-provider-todoist

default: install

build:
	go build -o ${BINARY}

install: build
	mv ${BINARY} ~/.terraform.d/plugins
