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
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var Timeout int
var SuppressCrit bool
var BlacklistServers = []string{
	"all.s5h.net",
	"b.barracudacentral.org",
	"bl.emailbasura.org",
	"bl.spamcannibal.org",
	"bl.spamcop.net",
	"blacklist.woody.ch",
	"bogons.cymru.com",
	"cbl.abuseat.org",
	"cdl.anti-spam.org.cn",
	"combined.abuse.ch",
	"db.wpbl.info",
	"dnsbl-1.uceprotect.net",
	"dnsbl-2.uceprotect.net",
	"dnsbl-3.uceprotect.net",
	"dnsbl.anticaptcha.net",
	"dnsbl.dronebl.org",
	"dnsbl.inps.de",
	"dnsbl.sorbs.net",
	"dnsbl.spfbl.net",
	"drone.abuse.ch",
	"duinv.aupads.org",
	"dul.dnsbl.sorbs.net",
	"dyna.spamrats.com",
	"dynip.rothen.com",
	"http.dnsbl.sorbs.net",
	"ips.backscatterer.org",
	"ix.dnsbl.manitu.net",
	"korea.services.net",
	"misc.dnsbl.sorbs.net",
	"noptr.spamrats.com",
	"orvedb.aupads.org",
	"pbl.spamhaus.org",
	"proxy.bl.gweep.ca",
	"psbl.surriel.com",
	"relays.bl.gweep.ca",
	"relays.nether.net",
	"sbl.spamhaus.org",
	"short.rbl.jp",
	"singular.ttk.pte.hu",
	"smtp.dnsbl.sorbs.net",
	"socks.dnsbl.sorbs.net",
	"spam.abuse.ch",
	"spam.dnsbl.anonmails.de",
	"spam.dnsbl.sorbs.net",
	"spam.spamrats.com",
	"spambot.bls.digibase.ca",
	"spamrbl.imp.ch",
	"spamsources.fabel.dk",
	"ubl.lashback.com",
	"ubl.unsubscore.com",
	"virus.rbl.jp",
	"web.dnsbl.sorbs.net",
	"wormrbl.imp.ch",
	"xbl.spamhaus.org",
	"z.mailspike.net",
	"zen.spamhaus.org",
	"zombie.dnsbl.sorbs.net",
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
		log.Println(err)
		os.Exit(UNKNOWN)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.nagios-dnsblklist.yaml)")
	RootCmd.PersistentFlags().IntVarP(&Timeout, "timeout", "t", 30, "Pick a timeout in seconds")
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
		SuppressCrit = viper.GetBool("suppresscrit")
	}
}
