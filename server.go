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
	HighlightCurrent
	FocusCurrent
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

	Scene.listItems(con)
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
		item, e := Scene.Items[string(input)]
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
					Scene.listItems(msg.from.con)
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
			selectedItem := Scene.Items[msg.from.editContext.modelName]
			// NOTE: after item selection
			switch msg.from.editContext.mode {
			case EditPos:
				fmt.Fprintf(msg.from.con, "Poisiton (%+v):", selectedItem.pos)
			case EditRot:
				rot := Scene.Items[msg.from.editContext.modelName].rot
				fmt.Fprintf(msg.from.con, "Rotation (X:%f, Y: %f, Z: %f):", rot.X*rl.Rad2deg, rot.Y*rl.Rad2deg, rot.Z*rl.Rad2deg)
			case EditScale:
				fmt.Fprintf(msg.from.con, "Scale (%+v):", selectedItem.scale)
			case ResetChanges:
				msg.from.con.Write([]byte("reseting unsaved changes...\n"))
				// TODO: reset the changes here
				msg.from.resetEditMode()
			case HighlightCurrent:
				selectedItem.highlight = !selectedItem.highlight
				msg.from.con.Write([]byte("element highlighted...\n"))
				msg.from.resetEditMode()
			case FocusCurrent:
				for _, item := range Scene.Items {
					if item.uid != msg.from.editContext.modelName && selectedItem.focus {
						selectedItem.focus = false
						GameState.editFocusedItemUid = item.uid
					} else {
						selectedItem.focus = true
						GameState.editFocusedItemUid = item.uid
					}
				}
				msg.from.con.Write([]byte("focused on element...\n"))
				msg.from.resetEditMode()
			case ToggleModel:
				selectedItem.hidden = !Scene.Items[msg.from.editContext.modelName].hidden
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
					Scene.Items[msg.from.editContext.modelName].pos = rl.NewVector3(float32(v1), float32(v2), float32(v3))
				case EditRot:
					if len(vals) != 3 {
						msg.from.con.Write([]byte("invalid number of arguments, need 3\n"))
						break
					}
					msg.from.con.Write([]byte("Rotation: "))
					v1, _ := strconv.ParseFloat(vals[0], 32)
					v2, _ := strconv.ParseFloat(vals[1], 32)
					v3, _ := strconv.ParseFloat(vals[2], 32)
					v := rl.NewVector3(float32(v1)*rl.Deg2rad, float32(v2)*rl.Deg2rad, float32(v3)*rl.Deg2rad)
					Scene.Items[msg.from.editContext.modelName].model.Transform = rl.MatrixRotateXYZ(v)
				case EditScale:
					if len(vals) != 1 {
						msg.from.con.Write([]byte("invalid number of arguments, need 1\n"))
						break
					}
					msg.from.con.Write([]byte("scale edit"))
				case ResetChanges:
					msg.from.con.Write([]byte("reset all changes"))
					v1, _ := strconv.ParseFloat(vals[0], 32)
					Scene.Items[msg.from.editContext.modelName].scale = float32(v1)
				case DeleteModel:
					if input[0] == 'y' || input == "yes" {
						delete(Scene.Items, msg.from.editContext.modelName)
						msg.from.editContext.mode = Cancel
						msg.from.editContext.modelName = ""
						fmt.Fprint(msg.from.con, "Model deleted... \n")
					} else {
						fmt.Fprint(msg.from.con, "Canceled...\n")
						msg.from.editContext.mode = Cancel
					}
				case Cancel:
					msg.from.con.Write([]byte("Canceled...\n"))
				default:
					msg.from.con.Write([]byte("\nif you see this message open an issue\n"))
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
	temp := []string{"rotation", "position", "scale", "reset", "hide model", "toggle highlight", "toggle focus", "delete model", "cancel"}
	el := Scene.Items[e.editContext.modelName]
	fmt.Fprintf(e.con, "[%s] Pos: %+v | Rot: %+v | Scale:%+v | Visible:%+v | %+v \n", e.editContext.modelName, el.pos, el.rot, el.scale, el.hidden, el.model.Materials.Shader)
	for i, cmd := range temp {
		fmt.Fprintf(e.con, "[%d] - %s \n", i+1, cmd)
	}
}

func (s *Scene3D) listItems(con net.Conn) {
	con.Write([]byte("\n-----Item List------\n"))
	for i := range Scene.Items {
		con.Write([]byte(i + "\n"))
	}
	con.Write([]byte("--------------------\n"))
}

func (editor *Editor) resetEditMode() {
	editor.editContext.mode = Cancel
	editor.askAttr()
}
