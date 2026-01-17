package services

import (
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"mime/multipart"
	"os"
	"path/filepath"
)

type FileService struct {
	UploadDir string
}

func NewFileService(uploadDir string) *FileService {
	err := os.MkdirAll(uploadDir, os.ModePerm)
	if err != nil {
		return nil
	}
	return &FileService{UploadDir: uploadDir}
}

func (fs *FileService) SaveFile(fileHeader *multipart.FileHeader) (string, error) {
	ext := filepath.Ext(fileHeader.Filename)
	uniqueName := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dst := filepath.Join(fs.UploadDir, uniqueName)

	src, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer func() {
		if cerr := src.Close(); cerr != nil {
			log.Fatalf("Error while closing source file: %v", cerr)
		}
	}()

	out, err := os.Create(dst)
	if err != nil {
		return "", err
	}
	defer func() {
		if cerr := out.Close(); cerr != nil {
			log.Fatalf("Error while closing destination file: %v", cerr)
		}
	}()

	_, err = io.Copy(out, src)
	if err != nil {
		return "", err
	}

	return uniqueName, nil
}

func (fs *FileService) GetFilePath(filename string) string {
	return fmt.Sprintf("%s/%s", fs.UploadDir, filename)
}
