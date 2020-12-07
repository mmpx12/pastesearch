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
   -f, --firefox                Open paste in firefox (if result < 20)
```

- If paste exist it will be print in blue in the output.
- If paste was deleted or expired it will be print in red and will not be saved or open in firefox.

![img](out.png) 
