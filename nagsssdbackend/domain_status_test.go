package nagsssdbackend

import (
	"testing"
)

func TestParseDomainStatus(t *testing.T) {

	type testCase_ struct {
		Expected     domainStatus
		TestFilename string
	}

	for testName, testCase := range map[string]testCase_{
		"domain-online": {
			Expected: domainStatus{
				Online: true,
			},
			TestFilename: "sssctl_domain_status_online.txt",
		},
		"domain-offline": {
			Expected: domainStatus{
				Online: false,
			}, TestFilename: "sssctl_domain_status_offline.txt",
		},
	} {

		t.Run(testName, func(t *testing.T) {
			fbytes := testHelperLoadTestdataFile(t, testCase.TestFilename)

			dStatus, err := parseDomainStatus(fbytes)
			if err != nil {
				t.Error(err)
			}

			assertDomainStatusEqual(t, dStatus, testCase.Expected)
		})
	}
}

func assertDomainStatusEqual(t *testing.T, a, b domainStatus) {
	if a.Online != b.Online {
		t.Errorf("Online status differs")
	}
}
