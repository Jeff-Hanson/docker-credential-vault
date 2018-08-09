// Copyright © 2017 Lee Briggs <lee@leebriggs.co.uk>
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
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"bytes"

	"github.com/mitchellh/mapstructure"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	v "github.com/Jeff-Hanson/docker-credential-vault/vault"
)

type DockerLoginCredentials struct {
	ServerURL string `mapstructure:"serverurl"`
	Username  string `mapstructure:"username"`
	Secret    string `mapstructure:"password"`
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Get credentials from vault",
	Run: func(cmd *cobra.Command, args []string) {

		// read in stdin
		urlBytes, err := ioutil.ReadAll(os.Stdin)

		// if stdin can't be read, bomb
		if err != nil {
			Debug.Println("Error reading stdin", err)
			return
		}

		Debug.Printf("Input for Get: %s", string(bytes.TrimSpace(urlBytes)))

		// create a base64 encoded url from stdin
		url := b64.StdEncoding.EncodeToString(bytes.TrimSpace(urlBytes))

		Debug.Printf("Base 64 of input: %s", url)

		vaultHost = viper.GetString("vault")
		vaultToken = viper.GetString("token")

		// create a vault client
		client, err := v.New(vaultToken, vaultHost, vaultPort)

		if err != nil {
			Debug.Println("Error creating vault client", err)
			return
		}

		secret, err := client.Logical().Read("secret/" + url)

		if err != nil {
			Debug.Println("Error reading vault creds", err)
			return
		}

		if secret == nil {
			Debug.Printf("Error reading secret from vault: %s / %s", string(bytes.TrimSpace(urlBytes)), url)
			return
		}

		var creds DockerLoginCredentials

		if err := mapstructure.Decode(secret.Data, &creds); err != nil {
			Debug.Println("Error parsing vault response: ", err)
			return
		}

		jsonCreds, _ := json.Marshal(creds)

		Debug.Printf("Creds being returned: %s", string(jsonCreds))
		fmt.Printf(string(jsonCreds))
	},
}

func init() {
	RootCmd.AddCommand(getCmd)
}
