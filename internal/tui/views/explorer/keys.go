package explorer

import "github.com/charmbracelet/bubbles/key"

type KeyMap struct {
	NextPane   key.Binding
	PrevPane   key.Binding
	RunQuery   key.Binding
	ScrollUp   key.Binding
	ScrollDown key.Binding
	PageUp     key.Binding
	PageDown   key.Binding
	ToggleHelp key.Binding
	Quit       key.Binding
}

func NewKeyMap() KeyMap {
	return KeyMap{
		NextPane: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next pane"),
		),
		PrevPane: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev pane"),
		),
		RunQuery: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "run query"),
		),
		ScrollUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("up/k", "scroll up"),
		),
		ScrollDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("down/j", "scroll down"),
		),
		PageUp: key.NewBinding(
			key.WithKeys("pgup", "b"),
			key.WithHelp("pgup/b", "page up"),
		),
		PageDown: key.NewBinding(
			key.WithKeys("pgdown", "f", " "),
			key.WithHelp("pgdn/f/space", "page down"),
		),
		ToggleHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "toggle help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c"),
			key.WithHelp("ctrl+c", "quit"),
		),
	}
}

func (k *KeyMap) SetFocusMode(pane PaneID) {
	inputFocused := pane == InputPane

	k.RunQuery.SetEnabled(inputFocused)
	k.ScrollUp.SetEnabled(!inputFocused)
	k.ScrollDown.SetEnabled(!inputFocused)
	k.PageUp.SetEnabled(!inputFocused)
	k.PageDown.SetEnabled(!inputFocused)
	k.ToggleHelp.SetEnabled(!inputFocused)
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		k.NextPane,
		k.PrevPane,
		k.RunQuery,
		k.ScrollUp,
		k.ScrollDown,
		k.ToggleHelp,
		k.Quit,
	}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.NextPane, k.PrevPane},
		{k.RunQuery, k.ToggleHelp, k.Quit},
		{k.ScrollUp, k.ScrollDown, k.PageUp, k.PageDown},
	}
}
