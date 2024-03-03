package zip

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const DirToExtractZip = "extracted_zip"

func ExtractZip(zipPath string) (string, error) {
	extractDir := filepath.Join(filepath.Dir(zipPath), DirToExtractZip)
	err := os.MkdirAll(extractDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", err
	}
	defer func(zipReader *zip.ReadCloser) {
		err = zipReader.Close()
		if err != nil {

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

		if strings.HasPrefix(file.Name, "__MACOSX") {
			continue
		}

		targetPath := filepath.Join(extractDir, strings.TrimPrefix(file.Name, prefix+string(filepath.Separator)))

		if file.FileInfo().IsDir() {
			err1 := os.MkdirAll(targetPath, os.ModePerm)
			if err1 != nil {
				return "", err1
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return "", err
		}

		srcFile, err := file.Open()
		if err != nil {
			err1 := srcFile.Close()
			if err1 != nil {
				fmt.Println("error")
			}
			return "", err
		}

		destFile, err1 := os.Create(targetPath)
		if err1 != nil {
			err1 = destFile.Close()
			if err1 != nil {
				fmt.Println("error")
			}
			return "", err
		}

		if _, err := io.Copy(destFile, srcFile); err != nil {
			return "", err
		}

		err1 = srcFile.Close()
		if err1 != nil {
			fmt.Println("error")
		}
		err1 = destFile.Close()
		if err1 != nil {
			fmt.Println("error")
		}

	}

	return extractDir, nil
}
