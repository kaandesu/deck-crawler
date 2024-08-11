package main

type key int32

const (
	BackKey key = iota
	QuitKey
	SaveKey
)

type Help struct {
	Key  string
	Desc string
}

type Binding struct {
	help     Help
	keys     []string
	disabled bool
}

type keyMap struct {
	Back         Binding
	Quit         Binding
	Save         Binding
	Help         Binding
	FullScreen   Binding
	CycleCamMode Binding
}

var keys = keyMap{
	Back: NewBinding(
		WithKeys("b", "back"),
		WithHelp("b", "go back"),
	),
	Quit: NewBinding(
		WithKeys("q", "quit"),
		WithHelp("q", "quit editor mode"),
	),
	Save: NewBinding(
		WithDisabled(),
		WithKeys("s", "save"),
		WithHelp("s", "[wip] save the changes"),
	),
	Help: NewBinding(
		WithDisabled(),
		WithKeys("h", "help"),
		WithHelp("h", "[wip] display the help menu"),
	),
	FullScreen: NewBinding(
		WithKeys("full", "f"),
		WithHelp("f", "toggle fullscreen"),
	),
	CycleCamMode: NewBinding(
		WithKeys("c", "cam"),
		WithHelp("c", "toggle cam mode"),
	),
}

type BindingOpt func(*Binding)

func Matches(key string, b ...Binding) bool {
	for _, binding := range b {
		for _, v := range binding.keys {
			if key == v && binding.Enabled() {
				return true
			}
		}
	}
	return false
}

func NewBinding(opts ...BindingOpt) Binding {
	b := &Binding{}
	for _, opt := range opts {
		opt(b)
	}
	return *b
}

func WithKeys(keys ...string) BindingOpt {
	return func(b *Binding) {
		b.keys = keys
	}
}

func (b Binding) Enabled() bool {
	return !b.disabled && b.keys != nil
}

func WithHelp(key, desc string) BindingOpt {
	return func(b *Binding) {
		b.help = Help{Key: key, Desc: desc}
	}
}

func WithDisabled() BindingOpt {
	return func(b *Binding) {
		b.disabled = true
	}
}
