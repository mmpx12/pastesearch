package main

import (
	"encoding/json"
	"fmt"
	"github.com/itchyny/gojq"
	"github.com/speedata/optionparser"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"sync"
	"time"
)

var search, searchv2, save, prefix, mail, domain string
var browser, slow bool
var wg sync.WaitGroup

func main() {
	ArgsParse()
	paste := Parse(Search())
	for _, i := range paste {
		if !slow {
			wg.Add(1)
			go GetPaste(i)
		} else {
			GetPaste(i)
			time.Sleep(300 * time.Millisecond)
		}
	}
	if !slow {
		wg.Wait()
	}
}

func Search() []byte {
	var res *http.Response
	switch {
	case search != "":
		res, _ = http.Get("https://psbdmp.ws/api/search/" + search)
	case searchv2 != "":
		res, _ = http.Get("https://psbdmp.ws/api/v2/search/" + searchv2)
	case mail != "":
		res, _ = http.Get("https://psbdmp.ws/api/search/email/" + mail)
	case domain != "":
		res, _ = http.Get("https://psbdmp.ws/api/search/domain/" + domain)
	default:
		os.Exit(1)
	}
	defer res.Body.Close()
	result, _ := ioutil.ReadAll(res.Body)
	return result
}

func Parse(raw []byte) []string {
	var input interface{}
	result := make([]string, 0)
	json.Unmarshal(raw, &input)
	parsed, _ := gojq.Parse(`.data[] | [.id +"@"+ if .tags == "" then "none" else .tags end + "@" + .time] |@tsv`)
	val := parsed.Run(input)
	for {
		v, ok := val.Next()
		if !ok {
			break
		}
		res := fmt.Sprintf("%v\n", v)
		result = append(result, res)
	}
	if len(result) == 0 {
		fmt.Println("No result")
		os.Exit(0)
	} else if len(result) > 20 {
		browser = false
	}
	return result
}

func SavePaste(paste_id string, content *http.Response) {
	f, _ := os.Create(save + "/" + prefix + paste_id + ".txt")
	defer f.Close()
	io.Copy(f, content.Body)
}

func OpenBrowser(url string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", url).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		exec.Command("open", url).Start()
	case "android":
		exec.Command("termux-open-url", url).Start()
	default:
		fmt.Println("Error: Can't open browser")
		browser = false
	}
}

func GetPaste(paste string) {
	if !slow {
		defer wg.Done()
	}
	detail := strings.Split(strings.TrimSuffix(paste, "\n"), "@")
	content, _ := http.Get("https://pastebin.com/raw/" + detail[0])
	defer content.Body.Close()
	if content.StatusCode <= 300 {
		result := "\033[32m[\033[36mhttps://pastebin.com/raw/" + detail[0] + "\033[32m]─→\033[32m[\033[33m" + detail[2] + "\033[32m]"
		if detail[1] != "none" {
			result += "─→\033[32m[\033[35m" + detail[1] + "\033[32m]"
		}
		fmt.Println(result)
		if save != "" {
			SavePaste(detail[0], content)
		}
		if browser {
			OpenBrowser("https://pastebin.com/raw/" + detail[0])
		}
	} else if content.StatusCode == 403 {
		result := "\033[32m[\033[35mhttps://pastebin.com/raw/" + detail[0] + "\033[32m]─→\033[32m[\033[33m" + detail[2] + "\033[32m]"
		if detail[1] != "none" {
			result += "─→\033[32m[\033[35m" + detail[1] + "\033[32m]"
		}
		result += " \033[31mCAPTCHA :("
		fmt.Println(result)
	} else {
		fmt.Println("\033[32m[\033[31mhttps://pastebin.com/raw/" + detail[0] + "\033[32m]─→\033[32m[\033[33m" + detail[2] + "\033[32m]")
	}
}

func ArgsParse() {
	switch len(os.Args) {
	case 1:
		help()
	case 2:
		search = os.Args[1]
		return
	}
	op := optionparser.NewOptionParser()
	op.On("-s", "--search search", "", &search)
	op.On("-S", "--searchv2 searchv2", "", &searchv2)
	op.On("-m", "--mail mail", "", &mail)
	op.On("-d", "--domain domain", "", &domain)
	op.On("-b", "--browser", "", &browser)
	op.On("-o", "--save save", "", &save)
	op.On("-p", "--prefix prefix", "", &prefix)
	op.On("-x", "--slow", "", &slow)
	op.On("-h", "--help", "", help)
	op.Parse()
	if save != "" {
		if _, err := os.Stat(save); os.IsNotExist(err) {
			os.Mkdir(save, 0777)
		}
	}

}

func help() {
	fmt.Println(`pastebin searcher

usage:
   -s, --search [QUERY]        General search on pastebin
   -S, --searchv2 [QUERY]      Same as -s but with api v2
   -m, --mail [EMAIL]          Search for emails 
   -d, --domain [DOMAIN]       Search specific domain
   -o, --save [DIRECTORY ]     Save paste into directory
   -p, --prefix [PREFIX]       Prefix when save paste
   -b, --browser               Open paste in browser (if result < 20)
   -x, --slow                  Avoid triggering captcha (lot slower)

output:`)
	fmt.Println("  \033[32m[\033[31m$LINK\033[32m]-→\033[32m[\033[33m$DATE\033[32m] : Paste was removed")
	fmt.Println("  \033[32m[\033[36m$LINK\033[32m]-→\033[32m[\033[33m$DATE\033[32m] : Paste exists")
	fmt.Println("  \033[32m[\033[35m$LINK\033[32m]-→\033[32m[\033[33m$DATE\033[32m] : blocked by Captcha")
	fmt.Println("Some paste also have tags:")
	fmt.Println("  -→\033[32m[\033[35m$TAGS\033[32m]")
	os.Exit(1)
}
