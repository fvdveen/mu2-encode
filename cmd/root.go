package cmd

import (
	"strings"

	"github.com/fvdveen/mu2-encode/encode"
	encodepb "github.com/fvdveen/mu2-proto/go/proto/encode"
	"github.com/hashicorp/consul/api"
	"github.com/micro/go-grpc"
	"github.com/micro/go-micro"
	"github.com/micro/go-micro/registry/consul"
	"github.com/micro/go-micro/transport"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logLvl string
	conf   struct {
		Consul struct {
			Address string `mapstructure:"address"`
		} `mapstructure:"consul"`
		Config struct {
			Path string `mapstructure:"path"`
			Type string `mapstructure:"type"`
		} `mapstructure:"config"`
	}
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "mu2-encode",
	Short: "Mu2 encode service",
	Long: `Mu2 is a discord music bot. 

This is the encode service for mu2.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cc := api.DefaultConfig()
		if conf.Consul.Address != "" {
			cc.Address = conf.Consul.Address
		}

		srv := grpc.NewService(
			micro.Name("mu2.service.encode"),
			micro.Version("latest"),
			micro.Registry(consul.NewRegistry(consul.Config(cc))),
			micro.Transport(
				transport.NewTransport(
					transport.Secure(true),
				),
			),
		)

		s := encode.NewService()

		encodepb.RegisterEncodeServiceHandler(srv.Server(), s)

		return srv.Run()
	},
	SilenceUsage: true,
}

// Execute runs the cli
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.SetEnvPrefix("MU2")

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&conf.Consul.Address, "consul-addr", "", "consul address")
	rootCmd.PersistentFlags().StringVar(&conf.Config.Path, "config-path", "encode/config", "config path on the kv store")
	rootCmd.PersistentFlags().StringVar(&conf.Config.Type, "config-type", "json", "config type on the kv store")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	viper.AutomaticEnv() // read in environment variables that match

	for _, key := range viper.AllKeys() {
		val := viper.Get(key)
		viper.Set(key, val)
	}

	if err := viper.Unmarshal(&conf); err != nil {
		logrus.WithField("type", "main").Fatalf("Unmarshalling config: %v", err)
		return
	}

	var lvl logrus.Level

	lvl = logrus.DebugLevel

	logrus.SetLevel(lvl)
}
