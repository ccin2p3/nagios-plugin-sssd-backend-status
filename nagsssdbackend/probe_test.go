package nagsssdbackend

import (
	"reflect"
	"testing"
)

func TestFetchAllDomains(t *testing.T) {

	type testCase_ struct {
		Expected     []string
		TestFilename string
	}

	for testName, testCase := range map[string]testCase_{
		"one-domain": {
			Expected:     []string{"thedomain.cc.in2p3.fr"},
			TestFilename: "sssctl_domain_list_one_domain.txt",
		},
		"two-domains": {
			Expected:     []string{"thedomain.cc.in2p3.fr", "theotherdomain.cc.in2p3.fr"},
			TestFilename: "sssctl_domain_list_two_domains.txt",
		},
	} {

		t.Run(testName, func(t *testing.T) {
			stubExecFuncWithOutputFromFile(t, testCase.TestFilename, nil)

			probe := sssdBackendStatusProbe{}
			domains, err := probe.fetchAllDomains()
			if err != nil {
				t.Error(err)
			}

			if !reflect.DeepEqual(domains, testCase.Expected) {
				t.Errorf("was expecting %s, got %s", testCase.Expected, domains)
			}
		})
	}
}
