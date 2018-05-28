package cmd

import (
	"fmt"

	"github.com/Its-Alex/go-postgresql-dhcp-server/database"
	"github.com/Its-Alex/go-postgresql-dhcp-server/dhcp"
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
		Use:   "go-postgresql-dhcp-server",
		Short: "DHCP4 reservation tool",
		Long:  `Reservation tool for ipv4 plugged with postgres`,
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			// Set logrus format to json
			logrus.SetFormatter(&logrus.JSONFormatter{})
			if viper.IsSet("verbose") {
				logrus.SetLevel(logrus.DebugLevel)
			}

			var err error
			database.Db, err = sqlx.Connect("postgres", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
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
		Run: func(cmd *cobra.Command, args []string) {
			dhcp.Start()
		},
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
	rootCmd.Flags().StringP("verbose", "v", "", "set verbosity to debug")
	viper.BindPFlag("interface", rootCmd.Flags().Lookup("interface"))
	viper.BindPFlag("server_ip", rootCmd.Flags().Lookup("server_ip"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("verbose", rootCmd.Flags().Lookup("verbose"))
	viper.BindEnv("interface")
	viper.BindEnv("server_ip")
	viper.BindEnv("port")
	viper.BindEnv("verbose")

	viper.SetDefault("interface", "en0")
	viper.SetDefault("port", "67")
	viper.SetDefault("server_ip", "192.168.0.254")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
