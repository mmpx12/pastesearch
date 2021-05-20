Pastesearch
============

Search on pastebin based on psdmp.ws api

### Usage:

usage:
======

```
   -s, --search [QUERY]         General search on pastebin
   -S, --searchv2 [QUERY]       Same as -s but with api v2
   -m, --mail  [EMAIL]          Search for emails 
   -d, --domain [DOMAIN]        Search specific domain
   -o, --save   [directory]     Save paste into directory
   -p, --prefix [prefix]        prefix when save paste
   -b, --browser                Open paste in browser (if result < 20)

   // only for go
   -x, --slow                   Don't use goroutine (slower / less risk of triggering captcha)
```

- If paste exist it will be print in blue in the output.
- If paste was deleted or expired it will be print in red and will not be saved or open in browser.
- If pastebin is not reachable cause of captcha, it will be print in purple

![img](out.png) 

# GO version 

### WARNINGS:

Since it's using goroutine for retrive the paste, dont using it to much in a short period of time or you will be block by captcha.

If you have lot of stuff to check use the bash version, it will be more than 30 time slower than the go version, but you will not have to worry about triggering captcha.

You can alse use the flag `-x|--slow` instead of the bash version


Build:
======

For Linux: 

```sh
make 
sudo make install
# or
sudo make all
```
For Termux:

```sh
make
make termux-install
# or
make termux-all
```

clean:
======

for linux:

```sh
sudo make clean
```

for termux

```sh
make termux-clean
```

### bash 


Lot slower than the go version but more safe (you will not trigger captcha).

### Install 

```sh 
sudo rm -f /usr/bin/pastesearch
sudo cp pastesearch.sh /usr/bin/pastesearch
sudo chmod +x !$
```
