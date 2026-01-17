package handlers

import (
	"github.com/tonysanin/brobar/pkg/helpers"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/tonysanin/brobar/file-service/internal/services"
	"github.com/tonysanin/brobar/pkg/response"
)

type FileHandler struct {
	fileService *services.FileService
}

func NewFileHandler(s *services.FileService) *FileHandler {
	return &FileHandler{fileService: s}
}

func (h *FileHandler) UploadFile(c fiber.Ctx) error {
	fileUrl := helpers.GetEnv("FILE_URL", "")

	file, err := c.FormFile("file")
	if err != nil {
		return response.BadRequest(c, err)
	}

	filename, err := h.fileService.SaveFile(file)
	if err != nil {
		return response.Error(c, fiber.StatusInternalServerError, err)
	}

	return response.Success(c, fiber.Map{
		"filename": filename,
		"path":     fileUrl + filename,
		"message":  "file uploaded successfully",
	})
}

func (h *FileHandler) GetFile(c fiber.Ctx) error {
	filename := c.Params("filename")
	path := h.fileService.GetFilePath(filename)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return response.NotFound(c)
	}

	return c.SendFile(path, fiber.SendFile{Download: false})
}
