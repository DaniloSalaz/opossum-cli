package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

var comandos []moduleCommand

var urlAbuseIP = "https://api.abuseipdb.com/api/v2"
var msgCommadn = "opossum"

type item struct {
	command     *flag.FlagSet
	name        string
	description string
}
type moduleCommand struct {
	command     *flag.FlagSet
	name        string
	description string
	items       []item
}

func (m *moduleCommand) addSubCLIs(item item) {
	m.items = append(m.items, item)
}

func (m *moduleCommand) isExistItem(nameItem string) bool {
	for _, i := range m.items {
		if i.name == nameItem {
			return true
		}
	}
	return false
}

func setComandos(c moduleCommand) {
	comandos = append(comandos, c)
}
func showHelpByNameItemCommand(nameCommand string, nameItem string) bool {
	for _, c := range comandos {
		if c.name == nameCommand {
			for _, i := range c.items {
				if i.name == nameItem {
					fmt.Printf("\n usage: %s %s  %s [options]\n\n opntions: \n", msgCommadn, c.name, i.name)
					i.command.PrintDefaults()
				}
			}
			return true
		}
	}
	return false
}
func showHelpByNameCommand(nameCommand string) bool {
	fmt.Printf("usage: %s %s <item> \n\n items: \n", msgCommadn, nameCommand)

	for _, c := range comandos {
		if c.name == nameCommand {
			for _, i := range c.items {
				fmt.Printf("\t %s \t %s \n", i.name, i.description)
			}
			return true
		}
	}
	return false
}
func showHelp() {
	fmt.Printf("usage: %s [command] \n\n commands: \n ", msgCommadn)
	for _, c := range comandos {
		fmt.Printf("\t %s \t %s \n", c.name, c.description)
	}
	fmt.Println()
}

func comprobarComand(arg []string) bool {
	var existFirtsCmd = false
	for _, c := range comandos {
		if len(arg) > 0 && c.name == arg[0] {
			if len(arg) > 1 && c.isExistItem(arg[1]) {
				return true
			}
			existFirtsCmd = true
		}
	}
	if existFirtsCmd {
		showHelpByNameCommand(arg[0])
	} else {
		showHelp()
	}
	return false
}

func printResponseJSON(body []byte) {
	var prettyJSON bytes.Buffer
	error := json.Indent(&prettyJSON, body, "", "\t")
	if error != nil {
		fmt.Println("JSON parse error: ", error)
		return
	}
	fmt.Println(string(prettyJSON.Bytes()))
}
func lauchHTTP(action string, apikey string, params map[string]string) {
	client := &http.Client{}
	urlAbuseIP = urlAbuseIP + action

	req, err := http.NewRequest("GET", urlAbuseIP, nil)
	req.Header.Add("key", apikey)
	req.Header.Add("Accept", "application/json")

	query := req.URL.Query()

	for key, value := range params {
		query.Add(key, value)
	}

	req.URL.RawQuery = query.Encode()
	resp, errHTTP := client.Do(req)

	if errHTTP != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()
	body, errBody := ioutil.ReadAll(resp.Body)

	if errBody != nil {
		fmt.Println(err)
		return
	}
	os.Stdout.Write(body)

}

func checkIP(key string, ip string, days int, verbose bool, help bool) {
	if help || key == "" || ip == "" {
		fmt.Println(ip, key)
		// showHelpByNameItemCommand("abuseip", "check")
		return
	}
	hashmap := make(map[string]string)
	hashmap["ipAddress"] = ip
	hashmap["maxAgeInDays"] = strconv.Itoa(days)
	if verbose {
		hashmap["verbose"] = "true"
	}

	lauchHTTP("/check", key, hashmap)

}
func blacklist(key string, min int, limit int, plaintext bool, desc bool, last bool, help bool) {
	if help || key == "" {
		showHelpByNameItemCommand("abuseip", "blacklist")
		return
	}
	hashmap := make(map[string]string)
	hashmap["confidenceMinimum"] = strconv.Itoa(min)
	hashmap["limit"] = strconv.Itoa(min)
	if plaintext {
		hashmap["plaintext"] = "true"
	}
	if desc {
		hashmap["abuseConfidenceScore"] = "true"
	}
	if last {
		hashmap["lastReportedAt"] = "true"
	}

	lauchHTTP("/blacklist", key, hashmap)
}
func main() {
	var key string
	var ip string
	var days int
	var verbose bool
	var help bool

	var min int
	var limit int
	var plaintext bool
	var desc bool
	var last bool

	a := flag.NewFlagSet("abuseip", flag.ContinueOnError)
	e := flag.NewFlagSet("abc", flag.ContinueOnError)

	c := flag.NewFlagSet("check", flag.ContinueOnError)
	c.StringVar(&key, "key", "", "ApiKey (Required)")
	c.StringVar(&ip, "ip", "", "IP to verify. (Required)")
	c.IntVar(&days, "maxDays", 30, "Determines how far back in time(days). The default is 30, min (1) max (365)")
	c.BoolVar(&verbose, "verbose", false, "Verbose")
	c.BoolVar(&help, "help", false, "Show help")

	b := flag.NewFlagSet("blacklist", flag.ContinueOnError)
	b.StringVar(&key, "key", "", "ApiKey (Required)")
	b.IntVar(&min, "min", 100, "Confidence minimum. The default is 30, min (25) max (100)")
	b.IntVar(&limit, "limit", 100000, "Limit return IPs. The default 10.000")
	b.BoolVar(&plaintext, "plaintext", false, "Response plain text")
	b.BoolVar(&desc, "desc", false, "Descending order")
	b.BoolVar(&last, "last", false, "Descending order of the last IPs")
	b.BoolVar(&help, "help", false, "Show help")

	moduloAbuse := moduleCommand{command: a, name: "abuseip", description: "Requests to https://www.abuseipdb.com/"}
	moduloAbc := moduleCommand{command: e, name: "abc", description: "Requests to https://www.abuseipdb.com/"}

	command := item{command: c, name: "check", description: "IP Address report"}
	command2 := item{command: b, name: "blacklist", description: "Blacklist of reported IPs"}

	moduloAbuse.addSubCLIs(command)
	moduloAbuse.addSubCLIs(command2)

	setComandos(moduloAbuse)
	setComandos(moduloAbc)

	c.Parse(os.Args[2:])
	fmt.Println(key)

	if comprobarComand(os.Args[1:]) {
		switch os.Args[2] {
		case "check":
			c.Parse(os.Args[2:])
			fmt.Println(key)
			checkIP(key, ip, days, verbose, help)
		case "blacklist":
			b.Parse(os.Args[2:])
			blacklist(key, min, limit, plaintext, desc, last, help)
		}

	}
}
