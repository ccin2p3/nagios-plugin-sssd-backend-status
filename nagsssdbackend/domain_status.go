package nagsssdbackend

import (
	"bytes"
	"fmt"

	"github.com/pkg/errors"
)

var (
	onlineStatusLinePrefix   = []byte("Online status: ")
	onlineStatusStatusBytes  = []byte("Online")
	offlineStatusStatusBytes = []byte("Offline")
)

type domainStatus struct {
	Online bool
}

func parseDomainStatus(bdStatus []byte) (domainStatus, error) {
	var dStatus domainStatus

	lines := bytes.Split(bdStatus, []byte("\n"))
	for _, line := range lines {
		if bytes.HasPrefix(line, onlineStatusLinePrefix) {
			bstatus := line[len(onlineStatusLinePrefix):]
			status, err := parseOnlineStatus(bstatus)
			if err != nil {
				return dStatus, errors.Wrap(err, "parsing online status")
			}
			dStatus.Online = status
		}
	}

	return dStatus, nil
}

func parseOnlineStatus(status []byte) (bool, error) {
	if bytes.Equal(status, onlineStatusStatusBytes) {
		return true, nil
	}

	if bytes.Equal(status, offlineStatusStatusBytes) {
		return false, nil
	}

	return false, fmt.Errorf("unexpected online status '%s'", status)
}
