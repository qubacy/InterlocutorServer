package server

import (
	"encoding/json"
	"fmt"
	"ilserver/domain"
	"ilserver/repository"
	"ilserver/transport/control"
	"ilserver/transport/overWs"
	"ilserver/transport/overWsDto"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

type App struct {
	httpServer *http.Server
}

func NewApp() *App {
	return &App{}
}

// -----------------------------------------------------------------------

func (s *App) Run() error {
	wsHandler := overWs.NewHandler()
	controlHandler := control.NewHandler()

	// ***

	overWs.BackgroundUpdateRooms(wsHandler)

	// ***

	var mux = http.NewServeMux()

	// *** websocket and control ***

	prepareWsServer(mux, wsHandler)
	prepareControlServer(mux, controlHandler)

	// *** simple debug http ***

	var useDebugSvr bool = viper.GetBool(
		"debug_server.use")
	if useDebugSvr {
		prepareDebugServer(mux, wsHandler)
	}

	// ***

	s.httpServer = &http.Server{
		Addr:           ":" + viper.GetString("port"),
		MaxHeaderBytes: 1 << 20, // 1 MB
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		Handler:        mux,
	}

	return s.httpServer.ListenAndServe()
}

// route preparation
// -----------------------------------------------------------------------

func prepareControlServer(mux *http.ServeMux, handler *control.Handler) {
	mux.HandleFunc("/sign-in", handler.SignIn)
	mux.HandleFunc("/topic", handler.PostTopic)
	mux.HandleFunc("/topics", handler.GetTopics)
}

func prepareWsServer(mux *http.ServeMux, handler *overWs.Handler) {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		websocket, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Upgrade err:", err)
			return
		}

		log.Println("Ws", websocket.RemoteAddr().String(), "connected")
		if err := listen(handler, websocket); err != nil {
			log.Println("One ws listen err:", err)
			websocket.Close()
			return
		}
	})
}

func prepareDebugServer(mux *http.ServeMux, handler *overWs.Handler) {
	// GET
	mux.HandleFunc("/debug/runtime/rooms",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				http.NotFound(w, r)
				return
			}

			bytes, _ := json.Marshal(handler.RoomService.Rooms)
			w.Header().Add("Content-Type", "application/json")
			w.Write(bytes)
		})
	// GET
	mux.HandleFunc("/debug/database/admin-count",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.Write([]byte("Method is not GET"))
				return
			}

			err, count := repository.Instance().RecordCountInTable("Admins")
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(
				strconv.Itoa(count)))
		})
	// GET
	mux.HandleFunc("/debug/database/has-admin-with-login",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.Write([]byte("Method is not GET"))
				return
			}

			login := r.URL.Query().Get("login")
			err, has := repository.Instance().HasAdminByLogin(login)

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(
				strconv.FormatBool(has)))
		})
	// POST
	mux.HandleFunc("/debug/database/insert-topic",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.Write([]byte("Method is not POST"))
				return
			}

			err := r.ParseForm()
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			// ***

			name := r.Form.Get("name")
			langAsStr := r.Form.Get("lang")

			lang, err := strconv.Atoi(langAsStr)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			// ***

			err, idr := repository.Instance().InsertTopic(
				domain.Topic{
					Lang: lang,
					Name: name,
				},
			)

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(
				strconv.FormatInt(idr, 10)))
		})
	// POST
	mux.HandleFunc("/debug/database/update-admin-pass",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.Write([]byte("Method is not POST"))
				return
			}

			err := r.ParseForm()
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			// ***

			login := r.Form.Get("login")
			newPass := r.Form.Get("newPass")

			err = repository.Instance().UpdateAdminPass(
				login, newPass)

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(
				strconv.FormatBool(true)))
		})
	// GET
	mux.HandleFunc("/debug/database/random-one-topic",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.Write([]byte("Method is not GET"))
				return
			}

			err := r.ParseForm()
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			// ***

			langAsStr := r.Form.Get("lang")
			lang, err := strconv.Atoi(langAsStr)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			// ***

			err, tc := repository.Instance().SelectRandomOneTopic(lang)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}

			// ***

			bytes, err := json.Marshal(tc)
			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Header().Add("Content-Type", "application/json")
			w.Write(bytes)
		})
	// GET
	mux.HandleFunc("/debug/database/has-admin-with-login-and-pass",
		func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				w.Write([]byte("Method is not GET"))
				return
			}

			pass := r.URL.Query().Get("pass")
			login := r.URL.Query().Get("login")
			err, has := repository.Instance().HasAdminWithLoginAndPass(
				login, pass)

			if err != nil {
				w.Write([]byte(err.Error()))
				return
			}
			w.Write([]byte(
				strconv.FormatBool(has)))
		})
}

// -----------------------------------------------------------------------

func listen(handler *overWs.Handler, conn *websocket.Conn) error {
	handler.AddConn(conn)

	for {
		messageType, messageContent, err := conn.ReadMessage()

		// TODO: изучить закрытие веб-сокета
		if err != nil {
			switch err.(type) {
			case *websocket.CloseError:
				concreteErr := err.(*websocket.CloseError)
				log.Printf("Unexpected read message, close err %v", concreteErr)
			case *websocket.HandshakeError:
				concreteErr := err.(*websocket.HandshakeError)
				log.Printf("Unexpected read message, handshake err %v", concreteErr)
			}
			return handler.RemoveConnAndClose(conn)
		}

		// ***

		log.Println(string(messageContent))
		log.Println(messageType)

		if messageType == websocket.CloseMessage {
			return handler.RemoveConnAndClose(conn)
		}

		if messageType != websocket.TextMessage {
			log.Println(conn.RemoteAddr(), "message type is not text")
			return handler.RemoveConnAndClose(conn)
		}

		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			return handler.RemoveConnAndClose(conn)
		}

		// ***

		var pack overWsDto.Pack
		err = json.Unmarshal(messageContent, &pack)
		if err != nil {
			log.Println(conn.RemoteAddr(), err)

			// TODO: отправить пакет с информацией об ошибки

			return handler.RemoveConnAndClose(conn)
		}
		log.Println(conn.RemoteAddr(), pack)

		// ***

		err = routeWsPack(handler, conn, pack)
		if err != nil {
			log.Println(conn.RemoteAddr(), err)
			return handler.RemoveConnAndClose(conn)
		}
	}
}

func routeWsPack(handler *overWs.Handler, conn *websocket.Conn, pack overWsDto.Pack) error {
	bytes, err := json.Marshal(pack.RawBody)
	if err != nil {
		return err
	}

	// ***

	if pack.Operation == overWsDto.SEARCHING_START {
		var reqDto overWsDto.CliSearchingStartBodyClient
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return err
		}

		return handler.SearchingStart(conn, reqDto)
	} else if pack.Operation == overWsDto.SEARCHING_STOP {
		var reqDto overWsDto.CliSearchingStopBody
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return err
		}

		return handler.SearchingStop(conn, reqDto)
	} else if pack.Operation == overWsDto.CHATTING_NEW_MESSAGE {
		var reqDto overWsDto.CliChattingNewMessageBody
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return err
		}

		return handler.ChattingNewMessage(conn, reqDto)
	} else if pack.Operation == overWsDto.CHOOSING_USERS_CHOSEN {
		var reqDto overWsDto.CliChoosingUsersChosenBody
		err = json.Unmarshal(bytes, &reqDto)
		if err != nil {
			return nil
		}

		return handler.ChoosingUsersChosen(conn, reqDto)
	}

	// ***

	return fmt.Errorf("routeWsPack, operation is unknown")
}
