package controller

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"product-crud/config"
	"product-crud/controller/response"
	"product-crud/util/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileController interface {
	Upload(c *gin.Context)
	Download(c *gin.Context)
}

type fileController struct {
}

func NewFileController() *fileController {
	logger.Info("Initializing file controller..")
	return &fileController{}
}

func (fc fileController) Upload(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic: %+v", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		panic(err)
	}
	newpath := filepath.Join(config.Env.FilePath, "public")
	err = os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logger.Info(`file size: %+v`, header.Size)
	ext := filepath.Ext(header.Filename)
	uuid, err := uuid.NewRandom()
	if err != nil {
		panic(err)
	}
	filename := uuid.String() + ext
	out, err := os.Create(newpath + "/" + filename)
	if err != nil {
		panic(err)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		panic(err)
	}
	filepath := filename
	response.Success(c, filepath, false)
}

func (fc fileController) Download(c *gin.Context) {
	defer func() {
		if r := recover(); r != nil {
			logger.Error("Recovered from panic: %+v", r)
			response.Fail(c, "Internal Server Error")
			return
		}
	}()
	newpath := filepath.Join(config.Env.FilePath, "public")
	filename := c.Param("name")
	if _, err := os.Stat(newpath + "/" + filename); errors.Is(err, os.ErrNotExist) {
		panic(err)
	}
	c.File(newpath + "/" + filename)

}

var _ FileController = (*fileController)(nil)
