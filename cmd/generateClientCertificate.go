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
	"crypto/x509"
	"fmt"
	"os"
	"path/filepath"

	"github.com/appscode/go/log"
	"github.com/appscode/go/term"
	"github.com/appscode/kutil/tools/certstore"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/cert"
)

// generateClientCertificateCmd represents the generateClientCertificate command
var generateClientCertificateCmd = &cobra.Command{
	Use:   "generateClientCertificate",
	Short: "generate certificate for client",
	Long: `It generate certificate for client`,
	Run: func(cmd *cobra.Command, args []string) {

		//check if client name is sent as flag
		if len(args)==0{
			log.Fatalln("Missing client name.")
		}
		if len(args)>1 {
			log.Fatalln("Multiple client name found")
		}

		cfg:=cert.Config{
			CommonName: args[0],
			Usages: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth},
		}

		// create certificate store to store client certificate
		store,err:=certstore.NewCertStore(afero.NewOsFs(),filepath.Join(getRootDir(),"pki"),cfg.Organization...)

		if err!=nil{
			log.Fatalf("Failed to create certificate store. Reason: %v",err)
		}

		//check if certificate already exist. If exist ask user if he/she want to overwrite.
		if store.IsExists(filename(cfg)){
			if !term.Ask(fmt.Sprintf("Client certificate found at %s. Do you want to overwrite?", store.Location()), false){
				os.Exit(1)
			}

		}
		// load CA cert and key
		err = store.LoadCA()

		if err!=nil{
			log.Fatalf("Failed to load ca certificate. Reason: %v.",err)
		}

		// generate certificate pair for server
		crt, key, err := store.NewClientCertPair(cfg.CommonName,cfg.Organization...)

		if err!=nil{
			log.Fatalf("Failed to generate certificate pair for client. Reason: %v.",err)
		}

		// write this certificate pair in file
		err = store.WriteBytes(filename(cfg),crt,key)

		if err!=nil{
			log.Fatalf("Failed to write certificate pair in file. Reason: %v.",err)
		}

		// certificate for server successfully create
		term.Successln("Certificate for client is successfully created in ",store.Location())
	},
}

func init() {
	RootCmd.AddCommand(generateClientCertificateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateClientCertificateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateClientCertificateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
