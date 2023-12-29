package utils

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func UnZipFromBuffer(buf []byte, dst string) error {
	archive, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))

	if err != nil {
		return fmt.Errorf("cannot construct a zip reader, reason: %s", err.Error())
	}

	for _, f := range archive.File {
		filePath := filepath.Join(dst, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
				return fmt.Errorf("cannot make directory of zip file, reason: %s", err.Error())
			}
			continue
		}

		dstFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return fmt.Errorf("cannot open file: %s on disk, reason: %s", filePath, err.Error())
		}

		fileInArchive, err := f.Open()
		if err != nil {
			return fmt.Errorf("cannot open file: %s in archive, reason: %s", filePath, err.Error())
		}

		if _, err := io.Copy(dstFile, fileInArchive); err != nil {
			return fmt.Errorf("cannot copy file: %s from arcive to disk, reason: %s", filePath, err.Error())
		}

		dstFile.Close()
		fileInArchive.Close()
	}

	return nil
}
