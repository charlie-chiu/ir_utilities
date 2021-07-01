package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func main() {
	log.SetFlags(log.Lshortfile)

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

	if resp.StatusCode != http.StatusOK {
		log.Fatal("response status error: expected 200, got ", resp.StatusCode)
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

	var ipList []string
	for i := 0; i < len(list.Data); i++ {
		ipList = append(ipList, list.Data[i].IP)
	}

	log.Print(strings.Join(ipList, ","))
}
