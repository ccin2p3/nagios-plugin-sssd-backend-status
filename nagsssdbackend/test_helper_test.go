package nagsssdbackend

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func testHelperLoadTestdataFile(t *testing.T, filename string) []byte {
	file, err := os.Open(filepath.Join("testdata", filename))
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	filebytes, err := ioutil.ReadAll(file)
	if err != nil {
		t.Fatal(err)
	}

	return filebytes
}

// stubExecFuncWithOutputFromFile allows to easily stub the exec function
// and return arbitrary content from a file
// It also takes an error as input to allow to test error reporting / handling
func stubExecFuncWithOutputFromFile(t *testing.T, filename string, err error) {
	execFnc = func(name string, args ...string) ([]byte, error) {
		var fcontent []byte
		if err != nil {
			return fcontent, err
		}

		return testHelperLoadTestdataFile(t, filename), nil
	}
}
