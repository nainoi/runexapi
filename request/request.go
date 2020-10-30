package request

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"time"
)

//GetRequest client
func GetRequest(urlS string, bearer string, params map[string]interface{}) ([]byte, error){
	req, err := http.NewRequest("GET", urlS, nil)
	req.Header.Add("Authorization", bearer)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded, charset=UTF-8")

	timeout := time.Duration(6 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	//client.CheckRedirect = checkRedirectFunc

	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 200 || resp.StatusCode < 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		return body, err
	}
	return []byte{}, err
}

//PostRequest client
func PostRequest(urlS string, bearer string, params map[string]interface{}) ([]byte, error){
	bytesRepresentation, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", urlS, bytes.NewBuffer(bytesRepresentation))
	req.Header.Add("Authorization", bearer)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded, charset=UTF-8")

	timeout := time.Duration(6 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	//client.CheckRedirect = checkRedirectFunc

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	// if err != nil {
	// 	log.Println(err)
	// }

	// var result map[string]interface{}

	// json.NewDecoder(resp.Body).Decode(&result)

	// log.Println(result)
	// log.Println(result["data"])

	defer resp.Body.Close()

	if resp.StatusCode >= 200 || resp.StatusCode < 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		return body, err
	}

	return []byte{}, err
}

//MultipartRequest client
func MultipartRequest(urlS string, bearer string, file multipart.File, params map[string]interface{}) ([]byte, error){
	// Open the file
	//file, err := os.Open("name.txt")
	// if err != nil {
	// 	log.Println(err)
	// }
	// Close the file later
	defer file.Close()

	// Buffer to store our request body as bytes
	var requestBody bytes.Buffer

	// Create a multipart writer
	multiPartWriter := multipart.NewWriter(&requestBody)

	// Initialize the file field
	fileWriter, err := multiPartWriter.CreateFormFile("file_field", "name.txt")
	if err != nil {
		log.Println(err)
	}

	// Copy the actual file content to the field field's writer
	_, err = io.Copy(fileWriter, file)
	if err != nil {
		log.Println(err)
	}

	// Populate other fields
	fieldWriter, err := multiPartWriter.CreateFormField("normal_field")
	if err != nil {
		log.Println(err)
	}

	_, err = fieldWriter.Write([]byte("Value"))
	if err != nil {
		log.Println(err)
	}

	// We completed adding the file and the fields, let's close the multipart writer
	// So it writes the ending boundary
	multiPartWriter.Close()
	bytesRepresentation, err := json.Marshal(params)
	if err != nil {
		log.Println(err)
	}

	req, err := http.NewRequest("POST", urlS, bytes.NewBuffer(bytesRepresentation))
	req.Header.Add("Authorization", bearer)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded, charset=UTF-8")

	timeout := time.Duration(6 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	//client.CheckRedirect = checkRedirectFunc

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}

	// resp, err := http.Post(url, "application/json", bytes.NewBuffer(bytesRepresentation))
	// if err != nil {
	// 	log.Println(err)
	// }

	// var result map[string]interface{}

	// json.NewDecoder(resp.Body).Decode(&result)

	// log.Println(result)
	// log.Println(result["data"])

	defer resp.Body.Close()

	if resp.StatusCode >= 200 || resp.StatusCode < 300 {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
		}
		return body, err
	}

	return []byte{}, err
}