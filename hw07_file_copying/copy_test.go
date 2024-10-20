package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"testing"
	"time"
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
			expectedError:    io.EOF,
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
	}

	fmt.Println(os.Getwd())

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := Copy(tc.fromPath, tc.toPath, tc.offset, tc.limit)
			if !errors.Is(err, tc.expectedError) {
				t.Errorf("Expected error %v, got %v", tc.expectedError, err)
			}
			if tc.expectedError == io.EOF {
				// Read actual file content
				actualFile, err := os.Open(tc.toPath)
				actualFileInfo, err := actualFile.Stat()
				actualContent := make([]byte, actualFileInfo.Size())
				_, err = actualFile.Read(actualContent)

				time.Sleep(2 * time.Second)

				// Read expected file content
				expectedFile, err := os.Open(tc.expectedFilePath)
				expectedFileInfo, err := expectedFile.Stat()
				expectedContent := make([]byte, expectedFileInfo.Size())
				_, err = expectedFile.Read(expectedContent)

				if err != nil {
					t.Fatalf("Failed to read output file: %v", err)
				}
				if !bytes.Equal(actualContent, expectedContent) {
					t.Errorf("Expected output file content %s, got %s", actualContent, expectedContent)
				}
				os.Remove(tc.toPath)
			}
		})
	}
}
