package controller

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"product-crud/config"
	"product-crud/constant/errorconstants"
	"product-crud/util/apiresponse"
	"product-crud/util/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type IFileController interface {
	Upload(c *gin.Context)
	Download(c *gin.Context)
}

type FileController struct {
}

func NewFileController() FileController {
	logger.Info("Initializing file controller..")
	return FileController{}
}

func (fc FileController) Upload(c *gin.Context) {
	logger.Info(`Upload file request`)
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		logger.Error("Error : %v", err)
		panic(errorconstants.INTERNAL_ERROR)
	}
	newpath := filepath.Join(config.Env.FilePath, "public")
	err = os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(errorconstants.INTERNAL_ERROR)
	}
	logger.Info(`file size: %+v`, header.Size)
	ext := filepath.Ext(header.Filename)
	uuid, err := uuid.NewRandom()
	if err != nil {
		logger.Error("Error : %v", err)
		panic(errorconstants.INTERNAL_ERROR)
	}
	filename := uuid.String() + ext
	out, err := os.Create(newpath + "/" + filename)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(errorconstants.INTERNAL_ERROR)
	}
	defer out.Close()
	_, err = io.Copy(out, file)
	if err != nil {
		logger.Error("Error : %v", err)
		panic(errorconstants.INTERNAL_ERROR)
	}
	filepath := filename
	apiresponse.Ok(c, filepath, false)
}

func (fc FileController) Download(c *gin.Context) {
	newpath := filepath.Join(config.Env.FilePath, "public")
	filename := c.Param("name")
	if _, err := os.Stat(newpath + "/" + filename); errors.Is(err, os.ErrNotExist) {
		logger.Error("Error : %v", err)
		panic(errorconstants.INTERNAL_ERROR)
	}
	c.File(newpath + "/" + filename)

}

var _ IFileController = (*FileController)(nil)
