package uploads

import (
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"

	//"bitbucket.org/suthisakch/evenex/pkg/upload"
	"github.com/gin-gonic/gin"
	guuid "github.com/google/uuid"
	"github.com/nfnt/resize"
	config "thinkdev.app/think/runex/runexapi/config"
	"thinkdev.app/think/runex/runexapi/pkg/app"
	"thinkdev.app/think/runex/runexapi/pkg/e"
	"thinkdev.app/think/runex/runexapi/utils"
)

// GetImageFullUrl get the full access path
func GetImageFullUrl(name string) string {
	return config.HTTP + "/" + GetImagePath() + name
}

// GetImageName get image name
func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = utils.EncodeMD5(fileName)

	return fileName + ext
}

// GetImagePath get save path
func GetImagePath() string {
	return config.ImageSavePath
}

// GetImageFullPath get full save path
func GetImageFullPath() string {
	return config.RuntimeRootPath + GetImagePath()
}

// CheckImageExt check image file ext
func CheckImageExt(fileName string) bool {
	ImageAllowExts := [3]string{".jpg", ".jpeg", ".png"}
	ext := utils.GetExt(fileName)
	for _, allowExt := range ImageAllowExts {
		if strings.ToUpper(allowExt) == strings.ToUpper(ext) {
			return true
		}
	}

	return false
}

// Uploads image by create event
func Uploads(c *gin.Context) {
	file, header, err := c.Request.FormFile("upload")
	filename := header.Filename
	fmt.Println(filename)

	uniqidFilename := guuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", uniqidFilename.String())

	pathDir := "." + config.UPLOAD_IMAGE
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		os.MkdirAll(pathDir, os.ModePerm)
	}

	out, err := os.Create(pathDir + uniqidFilename.String() + ".png")

	path := config.UPLOAD_IMAGE + uniqidFilename.String() + ".png"
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
		//log.Fatal(err)
	}

	res := gin.H{"msg": "success", "url": path}
	c.JSON(http.StatusOK, res)
}

func UploadCover(c *gin.Context) {
	file, header, err := c.Request.FormFile("upload")
	file2, _, _ := c.Request.FormFile("upload")
	filename := header.Filename
	fmt.Println(filename)

	uniqidFilename := guuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", uniqidFilename.String())

	pathDir := "." + config.UPLOAD_IMAGE
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		os.MkdirAll(pathDir, os.ModePerm)
	}

	out, err := os.Create(pathDir + uniqidFilename.String() + ".png")

	path := config.UPLOAD_IMAGE + uniqidFilename.String() + ".png"
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)

	//thumbnail
	pathDirThumbnail := "." + config.UPLOAD_IMAGE + "thumbnail/"
	if _, err := os.Stat(pathDirThumbnail); os.IsNotExist(err) {
		os.MkdirAll(pathDirThumbnail, os.ModePerm)
	}

	// decode jpeg into image.Image
	//img, err := jpeg.Decode(file2)
	img, str, err := image.Decode(file2)
	log.Println(str)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m := resize.Resize(280, 0, img, resize.Lanczos3)

	outThumbnail, err := os.Create(pathDirThumbnail + uniqidFilename.String() + ".png")

	pathThumbnail := config.UPLOAD_IMAGE + "thumbnail/" + uniqidFilename.String() + ".png"
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer outThumbnail.Close()

	// write new image to file
	jpeg.Encode(outThumbnail, m, nil)
	//_, err = io.Copy(outThumbnail, file)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
		//log.Fatal(err)
	}

	res := gin.H{
		"msg":  "success",
		"url":  pathThumbnail,
		"path": path,
	}
	c.JSON(http.StatusOK, res)
}

func UploadSlip(c *gin.Context) {
	appG := app.Gin{C: c}
	file, header, err := c.Request.FormFile("upload")
	filename := header.Filename
	imageName := GetImageName(filename)
	fmt.Println(filename)

	// decode jpeg into image.Image
	//img, err := jpeg.Decode(file)
	img, str, err := image.Decode(file)
	log.Println(str)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	m := resize.Resize(680, 0, img, resize.Lanczos3)

	uniqidFilename := guuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", uniqidFilename.String())

	pathDir := "./upload/image/slip/"
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		os.MkdirAll(pathDir, os.ModePerm)
	}

	out, err := os.Create(pathDir + imageName)

	path := "/upload/image/slip/" + imageName
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer out.Close()
	//_, err = io.Copy(out, file)
	// write new image to file
	jpeg.Encode(out, m, nil)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
		//log.Fatal(err)
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"image_save_url": path,
	})
}

func UploadWithFolder(c *gin.Context) {
	appG := app.Gin{C: c}
	file, header, err := c.Request.FormFile("upload")
	file2, _, _ := c.Request.FormFile("upload")
	foldername := c.Request.FormValue("folder")
	width := c.Request.FormValue("width")
	filename := header.Filename
	imageName := GetImageName(filename)
	fmt.Println(filename)
	// decode jpeg into image.Image
	img, err := jpeg.Decode(file)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	// resize to width 1000 using Lanczos resampling
	// and preserve aspect ratio
	u64, err := strconv.Atoi(width)
	img_width := uint(u64)
	//img_width, err := strconv.ParseUint(width, 10, 32)
	m := resize.Resize(img_width, 0, img, resize.Lanczos3)

	uniqidFilename := guuid.New()
	fmt.Printf("github.com/google/uuid:         %s\n", uniqidFilename.String())

	pathDir := "./upload/image/" + foldername + "/"
	if _, err := os.Stat(pathDir); os.IsNotExist(err) {
		os.MkdirAll(pathDir, os.ModePerm)
	}

	out, err := os.Create(pathDir + imageName)

	path := "/upload/image/" + foldername + "/" + imageName
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer out.Close()
	//_, err = io.Copy(out, file)
	// write new image to file
	jpeg.Encode(out, m, nil)

	//thumbnail
	pathDirThumbnail := "./upload/image/" + foldername + "/" + "thumbnail/"
	if _, err := os.Stat(pathDirThumbnail); os.IsNotExist(err) {
		os.MkdirAll(pathDirThumbnail, os.ModePerm)
	}

	// decode jpeg into image.Image
	img2, err := jpeg.Decode(file2)
	if err != nil {
		log.Fatal(err)
	}
	file.Close()

	m2 := resize.Resize(280, 0, img2, resize.Lanczos3)

	outThumbnail, err := os.Create(pathDirThumbnail + uniqidFilename.String() + ".jpg")

	pathThumbnail := "./upload/image/" + foldername + "/" + "thumbnail/" + uniqidFilename.String() + ".jpg"
	if err != nil {
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	defer outThumbnail.Close()

	// write new image to file
	jpeg.Encode(outThumbnail, m2, nil)
	//_, err = io.Copy(outThumbnail, file)

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_SAVE_IMAGE_FAIL, nil)
		return
		//log.Fatal(err)
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"image_save_url":      path,
		"image_thumbnail_url": pathThumbnail,
	})
}
