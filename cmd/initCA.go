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
	"fmt"
	"os"
	"path/filepath"

	"github.com/appscode/go/log"
	"github.com/appscode/go/term"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
	"github.com/appscode/kutil/tools/certstore"

	"k8s.io/client-go/util/homedir"
)
func getRootDir() string{
	homeDir := homedir.HomeDir()
	rootDir := filepath.Join(homeDir,".wccertificates")
	if _,err := os.Stat(rootDir); err!=nil{
			os.MkdirAll(rootDir,0755)
		}
	return rootDir
}


// initCACmd represents the initCA command
//This create a CA to use for sign server and client.

var initCACmd = &cobra.Command{
	Use:   "initCA",
	Short: "create CA",
	Long: `create CA for self sign`,
	Run: func(cmd *cobra.Command, args []string) {
		log.Infoln(getRootDir())

		//create a certificate store to store certificate
		store,err:=certstore.NewCertStore(afero.NewOsFs(),filepath.Join(getRootDir(),"pki"))

		if err!=nil{
			log.Fatalf("Failed to create certificate store. Reason: %v",err)
		}

		//if ca already exist then ask user if he/she want to overwrite it.
		if store.IsExists("ca"){
			if !term.Ask(fmt.Sprintf("CA certificate found at %s. Do you want to overwrite?",store.Location()),false){
				os.Exit(1)
			}
		}

		// create new CA
		err = store.NewCA()

		if err !=nil{
			log.Fatalf("Failed to create CA. Reason: %v.",err)
		}

		// CA successfully created
		term.Successln("CA successfully generated in",store.Location())
	},
}

func init() {
	RootCmd.AddCommand(initCACmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCACmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCACmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
