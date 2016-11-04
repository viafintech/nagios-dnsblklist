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
	"net"
	"os"
	"strings"
	"time"

	"github.com/go-playground/log"
	"github.com/spf13/cobra"
)

type dnsInfo struct {
	returnCode int
	Message    string
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Expects an ip-address(ipv4) to check if it is blacklisted.[127.0.0.1]",
	Long: `Checks the supplied ip-address(ipv4) and returns:
* 0: not blacklisted
* 1: a blacklist server timed out or timeout reached
* 2: it was found on a blacklist server
* 3: an unknown error occured`,
	Run: func(cmd *cobra.Command, args []string) {
		ip, error := isIpInputValid(args)
		if error != OK {
			log.Alert("Please specify a correct ip address.")
			os.Exit(UNKNOWN)
		}

		dnsInfoCollector := make(chan *dnsInfo)
		defer close(dnsInfoCollector)

		reversedIpString := reverseIpString(ip)

		isTimeOut := startTimer()
		defer close(isTimeOut)

		dnsBlacklistCheckerCount := len(BlacklistServers)
		log.Debug("Checking ", dnsBlacklistCheckerCount, " blacklist servers.")

		for _, blacklistServer := range BlacklistServers {
			go checkIpAgainstBlacklistDomain(dnsInfoCollector, blacklistServer, reversedIpString)
		}

		for {
			select {
			case dnsInfoOutput := <-dnsInfoCollector:
				switch dnsInfoOutput.returnCode {
				case WARNING:
					log.Warn(dnsInfoOutput.Message)
					os.Exit(WARNING)
				case UNKNOWN:
					log.Alert(dnsInfoOutput.Message)
					os.Exit(UNKNOWN)
				case CRITICAL:
					if SuppressCrit {
						log.Warn(dnsInfoOutput.Message)
						os.Exit(WARNING)
					} else {
						log.Error(dnsInfoOutput.Message)
						os.Exit(CRITICAL)
					}
				case OK:
					dnsBlacklistCheckerCount -= 1
				}
			case <-isTimeOut:
				log.Warn(fmt.Sprintf("Timeout is reached but %d dns blacklist crawler were still working.", dnsBlacklistCheckerCount))
				os.Exit(WARNING)
			default:
				if dnsBlacklistCheckerCount <= 1 {
					log.Info("The IP isn't blacklisted.")
					os.Exit(OK)
				}
			}
		}
	},
}

func startTimer() chan bool {
	isTimerOver := make(chan bool, 1)
	log.Debug("Starting timer with ", Timeout, " seconds.")
	go func() {
		time.Sleep(time.Duration(Timeout) * time.Second)
		log.Debug("The Timer reached it's end after ", Timeout, " and informs the main proc.")
		isTimerOver <- true
	}()
	return isTimerOver
}

func isIpInputValid(input []string) (net.IP, int) {
	if len(input) <= 0 {
		log.Debug("An argument wasn't supplied.")
		return nil, WARNING
	}
	parsedIp := net.ParseIP(input[0]).To4()
	if parsedIp == nil {
		log.Debug("The IP address turned out to be not valid.")
		return nil, WARNING
	}
	return parsedIp, OK
}

func reverseIpString(ip net.IP) string {
	ip_components := strings.Split(ip.String(), ".")
	reversedIpAddress := fmt.Sprintf("%s.%s.%s.%s", ip_components[3], ip_components[2], ip_components[1],
		ip_components[0])
	log.Debug("Reversed Ip Address for testing is: ", reversedIpAddress)
	return reversedIpAddress
}

func checkIpAgainstBlacklistDomain(ret chan *dnsInfo, blacklistDomain string, reversedIpAddress string) {
	log.Debug("Checking ", reversedIpAddress, " against ", blacklistDomain)
	ns_records, err := net.LookupHost(fmt.Sprintf("%s.%s", reversedIpAddress, blacklistDomain))

	log.Debug(reversedIpAddress, ".", blacklistDomain, ": Lookup host finished. Starting to interpret outcome.")
	nerr, ok := err.(*net.DNSError)

	if ok && nerr.Err == "no such host" && len(ns_records) == 0 {
		ret <- &dnsInfo{OK, fmt.Sprintf("%s is not listed on blacklistdomain:%s", reversedIpAddress, blacklistDomain)}
	} else if ok && nerr.Timeout() {
		ret <- &dnsInfo{WARNING, err.Error()}
	} else if ok && nerr.Temporary() {
		ret <- &dnsInfo{UNKNOWN, fmt.Sprintf("A temporary failure was detected with error: %s", err.Error())}
	} else {
		ret <- &dnsInfo{CRITICAL, fmt.Sprintf("%s is listed on the blacklist with domain %s by %s",
			reversedIpAddress, blacklistDomain, ns_records)}
	}

	log.Debug("Checking ", reversedIpAddress, ".", blacklistDomain, " finished.")
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
