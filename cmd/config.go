package cmd

import (
	"fmt"
	"log"

	"github.com/jbchouinard/mxremind/pkg/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func init() {
	rootCmd.AddCommand(configCmd)
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Print current configuration in YAML format",
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		bytes, err := yaml.Marshal(conf)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s", bytes)
	}}
