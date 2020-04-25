package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	_ "gopkg.in/yaml.v2"

	"github.com/hashicorp/consul/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	logger      = NewLogger("Info", "cmMigration")
	loggerErr   = NewLogger("Error", "cmMigration")
	loggerAudit = GetLogger("Audit")
	configFile  string
	// consul      string
	client     *api.KV
	ADD        map[string]interface{}
	UPDATE     map[string]interface{}
	DELETE     map[string]interface{}
	MIGRATIOIN map[string]interface{}
)

// cmMigrationCmd represents the cmMigration command
var cmMigrationCmd = &cobra.Command{
	Use:   "consulTool",
	Short: "consulTool commad",
	Long: `consulTool will handle central configuration add/update/delete/migration base
	on configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("cmMigration called")
		logger.Println("cm migration call")
		getMigrationData()
		migration()
	},
}

func init() {
	rootCmd.AddCommand(cmMigrationCmd)

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cmMigrationCmd.PersistentFlags().String("foo", "", "A help for foo")

	cmMigrationCmd.Flags().StringVarP(&configFile, "file", "f", "migration.yaml", "Migration yaml file")
	cmMigrationCmd.MarkFlagRequired("file")
	cmMigrationCmd.Flags().StringP("consul", "c", "10.120.115.125:31007", "consul-host:port localhost:8500")
	// cmMigrationCmd.Flags().StringVarP(&consul, "consul", "t", "localhost:8500", "consul-host:port localhost:8500")
	cmMigrationCmd.MarkFlagRequired("consul")
	viper.BindPFlag("consul", cmMigrationCmd.LocalFlags().Lookup("consul"))

}

func getMigrationData() {
	logger.Printf("Read configuration from file %s\n ", configFile)
	config := viper.New()
	fileDir := filepath.Dir(configFile)
	file := strings.Split(configFile, ".")
	config.SetConfigName(file[0])
	config.SetConfigType(file[1])
	// fmt.Printf("fileDir %s", fileDir)
	config.AddConfigPath(fileDir)
	config.AddConfigPath(".")
	if err := config.ReadInConfig(); err != nil {
		logger.Printf("Error reading Configurtion file %s  %s", configFile, err)
	}

	ADD = config.GetStringMap("add")
	logger.Printf("ADD key count %d", len(ADD))
	UPDATE = config.GetStringMap("delete")
	logger.Printf("UPDATE key count %d", len(UPDATE))
	MIGRATIOIN = config.GetStringMap("migration")
	logger.Printf("MIGRATIOIN key count %d", len(MIGRATIOIN))
	DELETE = config.GetStringMap("delete")
	logger.Printf("DELETE key count %d", len(DELETE))
}

func getConsulClient() (*api.KV, error) {
	consulServer := viper.GetString("consul")
	logger.Printf("Initialize consul to %s", consulServer)
	consulConfig := &api.Config{
		Scheme:  "http",
		Address: consulServer,
	}
	client, err := api.NewClient(consulConfig)
	if err != nil {
		logger.Panicf("Can not connect to consul server %s", consulServer)
	}
	// Get a handle to the KV API
	return client.KV(), nil
}

func handleADD(opt *api.QueryOptions) {
	logger.Println("Add key handling:")
	logger.Println(ADD)
	for k, v := range ADD {
		keypair := &api.KVPair{
			Key:   k,
			Value: []byte(v.(string)),
		}
		if keyPare, _, err := client.Get(k, opt); err == nil {
			logger.Printf("key %s exists, skipped", keyPare)
			continue
		}
		_, err := client.Put(keypair, nil)
		if err != nil {
			loggerErr.Panicf("Add key [%s] failed :%s", k, err)
		} else {
			loggerAudit.Printf("Add: %s", k)
		}

	}
}

func handleDELETE(opt *api.QueryOptions) {

}

func handleUpdate(opt *api.QueryOptions) {}

func handleMigration(opt *api.QueryOptions) {}

func migration() {
	client, err = getConsulClient()
	if err != nil {
		logger.Panicln("Consul client create failed")
	}
	queryOpt := &api.QueryOptions{}
	if len(ADD) > 0 {
		handleADD(queryOpt)
	}
}
