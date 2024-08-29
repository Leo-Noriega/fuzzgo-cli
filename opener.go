package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: Add choice to use 'bat' with files or other options
var choices = []string{"nvim", "code", "open", "yazi"}


// Model - Holds the state of the application
type model struct {
	cursor int
	choice string
}

// Init - Initializes the program
func (m model) Init() tea.Cmd {
	return nil
}

// Update - Handles messages and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = len(choices) - 1
			}
		case "down", "j":
			m.cursor++
			if m.cursor >= len(choices) {
				m.cursor = 0
			}
		case "enter":
			m.choice = choices[m.cursor]
			return m, tea.Quit
		}
	}
	return m, nil
}

// TODO: Add Glamour to beautify the output or LipGloss all by Charmbracelet

// View - Renders the UI
func (m model) View() string {
	s := strings.Builder{}
	s.WriteString("\nChoose an application to open the project with:\n\n")
	for i := 0; i < len(choices); i++ {
		if m.cursor == i {
			s.WriteString("(∆) ")
		} else {
			s.WriteString("( ) ")
		}
		s.WriteString(choices[i])
		s.WriteString("\n")
	}
	s.WriteString("\nPress 'q' or 'ctrl+c' to quit.\n")
	return s.String()
}

func runCommand(command string, args ...string) {
	cmd := exec.Command(command, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error executing command: %v\n", err)
		os.Exit(1)
	}
}

// TODO: Refactor main to clean the function in smaller parts
func main() {
	// Execute the command and capture its output
	cmd := exec.Command("sh", "-c", `fd . '/Users/noriega' -t d -E Library --hidden | fzf --height=40% --layout=reverse --info=inline --border --margin=1 --padding=1`)
	out, err := cmd.Output()

	if err != nil {
		fmt.Println("Error executing command:", err)
		os.Exit(1)
	}

	// Use the ouput to execute next command
	path := strings.TrimSpace(string(out))
	fmt.Println("\nWorking with the directory:", path)

	// Initialize the program with the model
	p := tea.NewProgram(model{})
	m, err := p.Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}

	// cd to the directory
	cmd = exec.Command("sh", "-c", "cd "+path)

	// Assert the final tea.Model to our local model and print the choice.
	if m, ok := m.(model); ok && m.choice != "" {
		runCommand(m.choice, path)
	}
}

