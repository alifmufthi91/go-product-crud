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

func NewFileController() FileController {
	logger.Info("Initializing file controller..")
	return fileController{}
}

func (fc fileController) Upload(c *gin.Context) {
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	newpath := filepath.Join(config.Env.FilePath, "public")
	err = os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	ext := filepath.Ext(header.Filename)
	uuid, err := uuid.NewRandom()
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	filename := uuid.String() + ext
	out, err := os.Create(newpath + "/" + filename)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	filepath := filename
	response.Success(c, filepath)
}

func (fc fileController) Download(c *gin.Context) {
	newpath := filepath.Join(config.Env.FilePath, "public")
	filename := c.Param("name")
	if _, err := os.Stat(newpath + "/" + filename); errors.Is(err, os.ErrNotExist) {
		logger.Error(err.Error())
		response.Fail(c, errors.New("something went wrong").Error())
		return
	}
	c.File(newpath + "/" + filename)

}
