#这个是Linux的版本，Windows没有 /usr/bin/dokcer 需要手动迁移
CMD=/usr/local/go/bin/go
BIN_PATH=bin
SRC_PATH=src

# make 会调用 all
all: clean build install

# make build 会调用 build
build:
	sudo env GOROOT=/usr/local/go PATH=$$PATH:/usr/local/go/bin $(CMD) build -o $(BIN_PATH)/docker $(SRC_PATH)/*

install:
	cp bin/docker /usr/bin/docker
	cp bin/docker /usr/local/bin/docker

uninstall:
	rm -rf /usr/bin/docker /usr/local/bin/docker
	rm -rf bin/docker

clean: uninstall
