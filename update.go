//Julie
package main

import (
	"time"
	"golang-addon/week-1/golang-clidle/wordle"
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case msgResetStatus:
		// If there is more than one pending status message, that means
		// something else is currently displaying a status message, so we don't
		// want to overwrite it.
		m.statusPending--
		if m.statusPending == 0 {
			m.handleResetStatus()
		}

	// Handle keypresses
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlD:
			return m, tea.Quit
		
		case tea.KeyBackspace:
			m.handleDeleteChar()

		case tea.KeyEnter:
			m.handleSubmitActiveGuess()

		case tea.KeyRunes:
			if len(msg.Runes) == 1 {
				m.handleSubmitChar(msg.Runes[0])
			}
		

		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// handleSetStatus sets the status message, and returns a tea.Cmd that restores the
// default status message after a delay.
func (m *model) handleSetStatus(msg string, duration time.Duration) tea.Cmd {
	m.status = msg
	if duration > 0 {
		m.statusPending++
		return tea.Tick(duration, func(time.Time) tea.Msg {
			return msgResetStatus{}
		})
	}
	return nil
}

// handleResetStatus immediately resets the status message to its default value.
func (m *model) handleResetStatus() {
	m.status = "Guess the word!"
}

// msgResetStatus is sent when the status line should be reset.
type msgResetStatus struct{}


func (m *model) handleSubmitChar(r rune) {
	if m.cursor < wordle.WordSize {
		if 'a' <= r && r <= 'z' {
			r -= 'a' - 'A'
		}

		if 'A' <= r && r <= 'Z' {
			m.activeGuess[m.cursor] = byte(r)
			m.cursor++
		}
	}
}

func (m *model) handleDeleteChar() {
	if m.cursor > 0 {
		m.cursor--

	}
}

func (m *model) handleSubmitActiveGuess() {
	ws := m.ws
	g := wordle.NewGuess(string(m.activeGuess[:]))
	g.UpdateLettersWithWord(ws.Word)

	err := ws.AppendGuess(g)
	if err != nil {
		m.handleSetStatus(err.Error(), 1*time.Second)
		return
	}

	m.handleResetActiveGuess()
}

func (m *model) handleResetActiveGuess() {
	copy(m.activeGuess[:], "")
	m.cursor = 0
}