package main

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type (
	Editor struct {
		con  net.Conn
		addr string
	}
	Message struct {
		from    *Editor
		payload []byte
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
		editorsOnline: make(map[string]*Editor),
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
		// WARN: temp solution
		GameState.editMode = false
		con.Close()
	}()

	buf := make([]byte, 1024)
	editor := &Editor{
		con:  con,
		addr: con.RemoteAddr().String(),
	}
	s.editorsOnline[editor.addr] = editor

	GameState.editMode = true

	con.Write([]byte("\n-----------\n"))
	for i := range ViewportState.Items {
		con.Write([]byte(i + "\n"))
	}
	con.Write([]byte("\n-----------\n"))

	for {
		n, err := con.Read(buf)
		if err != nil {
			fmt.Printf("[%s] has disconnected \n", editor.addr)
			// TODO: remove editor from the map
			// if there are no editors change the editorMode to false
			break
		}

		s.msgch <- Message{
			from:    editor,
			payload: buf[:n],
		}
	}
}

func (s *Server) awijfiweojf() {}

func (s *Server) handleMessage() {
	for msg := range s.msgch {
		modelName, values, found := bytes.Cut(msg.payload, []byte{'-'})
		_ = modelName
		_ = values

		if !found {
			continue
		}
		item, e := ViewportState.Items[string(modelName)]
		_ = item
		if !e {
			errMsg := fmt.Sprintf("name not found: %s \n", modelName)
			fmt.Print(errMsg)
			msg.from.con.Write([]byte(errMsg))
			continue
		}

		sanitizedName := strings.ReplaceAll(string(values), "\n", "")
		s, _ := strconv.Atoi(string(sanitizedName))

		fmt.Printf("Updating %s...\n", sanitizedName)
		item.model.Transform = rl.MatrixRotateXYZ(rl.NewVector3(0, float32(s), 0))
	}
}
