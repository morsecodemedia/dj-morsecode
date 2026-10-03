package main

import (
	tea "github.com/charmbracelet/bubbletea"
)

func (m model) updateWindowSize(
	msg tea.WindowSizeMsg,
) (tea.Model, tea.Cmd) {

	m.Width = msg.Width
	m.Height = msg.Height

	return m, nil

}
