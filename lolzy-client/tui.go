package main

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00D7FF")).
			MarginBottom(1)

	statBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1).
			MarginRight(2)

	winrateStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00"))
	titleStyle   = lipgloss.NewStyle().Italic(true).Foreground(lipgloss.Color("#888888"))
)

func RenderCounterMap(result map[string][]ChampCounter) {
	// 1. Kulcsok kinyerése és rendezése (hogy ne ugráljanak a kártyák)
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Stílusok
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1).
		MarginRight(1).
		MarginBottom(1).
		Width(35) // Fix szélesség a rácshoz

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFA500")). // Narancs fejléc
		Underline(true)

	var allCards []string

	// 2. Kártyák létrehozása
	for _, name := range keys {
		counters := result[name]

		// Táblázat építése az ellenfeleknek
		rows := [][]string{}
		// Csak az első 5 countert mutatjuk, hogy ne legyen túl hosszú a kártya
		maxRows := 5
		for i, c := range counters {
			if i >= maxRows {
				break
			}

			rows = append(rows, []string{
				c.Name,
				fmt.Sprintf("%.1f%%", c.LoseRate),
			})
		}

		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
			Rows(rows...).
			Render()

		// Kártya tartalmának összeállítása
		cardContent := fmt.Sprintf("%s\n%s\n%s",
			headerStyle.Render(strings.ToUpper(name)),
			lipgloss.NewStyle().Faint(true).Render("Legerősebb counterek:"),
			t,
		)

		allCards = append(allCards, cardStyle.Render(cardContent))
	}

	// 3. Rácsba rendezés (pl. 3 kártya egy sorban)
	cardsPerRow := 3
	var finalRows []string

	for i := 0; i < len(allCards); i += cardsPerRow {
		end := i + cardsPerRow
		if end > len(allCards) {
			end = len(allCards)
		}
		// Vízszintesen összefűzzük az aktuális sor kártyáit
		row := lipgloss.JoinHorizontal(lipgloss.Top, allCards[i:end]...)
		finalRows = append(finalRows, row)
	}

	// Függőlegesen összefűzzük a sorokat és kiírjuk
	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, finalRows...))
}

func GetChampCard(c Champ) string {
	// Statisztikák (színes winrate-el)
	wrColor := "#00FF00"
	if c.Winrate < 0.5 {
		wrColor = "#FF0000"
	}

	wrStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(wrColor)).Bold(true)

	stats := fmt.Sprintf("WR: %s | PR: %.1f%%",
		wrStyle.Render(fmt.Sprintf("%.1f%%", c.Winrate)),
		c.Pickrate)

	// Táblázat a countereknek (kompaktabb nézet)
	rows := [][]string{}
	for _, cnt := range c.Counters[:3] { // Csak az első 3 counter, hogy ne legyen túl hosszú
		rows = append(rows, []string{cnt.Name, fmt.Sprintf("%.1f%%", cnt.LoseRate)})
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		Rows(rows...).
		Render()

	// A teljes kártya összeállítása egy keretbe
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("63")).
		Padding(0, 1).
		MarginRight(2). // Térköz a kártyák között
		Width(30)       // Fix szélesség, hogy egyformák legyenek

	content := fmt.Sprintf("%s\n%s\n%s\n%s",
		lipgloss.NewStyle().Bold(true).Underline(true).Render(c.Name),
		stats,
		lipgloss.NewStyle().Faint(true).Render("Top Counters:"),
		t)

	return cardStyle.Render(content)
}

func RenderAllChamps(champs []Champ) {
	var rows []string
	var currentBatch []string

	cardsPerRow := 3 // Hány kártya legyen egymás mellett

	for i, c := range champs {
		currentBatch = append(currentBatch, GetChampCard(c))

		// Ha megtelt egy sor, vagy az utolsó elemnél tartunk
		if (i+1)%cardsPerRow == 0 || i == len(champs)-1 {
			// Vízszintesen összefűzzük a batch-et
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentBatch...))
			currentBatch = []string{} // Reset a következő sorhoz
		}
	}

	// A sorokat függőlegesen összefűzzük és kiírjuk
	fmt.Println(lipgloss.JoinVertical(lipgloss.Left, rows...))
}
