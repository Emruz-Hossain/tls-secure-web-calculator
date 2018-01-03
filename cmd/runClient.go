// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
	"io/ioutil"
)

// runClientCmd represents the runClient command
var runClientCmd = &cobra.Command{
	Use:   "runClient",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {

		// parse the flags.
		if len(args)<5{
			log.Fatalln("Argument Missing. Expected Argument: <ca.crt path> <client.crt path> <client.key path> <first operand> <second operand>")
		}
		caCRT	:= args[0]
		clientCRT := args[1]
		clientKey := args[2]

		FirstOperand,err:=strconv.Atoi(args[3])
		if err!=nil{
			log.Fatalln("String to Conversion Error for first argument.")
		}

		SeconOperand,err:=strconv.Atoi(args[4])
		if err!=nil{
			log.Fatalln("String to Conversion Error for second argument.")
		}

		//load client cert
		clientCert,err:=tls.LoadX509KeyPair(clientCRT,clientKey)

		if err!=nil{
			log.Fatalln("Can't load .crt and .key file for client. Reason: %v",err)
		}

		// load CA cert
		caCert, err:= ioutil.ReadFile(caCRT)

		if err!=nil{
			log.Fatalf("Failed to read ca.crt. Reason: %v.",err)
		}

		// create ca cert pool
		caCertPool := x509.NewCertPool()
		ok:=caCertPool.AppendCertsFromPEM(caCert)

		if !ok{
			log.Fatalf("Can't append caCert to caCertPool.")
		}

		// set up tls configuration
		tlsConfig:=&tls.Config{
			Certificates: []tls.Certificate{clientCert},
			RootCAs: caCertPool,
		}

		tlsConfig.BuildNameToCertificate()

		// create a transport for client
		transport:=&http.Transport{TLSClientConfig:tlsConfig}


		//create a client
		client:=http.Client{Transport:transport}

		url:="https://127.0.0.1:8080?FirstOperand="+strconv.Itoa(FirstOperand)+"&"+"SecondOperand="+strconv.Itoa(SeconOperand)

		// create a request
		req,err:=http.NewRequest("GET",url,nil)

		if err!=nil{
			log.Fatalln("Can't create request. Reason: ",err)
		}

		// add basic authentication credential at header
		req.Header.Add("Authorization","Basic "+"ZW1ydXo6MTIzNA==")

		// send GET request to server
		resp,err:=client.Do(req)

		if err!=nil{
			log.Fatalln("Can't send request. Reason: ",err)
		}

		// check if server send a successful response
		if resp.StatusCode==200{

			defer resp.Body.Close()
			bodyText,err:=ioutil.ReadAll(resp.Body)

			if err!=nil{
				log.Fatalln("Error in reading response body. Reason: ",err)
			}

			fmt.Println("Response: ",string(bodyText))
		}else{
			fmt.Println(resp.StatusCode,",",resp.Status)
		}
	},
}

func init() {
	RootCmd.AddCommand(runClientCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runClientCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runClientCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
