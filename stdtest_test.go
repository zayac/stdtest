package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestOpenTestConfig(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "config.json")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte(`
[{
	"name":"test",
	"input":"1 2 3",
	"output":["4", "5", "6"],
	"is_output_ordered": true
}]`)
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	got, err := openTestConfig(tmpFile.Name())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := []TestConfig{
		{Name: "test", Input: "1 2 3", Output: []string{"4", "5", "6"}, IsOutputOrdered: true},
	}
	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %v, got %v", want, got)
	}
}

func TestMatchOutputUnordered(t *testing.T) {
	for _, tst := range []struct {
		name    string
		inStr   string
		inStrs  []string
		wantErr bool
	}{
		{
			name:   "found",
			inStr:  "foo bar baz",
			inStrs: []string{"baz", "foo"},
		}, {
			name:    "not found",
			inStr:   "foo bar baz",
			inStrs:  []string{"foo", "foob"},
			wantErr: true,
		},
	} {
		if got := matchOutputUnordered(tst.inStr, tst.inStrs); (got != nil) != tst.wantErr {
			t.Errorf("%s: want error: %v, got %q", tst.name, tst.wantErr, got)
		}
	}
}

func TestMatchOutputOrdered(t *testing.T) {
	for _, tst := range []struct {
		name    string
		inStr   string
		inStrs  []string
		wantErr bool
	}{
		{
			name:    "wrong order",
			inStr:   "zurich riga moscow daugavpils",
			inStrs:  []string{"moscow", "zurich"},
			wantErr: true,
		}, {
			name:   "found",
			inStr:  "zurich riga moscow daugavpils",
			inStrs: []string{"zurich", "riga"},
		}, {
			name:    "not found",
			inStr:   "zurich riga moscow daugavpils",
			inStrs:  []string{"zurich", "liepaja"},
			wantErr: true,
		},
	} {
		if got := matchOutputOrdered(tst.inStr, tst.inStrs); (got != nil) != tst.wantErr {
			t.Errorf("%s: want error: %v, got %q", tst.name, tst.wantErr, got)
		}
	}
}

/*func TestTestProgram(t *testing.T) {
	tmpFile, err := ioutil.TempFile("", "sum.go")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte(`package main
import "fmt"
func main() {
	var sum, val int
	for i := 0; i < 5; i++ {
		fmt.Scan(&val)
		sum += val
	}
	fmt.Println(sum)
}`)
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	config := []TestConfig{
		{
			Name:   "sample",
			Input:  "1 2 3 4 5",
			Output: []string{"15"},
		},
	}
	if err := testProgram(tmpFile.Name(), config); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}*/
