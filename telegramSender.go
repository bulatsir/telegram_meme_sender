package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func errorCheck(err error) {
	if err != nil {
		log.Println(err)
	}
}

type ConfigJson struct {
	bot_token string
	chat_id   string
}

func main() {

	jsonConfigFile, _ := os.Open("config.json")
	decoder := json.NewDecoder(jsonConfigFile)
	config := new(ConfigJson)
	err := decoder.Decode(&config)
	errorCheck(err)
	fmt.Println(config.bot_token)
	postUriForJpg := "https://api.telegram.org/" + config.bot_token + "/sendPhoto"
	fmt.Println(postUriForJpg)

	tick := time.NewTicker(time.Second * 15)
	for _ = range tick.C {
		//get list of files
		fileList, err := ioutil.ReadDir("./")
		errorCheck(err)
		//shuffle fileList slice
		for _, f := range fileList {
			rand.Shuffle(len(fileList), func(i, j int) {
				fileList[i], fileList[j] = fileList[j], fileList[i]
			})

			fileInDirectory := f.Name()

			extension := filepath.Ext(fileInDirectory)
			if extension == ".jpg" {
				//open file

				openedFile, err := os.Open(fileInDirectory)
				errorCheck(err)

				err = openedFile.Close()
				errorCheck(err)
				//create buffer
				body := new(bytes.Buffer)

				//create writer from body
				writer := multipart.NewWriter(body)

				//initializing field for file
				file_part, err := writer.CreateFormFile("photo", fileInDirectory)
				errorCheck(err)

				//copy openedFile to file_part?
				_, err = io.Copy(file_part, openedFile)
				errorCheck(err)

				//initializing field for another parameter
				field_part, err := writer.CreateFormField("chat_id")

				//write value for parameter
				_, err = field_part.Write([]byte(config.chat_id))
				errorCheck(err)

				writer.Close()

				req, err := http.NewRequest("POST", postUriForJpg, body)

				//add header
				req.Header.Set("Content-Type", writer.FormDataContentType())
				client := &http.Client{}
				resp, err := client.Do(req)
				errorCheck(err)

				log.Println(resp)

			}

		}
	}
}
