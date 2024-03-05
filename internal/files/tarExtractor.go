package files

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/Checkmarx-Containers/containers-resolver/internal/logger"
	"io"
	"os"
	"path/filepath"
)

const DirToExtractTar = "extracted_tar"

func extractTar(l *logger.Logger, tarPath string) (string, error) {
	r, err := os.Open(tarPath)
	if err != nil {
		fmt.Println("error")
	}

	extractDir := filepath.Join(filepath.Dir(tarPath), DirToExtractTar)
	err = os.MkdirAll(extractDir, os.ModePerm)
	if err != nil {
		l.Error("Could not create directory `%s`", extractDir, err)
		return "", err
	}

	gzr, err := gzip.NewReader(r)
	if err != nil {
		l.Error("Could not create gzip reader `%s`", extractDir, err)
		return "", err
	}
	defer func(gzr *gzip.Reader) {
		err = gzr.Close()
		if err != nil {
			l.Warn("error whole closing gzip reader", err)
		}
	}(gzr)

	tr := tar.NewReader(gzr)

	for {
		header, headerErr := tr.Next()
		switch {

		case headerErr == io.EOF:
			return extractDir, nil

		case headerErr != nil:
			l.Error("could not get next header. error: ", headerErr)
			return "", headerErr

		case header == nil:
			continue
		}

		target := filepath.Join(extractDir, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if _, dirFileErr := os.Stat(target); dirFileErr != nil {
				if dirFileErr = os.MkdirAll(target, 0755); dirFileErr != nil {
					l.Error("could not create dir: %s, err: %v", target, dirFileErr)
					return "", dirFileErr
				}
			}

		case tar.TypeReg:
			f, regFileErr := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if regFileErr != nil {
				l.Error("could not open file: %s, err: %v", target, regFileErr)
				return "", regFileErr
			}

			if _, regFileErr = io.Copy(f, tr); regFileErr != nil {
				l.Error("could not copy file: %s, err: %v", f.Name(), regFileErr)
				return "", regFileErr
			}

			regFileErr = f.Close()
			if regFileErr != nil {
				l.Warn("error whole closing file %s", f.Name(), regFileErr)
				return "", regFileErr
			}
		}
	}
}
