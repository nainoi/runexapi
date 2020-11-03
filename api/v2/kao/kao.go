package kao

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"thinkdev.app/think/runex/runexapi/api/v2/response"
	"thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/model"
	"thinkdev.app/think/runex/runexapi/repository"
)

//TK is token key koa
var (
	TK = "e35Sa9MvZJ1fA0PV"
)

// KoaAPI struct repo
type KoaAPI struct {
	KaoRepository repository.KaoRepository
}

// GetKaoActivity godoc
// @Summary get Kao event detail
// @Description GetKaoActivity get Kao event detail open id API calls
// @Consume application/x-www-form-urlencoded
// @Tags kao
// @Accept  application/json
// @Produce application/json
// @Param payload body model.GetKaoActivityRequest true "payload"
// @Success 200 {object} response.ResponseOAuth
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /kao [post]
func (api KoaAPI) GetKaoActivity(c *gin.Context) {
	var (
		res = response.Gin{C: c}
	)
	var form model.GetKaoActivityRequest
	if err := c.ShouldBind(&form); err != nil {
		res.Response(http.StatusBadRequest, err.Error(), nil)
		return
	}
	//userID := "5d772660c8a56133c2d7c5ba"
	//userID, _ := oauth.GetValuesToken(c)
	urlS := fmt.Sprintf("https://kaokonlakao-www-tabshier.azurewebsites.net/api/%s/bib/%s", form.Slug, form.EBIB)
	var bearer = "Bearer olcgZVpqDXQikRDG"
	//reqURL, _ := url.Parse(urlS)
	req, err := http.NewRequest("GET", urlS, nil)
	req.Header.Add("Authorization", bearer)
	//req.Header.Add("Content-Type", "application/x-www-form-urlencoded, charset=UTF-8")

	timeout := time.Duration(6 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	client.CheckRedirect = checkRedirectFunc

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
		var koaObject model.KaoObject
		err = json.Unmarshal(body, &koaObject)
		if err != nil {
			log.Println(err)
		}
		res.Response(http.StatusOK, "success", koaObject)
		c.Abort()
		return
	}

	res.Response(http.StatusInternalServerError, err.Error(), nil)
}

//SendKaoActivity api
func (api KoaAPI) SendKaoActivity(c *gin.Context) {
	url := fmt.Sprintf("https://kaokonlakao-www-tabshier.azurewebsites.net/api/%s/bib/%s/submit")
	var (
		res = response.Gin{C: c}
	)

	file, _, err := c.Request.FormFile("image")
	if err != nil {
		res.Response(http.StatusBadRequest, "Image is require", nil)
		c.Abort()
		return
	}

	time := c.Request.FormValue("time")
	if time == "" {
		res.Response(http.StatusBadRequest, "Time is require", nil)
		c.Abort()
		return
	}

	distance := c.Request.FormValue("distance")
	if distance == "" {
		res.Response(http.StatusBadRequest, "Time is require", nil)
		c.Abort()
		return
	}

	uniqidFilename := uuid.New()

	pathDir := "." + config.UPLOAD_KAO
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		os.MkdirAll(pathDir, os.ModePerm)
	}

	out, err := os.Create(pathDir + uniqidFilename.String() + ".png")

	//path := pathDir + uniqidFilename.String() + ".png"
	if err != nil {
		log.Println(err)
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		c.Abort()
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)
	bodyWriter.WriteField("distance", distance)
	bodyWriter.WriteField("time", time)
	fileWriter, err := bodyWriter.CreateFormFile("imageData", uniqidFilename.String()+".png")
	if err != nil {
		fmt.Println(err)
		//fmt.Println("Create form file error: ", error)
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		c.Abort()
		return
	}
	buf := bytes.NewBuffer(nil)
	_, err = io.Copy(buf, file)
	fileWriter.Write(buf.Bytes())
	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()
	resp, err := http.Post(url, contentType, bodyBuf)
	if err != nil {
		res.Response(http.StatusInternalServerError, err.Error(), nil)
		c.Abort()
		return
	}
	resp.Header.Add("Authorization", TK)
	defer resp.Body.Close()
	fmt.Println(resp)

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		errString := fmt.Errorf("[%d %s]%s", resp.StatusCode, resp.Status, string(b))
		log.Println(errString)
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	fmt.Println(string(respData))
	// if err == nil {
	// 	req, err := http.NewRequest("POST", url, bytes.NewBuffer(b))
	// 	req.Header.Add("Authorization", TK)
	// 	req.Header.Set("Content-Type", "application/json")
	// 	client := &http.Client{}
	// 	resp, err := client.Do(req)
	// 	if err != nil {
	// 		log.Println("Error on response.\n[ERRO] -", err)
	// 	} else {
	// 		defer resp.Body.Close()
	// 	}
	// }
}

func checkRedirectFunc(req *http.Request, via []*http.Request) error {
	req.Header.Add("Authorization", via[0].Header.Get("Authorization"))
	return nil
}

// KaoActivity send activity data to Kao
func KaoActivity(path string, distance float64, timeRun int64, slug string, ebib string, orderID string, phone string) ([]byte, error) {
	urls := fmt.Sprintf("https://kaokonlakao.com/api/%s/item/%s/submit", slug, orderID)
	pathAbs, err := filepath.Abs(fmt.Sprintf(".%s",path))
	if err != nil {
		return []byte{}, err
	}
	file, err := os.Open(pathAbs)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}

	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile("imageData", filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	writer.WriteField("distance", fmt.Sprintf("%.2f", distance))
	writer.WriteField("time", fmt.Sprintf("%d", timeRun))
	writer.WriteField("submitterMobileNumber", phone)

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	var bearer = "Bearer e35Sa9MvZJ1fA0PV"
	req, err := http.NewRequest("POST", urls, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Add("Authorization", bearer)
	
	timeout := time.Duration(6 * 3 * time.Second)
	client := &http.Client{
		Timeout: timeout,
	}
	//client.CheckRedirect = checkRedirectFunc

	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	// resp, err := http.Post(urls, "multipart/form-data", bodyBuf)
	// if err != nil {
	// 	return []byte{}, err
	// }
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, _ := ioutil.ReadAll(resp.Body)
		errString := fmt.Errorf("[%d %s]%s", resp.StatusCode, resp.Status, string(b))
		log.Println(errString)
		return []byte{}, errString
	}
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	return respData, err
}
