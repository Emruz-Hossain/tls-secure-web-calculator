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
	"log"
	"net"
	"os"
	"path/filepath"

	"github.com/appscode/go/term"
	"github.com/appscode/kutil/tools/certstore"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"k8s.io/client-go/util/cert"
)

// generateServerCertificateCmd represents the generateServerCertificate command
var generateServerCertificateCmd = &cobra.Command{
	Use:   "generateServerCertificate",
	Short: "generate certificate for server",
	Long: `This command generate certificate for server`,
	Run: func(cmd *cobra.Command, args []string) {
		sans:=cert.AltNames{
			IPs: []net.IP{net.ParseIP("127.0.0.1")},
		}
		cfg:=cert.Config{
			CommonName: "server",
			AltNames:	sans,
			Usages:	[]x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		}

		// create certificate store to store certificate of server
		store,err:=certstore.NewCertStore(afero.NewOsFs(),filepath.Join(getRootDir(),"pki"),cfg.Organization...)

		if err!=nil{
			log.Fatalf("Failed to create certificate store. Reason: %v",err)
		}

		//check if certificate already exist. If exist ask user if he/she want to overwrite.
		if store.IsExists(filename(cfg)){
			if !term.Ask(fmt.Sprintf("Server certificate found at %s. Do you want to overwrite?", store.Location()), false){
				os.Exit(1)
			}

		}
		// load CA cert and key
		err = store.LoadCA()

		if err!=nil{
			log.Fatalf("Failed to load ca certificate. Reason: %v.",err)
		}

		// generate certificate pair for server
		crt, key, err := store.NewServerCertPair(cfg.CommonName, cfg.AltNames)

		if err!=nil{
			log.Fatalf("Failed to generate certificate pair for server. Reason: %v.",err)
		}

		// write this certificate pair in file
		err = store.WriteBytes(filename(cfg),crt,key)

		if err!=nil{
			log.Fatalf("Failed to write certificate pair in file. Reason: %v.",err)
		}

		// certificate for server successfully create
		term.Successln("Certificate for server is successfully created in ",store.Location())


	},
}

func init() {
	RootCmd.AddCommand(generateServerCertificateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// generateServerCertificateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// generateServerCertificateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
