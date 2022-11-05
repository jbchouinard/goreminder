package config

import (
	"crypto/tls"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("database.migrate", true)
	viper.SetDefault("smtp.port", 587)
	viper.SetDefault("smtp.authenticated", true)
	viper.SetDefault("smtp.tls.enabled", true)
	viper.SetDefault("smtp.tls.insecure", false)
	viper.SetDefault("imap.port", 993)
	viper.SetDefault("imap.tls.enabled", true)
	viper.SetDefault("imap.tls.insecure", false)
	viper.SetDefault("imap.authenticated", true)
}

func assertKeys(required []string) {
	for _, key := range required {
		if !viper.IsSet(key) {
			log.Fatal().Msgf("missing configuration key %q", key)
		}
	}
}

type TlsConfig struct {
	Enabled  bool `yaml:"enabled"`
	Insecure bool `yaml:"insecure"`
}

type ServerConfig struct {
	Address       string     `yaml:"address"`
	Password      string     `yaml:"password"`
	Authenticated bool       `yaml:"authenticated"`
	Host          string     `yaml:"host"`
	Port          uint16     `yaml:"port"`
	Tls           *TlsConfig `yaml:"tls"`
}

func (sc *ServerConfig) TlsConfig() *tls.Config {
	return &tls.Config{
		InsecureSkipVerify: sc.Tls.Insecure,
		ServerName:         sc.Host,
	}
}

func GetServerConfig(prefix string) *ServerConfig {
	addressKey := prefix + ".address"
	passwordKey := prefix + ".password"
	authenticatedKey := prefix + ".authenticated"
	hostKey := prefix + ".host"
	portKey := prefix + ".port"
	tlsKey := prefix + ".tls.enabled"
	insecureKey := prefix + ".tls.insecure"
	assertKeys([]string{addressKey, passwordKey, authenticatedKey, hostKey, portKey, tlsKey, insecureKey})
	return &ServerConfig{
		Address:       viper.GetString(addressKey),
		Password:      viper.GetString(passwordKey),
		Authenticated: viper.GetBool(authenticatedKey),
		Host:          viper.GetString(hostKey),
		Port:          viper.GetUint16(portKey),
		Tls: &TlsConfig{
			Enabled:  viper.GetBool(tlsKey),
			Insecure: viper.GetBool(insecureKey),
		},
	}
}

type MailboxConfig struct {
	In        string `yaml:"in"`
	Processed string `yaml:"processed"`
}

func GetMailboxConfig(prefix string) *MailboxConfig {
	inKey := prefix + ".in"
	processedKey := prefix + ".processed"
	assertKeys([]string{inKey, processedKey})
	return &MailboxConfig{
		In:        viper.GetString(inKey),
		Processed: viper.GetString(processedKey),
	}
}

type DatabaseConfig struct {
	Migrate bool   `yaml:"migrate"`
	URL     string `yaml:"url"`
}

func GetDatabaseConfig(prefix string) *DatabaseConfig {
	migrateKey := prefix + ".migrate"
	urlKey := prefix + ".url"
	assertKeys([]string{migrateKey, urlKey})
	return &DatabaseConfig{
		Migrate: viper.GetBool(migrateKey),
		URL:     viper.GetString(urlKey),
	}
}

type Config struct {
	Timezone      string          `yaml:"timezone"`
	SendInterval  uint16          `yaml:"send_interval"`
	FetchInterval uint16          `yaml:"fetch_interval"`
	Database      *DatabaseConfig `yaml:"database"`
	Mailbox       *MailboxConfig  `yaml:"mailbox"`
	IMAP          *ServerConfig   `yaml:"imap"`
	SMTP          *ServerConfig   `yaml:"smtp"`
}

func (conf *Config) Location() *time.Location {
	loc, err := time.LoadLocation(conf.Timezone)
	if err != nil {
		log.Fatal().Err(err).Msg("error loading location")
	}
	return loc
}

func GetConfig() *Config {
	assertKeys([]string{"timezone", "send_interval", "fetch_interval"})
	return &Config{
		Timezone:      viper.GetString("timezone"),
		SendInterval:  viper.GetUint16("send_interval"),
		FetchInterval: viper.GetUint16("fetch_interval"),
		Database:      GetDatabaseConfig("database"),
		Mailbox:       GetMailboxConfig("mailbox"),
		SMTP:          GetServerConfig("smtp"),
		IMAP:          GetServerConfig("imap"),
	}
}

func (mc *Config) Describe() string {
	return fmt.Sprintf("%s:%s", mc.IMAP.Address, mc.Mailbox.In)
}

func (mc *Config) Log(message string) {
	log.Info().Msgf("%s: %s", mc.Describe(), message)
}

func (mc *Config) Logf(format string, p ...any) {
	mc.Log(fmt.Sprintf(format, p...))
}
