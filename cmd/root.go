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

package cmd

import (
	"fmt"
	"os"
	"strings"

	cmderrors "github.com/ccin2p3/nagios-plugin-sssd-backend-status/cmd/errors"
	"github.com/ccin2p3/nagios-plugin-sssd-backend-status/nagsssdbackend"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	allDomainsKeyword      = "ALL"
	domainsStringSeparator = ","
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:           "check_sssd_backend_status",
	Short:         "Checks that sssd backends are Online",
	Long:          `Checks that sssd backends are Online. This tool rely on sssctl.`,
	SilenceUsage:  true,
	SilenceErrors: true,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	RunE: runE,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		initLogging()

		if err := processArgs(cmd, args); err != nil {
			return err
		}

		return nil
	},
}

func processArgs(cmd *cobra.Command, args []string) error {

	domainsStr, err := cmd.Flags().GetString("domains")
	if err != nil {
		return err
	}

	var domains []string
	if domainsStr != "" {
		domains = strings.Split(domainsStr, domainsStringSeparator)
	}

	viper.Set("domains", domains)

	return nil
}

func runE(cmd *cobra.Command, args []string) error {

	domains := viper.GetStringSlice("domains")

	probe := nagsssdbackend.NewSSSdBackendStatusProbe(domains)

	err := probe.Execute()
	if err != nil {
		return cmderrors.NewCmdError(err, 3 /* Nagios unknown */)
	}

	return nil
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("domains", "", "domains to check (comma separated). Defaults to check all domains")

	rootCmd.PersistentFlags().Bool("debug", false, "enable debug")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		rerr := err
		rc := 1
		if cerr, ok := err.(cmderrors.CmdError); ok {
			rerr = cerr.Err
			rc = cerr.Rc
		}

		fmt.Fprintf(os.Stderr, "%s\n", rerr)
		os.Exit(rc)
	}

	os.Exit(0)
}

func initLogging() {
	enableDebug := viper.GetBool("debug")

	if enableDebug {
		log.SetLevel(log.DebugLevel)
	}

	log.SetOutput(os.Stdout)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
