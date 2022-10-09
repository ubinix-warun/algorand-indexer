package main

import (
	"fmt"
	// "io"
	"os"
	// "strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	// "github.com/spf13/cobra/doc"
	// "github.com/spf13/viper"

	// bg "github.com/algorand/indexer/cmd/block-generator/core"
	// v "github.com/algorand/indexer/cmd/validator/core"
	"github.com/algorand/indexer/config"
	"github.com/algorand/indexer/idb"
	"github.com/algorand/indexer/idb/dummy"
	// _ "github.com/algorand/indexer/idb/postgres"
	// _ "github.com/algorand/indexer/util/disabledeadlock"
	// "github.com/algorand/indexer/util/metrics"
	// "github.com/algorand/indexer/version"

)

const autoLoadIndexerConfigFileName = config.FileName
const autoLoadParameterConfigFileName = "api_config"

// Calling os.Exit() directly will not honor any defer'd statements.
// Instead, we will create an exit type and handler so that we may panic
// and handle any exit specific errors
type exit struct {
	RC int // The exit code
}

// exitHandler will handle a panic with type of exit (see above)
func exitHandler() {
	if err := recover(); err != nil {
		if exit, ok := err.(exit); ok {
			os.Exit(exit.RC)
		}

		// It's not actually an exit type, restore panic
		panic(err)
	}
}

// Requires that main (and every go-routine that this is used)
// have defer exitHandler() called first
func maybeFail(err error, errfmt string, params ...interface{}) {
	if err == nil {
		return
	}
	logger.WithError(err).Errorf(errfmt, params...)
	panic(exit{1})
}

var rootCmd = &cobra.Command{
	Use:   "webhook",
	Short: "Algorand Webhook",
	Long:  `Webhook ...`,
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		//If no arguments passed, we should fallback to help
		cmd.HelpFunc()(cmd, args)
	},
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if doVersion {
			fmt.Printf("%s\n", "0.0.1")
			panic(exit{0})
		}
	},
}

var (
	// postgresAddr   string
	dummyIndexerDb bool
	doVersion       bool
	logLevel       string
	logFile        string
	logger         *log.Logger
)

func indexerDbFromFlags(opts idb.IndexerDbOptions) (idb.IndexerDb, chan struct{}) {
	// if postgresAddr != "" {
	// 	db, ch, err := idb.IndexerDbByName("postgres", postgresAddr, opts, logger)
	// 	maybeFail(err, "could not init db, %v", err)
	// 	return db, ch
	// }
	if dummyIndexerDb {
		return dummy.IndexerDb(), nil
	}
	logger.Errorf("no import db set")
	panic(exit{1})
}

func init() {
	

	logger = log.New()
	logger.SetFormatter(&log.JSONFormatter{
		DisableHTMLEscape: true,
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(log.InfoLevel)

	daemonCmd := DaemonCmd()
	rootCmd.AddCommand(daemonCmd)

	// Version should be available globally
	rootCmd.Flags().BoolVarP(&doVersion, "version", "v", false, "print version and exit")
	addFlags := func(cmd *cobra.Command) {
		cmd.Flags().StringVarP(&logLevel, "loglevel", "l", "info", "verbosity of logs: [error, warn, info, debug, trace]")
		cmd.Flags().StringVarP(&logFile, "logfile", "f", "", "file to write logs to, if unset logs are written to standard out")
		cmd.Flags().BoolVarP(&dummyIndexerDb, "dummydb", "n", false, "use dummy indexer db")
		cmd.Flags().BoolVarP(&doVersion, "version", "v", false, "print version and exit")
	}
	addFlags(daemonCmd)

}

func configureLogger() error {
	if logLevel != "" {
		level, err := log.ParseLevel(logLevel)
		if err != nil {
			return err
		}
		logger.SetLevel(level)
	}

	if logFile == "-" {
		logger.SetOutput(os.Stdout)
	} else if logFile != "" {
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		logger.SetOutput(f)
	}

	return nil
}


func main() {

	// Setup our exit handler for maybeFail() and other exit panics
	defer exitHandler()

	if err := rootCmd.Execute(); err != nil {
		logger.WithError(err).Error("an error occurred running webhook")
		os.Exit(1)
	}
	
}