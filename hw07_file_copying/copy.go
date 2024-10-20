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
	var err error

	// Подсчет и коррекция /r до оффсета
	*offset, err = correctParam(f, offset)
	if err != nil {
		return err
	}

	_, err = f.Seek(*offset, 0)
	// Подсчет и коррекция /r до лимита
	*limit, err = correctParam(f, limit)
	if err != nil {
		return err
	}

	return err
}

func correctParam(f *os.File, param *int64) (int64, error) {
	var countR int64 = 0
	var curParam int64 = 0
	reader := bufio.NewReader(f)

	// Подсчет и коррекция /r до лимита
	for curParam < *param+countR {
		b, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		curParam++
		if b == '\r' {
			countR++
		}
	}
	_, err := f.Seek(0, 0)

	return curParam, err
}
