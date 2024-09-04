package main

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"testing"

	"contribs-go/model"
)

func TestExtractor_Extract(t *testing.T) {
	for _, tt := range getTests(t) {
		t.Run(tt.name, func(t *testing.T) {
			ex := NewExtractor([]byte(tt.src))
			if ex.Error != nil {
				t.Fatal(ex.Error)
			}

			if got := ex.Extract(); !reflect.DeepEqual(got, tt.want) {
				if ex.Error != nil {
					t.Fatal(ex.Error)
				}

				t.Errorf("Extractor.Extract()\ngot 	= %v\nwant 	= %v", got, tt.want)
			}
		})
	}
}

func getTests(t *testing.T) []struct {
	name string
	src  string
	want map[model.Locus]struct{}
} {
	t.Helper()

	return []struct {
		name string
		src  string
		want map[model.Locus]struct{}
	}{
		{
			name: "moby/moby/pkg/system/chtimes.go",
			src:  openTest(t, "chtimes"),
			want: map[model.Locus]struct{}{
				{
					Ident: "os.Chtimes",
					Line:  42,
				}: {},
				{
					Ident: "syscall.Timespec",
					Line:  15,
				}: {},
				{
					Ident: "time.Time",
					Line:  11,
				}: {},
				{
					Ident: "time.Unix",
					Line:  14,
				}: {},
				{
					Ident: "time.Unix",
					Line:  22,
				}: {},
				{
					Ident: "unsafe.Sizeof",
					Line:  15,
				}: {},
				{
					Ident: "time.Unix",
					Line:  25,
				}: {},
				{
					Ident: "time.Time",
					Line:  33,
				}: {},
				{
					Ident: "time.Time",
					Line:  53,
				}: {},
			},
		},
		{
			name: "kubernetes/kubernetes/pkg/util/filesystem/defaultfs.go",
			src:  openTest(t, "defaultfs"),
			want: map[model.Locus]struct{}{
				{
					Ident: "os.Chtimes",
					Line:  109,
				}: {},
				{
					Ident: "os.Create",
					Line:  84,
				}: {},
				{
					Ident: "os.CreateTemp",
					Line:  134,
				}: {},
				{
					Ident: "path/filepath.Join",
					Line:  74,
				}: {},
				{
					Ident: "path/filepath.Walk",
					Line:  148,
				}: {},
				{
					Ident: "path/filepath.WalkFunc",
					Line:  49,
				}: {},
				{
					Ident: "path/filepath.WalkFunc",
					Line:  147,
				}: {},
				{
					Ident: "os.DirEntry",
					Line:  48,
				}: {},
				{
					Ident: "os.DirEntry",
					Line:  142,
				}: {},
				{
					Ident: "os.File",
					Line:  153,
				}: {},
				{
					Ident: "os.FileInfo",
					Line:  36,
				}: {},
				{
					Ident: "os.FileInfo",
					Line:  78,
				}: {},
				{
					Ident: "os.FileMode",
					Line:  39,
				}: {},
				{
					Ident: "os.FileMode",
					Line:  103,
				}: {},
				{
					Ident: "os.MkdirAll",
					Line:  104,
				}: {},
				{
					Ident: "os.MkdirTemp",
					Line:  64,
				}: {},
				{
					Ident: "os.MkdirTemp",
					Line:  129,
				}: {},
				{
					Ident: "os.ReadDir",
					Line:  143,
				}: {},
				{
					Ident: "os.ReadFile",
					Line:  124,
				}: {},
				{
					Ident: "os.Remove",
					Line:  119,
				}: {},
				{
					Ident: "os.RemoveAll",
					Line:  114,
				}: {},
				{
					Ident: "os.Rename",
					Line:  99,
				}: {},
				{
					Ident: "os.Stat",
					Line:  79,
				}: {},
				{
					Ident: "strings.HasPrefix",
					Line:  93,
				}: {},
				{
					Ident: "strings.HasPrefix",
					Line:  96,
				}: {},
				{
					Ident: "time.Time",
					Line:  40,
				}: {},
				{
					Ident: "time.Time",
					Line:  108,
				}: {},
			},
		},
		{
			name: "moby/moby/pkg/tarsum/fileinfosums.go",
			src:  openTest(t, "fileinfosums"),
			want: map[model.Locus]struct{}{
				{
					Ident: "runtime.GOOS",
					Line:  47,
				}: {},
				{
					Ident: "runtime.GOOS",
					Line:  48,
				}: {},
				{
					Ident: "sort.Sort",
					Line:  88,
				}: {},
				{
					Ident: "sort.Sort",
					Line:  93,
				}: {},
				{
					Ident: "sort.Sort",
					Line:  100,
				}: {},
				{
					Ident: "sort.Sort",
					Line:  102,
				}: {},
				{
					Ident: "strings.EqualFold",
					Line:  47,
				}: {},
			},
		},
		{
			name: "golang/go/src/context/context.go",
			src:  openTest(t, "context"),
			want: map[model.Locus]struct{}{
				{
					Ident: "sync/atomic.Int32",
					Line:  302,
				}: {},
				{
					Ident: "sync/atomic.Value",
					Line:  366,
				}: {},
				{
					Ident: "errors.New",
					Line:  109,
				}: {},
				{
					Ident: "sync.Mutex",
					Line:  365,
				}: {},
				{
					Ident: "sync.Once",
					Line:  279,
				}: {},
				{
					Ident: "time.AfterFunc",
					Line:  579,
				}: {},
				{
					Ident: "time.AfterFunc",
					Line:  579,
				}: {},
				{
					Ident: "time.Duration",
					Line:  630,
				}: {},
				{
					Ident: "time.Duration",
					Line:  637,
				}: {},
				{
					Ident: "time.Now",
					Line:  631,
				}: {},
				{
					Ident: "time.Now",
					Line:  638,
				}: {},
				{
					Ident: "time.Time",
					Line:  18,
				}: {},
				{
					Ident: "time.Time",
					Line:  125,
				}: {},
				{
					Ident: "time.Time",
					Line:  523,
				}: {},
				{
					Ident: "time.Time",
					Line:  552,
				}: {},
				{
					Ident: "time.Time",
					Line:  559,
				}: {},
				{
					Ident: "time.Time",
					Line:  593,
				}: {},
				{
					Ident: "time.Time",
					Line:  596,
				}: {},
				{
					Ident: "time.Timer",
					Line:  591,
				}: {},
				{
					Ident: "time.Until",
					Line:  571,
				}: {},
				{
					Ident: "time.Until",
					Line:  603,
				}: {},
			},
		},
		{
			name: "golang/go/src/strconv/ctoa_test.go",
			src:  openTest(t, "ctoa_test"),
			want: map[model.Locus]struct{}{
				{
					Ident: "testing.T",
					Line:  12,
				}: {},
				{
					Ident: "testing.T",
					Line:  46,
				}: {},
			},
		},
		{
			name: "golang/go/test/fixedbugs/bug233.go",
			src:  openTest(t, "bug233"),
			want: map[model.Locus]struct{}{
				{
					Ident: "fmt.Print",
					Line:  9,
				}: {},
			},
		},
		{
			name: "cilium/test/bpf_tests/trf.pb.go",
			src:  openTest(t, "trf.pb"),
			want: map[model.Locus]struct{}{
				{
					Ident: "reflect.TypeOf",
					Line:  364,
				}: {},
				{
					Ident: "sync.Once",
					Line:  287,
				}: {},
			},
		},
	}
}

func openTest(t *testing.T, name string) string {
	t.Helper()

	name = fmt.Sprintf("testfiles/%s.go.txt", name)
	file, err := os.Open(name)
	if err != nil {
		t.Fatal(err)
	}
	bs, err := io.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	return string(bs)
}
