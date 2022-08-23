package helper

import "testing"

func TestStdin(t *testing.T) {
	t.Parallel()
	result, p := StdinPipeOrMany("enter 1 > ", "enter 2 > ", "enter x ")
	if p {
		t.Error("p should not be true")
	}
	if len(result) != 3 {
		t.Error("result should have len 3")
	}
	//want := [][]byte{[]byte{'1'}, []byte{'2'}, []byte{'3'}}
	if result[0][0] != '1' && result[1][0] != '1' && result[2][0] != '1' {
		t.Error("sumting wong")
	}
}
