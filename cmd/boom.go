/*
Copyright © 2026 gitstick

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
	"time"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/spf13/cobra"
)

// boomCmd represents the boom command
var boomCmd = &cobra.Command{
	Use:   "boom",
	Short: "Animate an ASCII explosion",
	Long:  "Explode the selected character art across the terminal. Press q or ctrl+c to quit.",
	Run: func(cmd *cobra.Command, args []string) {
		character, _ := cmd.Flags().GetString("character")
		runBoom(character)
	},
}

func init() {
	rootCmd.AddCommand(boomCmd)
}

var boomStatusStyle = lipgloss.NewStyle().Foreground(lipgloss.White)

const boomMaxFrames = 55

type boomModel struct {
	width     int
	height    int
	frame     int
	art       []string
	artWidth  int
	artHeight int
	pieces    []boomPiece
}

type boomTickMsg time.Time

type boomPiece struct {
	x  float64
	y  float64
	vx float64
	vy float64
}

func (m boomModel) Init() tea.Cmd {
	return boomTick
}

func (m boomModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case boomTickMsg:
		m.frame++
		m.updatePieces()
		if m.frame >= boomMaxFrames {
			return m, tea.Quit
		}
		return m, boomTick
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ensurePieces()
	}

	return m, nil
}

func (m boomModel) View() tea.View {
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

	for _, piece := range m.pieces {
		m.drawPiece(rows, piece)
	}

	var s strings.Builder
	for y, row := range rows {
		s.WriteString(string(row))
		if y < len(rows)-1 {
			s.WriteByte('\n')
		}
	}
	s.WriteByte('\n')
	s.WriteString(boomStatusStyle.Render("Press q or ctrl+c to quit."))

	v := tea.NewView(s.String())
	v.AltScreen = true
	return v
}

func (m boomModel) drawPiece(rows [][]rune, piece boomPiece) {
	x := int(math.Round(piece.x))
	y := int(math.Round(piece.y))

	for artY, line := range m.art {
		screenY := y + artY
		if screenY < 0 || screenY >= len(rows) {
			continue
		}

		for artX, char := range []rune(line) {
			screenX := x + artX
			if screenX < 0 || screenX >= m.width || char == ' ' || char == '⠀' {
				continue
			}
			rows[screenY][screenX] = char
		}
	}
}

func (m *boomModel) updatePieces() {
	m.ensurePieces()

	for i := range m.pieces {
		m.pieces[i].x += m.pieces[i].vx
		m.pieces[i].y += m.pieces[i].vy
	}
}

func (m *boomModel) ensurePieces() {
	if len(m.pieces) > 0 || m.width == 0 || m.height == 0 || m.artWidth == 0 || m.artHeight == 0 {
		return
	}

	count := 10
	centerX := float64(m.width-m.artWidth) / 2
	centerY := float64(max(m.height-1-m.artHeight, 0)) / 2

	for i := range count {
		angle := (math.Pi * 2 * float64(i)) / float64(count)
		angle += (rand.Float64() - 0.5) * 0.45
		speed := 0.7 + rand.Float64()*1.2

		m.pieces = append(m.pieces, boomPiece{
			x:  centerX,
			y:  centerY,
			vx: math.Cos(angle) * speed * 2.8,
			vy: math.Sin(angle) * speed * 1.1,
		})
	}
}

func boomTick() tea.Msg {
	time.Sleep(time.Millisecond * 95)
	return boomTickMsg(time.Now())
}

func initialBoomModel(character string) boomModel {
	art := strings.Split(characterArt(character), "\n")
	artWidth := 0
	for _, line := range art {
		artWidth = max(artWidth, len([]rune(line)))
	}

	return boomModel{
		art:       art,
		artWidth:  artWidth,
		artHeight: len(art),
	}
}

func runBoom(character string) {
	p := tea.NewProgram(initialBoomModel(character))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v", err)
	}
}
