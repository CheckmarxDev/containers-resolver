package zip

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const DirToExtractTar = "extracted_tar"

func ExtractTar(tarPath string) (string, error) {
	r, err := os.Open(tarPath)
	if err != nil {
		fmt.Println("error")
	}

	extractDir := filepath.Join(filepath.Dir(tarPath), DirToExtractTar)
	err = os.MkdirAll(extractDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	gzr, err := gzip.NewReader(r)
	if err != nil {
		return "", err
	}
	defer func(gzr *gzip.Reader) {
		err = gzr.Close()
		if err != nil {
			fmt.Println("error")
		}
	}(gzr)

	tr := tar.NewReader(gzr)

	for {
		header, headerErr := tr.Next()
		switch {

		// if no more files are found return extracted files dir
		case headerErr == io.EOF:
			return extractDir, nil

		case headerErr != nil:
			return "", headerErr

		case header == nil:
			continue
		}

		target := filepath.Join(extractDir, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if _, dirFileErr := os.Stat(target); dirFileErr != nil {
				if dirFileErr = os.MkdirAll(target, 0755); dirFileErr != nil {
					return "", dirFileErr
				}
			}

		case tar.TypeReg:
			f, regFileErr := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if regFileErr != nil {
				return "", regFileErr
			}

			if _, regFileErr = io.Copy(f, tr); err != nil {
				return "", regFileErr
			}

			regFileErr = f.Close()
			if regFileErr != nil {
				return "", regFileErr
			}
		}
	}
}
