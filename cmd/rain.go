/*
Copyright © 2026 gitstick
*/
package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
)

var rainCmd = &cobra.Command{
	Use:   "rain",
	Short: "Make the haha face rain from the sky",
	Long:  "Make the haha face rain from the sky. Press q or ctrl+c to quit.",
	Run: func(cmd *cobra.Command, args []string) {
		runRain()
	},
}

func init() {
	rootCmd.AddCommand(rainCmd)
}

var rainStatusStyle = lipgloss.NewStyle().Foreground(lipgloss.White)

type rainDrop struct {
	x     int
	y     int
	speed int
	delay int
}

type model struct {
	width     int
	height    int
	art       []string
	artWidth  int
	artHeight int
	drops     []rainDrop
	frame     int
}

func (m model) Init() tea.Cmd {
	return tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tickMsg:
		m.frame++
		m.updateDrops()
		return m, tick
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ensureDrops()
	}

	return m, nil
}

func (m model) View() tea.View {
	if m.width == 0 || m.height == 0 {
		v := tea.NewView("Initializing...")
		v.AltScreen = true
		return v
	}

	screenHeight := max(m.height-1, 1)
	rows := make([][]rune, screenHeight)
	for y := range rows {
		rows[y] = make([]rune, m.width)
		for x := range rows[y] {
			rows[y][x] = ' '
		}
	}

	for _, drop := range m.drops {
		if drop.delay > 0 {
			continue
		}
		m.drawDrop(rows, drop)
	}

	var s strings.Builder
	for y, row := range rows {
		s.WriteString(string(row))
		if y < len(rows)-1 {
			s.WriteByte('\n')
		}
	}
	s.WriteByte('\n')
	s.WriteString(rainStatusStyle.Render("Press q or ctrl+c to quit."))

	v := tea.NewView(s.String())
	v.AltScreen = true
	return v
}

func (m model) drawDrop(rows [][]rune, drop rainDrop) {
	for artY, line := range m.art {
		screenY := drop.y + artY
		if screenY < 0 || screenY >= len(rows) {
			continue
		}

		for artX, char := range []rune(line) {
			screenX := drop.x + artX
			if screenX < 0 || screenX >= m.width || char == ' ' {
				continue
			}
			rows[screenY][screenX] = char
		}
	}
}

func (m *model) updateDrops() {
	m.ensureDrops()

	for i := range m.drops {
		if m.drops[i].delay > 0 {
			m.drops[i].delay--
			continue
		}

		if m.frame%m.drops[i].speed == 0 {
			m.drops[i].y++
		}

		if m.drops[i].y > m.height {
			m.drops[i] = m.newDrop(true)
		}
	}
}

func (m *model) ensureDrops() {
	if m.width == 0 || m.height == 0 || m.artWidth == 0 {
		m.drops = nil
		return
	}

	target := max(1, min(m.width/max(m.artWidth, 1), 4))
	for len(m.drops) < target {
		m.drops = append(m.drops, m.newDrop(false))
	}
	if len(m.drops) > target {
		m.drops = m.drops[:target]
	}
}

func (m model) newDrop(fromTop bool) rainDrop {
	maxX := max(m.width-m.artWidth, 0)
	y := -rand.Intn(max(m.artHeight, 1))
	delay := rand.Intn(20)
	if !fromTop {
		y = rand.Intn(max(m.height, 1)) - m.artHeight
		delay = rand.Intn(8)
	}

	return rainDrop{
		x:     rand.Intn(maxX + 1),
		y:     y,
		speed: rand.Intn(3) + 1,
		delay: delay,
	}
}

type tickMsg time.Time

func tick() tea.Msg {
	time.Sleep(time.Millisecond * 70)
	return tickMsg(time.Now())
}

func initialModel() model {
	art := strings.Split(trollface, "\n")
	artWidth := 0
	for _, line := range art {
		artWidth = max(artWidth, len([]rune(line)))
	}

	return model{
		art:       art,
		artWidth:  artWidth,
		artHeight: len(art),
	}
}

func runRain() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}
