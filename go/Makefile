build:
	go build -ldflags="-w -s"

install:
	mv pastesearch /usr/bin/pastesearch

termux-install:
	mv pastesearch /data/data/com.termux/files/usr/bin/pastesearch

all: build install

termux-all: build termux-install

clean:
	rm -f pastesearch /usr/bin/pastesearch

termux-clean:
	rm -f pastesearch /data/data/com.termux/files/usr/bin/pastesearch
