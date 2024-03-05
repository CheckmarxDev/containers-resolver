package files

import (
	"archive/zip"
	"github.com/Checkmarx-Containers/containers-resolver/internal/logger"
	"io"
	"os"
	"path/filepath"
	"strings"
)

const DirToExtractZip = "extracted_zip"

func extractZip(l *logger.Logger, zipPath string) (string, error) {
	extractDir := filepath.Join(filepath.Dir(zipPath), DirToExtractZip)
	err := os.MkdirAll(extractDir, os.ModePerm)
	if err != nil {
		l.Error("Could not create directory `%s`", extractDir, err)
		return "", err
	}

	zipReader, err := zip.OpenReader(zipPath)
	if err != nil {
		l.Error("Could not create zip reader `%v`", err)

		return "", err
	}
	defer func(zipReader *zip.ReadCloser) {
		err = zipReader.Close()
		if err != nil {
			l.Warn("error whole closing zip reader", err)
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
				l.Error("Could not create new directory `%s`", targetPath, err)
				return "", fileErr
			}
			continue
		}

		srcFile, fileErr = file.Open()
		if fileErr != nil {
			fileErr = srcFile.Close()
			if fileErr != nil {
				l.Warn("Could not close src file `%s`", srcFile, err)
			}
			return "", fileErr
		}

		destFile, fileErr = os.Create(targetPath)
		if fileErr != nil {
			fileErr = destFile.Close()
			if fileErr != nil {
				l.Warn("Could not close dest file `%s`", destFile, err)
			}
			return "", fileErr
		}

		if _, fileErr = io.Copy(destFile, srcFile); fileErr != nil {
			l.Error("Could not close dest file `%s`", destFile, err)
			return "", fileErr
		}

		fileErr = srcFile.Close()
		if fileErr != nil {
			l.Warn("error while closing src file", err)
		}
		fileErr = destFile.Close()
		if fileErr != nil {
			l.Warn("error while closing dest file", err)
		}

	}
	l.Debug("Successfully extracts zip folder to: %s", extractDir)
	return extractDir, nil
}
