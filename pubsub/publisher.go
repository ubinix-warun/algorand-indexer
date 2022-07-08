package pubsub

import(
	"fmt"
	
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	
	"github.com/algorand/go-algorand/data/transactions"

    "encoding/json"
)

var (
	upgrader = websocket.Upgrader{}
)

func (imp *Publisher) serveSub(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	defer ws.Close()

	for {

		select {
		case pub := <-(*imp).Pubs:

			pubMsg, err := json.Marshal(pub)
			if err != nil {
				// panic(err)
				c.Logger().Debugf("%v", err)
			} else {
				c.Logger().Debugf("pub >> %s",string(pubMsg)) // {"full_name":"Bob"}
			}

			err = ws.WriteMessage(websocket.BinaryMessage, pubMsg)
			if err != nil {
				c.Logger().Debugf("%v", err)
				return err
			}

		// case <-time.After(5*time.Second):
			// ws.
		}

	}
}

// Publisher - ...
type Publisher struct {
	// Log *log.Logger
	Pubs chan transactions.SignedTxnInBlock
}


// NewPublisher builds an Publisher
// func NewPublisher(l *log.Logger) *Publisher {
func NewPublisher() *Publisher {

	pub := Publisher{
		// Log:             l,
		Pubs: 	make(chan transactions.SignedTxnInBlock, 100),
	}

	e := echo.New()

	go func() {

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.GET("/ws", pub.serveSub)

		fmt.Printf("publishing ws on %s\n", ":1323")
		e.Logger.Infof("publishing ws on %s", ":1323")
	
		e.Logger.Fatal(e.Start(":1323"))

	}()


	return &pub
}

