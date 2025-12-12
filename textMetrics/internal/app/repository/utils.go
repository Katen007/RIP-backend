package repository

import (
	"fmt"
	"math"

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

func calcFlesch(words, sentences, syllables int) int {
	if words <= 0 || sentences <= 0 {
		return 0
	}
	fre := 206.835 - 1.015*(float64(words)/float64(sentences)) - 84.6*(float64(syllables)/float64(words))
	if fre < 0 {
		fre = 0
	}
	if fre > 206.835 {
		fre = 206.835
	}
	return int(math.Round(fre))
}
