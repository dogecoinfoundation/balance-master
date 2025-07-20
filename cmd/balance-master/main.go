package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	balance_master "dogecoin.org/balance-master/pkg/balance-master"
	bmcfg "dogecoin.org/balance-master/pkg/config"
)

func main() {
	var dogeSchema, dogeHost, dogePort, dogeUser, dogePassword, databaseURL, migrationsPath, rpcServerHost, rpcServerPort string

	flags := flag.NewFlagSet("balance-master", flag.ExitOnError)
	flags.StringVar(&dogeSchema, "doge-schema", "http", "Dogecoin schema")
	flags.StringVar(&dogeHost, "doge-host", "localhost", "Dogecoin host")
	flags.StringVar(&dogePort, "doge-port", "22556", "Dogecoin port")
	flags.StringVar(&dogeUser, "doge-user", "test", "Dogecoin user")
	flags.StringVar(&dogePassword, "doge-password", "test", "Dogecoin password")
	flags.StringVar(&databaseURL, "database-url", "sqlite://balance-master.db", "Database URL")
	flags.StringVar(&migrationsPath, "migrations-path", "file://../db/migrations", "Migrations path")
	flags.StringVar(&rpcServerHost, "rpc-server-host", "localhost", "RPC server host")
	flags.StringVar(&rpcServerPort, "rpc-server-port", "8899", "RPC server port")

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-signalChan
		os.Exit(0)
	}()

	bm := balance_master.NewBalanceMaster(&bmcfg.Config{
		RpcServerHost:  rpcServerHost,
		RpcServerPort:  rpcServerPort,
		DogeScheme:     dogeSchema,
		DogeHost:       dogeHost,
		DogePort:       dogePort,
		DogeUser:       dogeUser,
		DogePassword:   dogePassword,
		DatabaseURL:    databaseURL,
		MigrationsPath: migrationsPath,
	})

	bm.Start()
}
