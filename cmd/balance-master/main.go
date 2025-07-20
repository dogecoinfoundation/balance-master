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

	flag.StringVar(&dogeSchema, "doge-schema", "http", "Dogecoin schema")
	flag.StringVar(&dogeHost, "doge-host", "localhost", "Dogecoin host")
	flag.StringVar(&dogePort, "doge-port", "22556", "Dogecoin port")
	flag.StringVar(&dogeUser, "doge-user", "test", "Dogecoin user")
	flag.StringVar(&dogePassword, "doge-password", "test", "Dogecoin password")
	flag.StringVar(&databaseURL, "database-url", "sqlite://balance-master.db", "Database URL")
	flag.StringVar(&migrationsPath, "migrations-path", "db/migrations", "Migrations Path")
	flag.StringVar(&rpcServerHost, "rpc-server-host", "localhost", "RPC server host")
	flag.StringVar(&rpcServerPort, "rpc-server-port", "8899", "RPC server port")

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
