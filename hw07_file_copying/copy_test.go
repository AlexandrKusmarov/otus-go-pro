package main

import (
	"bytes"
	"errors"
	"io"
	"os"
	"testing"
)

func TestCopy(t *testing.T) {
	type test struct {
		name             string
		fromPath         string
		toPath           string
		offset           int64
		limit            int64
		expectedError    error
		expectedFilePath string
	}

	tests := []test{
		{
			name:             "offset_0_limit_0",
			fromPath:         "testdata/input.txt",
			toPath:           "out.txt",
			offset:           0,
			limit:            0,
			expectedError:    nil,
			expectedFilePath: "testdata/out_offset0_limit0.txt",
		},
		{
			name:             "offset_0_limit_10",
			fromPath:         "testdata/input.txt",
			toPath:           "out.txt",
			offset:           0,
			limit:            10,
			expectedError:    nil,
			expectedFilePath: "testdata/out_offset0_limit10.txt",
		},
		{
			name:             "offset_0_limit_1000",
			fromPath:         "testdata/input.txt",
			toPath:           "out.txt",
			offset:           0,
			limit:            1000,
			expectedError:    nil,
			expectedFilePath: "testdata/out_offset0_limit1000.txt",
		},
		{
			name:             "offset_0_limit_10000",
			fromPath:         "testdata/input.txt",
			toPath:           "out.txt",
			offset:           0,
			limit:            10000,
			expectedError:    nil,
			expectedFilePath: "testdata/out_offset0_limit10000.txt",
		},
		{
			name:             "offset_0_limit_100000",
			fromPath:         "testdata/input.txt",
			toPath:           "out.txt",
			offset:           0,
			limit:            100000,
			expectedError:    nil,
			expectedFilePath: "testdata/out_offset0_limit100000.txt",
		},
		{
			name:             "offset_100_limit_1000",
			fromPath:         "testdata/input.txt",
			toPath:           "out.txt",
			offset:           100,
			limit:            1000,
			expectedError:    nil,
			expectedFilePath: "testdata/out_offset100_limit1000.txt",
		},
		{
			name:             "offset_6000_limit_1000",
			fromPath:         "testdata/input.txt",
			toPath:           "out.txt",
			offset:           6000,
			limit:            1000,
			expectedError:    nil,
			expectedFilePath: "testdata/out_offset6000_limit1000.txt",
		},
		{
			name:             "offset_600000_limit_100000",
			fromPath:         "testdata/input.txt",
			toPath:           "out.txt",
			offset:           600000,
			limit:            100000,
			expectedError:    ErrOffsetExceedsFileSize,
			expectedFilePath: "",
		},
		{
			name:             "input_equals_output_filename",
			fromPath:         "testdata/out.txt",
			toPath:           "testdata/out.txt",
			offset:           600000,
			limit:            100000,
			expectedError:    ErrUnsupportedFile,
			expectedFilePath: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			defer os.Remove(tc.toPath)
			err := Copy(tc.fromPath, tc.toPath, tc.offset, tc.limit)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}
			if errors.Is(tc.expectedError, io.EOF) {
				// Read actual file content
				actualContent := readData(t, tc.toPath)

				// Read expected file content
				expectedContent := readData(t, tc.expectedFilePath)
				if !bytes.Equal(actualContent, expectedContent) {
					t.Errorf("Expected output file content %s, got %s", actualContent, expectedContent)
				}
			}
		})
	}
}

func readData(t *testing.T, path string) []byte {
	t.Helper()
	actualFile, err := os.Open(path)
	if err != nil {
		t.Fatalf("Failed to open file: %v", err)
	}

	actualFileInfo, err := actualFile.Stat()
	if err != nil {
		t.Fatalf("Failed to stat file: %v", err)
	}

	actualContent := make([]byte, actualFileInfo.Size())
	_, err = actualFile.Read(actualContent)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}
	return actualContent
}
