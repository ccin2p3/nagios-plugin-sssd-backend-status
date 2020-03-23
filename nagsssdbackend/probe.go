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
