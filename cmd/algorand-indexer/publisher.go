package main

import (
	// "context"
	"fmt"
	// "os"
	// "os/signal"
	// "strings"
	// "sync"
	// "syscall"
	// "time"

	// "github.com/algorand/go-algorand/rpcs"
	// "github.com/algorand/go-algorand/data/transactions"
	// "github.com/spf13/cobra"
	// "github.com/spf13/viper"

	// "github.com/algorand/indexer/config"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	upgrader = websocket.Upgrader{}
)

// var publisherCmd = &cobra.Command{
// 	Use:   "publisher",
// 	Short: "run publisher",
// 	Long:  "run publisher. Serve block on WS.",
// 	//Args:
// 	Run: func(cmd *cobra.Command, args []string) {
// 		var err error
// 		config.BindFlags(cmd)
// 		err = configureLogger()
// 		if err != nil {
// 			fmt.Fprintf(os.Stderr, "failed to configure logger: %v", err)
// 			os.Exit(1)
// 		}

// 		// ctx, cf := context.WithCancel(context.Background())
// 		// defer cf()
// 		// {
// 		// 	cancelCh := make(chan os.Signal, 1)
// 		// 	signal.Notify(cancelCh, syscall.SIGTERM, syscall.SIGINT)
// 		// 	go func() {
// 		// 		<-cancelCh
// 		// 		logger.Println("Stopping Publisher.")
// 		// 		cf()
// 		// 	}()
// 		// }

// 		// var wg sync.WaitGroup
// 		// bot!
// 		// if bot != nil {
// 		// 	wg.Add(1)
// 		// 	go func() {
// 		// 		defer wg.Done()

// 		// 		// Wait until the something is available.
// 		// 		<-availableCh

// 		// } else {
// 		// 	logger.Info("No block importer configured.")
// 		// }

// 		e := echo.New()
// 		e.Use(middleware.Logger())
// 		e.Use(middleware.Recover())
// 		// e.Static("/", "../public")
// 		e.GET("/ws", hello)
	
// 		fmt.Printf("serving ws on %s\n", ":1323")
// 		logger.Infof("serving ws on %s", ":1323")
	
// 		e.Logger.Fatal(e.Start(":1323"))

// 		// wg.Wait()
// 	},
// }

func publish(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {
		// Write
		err := ws.WriteMessage(websocket.TextMessage, []byte("Hello, Client!"))
		if err != nil {
			c.Logger().Error(err)
		}

		// // Read
		// _, msg, err := ws.ReadMessage()
		// if err != nil {
		// 	c.Logger().Error(err)
		// }
		// fmt.Printf("%s\n", msg)
	}
}

func publisher() {
	
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.GET("/ws", publish)

	fmt.Printf("publishing ws on %s\n", ":1323")
	logger.Infof("publishing ws on %s", ":1323")

	e.Logger.Fatal(e.Start(":1323"))
}
