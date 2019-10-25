// Copyright © 2018-2019 Cash Payment Solutions GmbH
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
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type dnsInfo struct {
	returnCode int
	Message    string
}

type cloudflareDNSResponse struct {
	Status int
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
		ip, error := isIPInputValid(args)
		if error != OK {
			log.Println("Unknown: Please specify a correct ip address.")
			os.Exit(UNKNOWN)
		}

		dnsInfoCollector := make(chan *dnsInfo)
		defer close(dnsInfoCollector)

		reversedIPString := reverseIPString(ip)

		isTimeOut := startTimer()
		defer close(isTimeOut)

		dnsBlacklistCheckerCount := len(BlacklistServers)

		for _, blacklistServer := range BlacklistServers {
			go checkIPAgainstBlacklistDomain(
				dnsInfoCollector,
				blacklistServer,
				reversedIPString,
			)
		}

		for {
			select {
			case dnsInfoOutput := <-dnsInfoCollector:
				switch dnsInfoOutput.returnCode {
				case WARNING:
					log.Println("Warning: ", dnsInfoOutput.Message)
					os.Exit(WARNING)
				case UNKNOWN:
					log.Println("Unknown: ", dnsInfoOutput.Message)
					os.Exit(UNKNOWN)
				case CRITICAL:
					if SuppressCrit {
						log.Println("Warning: ", dnsInfoOutput.Message)
						os.Exit(WARNING)
					} else {
						log.Println("Critical: ", dnsInfoOutput.Message)
						os.Exit(CRITICAL)
					}
				case OK:
					dnsBlacklistCheckerCount--
				}
			case <-isTimeOut:
				log.Printf(
					"Warning: Timeout is reached but %d dns blacklist crawler were still working.\n",
					dnsBlacklistCheckerCount,
				)
				os.Exit(WARNING)
			default:
				if dnsBlacklistCheckerCount <= 1 {
					log.Println("Ok: The IP isn't blacklisted.")
					os.Exit(OK)
				}
			}
		}
	},
}

func startTimer() chan bool {
	isTimerOver := make(chan bool, 1)
	go func() {
		time.Sleep(time.Duration(Timeout) * time.Second)
		isTimerOver <- true
	}()
	return isTimerOver
}

func isIPInputValid(input []string) (net.IP, int) {
	if len(input) <= 0 {
		return nil, WARNING
	}
	parsedIP := net.ParseIP(input[0]).To4()
	if parsedIP == nil {
		return nil, WARNING
	}
	return parsedIP, OK
}

func reverseIPString(ip net.IP) string {
	ipComponents := strings.Split(ip.String(), ".")
	reversedIPAddress := fmt.Sprintf(
		"%s.%s.%s.%s",
		ipComponents[3],
		ipComponents[2],
		ipComponents[1],
		ipComponents[0],
	)
	return reversedIPAddress
}

func checkIPAgainstBlacklistDomain(ret chan *dnsInfo, blacklistDomain string, reversedIPAddress string) {
	client := &http.Client{}

	url := fmt.Sprintf(
		"https://cloudflare-dns.com/dns-query?name=%s.%s&type=A",
		reversedIPAddress,
		blacklistDomain,
	)

	dnsReq, err := http.NewRequest("GET", url, nil)

	if err != nil {
		ret <- &dnsInfo{
			WARNING,
			fmt.Sprintf(
				"Failed to create new dns request: %s",
				err.Error(),
			),
		}
	}

	dnsReq.Header.Add("Accept", "application/dns-json")

	dnsResp, err := client.Do(dnsReq)

	if err != nil {
		ret <- &dnsInfo{
			WARNING,
			fmt.Sprintf(
				"The dnssec request failed with: %s",
				err.Error(),
			),
		}
	}

	defer dnsResp.Body.Close()

	respBody, err := ioutil.ReadAll(dnsResp.Body)

	if err != nil {
		ret <- &dnsInfo{
			WARNING,
			fmt.Sprintf(
				"Parsing the dns body failed: %s",
				err.Error(),
			),
		}
	}

	var dnsData cloudflareDNSResponse
	err = json.Unmarshal(respBody, &dnsData)

	if err != nil {
		ret <- &dnsInfo{
			WARNING,
			fmt.Sprintf(
				"Unmarshalling the dns request failed: %s",
				err.Error(),
			),
		}
	}

	switch dnsData.Status {
	case 0:
		ret <- &dnsInfo{
			CRITICAL,
			fmt.Sprintf(
				"%s is listed on the blacklist with domain %s",
				reversedIPAddress,
				blacklistDomain,
			),
		}
	case 2, 3:
		ret <- &dnsInfo{
			OK,
			fmt.Sprintf(
				"%s is not listed on blacklistdomain:%s",
				reversedIPAddress,
				blacklistDomain,
			),
		}
	default:
		ret <- &dnsInfo{
			UNKNOWN,
			fmt.Sprintf(
				"Check the official RCODE's of DNS Requests: %d",
				dnsData.Status,
			),
		}
	}
}

func init() {
	RootCmd.AddCommand(checkCmd)
}
