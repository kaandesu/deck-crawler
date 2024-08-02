package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type EditMode int32

const (
	EditRot EditMode = iota
	EditPos
	EditScale
	ResetChanges
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
		msgch:         make(chan Message, 100),
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

	ViewportState.listItems(con)
	con.Write([]byte("enter item name: "))

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

		input := string(msg.payload[:len(msg.payload)-1])
		item, e := ViewportState.Items[string(input)]
		_ = item

		switch {
		case Matches(msg.payload.String(), keys.Back):
			if msg.from.editContext.modelName == "" {
				msg.from.con.Write([]byte("to quite type 'q' or 'quit' \n"))
			}

			if msg.from.editContext.modelName != "" {
				if msg.from.editContext.mode != Cancel {
					msg.from.resetEditMode()
				} else {
					msg.from.editContext.modelName = ""
					ViewportState.listItems(msg.from.con)
					msg.from.con.Write([]byte("enter item name: "))
				}
			}
			fmt.Println("editor - back")
			continue
		case Matches(msg.payload.String(), keys.Quit):
			msg.from.con.Close()
			fmt.Println("editor - quit")
			continue
		case Matches(msg.payload.String(), keys.Save):
			fmt.Println("save is not implemented")
			continue
		case Matches(msg.payload.String(), keys.Help):
			fmt.Println("help is not implemented")
			continue
		case Matches(msg.payload.String(), keys.FullScreen):
			GameState.editFull = !GameState.editFull
			continue
		case Matches(msg.payload.String(), keys.CycleCamMode):
			if GameState.camMode == rl.CameraFirstPerson {
				GameState.camMode = rl.CameraThirdPerson
				fmt.Println("switched cam mode: third")
			} else {
				GameState.camMode = rl.CameraFirstPerson
				fmt.Println("switched cam mode: first")
			}
			continue
		}

		switch {
		case msg.from.editContext.modelName == "" && !e:
			errMsg := fmt.Sprintf("name not found: %s \n", input)
			fmt.Print(errMsg)
			msg.from.con.Write([]byte(errMsg))

		case msg.from.editContext.modelName == "":
			msg.from.editContext.modelName = string(input)
			msg.from.askAttr()

		case msg.from.editContext.mode == Cancel:
			mode, err := msg.payload.parseEditMode()
			if err != nil {
				fmt.Println(err)
				msg.from.con.Write([]byte("invalid mode"))
				break
			}

			msg.from.editContext.mode = mode
			// NOTE: after item selection
			switch msg.from.editContext.mode {
			case EditPos:
				fmt.Fprintf(msg.from.con, "Poisiton (%+v):", ViewportState.Items[msg.from.editContext.modelName].pos)
			case EditRot:
				fmt.Fprintf(msg.from.con, "Rotation (%+v):", ViewportState.Items[msg.from.editContext.modelName].rot)
			case EditScale:
				fmt.Fprintf(msg.from.con, "Scale (%+v):", ViewportState.Items[msg.from.editContext.modelName].scale)
			case ResetChanges:
				msg.from.con.Write([]byte("reseting unsaved changes...\n"))
				// TODO: reset the changes here
				msg.from.resetEditMode()
			case ToggleModel:
				ViewportState.Items[msg.from.editContext.modelName].hidden = !ViewportState.Items[msg.from.editContext.modelName].hidden
				fmt.Fprint(msg.from.con, "Item toggled \n")
				msg.from.resetEditMode()
			case DeleteModel:
				msg.from.con.Write([]byte("Are you you want to DELETE the model? (y/n):"))
			case Cancel:
				msg.from.con.Write([]byte("canceled"))
			}
		default:
			if msg.from.editContext.mode != Cancel {
				vals := strings.Split(input, " ")
				switch msg.from.editContext.mode {
				case EditPos:
					if len(vals) != 3 {
						msg.from.con.Write([]byte("not enough arguments, need 3\n"))
						break
					}
					v1, _ := strconv.ParseFloat(vals[0], 32)
					v2, _ := strconv.ParseFloat(vals[1], 32)
					v3, _ := strconv.ParseFloat(vals[2], 32)
					msg.from.con.Write([]byte("Position: "))
					ViewportState.Items[msg.from.editContext.modelName].pos = rl.NewVector3(float32(v1), float32(v2), float32(v3))
				case EditRot:
					if len(vals) != 3 {
						msg.from.con.Write([]byte("not enough arguments, need 3\n"))
						break
					}
					msg.from.con.Write([]byte("Rotation: "))
					v1, _ := strconv.ParseFloat(vals[0], 32)
					v2, _ := strconv.ParseFloat(vals[1], 32)
					v3, _ := strconv.ParseFloat(vals[2], 32)
					v := rl.NewVector3(float32(v1), float32(v2), float32(v3))
					ViewportState.Items[msg.from.editContext.modelName].model.Transform = rl.MatrixRotateXYZ(v)
				case EditScale:
					msg.from.con.Write([]byte("scale edit"))
				case ResetChanges:
					msg.from.con.Write([]byte("reset all changes"))
				case ToggleModel:
					msg.from.con.Write([]byte("toggle model"))
				case DeleteModel:
					msg.from.con.Write([]byte("delete model"))
					if input[0] == 'y' || input == "yes" {
						// TODO: delete model here
						fmt.Fprint(msg.from.con, "[DELETE NOT IMPLEMENTED] \n")
					} else {
						fmt.Fprint(msg.from.con, "Canceled...\n")
					}
				case Cancel:
					msg.from.con.Write([]byte("canceled"))
				default:
					msg.from.con.Write([]byte("you shouldn't see this message i think"))
				}
			}
		}

	}
}

func (payload msgPayload) parseEditMode() (EditMode, error) {
	if len(payload) == 0 {
		return -1, fmt.Errorf("payload is empty")
	}
	b := payload[0] - 49
	if b >= byte(EditRot) && b <= byte(Cancel) {
		return EditMode(b), nil
	}
	return -1, fmt.Errorf("invalid byte for EditMode: %v", b)
}

func (data msgPayload) String() string {
	return string(data[:len(data)-1])
}

func (input EditMode) validate() error {
	if input < EditRot || input > Cancel {
		return errors.New("invalid EditMode value")
	}
	return nil
}

func (e *Editor) askAttr() {
	temp := []string{"rotation", "position", "scale", "reset", "hide model", "delete model", "cancel"}
	el := ViewportState.Items[e.editContext.modelName]
	fmt.Fprintf(e.con, "[%s] Pos: %+v | Rot: %+v | Scale:%+v | Hidden:%+v\n", e.editContext.modelName, el.pos, el.pos, el.scale, el.hidden)
	for i, cmd := range temp {
		fmt.Fprintf(e.con, "[%d] - %s \n", i+1, cmd)
	}
}

func (s *Scene3D) listItems(con net.Conn) {
	con.Write([]byte("\n-----Item List------\n"))
	for i := range ViewportState.Items {
		con.Write([]byte(i + "\n"))
	}
	con.Write([]byte("--------------------\n"))
}

func (editor *Editor) resetEditMode() {
	editor.editContext.mode = Cancel
	editor.askAttr()
}