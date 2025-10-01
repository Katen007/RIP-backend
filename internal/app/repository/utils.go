package repository

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateFileName() (string, error) {
	fileName, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return fileName.String(), nil
}

func CreateNewFilePath(filePath string) (string, error) {
	fileName, err := GenerateFileName()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", filePath, fileName), nil
}
