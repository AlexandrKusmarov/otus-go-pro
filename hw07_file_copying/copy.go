package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath string, toPath string, offset, limit int64) error {
	if err := checkFilePath(fromPath, toPath); err != nil {
		return err
	}

	special, err := isSpecialFile(fromPath)
	if err != nil || special {
		return ErrUnsupportedFile
	}

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

func isSpecialFile(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	// Проверяем, является ли файл устройством (char или block device)
	mode := fileInfo.Mode()
	if mode&os.ModeCharDevice != 0 || mode&os.ModeDevice != 0 {
		return true, nil
	}

	return false, nil
}

func checkFilePath(fromPath string, toPath string) error {
	abs1, err := filepath.Abs(fromPath)
	if err != nil {
		return ErrUnsupportedFile
	}
	abs2, err := filepath.Abs(toPath)
	if err != nil {
		return ErrUnsupportedFile
	}

	if abs1 == abs2 {
		return ErrUnsupportedFile
	}
	return nil
}
