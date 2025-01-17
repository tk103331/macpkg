package cpio

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"io"
	"os"
	"testing"
)

func TestNewcFormat(t *testing.T) {
	tests := []struct {
		name string
		path string
		want []string
	}{
		{
			name: "odc-test",
			path: "testdata/odc-test.cpio",
			want: []string{"test.sh", "test.txt", "test", "/dev/nvme0n1", "/run/dmeventd-client", "/run/docker.sock"},
		}, {
			name: "test_odc",
			path: "testdata/test_odc.cpio",
			want: []string{"gophers.txt", "readme.txt", "todo.txt", "checklist.txt"},
		}, {
			name: "test_svr4",
			path: "testdata/test_svr4.cpio",
			want: []string{"gophers.txt"},
		}, {
			name: "test_svr4_crc",
			path: "testdata/test_svr4_crc.cpio",
			want: []string{"gophers.txt"},
		},
	}
	for i := 0; i < len(tests); i++ {
		test := tests[i]
		t.Run(test.name, func(tt *testing.T) {
			reader, err := os.Open(test.path)
			if err != nil {
				tt.Errorf("Open() error = %v", err)
				return
			}
			defer reader.Close()
			cpioReader := NewReader(reader)
			files := make([]string, 0)
			for {
				file, err := cpioReader.Next()
				if errors.Is(err, io.EOF) {
					break
				} else if err != nil {
					tt.Errorf("Read() error = %v", err)
				}
				files = append(files, file.Name())
			}
			if !assert.Equal(tt, test.want, files) {
				tt.Logf("Read() files = %v, want %v", files, test.want)
			}
		})
	}
}
