package main

import (
	"DC-homework-1/uploader/mqClient"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"lib/authClient"
	"lib/dbClient"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type Config struct {
	BatchSize int
}

var config Config
var mq mqClient.MQClient
var ac authClient.AuthClient

type errorResponse struct {
	Error string `json:"error"`
}

func responseError(w http.ResponseWriter, err error, statusCode int) {
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: fmt.Sprintln(err)})
}

func (c *Config) Init() (err error) {
	var size int64
	if size, err = strconv.ParseInt(os.Getenv("BATCH_SIZE"), 10, 32); err != nil {
		return
	}
	c.BatchSize = int(size)
	return
}

func authorize(w http.ResponseWriter, r *http.Request, adminRequired bool) bool {
	w.Header().Set("Content-Type", "application/json")
	token := r.Header.Get("auth")
	if err := ac.EnsurePermission(token, adminRequired); err != nil {
		if errResp, ok := err.(*authClient.ErrorRespStatus); ok {
			responseError(w, errResp, errResp.StatusCode)
		} else {
			responseError(w, err, http.StatusInternalServerError)
		}
		return false
	}
	return true
}

func UploadHandler(w http.ResponseWriter, r *http.Request) {
	if !authorize(w, r, true) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var reader *multipart.Reader
	var err error
	if reader, err = r.MultipartReader(); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	var readerPart *multipart.Part
	if readerPart, err = reader.NextPart(); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	if readerPart.FormName() != "file" {
		responseError(w, errors.New(fmt.Sprintf("File is expected, %s is provided", readerPart.FormName())), http.StatusBadRequest)
		return
	}
	csvReader := csv.NewReader(readerPart)
	var header []string
	if header, err = csvReader.Read(); err != nil {
		responseError(w, err, http.StatusBadRequest)
		return
	}
	columnToId := make(map[string]int)
	for id, column := range header {
		columnToId[column] = id
	}
	requiredFields := []string{"title", "category"}
	for _, field := range requiredFields {
		if _, ok := columnToId[field]; !ok {
			responseError(w, errors.New(fmt.Sprintf("Missing required field: %s", field)), http.StatusBadRequest)
			return
		}
	}
	var batch []dbClient.Product
	for {
		var values []string
		if values, err = csvReader.Read(); err != nil {
			if err == io.EOF {
				break
			} else {
				responseError(w, err, http.StatusBadRequest)
				return
			}
		}
		product := dbClient.Product{
			Title:    values[columnToId["title"]],
			UniqueID: "",
			Category: values[columnToId["category"]],
		}
		if _, ok := columnToId["uniqueId"]; ok {
			product.UniqueID = values[columnToId["uniqueId"]]
		}
		batch = append(batch, product)
		if len(batch) == config.BatchSize {
			if err = mq.SendBatch(batch); err != nil {
				log.Fatal(err)
			}
			batch = batch[0:]
		}
	}
	if err = mq.SendBatch(batch); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var err error
	if err = config.Init(); err != nil {
		log.Fatal(err)
	}
	if ac, err = authClient.InitAuthClient(); err != nil {
		log.Fatal(err)
	}
	if mq, err = mqClient.InitMQClient(); err != nil {
		log.Fatal(err)
	}
	r := mux.NewRouter()
	r.HandleFunc("/upload", UploadHandler).Methods("POST")
	log.Println("Upload server started")
	log.Fatal(http.ListenAndServe(":8000", r))
}
