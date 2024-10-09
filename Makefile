.POSIX:

include config.mk

BIN_NAME = netac

all: build

build:
	go build -o $(NAME) $(CMD_PATH)

clean:
	rm -rf $(NAME)

install: $(NAME)
	mkdir -p $(PREFIX)/bin
	cp $(NAME) $(PREFIX)/bin/$(NAME)
	chmod 755 $(PREFIX)/bin/$(NAME)

uninstall:
	rm $(PREFIX)/bin/$(NAME)

.PHONY: all build clean install uninstall
