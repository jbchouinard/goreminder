/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "mxremind",
	Short: "MxRemind mail reminder server and utilities.",
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

var cfgFile string
var jsonLog bool

func init() {
	cobra.OnInitialize(initLogging)
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default \"./mxremind.yaml\")")

	startCmd.Flags().BoolVar(&jsonLog, "jsonlog", false, "output log in JSON format")

	rootCmd.PersistentFlags().String("timezone", "America/Montreal", "timezone location")
	viper.BindPFlag("timezone", rootCmd.PersistentFlags().Lookup("timezone"))

	rootCmd.PersistentFlags().String("db", "", "database connection URL")
	viper.BindPFlag("database.url", rootCmd.PersistentFlags().Lookup("db"))
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		cwd, err := os.Getwd()
		cobra.CheckErr(err)

		viper.AddConfigPath(cwd)
		viper.SetConfigType("yaml")
		viper.SetConfigName("mxremind")
	}

	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("MXREMIND")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Info().Msg("no config file found")
		} else {
			log.Fatal().Err(err).Msg("error reading configuration")
		}
	} else {
		log.Info().Msgf("using config file %s", viper.ConfigFileUsed())
	}

}

func initLogging() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if !jsonLog {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
