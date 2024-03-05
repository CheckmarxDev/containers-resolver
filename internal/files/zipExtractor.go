package files

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const DirToExtractZip = "extracted_zip"

func extractZip(zipPath string) (string, error) {
	extractDir := filepath.Join(filepath.Dir(zipPath), DirToExtractZip)
	err := os.MkdirAll(extractDir, os.ModePerm)
	if err != nil {
		log.Printf("Could not create directory `%s`", extractDir)
		return "", err
	}

	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		log.Printf("Could not create zip reader `%v`", err)
		return "", err
	}
	defer func(zipReader *zip.ReadCloser) {
		err = zipReader.Close()
		if err != nil {
			log.Println("error whole closing zip reader", err)
		}
	}(zipReader)

	prefix := ""
	if len(zipReader.File) > 0 {
		prefixParts := strings.Split(zipReader.File[0].Name, string(filepath.Separator))
		if len(prefixParts) > 1 {
			prefix = prefixParts[0]
		}
	}

	for _, file := range zipReader.File {

		var fileErr error
		var srcFile io.ReadCloser
		var destFile *os.File

		if strings.HasPrefix(file.Name, "__MACOSX") {
			continue
		}

		targetPath := filepath.Join(extractDir, strings.TrimPrefix(file.Name, prefix+string(filepath.Separator)))

		if file.FileInfo().IsDir() {
			fileErr = os.MkdirAll(targetPath, os.ModePerm)
			if fileErr != nil {
				log.Printf("Could not create new directory `%s`", targetPath)
				return "", fileErr
			}
			continue
		}

		if fileErr = os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); fileErr != nil {
			return "", fileErr
		}

		srcFile, fileErr = file.Open()
		if fileErr != nil {
			fileErr = srcFile.Close()
			if fileErr != nil {
				fmt.Println("error")
			}
			return "", fileErr
		}

		destFile, fileErr = os.Create(targetPath)
		if fileErr != nil {
			fileErr = destFile.Close()
			if fileErr != nil {
				fmt.Println("error")
			}
			return "", fileErr
		}

		if _, fileErr = io.Copy(destFile, srcFile); fileErr != nil {
			return "", fileErr
		}

		fileErr = srcFile.Close()
		if fileErr != nil {
			fmt.Println("error while closing src file")
		}
		fileErr = destFile.Close()
		if fileErr != nil {
			fmt.Println("error while closing dest file")

		}

	}

	return extractDir, nil
}
