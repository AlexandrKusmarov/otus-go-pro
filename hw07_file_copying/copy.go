package main

import (
	"bufio"
	"errors"
	"io"
	"os"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	defer fromFile.Close()

	stat, err := fromFile.Stat()
	if err != nil {
		return err
	}

	if offset > stat.Size() {
		return ErrOffsetExceedsFileSize
	}

	err = limitCorrector(fromFile, &limit, &offset)
	if err != nil {
		return err
	}

	if limit == 0 || limit > stat.Size() {
		limit = stat.Size() - offset
	}

	toFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer toFile.Close()

	_, err = fromFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	bar := progressbar.DefaultBytes(
		limit,
		"copying",
	)

	_, err = io.CopyN(io.MultiWriter(toFile, bar), fromFile, limit)
	if err != nil {
		return err
	}
	bar.Finish()

	return nil
}

func limitCorrector(f *os.File, limit *int64, offset *int64) error {
	reader := bufio.NewReader(f)
	countR := int64(0)
	curOffset := int64(0)
	curLimit := int64(0)
	var err error

	// Подсчет и коррекция /r до оффсета
	for curOffset < *offset+countR {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		curOffset++
		if b == '\r' {
			countR++
		}
	}

	*offset = curOffset
	countR = 0

	// Подсчет и коррекция /r до лимита
	for curLimit < *limit+countR {
		b, err := reader.ReadByte()
		if err != nil {
			break
		}
		curLimit++
		if b == '\r' {
			countR++
		}
	}
	*limit = curLimit

	_, err = f.Seek(0, 0)

	return err
}
