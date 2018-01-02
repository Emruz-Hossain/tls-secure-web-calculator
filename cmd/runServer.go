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
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/go/term"
	"github.com/spf13/cobra"
	"strings"
	"strconv"
)


type result struct {
	Sum int
	Sub int
	Mul int
	Div float64
}
type operands struct {
	FirstOperand  int
	SecondOperand int
}
type tests struct {
	F int
	S int
}


// runServerCmd represents the runServer command
var runServerCmd = &cobra.Command{
	Use:   "runServer",
	Short: "start wc-server",
	Long: `start web-calculator server`,
	Run: func(cmd *cobra.Command, args []string) {

		// get location of key file and crt file from commandline argument
		caCRT	:= args[0]
		serverCRT := args[1]
		serverKey := args[2]

		// read ca.crt
		pem, err:= ioutil.ReadFile(caCRT)

		if err!=nil{
			log.Fatalf("Failed to read ca.crt. Reason: %v.",err)
		}

		// create ca cert pool
		caCertPool := x509.NewCertPool()
		ok:=caCertPool.AppendCertsFromPEM(pem)

		if !ok{
			log.Fatalf("Can't append pem to caCertPool.")
		}

		// create tls configuration
		tlsConfig := &tls.Config{

			//Reject any TLS certification that can not be validated
			ClientAuth:	tls.RequireAndVerifyClientCert,

			//Ensure that we only use our "CA" to validate certificate
			ClientCAs:	caCertPool,

			// Use TLS 1.2 rather than 1.1
			MinVersion:	tls.VersionTLS12,

			// Don't allow session resumption
			SessionTicketsDisabled: true,

			// Force to use our preferred Cipher Suites.
			PreferServerCipherSuites: 	true,

			// Specify the cipher suits we prefer.
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305, // Go 1.8 only
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,   // Go 1.8 only
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
				tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			},
		}

		tlsConfig.BuildNameToCertificate()

		// Register handler function for given pattern in DefaultServeMux
		http.HandleFunc("/",handler)

		srv := &http.Server{
			Addr: ":8080",
			ReadTimeout: 5*time.Second,
			WriteTimeout: 10*time.Second,
			TLSConfig:tlsConfig,
		}

		term.Successln("Server is running.....")

		log.Fatalln(srv.ListenAndServeTLS(serverCRT,serverKey))

		//http.HandleFunc("/", handler)
		//http.ListenAndServe(":9000", nil)
	},
}

func init() {
	RootCmd.AddCommand(runServerCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// runServerCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// runServerCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}


func calculate(values operands) result {
	var response result
	response.Sum = values.FirstOperand + values.SecondOperand
	response.Sub = values.FirstOperand - values.SecondOperand
	response.Mul = values.FirstOperand * values.SecondOperand
	response.Div = float64(values.FirstOperand) / float64(values.SecondOperand)
	return response

}
func handler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Recieved new request....")
	// fmt.Println(request)
	if isAuthorised(writer, request) == false {
		writer.Header().Add("WWW-Authenticate", `Basic realm="Authorization Required"`)
		http.Error(writer, "401 Unauthorized", http.StatusUnauthorized)
	} else {
		if request.Method == "GET" {
			var values operands
			//fmt.Println(request.URL.Query())

			A, existA := request.URL.Query()["FirstOperand"]
			if existA {

				valA, errA := strconv.Atoi(A[0])
				if errA != nil {
					http.Error(writer, "FirstOperand not found", http.StatusBadRequest)
					return
				}
				values.FirstOperand = valA
			} else {
				http.Error(writer, "FirstOperand not found", http.StatusBadRequest)
				return
			}
			B, existB := request.URL.Query()["SecondOperand"]
			if existB {
				valB, errB := strconv.Atoi(B[0])
				if errB != nil {
					http.Error(writer, "SecondOperand not found", http.StatusBadRequest)
					return
				}
				values.SecondOperand = valB
			} else {
				http.Error(writer, "SecondOperand not found", http.StatusBadRequest)
				return
			}
			response := calculate(values)
			fmt.Println(response)
			responseJSON, err := json.MarshalIndent(response, "", " ")
			if err != nil {
				http.Error(writer, "Conversion Error", http.StatusInternalServerError)
			} else {
				fmt.Fprintln(writer, string(responseJSON))
			}

		}
		if request.Method == "POST" {
			//fmt.Println("POST method called")
			defer request.Body.Close()
			decoder := json.NewDecoder(request.Body)
			var values operands
			err := decoder.Decode(&values)
			if err != nil {
				fmt.Println(err)
			} else {
				response := calculate(values)
				fmt.Println(response)
				responseJSON, err := json.MarshalIndent(response, "", " ")
				if err != nil {
					http.Error(writer, "Conversion Error", http.StatusInternalServerError)
				} else {
					fmt.Fprintln(writer, string(responseJSON))
				}
			}
		}
	}

}

func isAuthorised(writer http.ResponseWriter, request *http.Request) bool {
	authorizationHeader := strings.SplitN(request.Header.Get("Authorization"), " ", 2)
	fmt.Println(authorizationHeader)
	if len(authorizationHeader) != 2 {
		return false
	}
	baseCredential, err := base64.StdEncoding.DecodeString(authorizationHeader[1])
	if err != nil {
		return false
	} else {
		credential := strings.SplitN(string(baseCredential), ":", 2)
		fmt.Println(credential)
		if credential[0] == "emruz" && credential[1] == "1234" {
			return true
		} else {
			return false
		}
	}
	return false
}
