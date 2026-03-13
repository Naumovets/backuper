package compressor

import (
	"compress/gzip"
	"io"
	"os"
)

func CompressFile(filename string, compressLevel int) (string, error) {
	newFilename := filename + ".gz"
	f, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer f.Close()

	out, err := os.Create(newFilename)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Use specified compression level if valid, otherwise default
	level := compressLevel
	if level < 0 || level > 9 {
		level = gzip.DefaultCompression
	}

	gw, err := gzip.NewWriterLevel(out, level)
	if err != nil {
		return "", err
	}
	defer gw.Close()

	if _, err := io.Copy(gw, f); err != nil {
		return "", err
	}

	return newFilename, nil
}
