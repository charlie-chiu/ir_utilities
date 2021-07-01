package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"

	"github.com/joho/godotenv"
)

//todo: handle indent

func main() {
	log.SetFlags(log.Lshortfile)

	hostIPList := fetchHostIPList()
	//hostIPList := []string{"103.145.22.84","182.255.60.85"}

	panels := generatePanels(hostIPList)
	//log.Println(panels)

	jsonModel := generateJSONModel(panels)
	//log.Println(jsonModel)

	file, err := os.Create("JSONModel.json")
	if err != nil {
		log.Fatal("create JSON Model file error, ", err)
	}
	_, err = file.WriteString(jsonModel)
	if err != nil {
		log.Fatal("write to JSON Model file error, ", err)
	}

	//log.Print(strings.Join(hostIPList, ","))
}

func generateJSONModel(panels string) string {
	tpl, err := template.ParseFiles("./JSONModel.tmpl")
	if err != nil {
		log.Fatal("ParseFile error, ", err)
	}

	buf := new(bytes.Buffer)
	err = tpl.Execute(buf, struct {
		Panels string
	}{
		Panels: panels,
	})

	if err != nil {
		log.Fatal("JSON Model template execute error, ", err)
	}

	return buf.String()
}

func generatePanels(hostIPList []string) string {
	tpl, err := template.ParseFiles("./panels.tmpl")
	if err != nil {
		log.Fatal("ParseFile error, ", err)
	}

	sb := strings.Builder{}
	buf := new(bytes.Buffer)
	for i, IP := range hostIPList {
		//log.Println(i+1, IP)
		err := tpl.Execute(buf, struct {
			ID   int
			Host string
		}{
			ID:   i + 1,
			Host: IP,
		})

		if err != nil {
			log.Fatal("alert panels template execute error, ", err)
		}
		sb.Write(buf.Bytes())
		if i < len(hostIPList)-1 {
			sb.WriteString(",\n")
		}
		buf.Reset()
	}

	panels := sb.String()
	return panels
}

func fetchHostIPList() []string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading .env file", err)
	}

	url := os.Getenv("LEA_API") + "/api/v1/ops/proxies"
	token := os.Getenv("API_SIMPLE_TOKEN")

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("client.Do error, ", err)
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("read response body error, ", err)
	}

	//log.Printf(string(bytes))

	type opsFormat struct {
		Data []struct {
			IP string `json:"{#HQS_OL}"`
		} `json:"data"`
	}

	list := &opsFormat{}

	err = json.Unmarshal(bytes, list)
	if err != nil {
		log.Fatal(err)
	}

	var hostIPList []string
	for i := 0; i < len(list.Data); i++ {
		hostIPList = append(hostIPList, list.Data[i].IP)
	}
	return hostIPList
}
