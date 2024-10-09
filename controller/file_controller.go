package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/swiftchat-backend/service"
	"github.com/kwa0x2/swiftchat-backend/utils"
	"net/http"
)

type IFileController interface {
	UploadFile(c *gin.Context)
}

type fileController struct {
	S3Service service.IS3Service
}

func NewFileController(s3Service service.IS3Service) IFileController {
	return &fileController{S3Service: s3Service}
}

// region "UploadFile" handles the file upload process.
func (ctrl *fileController) UploadFile(ctx *gin.Context) {
	// Extract the file from the form data
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Form File Error", err.Error()))
		return
	}

	// Ensure the file is closed after processing
	defer file.Close()

	// Upload the file to the S3 bucket and retrieve the file URL.
	fileURL, fileErr := ctrl.S3Service.UploadFile(file, header)
	if fileErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Error uploading file to S3 bucket"))
		return
	}

	ctx.JSON(http.StatusOK, fileURL)
}

// endregion
