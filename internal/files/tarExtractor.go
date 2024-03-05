package files

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const DirToExtractTar = "extracted_tar"

func extractTar(tarPath string) (string, error) {
	r, err := os.Open(tarPath)
	if err != nil {
		fmt.Println("error")
	}

	extractDir := filepath.Join(filepath.Dir(tarPath), DirToExtractTar)
	err = os.MkdirAll(extractDir, os.ModePerm)
	if err != nil {
		log.Printf("Could not create directory `%s`", extractDir)
		return "", err
	}

	gzr, err := gzip.NewReader(r)
	if err != nil {
		log.Printf("Could not create gzip reader `%s`", extractDir)
		return "", err
	}
	defer func(gzr *gzip.Reader) {
		err = gzr.Close()
		if err != nil {
			log.Println("error whole closing gzip reader", err)
		}
	}(gzr)

	tr := tar.NewReader(gzr)

	for {
		header, headerErr := tr.Next()
		switch {

		case headerErr == io.EOF:
			return extractDir, nil

		case headerErr != nil:
			log.Println("could not get next header", headerErr)
			return "", headerErr

		case header == nil:
			continue
		}

		target := filepath.Join(extractDir, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if _, dirFileErr := os.Stat(target); dirFileErr != nil {
				if dirFileErr = os.MkdirAll(target, 0755); dirFileErr != nil {
					log.Printf("could not create dir: %s, err: %v", target, dirFileErr)
					return "", dirFileErr
				}
			}

		case tar.TypeReg:
			f, regFileErr := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if regFileErr != nil {
				log.Printf("could not open file: %s, err: %v", target, regFileErr)
				return "", regFileErr
			}

			if _, regFileErr = io.Copy(f, tr); regFileErr != nil {
				log.Printf("could not copy file: %s, err: %v", f.Name(), regFileErr)
				return "", regFileErr
			}

			regFileErr = f.Close()
			if regFileErr != nil {
				log.Println("error whole closing file", f.Name(), regFileErr)
				return "", regFileErr
			}
		}
	}
}
