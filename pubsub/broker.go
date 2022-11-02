package pubsub

import(
	"fmt"
  	"net/http"

  	"github.com/go-playground/validator"
	
	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	
	"github.com/algorand/go-algorand/data/transactions"

	"bytes"
    "encoding/json"
	"strconv"
)

const MAX_CH_BUFFER = 100

// 	MAX_FEEDER >= 
//  	(MAX_SUBSCRIBER + MAX_HOOKER)

const MAX_FEEDER = 10
const MAX_SUBSCRIBER = 5
const MAX_HOOKER = 5

var (
	upgrader = websocket.Upgrader{}
)

type (
	
	Topic struct {
		Type  string `json:"type" validate:"required"`
		Param string `json:"param"`
	}

	MsgCreateTopic struct {
		Topic  Topic `json:"topic" validate:"required"`
		Target string
		Token  string
	}

	PubSubValidator struct {
		validator *validator.Validate
	}

)

func (cv *PubSubValidator) Validate(i interface{}) error {
	if err := cv.validator.Struct(i); err != nil {
		// Optionally, you could return the error to give each route more control over the status code
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return nil
}



func (imp *PubSocket) serveGen(c echo.Context) (err error) {

	typeGen := c.Param("type")

	if typeGen == "websocket" && len((*imp).Pubs) >= MAX_SUBSCRIBER {
		return echo.NewHTTPError(http.StatusBadRequest, "PubSocket: Too many subscriber.")
	}

	if typeGen == "webhook" && len((*imp).Hooks) >= MAX_HOOKER {
		return echo.NewHTTPError(http.StatusBadRequest, "PubSocket: Too many hooker.")
	}

	t := new(MsgCreateTopic)
	if err = c.Bind(t); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err = c.Validate(t); err != nil {
		return err
	}

	newToken := RandomString(16)

	if typeGen == "websocket" {

		// Setup Socket channel.
		sStream := make(TxStream, MAX_CH_BUFFER)
		sInfo := SockInfo{ 
			Token: newToken,
			Target: t.Target,
			s: sStream, 
			t: t.Topic,
		}

		(*imp).Pubs[newToken] = sInfo

	} else if typeGen == "webhook" {

		// Setup Hooker channel.
		hStream := make(HookStream, MAX_CH_BUFFER)
		hInfo := HookInfo{ 
			Token: newToken,
			Target: t.Target,
			s: hStream,
			t: t.Topic,
		}

		// Create Hooker routine.
		go (*imp).Hooker(hInfo, hStream) 

		(*imp).Hooks[newToken] = hInfo

	} else {

		return echo.NewHTTPError(http.StatusBadRequest, "PubSocket: Unknown Type.")
	}
	
	t.Token = newToken;

	return c.JSON(http.StatusOK, t)
}

func (imp *PubSocket) serveSock(c echo.Context) error {

	token := c.Param("token")

	if v, exists := (*imp).Pubs[token]; exists {

		ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}
		defer ws.Close()

		for {

			select {
			case pub := <-v.s: // get from Channel

				pubMsg, err := json.Marshal(pub)
				if err != nil {
					fmt.Printf("PubSocket.serveSock: %v\r\n", err)
				} else {

					err = ws.WriteMessage(websocket.TextMessage, (pubMsg))
					if err != nil {
						c.Logger().Debugf("%v", err)
						return err
					}

				}
				
			// case <-time.After(5*time.Second):
				// -_-
			}
		}
	
	}

	return echo.NewHTTPError(http.StatusBadRequest, "PubSocket: Unknown token")
}

func (imp *PubSocket) CheckRules(t Topic, tx transactions.SignedTxnInBlock) bool {

	if t.Type == "*" { // ALL
		return true
	}
	if t.Type == "SENDER" && t.Param == tx.Txn.Sender.String() {
		return true
	}
	if t.Type == "RECEIVER" && t.Param == tx.Txn.Receiver.String() {
		return true
	}
	if t.Type == "APPID" && (t.Param != "0" && len(t.Param) > 0) {
		appId, err := strconv.Atoi(t.Param)
		if err != nil {
			// false
			fmt.Printf("PubSocket.CheckRules: %v\r\n", err)
		} else {
			if uint64(appId) == uint64(tx.ApplicationID) {
				return true
			}
		}
	}

	return false;

}

func (imp *PubSocket) FeedTx() {

	for {

		select {
		case tx := <-(*imp).Tx:

			if len((*imp).Pubs) > 0 {
				
				for _ , s := range (*imp).Pubs {
					if (*imp).CheckRules(s.t, tx) {
						s.s <- tx
					}
				}

			}


			if len((*imp).Hooks) > 0 {

				for _ , h := range (*imp).Hooks {
					if (*imp).CheckRules(h.t, tx) {
						h.s <- tx
					}
				}

			}

		// case <-time.After(5*time.Second):
			// -_-
		}

	}

}

func (imp *PubSocket) Hooker(i HookInfo, h HookStream) {

	for {

		select {
		case hook := <-h:
			hookMsg, err := json.Marshal(hook)
			if err != nil {
				fmt.Printf("PubSocket.Hooker: %v\r\n", err)
			} else {
				hookBody := bytes.NewBuffer(hookMsg)
				resp, err := http.Post(i.Target, "application/json", hookBody)
				if err != nil {
					fmt.Printf("PubSocket.Hooker: %v\r\n", err)
				} else {

					resp.Body.Close()
				}

			}
			
		// case <-time.After(5*time.Second):
			// -_-
		}
	}
}

type TxStream chan transactions.SignedTxnInBlock
type HookStream chan transactions.SignedTxnInBlock

type SockInfo struct {
	Token	string
	Target 	string
	s		TxStream
	t		Topic
}
type HookInfo struct {
	Token	string
	Target 	string
	s		HookStream
	t		Topic
}

// PubSocket - ...
type PubSocket struct {
	Tx 		TxStream
	// -----------------
	Pubs 	map[string]SockInfo
	Hooks 	map[string]HookInfo
	// -----------------
	e		*echo.Echo
	token	string
}

var pubSocket PubSocket

// GetPubSocket get an PubSocket
func GetPubSocket() *PubSocket {

	return &pubSocket
}

// NewPubSocket builds an PubSocket
func NewPubSocket(addr string, token string) *PubSocket {

	pubSocket = PubSocket{
		Tx: 		make(TxStream, MAX_CH_BUFFER),
		// ---------------------------------
		Pubs:   	make(map[string]SockInfo),
		Hooks:		make(map[string]HookInfo),
		// ---------------------------------
		e: 			echo.New(),
		token: 		token,
	}

	// FeedTx proc.
	for i := 0; i < MAX_FEEDER; i++ {
		go pubSocket.FeedTx()
	}

	// Serve PubSocket with WebSocket port.
	e := pubSocket.e
  	e.Validator = &PubSubValidator{validator: validator.New()}

	go func() { 

		e.Use(middleware.Logger())
		e.Use(middleware.Recover())
		e.Use(middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
			KeyLookup: "header:X-API-Key",
			Validator: func(key string, c echo.Context) (bool, error) {
				return key == pubSocket.token, nil
			},
		}))

		e.POST("/generate/:type", pubSocket.serveGen)
		e.GET("/ws/:token", pubSocket.serveSock) // /ws/<token>

		e.Logger.Fatal(e.Start(addr))

	}()


	return &pubSocket
}
