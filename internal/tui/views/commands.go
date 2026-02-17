package views

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/vieitesss/jocq/internal/buffer"
)

func FetchDecodedData(data *buffer.Data) tea.Cmd {
	return func() tea.Msg {
		content := data.Decoded()

		return DecodedDataFetchedMsg{
			Content: content,
		}
	}
}
