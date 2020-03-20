package nagsssdbackend

import (
	"bytes"
	"fmt"

	nerrors "github.com/ccin2p3/nagios-plugin-sssd-backend-status/nagsssdbackend/errors"
	"github.com/ccin2p3/nagios-plugin-sssd-backend-status/nagsssdbackend/nagios"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type sssdBackendStatusProbe struct {
	domains []string
}

func NewSSSdBackendStatusProbe(domains []string) sssdBackendStatusProbe {
	return sssdBackendStatusProbe{
		domains: domains,
	}
}

func (p sssdBackendStatusProbe) Execute() error {
	domains := p.domains

	if len(domains) == 0 {
		var err error
		domains, err = p.fetchAllDomains()
		if err != nil {
			return errors.Wrap(err, "fetching all domains")
		}

		log.Debugf("discovered domains %s", domains)
	}

	var errs []error
	for _, domain := range domains {
		err := p.checkDomain(domain)
		if err != nil {
			errs = append(errs, nerrors.DomainStatusError{
				Err:    err,
				Domain: domain,
			})
		}
	}

	// this function will format errors to Nagios output
	// and exits the program with Nagios specific exit codes
	// This function will terminate the program
	nagios.ToNagiosOutput(errs)

	// This code will never be reached
	return nil
}

func (p sssdBackendStatusProbe) checkDomain(domain string) error {

	dStatus, err := p.fetchDomainStatus(domain)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("fetching domain '%s' status", domain))
	}

	if !dStatus.Online {
		return fmt.Errorf("domain is offline")
	}

	return nil
}

func (p sssdBackendStatusProbe) fetchDomainStatus(domain string) (domainStatus, error) {
	var dStatus domainStatus

	cmd := []string{"sssctl", "domain-status", domain}

	log.Debugf("executing command %s", cmd)

	rawOutput, err := execFnc(cmd[0], cmd[1:]...)
	if err != nil {
		return dStatus, errors.Wrap(err, fmt.Sprintf("fetching domain status: %s", rawOutput))
	}

	dStatus, err = parseDomainStatus(rawOutput)
	if err != nil {
		return dStatus, errors.Wrap(err, "parsing domain status")
	}

	return dStatus, nil
}

func (p sssdBackendStatusProbe) fetchAllDomains() ([]string, error) {
	var domains []string

	cmd := []string{"sssctl", "domain-list"}

	log.Debugf("executing command %s", cmd)

	rawOutput, err := execFnc(cmd[0], cmd[1:]...)
	if err != nil {
		return domains, errors.Wrap(err, "listing domains")
	}

	bdomains := bytes.Split(rawOutput, []byte("\n"))
	for _, bdomain := range bdomains {
		if len(bdomain) == 0 {
			continue
		}
		domains = append(domains, string(bdomain))
	}

	return domains, nil
}
