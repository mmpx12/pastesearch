#!/usr/bin/bash

usage(){
  echo -e """pastebin searcher

usage:
   -s, --search [QUERY]         General search on pastebin
   -S, --searchv2 [QUERY]       Same as -s but with api v2
   -m, --mail  [EMAIL]          Search for emails 
   -d, --domain [DOMAIN]        Search specific domain
   -o, --save   [directory]     Save paste into directory
   -p, --prefix [prefix]        prefix when save paste
   -f, --firefox                Open paste in firefox (if result < 20)

output:

\e[32m[\e[31m\$link\e[32m]-→\e[32m[\e[33m\$date\e[32m] : Paste was removed
\e[32m[\e[36m\$link\e[32m]-→\e[32m[\e[33m\$date\e[32m] : Paste exists
\nSome paste also have tags:
-→\e[32m[\e[35m\$tags\e[32m]\n
"""
exit 1
}

res=""
if [[ ${#@} > 1 ]]; then
  while [ "$1" != "" ]; do
    case $1 in
      -s | --search )
        shift
        search=true
        search_query="$1"
        ;;
      -S | --searchv2 )
        shift
        searchv2=true
        searchv2_query="$1"
        ;;
      -m | --mail )
        shift
        mail=true
        mail_query="$1"
        ;;
      -d | --domain )
        shift
        domain=true
        domain_query="$1"
        ;;
      -f | --firefox )
        firefox=true
        ;;
      -o | --save )
        shift
        save=true
        save_dir="$1"
        ;;
      -p | --prefix)
        shift
        prefix="$1"
        ;;
      * )
        usage
        ;;
    esac 
    shift
  done
elif [[ ${#@} == 1 ]]; then
  res=$(curl -s "https://psbdmp.ws/api/search/$1" | jq -r '.data[] | [.id +"@"+ if .tags == "" then "none" else .tags end + "@" + .time] |@tsv')
else 
  usage
fi


[[ $search == true ]] && res="$(curl -s "https://psbdmp.ws/api/search/$search_query" | jq -r '.data[] | [.id +"@"+ if .tags == "" then "none" else .tags end + "@" + .time] |@tsv')"

[[ $searchv2 == true ]] && res="${res}$(echo -e "" && curl -s "https://psbdmp.ws/api/v2/search/$searchv2_query" | jq -r '.data[] | [.id +"@"+ if .tags == "" then "none" else .tags end + "@" + .time] |@tsv')"

[[ $mail == true ]] && res="${res}$(echo -e "" && curl -s "https://psbdmp.ws/api/search/email/$mail_query" | jq -r '.data[] | [.id +"@"+ if .tags == "" then "none" else .tags end + "@" + .time] |@tsv')"
#pwndb --target "@$mail_query" | sed -r "s/[[:cntrl:]]\[[0-9]{1,3}m//g" | sed 's/\[+\]\t//g

[[ $domain == true ]] && res="${res}$(echo -e "" && curl -s "https://psbdmp.ws/api/search/domain/$domain_query" | jq -r '.data[] | [.id +"@"+ if .tags == "" then "none" else .tags end + "@" + .time] |@tsv')"


[[ $firefox == true && $(wc -l <<<"$res") -gt 15 ]] && echo "Too much result for firefox" && ff=""
[[ $firefox == true && $(wc -l <<<"$res") -lt  15 ]] && ff="firefox"
[[ $firefox == true ]] || ff=""
[[ $save == true && ! -d "$save_dir" ]] && mkdir -p "$save_dir"

if [ -z "$res" ]; then 
  echo -e "\e[31mNo paste found ..."
  exit 0
fi

while IFS= read -r line; do
  [[ ${#line} == 0 ]] && continue
  link="https://pastebin.com/raw/$(cut -d "@" -f1 <<<"$line")"
  tags=$(cut -d "@" -f2 <<<"$line")
  date=$(cut -d "@" -f3 <<<"$line" | cut -d " " -f1)
  sc=$(curl -sI "$link" | head -n 1|cut -d$' ' -f2)
  [[ $sc != 200 ]] && echo -e "\e[32m[\e[31m$link\e[32m]-→\e[32m[\e[33m$date\e[32m]" && continue
  [[ $sc == 200 ]] && echo -ne "\e[32m[\e[36m$link\e[32m]-→\e[32m[\e[33m$date\e[32m]"
  [[ $tags != "none" ]] && echo -e "-→\e[32m[\e[35m$tags\e[32m]" || echo ""
  eval $ff "$link" 2> /dev/null
  [[ $save == true ]] && curl -sk "$link" > "$save_dir/$prefix-$(cut -d "@" -f1 <<<"$line").txt"
done < <(sort -r -t"@" -k3 <<<"$res" | uniq)

