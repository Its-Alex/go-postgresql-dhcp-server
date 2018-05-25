package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	version = "0.0.0"
)

var (
	rootCmd = &cobra.Command{
		Use:   "dhcp4-reservation",
		Short: "DHCP4 reservation tool",
		Long:  `Reservation tool for ipv4 plugged with postgres`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			var err error
			Db, err = sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
				viper.GetString("psql_user"),
				viper.GetString("psql_password"),
				viper.GetString("psql_addr"),
				viper.GetString("psql_port"),
				viper.GetString("psql_db"),
				viper.GetString("psql_ssl"),
			))
			if err != nil {
				logrus.WithFields(logrus.Fields{
					"type": "database",
				}).Fatal(err)
			}
		},
		Run: dhcpStart,
	}
)

func init() {
	viper.SetEnvPrefix("dhcp4")
	viper.BindEnv("psql_addr")
	viper.BindEnv("psql_port")
	viper.BindEnv("psql_user")
	viper.BindEnv("psql_password")
	viper.BindEnv("psql_databse")
	viper.BindEnv("psql_ssl")

	viper.SetDefault("psql_addr", "localhost")
	viper.SetDefault("psql_port", "5432")
	viper.SetDefault("psql_user", "dhcp4")
	viper.SetDefault("psql_password", "dhcp4")
	viper.SetDefault("psql_databse", "dhcp4")
	viper.SetDefault("psql_ssl", "disable")

	rootCmd.Flags().String("interface", "", "network interface used by server")
	rootCmd.Flags().String("port", "", "port to start server")
	rootCmd.Flags().String("server_ip", "", "ip of server")
	viper.BindPFlag("interface", rootCmd.Flags().Lookup("interface"))
	viper.BindPFlag("server_ip", rootCmd.Flags().Lookup("server_ip"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindEnv("interface")
	viper.BindEnv("server_ip")
	viper.BindEnv("port")

	viper.SetDefault("interface", "en0")
	viper.SetDefault("port", "67")
	viper.SetDefault("server_ip", "192.168.0.254")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
