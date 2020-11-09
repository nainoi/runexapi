package upload

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"

	//"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/nfnt/resize"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/model"
)

// Uploads image by create event
func Uploads(c *gin.Context) {
	var res = response.Gin{C: c}

	// Maximum upload of 10 MB files
	c.Request.ParseMultipartForm(10 << 20)

	// Get handler for filename, size and headers
	file, handler, err := c.Request.FormFile("upload")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	// fmt.Printf("File Size: %+v\n", handler.Size)
	// fmt.Printf("MIME Header: %+v\n", handler.Header)

	uniqidFilename := uuid.New()
	name := fmt.Sprintf(uniqidFilename.String() + ".png")

	byteFile, err := ioutil.ReadAll(file)
	if err != nil {
		log.Println("error read")
		fmt.Println(err.Error())
	}

	data := uploadToS3(&byteFile, name)

	res.Response(http.StatusOK, "success", data)
}

// UploadCover cover upload to s3
func UploadCover(c *gin.Context) {
	file, header, err := c.Request.FormFile("upload")
	file2, _, _ := c.Request.FormFile("upload")
	filename := header.Filename
	fmt.Println(filename)

	var res = response.Gin{C: c}

	uniqidFilename := uuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", uniqidFilename.String())

	pathDir := "." + config.UPLOAD_IMAGE
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		os.MkdirAll(pathDir, os.ModePerm)
	}

	// decode jpeg into image.Image
	//img, err := jpeg.Decode(file2)
	img, str, err := image.Decode(file2)
	log.Println(str)
	if err != nil {
		log.Fatal(err)
	}

	path := pathDir + uniqidFilename.String() + "." + str
	out, err := os.Create(path)

	if err != nil {
		log.Println(err)
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	_, err = io.Copy(out, file)

	//thumbnail
	pathDirThumbnail := "." + config.UPLOAD_IMAGE + "thumbnail/"
	if _, err := os.Stat(pathDirThumbnail); os.IsNotExist(err) {
		os.MkdirAll(pathDirThumbnail, os.ModePerm)
	}

	defer file.Close()

	m := resize.Resize(540, 0, img, resize.Lanczos3)

	thumbPath := pathDirThumbnail + uniqidFilename.String() + "." + str

	outThumbnail, err := os.Create(thumbPath)
	if err != nil {
		log.Println(err)
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
	}

	// write new image to file
	jpeg.Encode(outThumbnail, m, nil)
	//_, err = io.Copy(outThumbnail, file)

	if err != nil {
		res.Response(http.StatusInternalServerError, err.Error(), gin.H{"message": err.Error()})
		return
		//log.Fatal(err)
	}

	thumbObject := UploadWithFolderToS3(thumbPath, "event", uniqidFilename.String()+"."+str)
	coverObject := UploadWithFolderToS3(path, "thumbnail", uniqidFilename.String()+"."+str)

	defer outThumbnail.Close()
	defer out.Close()
	defer os.Remove(path)
	defer os.Remove(thumbPath)

	resCover := model.CoverUploadResponse{
		ThumbURL: thumbObject.URL,
		CoverURL: coverObject.URL,
		MSG:      "success",
	}
	res.Response(http.StatusOK, "success", resCover)
}

func uploadToS3(data *[]byte, name string) model.UploadResponse {
	url := "https://storage.runex.co/upload"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)

	//file, errFile1 := os.Open(fmt.Sprintf("./%s", name))
	part1, errFile1 := writer.CreateFormFile("file", fmt.Sprintf("%s", name))
	//_, errFile1 = io.Copy(part1, file)

	// file, errFile1 := os.Open("/C:/Users/frogconn/Downloads/5f7324f2da1d9600135ed041vid1601548947586.mp4")
	//defer file.Close()
	//part1, errFile1 := writer.CreateFormField("file")
	var resObject = model.UploadResponse{}
	_, errFile1 = part1.Write(*data)
	if errFile1 != nil {
		fmt.Println(errFile1)
		log.Println("error create field")
		return resObject
	}
	err := writer.Close()
	if err != nil {
		log.Println("error close")
		fmt.Println(err)
		return resObject
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println("error req")
		fmt.Println(err)
		return resObject
	}
	req.Header.Add("token", "5Dk2o03a4hVjQPglSueFEah577fCGQfM")
	req.Header.Add("path", "runex/photo/")
	req.Header.Add("Cookie", "__cfduid=dd42cd8b41a9c49d5b75f756dc64e01451604633984")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println("error do req")
		fmt.Println(err)
		return resObject
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error read upload")
		fmt.Println(err)
		return resObject
	}

	if res.StatusCode == 200 {
		err = json.Unmarshal(body, &resObject)
		if err != nil {
			log.Println(err)
		}
	}
	return resObject
}

// UploadWithFolderToS3 upload to S3 with path and name
func UploadWithFolderToS3(path string, folder string, name string) model.UploadResponse {
	url := "https://storage.runex.co/upload"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	var resObject = model.UploadResponse{}
	file, errFile1 := os.Open(path)
	if errFile1 != nil {
		fmt.Println(errFile1)
		log.Println("error read file")
		return resObject
	}
	part1, errFile1 := writer.CreateFormFile("file", fmt.Sprintf("%s", name))
	_, errFile1 = io.Copy(part1, file)

	// file, errFile1 := os.Open("/C:/Users/frogconn/Downloads/5f7324f2da1d9600135ed041vid1601548947586.mp4")
	//defer file.Close()
	//part1, errFile1 := writer.CreateFormField("file")

	if errFile1 != nil {
		fmt.Println(errFile1)
		log.Println("error create field")
		return resObject
	}
	err := writer.Close()
	if err != nil {
		log.Println("error close")
		fmt.Println(err)
		return resObject
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, url, payload)

	if err != nil {
		log.Println("error req")
		fmt.Println(err)
		return resObject
	}
	req.Header.Add("token", "5Dk2o03a4hVjQPglSueFEah577fCGQfM")
	req.Header.Add("path", fmt.Sprintf("runex/photo/%s/", folder))
	req.Header.Add("Cookie", "__cfduid=dd42cd8b41a9c49d5b75f756dc64e01451604633984")

	req.Header.Set("Content-Type", writer.FormDataContentType())
	res, err := client.Do(req)
	if err != nil {
		log.Println("error do req")
		fmt.Println(err)
		return resObject
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println("error read upload")
		fmt.Println(err)
		return resObject
	}

	if res.StatusCode == 200 {
		err = json.Unmarshal(body, &resObject)
		if err != nil {
			log.Println(err)
		}
	}
	return resObject
}
