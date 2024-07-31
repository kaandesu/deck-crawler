package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
)

type EditMode int32

const (
	EditRot EditMode = iota
	EditPos
	EditScale
	ToggleModel
	DeleteModel
	Cancel
)

type (
	msgPayload []byte
	Editor     struct {
		con         net.Conn
		addr        string
		editContext struct {
			modelName string
			mode      EditMode
		}
	}
	Message struct {
		from    *Editor
		payload msgPayload
	}
	Server struct {
		quitch        chan struct{}
		msgch         chan Message
		ln            net.Listener
		logger        *log.Logger
		editorsOnline map[string]*Editor
		addr          string
	}
)

func NewLogger(filename string) *log.Logger {
	logFile, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		panic(err)
	}
	return log.New(logFile, "[deck_crawler_editor]", log.Ldate|log.Ltime)
}

func NewServer(addr string) *Server {
	return &Server{
		addr:          addr,
		msgch:         make(chan Message, 10),
		editorsOnline: make(map[string]*Editor, 10),
		quitch:        make(chan struct{}),
		logger:        NewLogger("./editor.log"),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	s.ln = ln

	go s.acceptLoop()
	go s.handleMessage()

	<-s.quitch
	close(s.msgch)

	return nil
}

func (s *Server) acceptLoop() {
	for {
		con, err := s.ln.Accept()
		if err != nil {
			break
		}
		go s.handleConn(con)
	}
}

func (s *Server) handleConn(con net.Conn) {
	defer func() {
		delete(s.editorsOnline, con.RemoteAddr().String())
		fmt.Printf("[%s] has disconnected \n", con.RemoteAddr().String())
		if len(s.editorsOnline) == 0 {
			GameState.editMode = false
		}
		con.Close()
	}()

	buf := make([]byte, 1024)
	editor := &Editor{
		con:  con,
		addr: con.RemoteAddr().String(),
	}
	editor.editContext.mode = Cancel
	s.editorsOnline[editor.addr] = editor

	GameState.editMode = true

	con.Write([]byte("\n-----Item List------\n"))
	for i := range ViewportState.Items {
		con.Write([]byte(i + "\n"))
	}
	con.Write([]byte("--------------------\n"))

	for {
		n, err := con.Read(buf)
		if err != nil {
			break
		}

		s.msgch <- Message{
			from:    editor,
			payload: buf[:n],
		}
	}
}

func (s *Server) handleMessage() {
	for msg := range s.msgch {

		modelName := string(msg.payload[:len(msg.payload)-1])
		item, e := ViewportState.Items[string(modelName)]
		_ = item

		if !e {
			errMsg := fmt.Sprintf("name not found: %s \n", modelName)
			fmt.Print(errMsg)
			msg.from.con.Write([]byte(errMsg))
			continue
		}

		if msg.from.editContext.modelName == "" {
			msg.from.editContext.modelName = string(modelName)
			msg.from.askAttr()
			return
		}

		if msg.from.editContext.mode == Cancel {

			mode, err := msg.payload.parseEditMode()
			if err != nil {
				fmt.Println(err)
				msg.from.con.Write([]byte("invalid mode"))
				return
			}

			fmt.Println(mode)

			switch mode {
			case EditPos:
				msg.from.con.Write([]byte("edit"))
			case EditRot:
			case EditScale:
			case ToggleModel:
			case DeleteModel:
			case Cancel:
			default:

			}
			return
		}

		// item.model.Transform = rl.MatrixRotateXYZ(rl.NewVector3(0, float32(s), 0))
	}
}

func (payload msgPayload) parseEditMode() (EditMode, error) {
	if len(payload) == 0 {
		return -1, fmt.Errorf("payload is empty")
	}
	b := payload[0]
	if b >= byte(EditRot) && b <= byte(Cancel) {
		return EditMode(b), nil
	}
	return -1, fmt.Errorf("invalid byte for EditMode: %v", b)
}

func sanitize(data []byte) string {
	return string(data[:len(data)-1])
}

func (input EditMode) validate() error {
	if input < EditRot || input > Cancel {
		return errors.New("invalid EditMode value")
	}
	return nil
}

func (e *Editor) askAttr() {
	temp := []string{"rotation", "position", "scale", "hide model", "delete model"}
	for i, cmd := range temp {
		fmt.Fprintf(e.con, "[%d] - %s \n", i+1, cmd)
	}
}
