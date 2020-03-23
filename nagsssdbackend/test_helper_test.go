/*
Copyright © 2020 IN2P3 Computing Centre, CNRS
Author(s): Remi Ferrand <remi.ferrand_at_cc.in2p3.fr>, 2020

This software is governed by the CeCILL-B license under French law and
abiding by the rules of distribution of free software.  You can  use,
modify and/ or redistribute the software under the terms of the CeCILL-B
license as circulated by CEA, CNRS and INRIA at the following URL
"http://www.cecill.info".

As a counterpart to the access to the source code and  rights to copy,
modify and redistribute granted by the license, users are provided only
with a limited warranty  and the software's author,  the holder of the
economic rights,  and the successive licensors  have only  limited
liability.

In this respect, the user's attention is drawn to the risks associated
with loading,  using,  modifying and/or developing or reproducing the
software by the user in light of its specific status of free software,
that may mean  that it is complicated to manipulate,  and  that  also
therefore means  that it is reserved for developers  and  experienced
professionals having in-depth computer knowledge. Users are therefore
encouraged to load and test the software's suitability as regards their
requirements in conditions enabling the security of their systems and/or
data to be ensured and,  more generally, to use and operate it in the
same conditions as regards security.

The fact that you are presently reading this means that you have had
knowledge of the CeCILL-B license and that you accept its terms.
*/
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
