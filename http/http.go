package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	url    string
	status int
	err    error
}

type statusMsg int

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

// Cmds are functions that perform some I/O and then return a Msg
func checkSomeURL(url string) tea.Cmd {
	return func() tea.Msg {
		// Create an HTTP client and make a GET request
		c := &http.Client{Timeout: 10 * time.Second}
		res, err := c.Get(url)
		if err != nil {
			// Wrap the error in a message and return it
			return errMsg{err}
		}

		// Return the HTTP status code as a message
		return statusMsg(res.StatusCode)
	}
}

func (m model) Init() tea.Cmd {
	// Note we don't call the function...
	// The Bubble Tea runtime will do that when the time is right
	return checkSomeURL(m.url)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case statusMsg:
		// The server returned a status message. Save it to our model.
		// We tell the Bubble Tea runtime that we want to exit but that
		// still let us render a final view
		m.status = int(msg)
		return m, tea.Quit
	case errMsg:
		// There was an error. Note it in the model and tell the runtime
		// that we want to quit.
		m.err = msg
		return m, tea.Quit
	case tea.KeyMsg:
		// Give users a way to exit if they want to.
		if msg.Type == tea.KeyCtrlC {
			return m, tea.Quit
		}
	}

	// Don't do anything else on other messages
	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n", m.err)
	}

	s := fmt.Sprintf("Checking %s ...", m.url)

	if m.status > 0 {
		s += fmt.Sprintf("%d %s!", m.status, http.StatusText(m.status))
	}

	return "\n" + s + "\n"
}

func main() {
	url := os.Args[1]
	if err := tea.NewProgram(model{url: url}).Start(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}
