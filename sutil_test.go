package sutil

import (
	"testing"

)


func TestWriteFile(t *testing.T) {

	err := WriteFile("aa", []byte("abcde\n"), 0600)

	if err != nil {
		t.Errorf("%s", err)
	}


	err = WriteFile("log/aa", []byte("abcde\n"), 0600)

	if err != nil {
		t.Errorf("%s", err)
	}


	err = WriteFile("log/log/aa", []byte("abcde\n"), 0600)

	if err != nil {
		t.Errorf("%s", err)
	}


	err = WriteFile("log/", []byte("abcde\n"), 0600)
	if err == nil {
		t.Errorf("%s", err)
	}

}

