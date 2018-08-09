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
	"fmt"
	"os"
	"strconv"
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var vaultHost string
var vaultPort int
var vaultToken string

// Debug log handle
var Debug *log.Logger

var logFile *os.File

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "docker-credential-vault",
	Short: "A credential helper for vault",
	Long:  `Stores, retrieves and erases your docker registry credentials from Hashcorp Vault`,
}

// Execute the command 
func Execute() {
	defer logFile.Close()
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	logFile, err := os.OpenFile("/tmp/vault-creds.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open log file", err)
	}
	Debug = log.New(logFile, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug.Println("Initializing...")

	cobra.OnInitialize(initConfig)
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.docker-credential-vault.yaml)")
	RootCmd.PersistentFlags().StringVarP(&vaultHost, "vault", "v", "vault.service.discover", "vault host to authenticate against")
	RootCmd.PersistentFlags().IntVarP(&vaultPort, "port", "p", 8200, "port of vault server to authenticate against")
	RootCmd.PersistentFlags().StringVarP(&vaultToken, "token", "t", "", "vault token to authenticate with")
	viper.BindPFlag("vault", RootCmd.PersistentFlags().Lookup("vault"))
	viper.BindPFlag("token", RootCmd.PersistentFlags().Lookup("token"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" { // enable ability to specify config file via flag
		viper.SetConfigFile(cfgFile)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath("/etc/docker/vault-credential-helper")
	viper.AddConfigPath("$HOME/.docker/vault-credential-helper")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		Debug.Println("Error reading config file", err)
	}
	
	v := viper.GetViper()
	Debug.Printf("Read in config: %v", v.AllSettings())

	var err error
	vaultHost = viper.Get("vault").(string)
	vaultPort, err = strconv.Atoi(viper.Get("port").(string))
	if err != nil {
		fmt.Println("Unable to convert port to int type. Value:", viper.Get("port").(string))
		return
	}
	vaultToken = viper.Get("token").(string)
	Debug.Printf("Vault params: Host-%s Port-%v Token-%v", vaultHost, vaultPort, vaultToken)

}
