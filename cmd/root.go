// Copyright © 2016 David Leib <david.leib@barzahlen.de>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"fmt"
	"os"

	"github.com/go-playground/log"
	"github.com/go-playground/log/handlers/console"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/Barzahlen/nagios-dnsblklist/logger"
)

var cfgFile string
var Verbosity int
var Timeout int
var SuppressCrit bool
var BlacklistServers = []string{
	"bl.spamcop.net",
	"cbl.abuseat.org",
	"b.barracudacentral.org",
	"dnsbl.sorbs.net",
	"http.dnsbl.sorbs.net",
	"dul.dnsbl.sorbs.net",
	"misc.dnsbl.sorbs.net",
	"smtp.dnsbl.sorbs.net",
	"socks.dnsbl.sorbs.net",
	"spam.dnsbl.sorbs.net",
	"web.dnsbl.sorbs.net",
	"zombie.dnsbl.sorbs.net",
	"dnsbl-1.uceprotect.net",
	"dnsbl-2.uceprotect.net",
	"dnsbl-3.uceprotect.net",
	"pbl.spamhaus.org",
	"sbl.spamhaus.org",
	"xbl.spamhaus.org",
	"zen.spamhaus.org",
	"psbl.surriel.com",
	"ubl.unsubscore.com",
	"dnsbl.njabl.org",
	"combined.njabl.org",
	"rbl.spamlab.com",
	"dyna.spamrats.com",
	"noptr.spamrats.com",
	"spam.spamrats.com",
	"cbl.anti-spam.org.cn",
	"cdl.anti-spam.org.cn",
	"dnsbl.inps.de",
	"drone.abuse.ch",
	"httpbl.abuse.ch",
	"korea.services.net",
	"short.rbl.jp",
	"virus.rbl.jp",
	"spamrbl.imp.ch",
	"wormrbl.imp.ch",
	"virbl.bit.nl",
	"dsn.rfc-ignorant.org",
	"ips.backscatterer.org",
	"spamguard.leadmon.net",
	"opm.tornevall.org",
	"netblock.pedantic.org",
	"multi.surbl.org",
	"ix.dnsbl.manitu.net",
	"tor.dan.me.uk",
	"relays.mail-abuse.org",
	"blackholes.mail-abuse.org",
	"rbl-plus.mail-abuse.org",
	"dnsbl.dronebl.org",
	"access.redhawk.org",
	"db.wpbl.info",
	"rbl.interserver.net",
	"query.senderbase.org",
	"bogons.cymru.com",
	"csi.cloudmark.com",
}

const OK = 0
const WARNING = 1
const CRITICAL = 2
const UNKNOWN = 3

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "nagios-dnsblklist",
	Short: "This tool is checking if a specific ip was blacklisted",
	Long: `This tool is performing a reverse dns lookup to check if the supplied
ip-address is listed on a blacklist server.

The ip-address is checked against a list of all relevant blacklist servers.

Current Version: 1.0.0`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nagios-dnsblklist.yaml)")
	RootCmd.PersistentFlags().IntVarP(&Timeout, "timeout", "t", 30, "Pick a timeout in seconds")
	RootCmd.PersistentFlags().IntVarP(&Verbosity, "verbosity", "v", 0, "Pick a verbosity level 0 = normal (default) and 1 = debug")
	RootCmd.PersistentFlags().BoolVarP(&SuppressCrit, "suppresscrit", "s", false,
		"Suppress critical message from the system and send warning instead.")
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName(".nagios-dnsblklist") // name of config file (without extension)
	viper.SetConfigType("yaml")               // default config file type is yaml
	viper.AddConfigPath("$HOME")              // adding home directory as first search path
	viper.AutomaticEnv()                      // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		BlacklistServers = viper.GetStringSlice("blacklistServers")
		Timeout = viper.GetInt("timeout")
		Verbosity = viper.GetInt("verbosity")
		SuppressCrit = viper.GetBool("suppresscrit")
	}

	consoleLog := console.New()
	consoleLog.SetFormatFunc(logger.NagiosFormatFunc)
	consoleLog.SetDisplayColor(false)

	switch Verbosity {
	case 0:
		log.RegisterHandler(consoleLog, log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.AlertLevel,
			log.PanicLevel, log.FatalLevel)
	case 1:
		log.RegisterHandler(consoleLog, log.AllLevels...)
	default:
		log.RegisterHandler(consoleLog, log.InfoLevel, log.WarnLevel, log.ErrorLevel, log.PanicLevel,
			log.FatalLevel)
	}
}
