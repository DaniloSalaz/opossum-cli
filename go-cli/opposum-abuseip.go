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

var comandos []cli
var urlAbuseIP = "https://api.abuseipdb.com/api/v2"

type cli struct {
	subcommand  *flag.FlagSet
	name        string
	description string
}

func setComandos(c cli) {
	comandos = append(comandos, c)
}
func showHelpByCommand(cli cli) {
	fmt.Println(cli.name)
	cli.subcommand.PrintDefaults()
	fmt.Println()
}
func showHelpByNameCommand(name string) bool {
	for _, c := range comandos {
		if c.name == name {
			fmt.Println(c.name)
			c.subcommand.PrintDefaults()
			return true
		}
	}
	return false
}
func showHelp() {
	fmt.Println("Comandos:")
	for _, c := range comandos {
		fmt.Printf("\t %s \t %s \n", c.name, c.description)
	}
}

func comprobarComand(arg []string) {
	if len(arg) < 1 {
		showHelp()
		os.Exit(1)
	} else if len(arg) == 0 {
		if !showHelpByNameCommand(arg[0]) {
			showHelp()
			os.Exit(1)
		}
	} else if arg[0] != "check" && arg[0] != "blacklist" {
		showHelp()
		os.Exit(1)
	}

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
		showHelpByNameCommand("check")
		os.Exit(1)
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
		showHelpByNameCommand("blacklist")
		os.Exit(1)
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

	c := flag.NewFlagSet("check", flag.ContinueOnError)
	c.StringVar(&key, "key", "", "ApiKey (Required)")
	c.StringVar(&ip, "ip", "", "IP a verificar. (Required)")
	c.IntVar(&days, "maxDays", 30, "Determina el tiempo (en días) máximo del reporte. Valor por defecto 30 días, mínimo (1) máximo (365)")
	c.BoolVar(&verbose, "verbose", false, "Verbose, retorna la toda la infomación")
	c.BoolVar(&help, "help", false, "Muestra ayuda")

	b := flag.NewFlagSet("blacklist", flag.ContinueOnError)
	b.StringVar(&key, "key", "", "ApiKey (Required)")
	b.IntVar(&min, "min", 100, "Confidence minimum. Valor por defecto 100, mínumo(25) máximo (100)")
	b.IntVar(&limit, "limit", 100000, "Límite de IPs. Valor por defecto 10.000")
	b.BoolVar(&plaintext, "plaintext", false, "Response en texto plano")
	b.BoolVar(&desc, "desc", false, "Ordenar descendente")
	b.BoolVar(&last, "last", false, "Ordenar descendente las IPs mas recientes")
	b.BoolVar(&help, "help", false, "Muestra ayuda")

	command := cli{subcommand: c, name: "check", description: "Retorna los detalles de una IP address"}
	command2 := cli{subcommand: b, name: "blacklist", description: "Retorna la lista negra de todas las IPs reportadas"}

	setComandos(command)
	setComandos(command2)

	comprobarComand(os.Args[1:])

	if os.Args[1] == "check" {
		c.Parse(os.Args[2:])
		checkIP(key, ip, days, verbose, help)
	}
	if os.Args[1] == "blacklist" {
		b.Parse(os.Args[2:])
		blacklist(key, min, limit, plaintext, desc, last, help)
	}
}
