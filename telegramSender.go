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
	Bot_token      string `json:"bot_token"`
	Chat_id        string `json:"chat_id"`
	Period_minutes int64  `json: "period_minutes"`
}

func main() {
	//json decode
	jsonConfigFile, _ := os.Open("config.json")
	decoder := json.NewDecoder(jsonConfigFile)
	config := new(ConfigJson)
	err := decoder.Decode(&config)
	errorCheck(err)
	fmt.Println(config)
	postUriForJpg := "https://api.telegram.org/bot" + config.Bot_token + "/sendPhoto"
	fmt.Println(postUriForJpg)

	tick := time.NewTicker(time.Minute * time.Duration(config.Period_minutes))
	for ; true; <-tick.C {

		//get list of files
		fileList, err := ioutil.ReadDir("./files")
		errorCheck(err)
		//shuffle fileList slice
		for _, f := range fileList {
			rand.Shuffle(len(fileList), func(i, j int) {
				fileList[i], fileList[j] = fileList[j], fileList[i]
			})

			fileInDirectory := f.Name()

			extension := filepath.Ext(fileInDirectory)
			if extension == ".jpg" || extension == ".png" {
				//open file
				openedFile, err := os.Open(fileInDirectory)
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
				_, err = field_part.Write([]byte(config.Chat_id))
				errorCheck(err)

				writer.Close()

				req, err := http.NewRequest("POST", postUriForJpg, body)

				//add header
				req.Header.Set("Content-Type", writer.FormDataContentType())
				client := &http.Client{}
				resp, err := client.Do(req)
				errorCheck(err)

				err = openedFile.Close()
				errorCheck(err)

				err = os.Remove(fileInDirectory)
				errorCheck(err)

				log.Println(resp)

			}
			if extension == ".mp4" || extension == ".gif" {
				postUriForJpg := "https://api.telegram.org/bot" + config.Bot_token + "/sendVideo"
				//open file
				openedFile, err := os.Open(fileInDirectory)
				errorCheck(err)

				//create buffer
				body := new(bytes.Buffer)

				//create writer from body
				writer := multipart.NewWriter(body)

				//initializing field for file
				file_part, err := writer.CreateFormFile("video", fileInDirectory)
				errorCheck(err)

				//copy openedFile to file_part?
				_, err = io.Copy(file_part, openedFile)
				errorCheck(err)

				//initializing field for another parameter
				field_part, err := writer.CreateFormField("chat_id")

				//write value for parameter
				_, err = field_part.Write([]byte(config.Chat_id))
				errorCheck(err)

				writer.Close()

				req, err := http.NewRequest("POST", postUriForJpg, body)

				//add header
				req.Header.Set("Content-Type", writer.FormDataContentType())
				client := &http.Client{}
				resp, err := client.Do(req)
				errorCheck(err)

				err = openedFile.Close()
				errorCheck(err)

				err = os.Remove(fileInDirectory)
				errorCheck(err)

				log.Println(resp)

			}
		}
	}
}
