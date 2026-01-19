package main

import (
	"fmt"
	"math"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"golang.org/x/term"
)

// Stílusok és Konstansok
const (
	MAP_CARD_WIDTH   = 40 // A map kártya szélessége
	CHAMP_CARD_WIDTH = 40 // A Champ kártya szélessége
)

var (
	// Színstílusok
	greenStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#00FF00")).Bold(true)
	redStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF0000")).Bold(true)
	faintStyle = lipgloss.NewStyle().Faint(true)

	baseCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("63")).
			Padding(0, 1).
			MarginRight(1).
			MarginBottom(1)
)

// --- SEGÉDFÜGGVÉNYEK ---

// Kiszámolja, hány kártya fér el a terminálban
func getCardsPerRow(cardWidth int) int {
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return 2 // Biztonsági tartalék, ha nem sikerül lekérni
	}
	// A kártya szélessége + a MarginRight(1)
	count := width / (cardWidth + 1)
	if count < 1 {
		return 1
	}
	return count
}

// Színkódolás érték alapján (isGood: a magas érték jó-e?)
func formatPercent(val float64, isGood bool) string {
	str := fmt.Sprintf("%.1f%%", val)
	// Ha a LoseRate magas (>50%), az "rossz", tehát piros. Ha a WR magas, az zöld.
	if (isGood && val >= 50) || (!isGood && val < 50) {
		return greenStyle.Render(str)
	}
	return redStyle.Render(str)
}

// --- RENDERELŐK ---

func RenderCounterMap(result map[string][]ChampCounter) {
	keys := make([]string, 0, len(result))
	for k := range result {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var allCards []string
	cardsPerRow := getCardsPerRow(MAP_CARD_WIDTH)

	for _, name := range keys {
		counters := result[name]
		rows := [][]string{}

		limit := int(math.Min(float64(len(counters)), 3))
		for i := 0; i < limit; i++ {
			c := counters[i]
			rows = append(rows, []string{
				c.Name,
				formatPercent(c.LoseRate, false),
			})
		}

		// Táblázat létrehozása fejléccel
		t := table.New().
			Border(lipgloss.RoundedBorder()).
			BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color("240"))).
			// ITT VAN A KÉRT MÓDOSÍTÁS:
			Headers("CHAMP", fmt.Sprintf("WR vs %s", strings.ToUpper(name))).
			Rows(rows...).
			Render()

		cardContent := fmt.Sprintf("%s\n%s",
			lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFA500")).Render("ROLE: "+strings.ToUpper(name)),
			t,
		)

		allCards = append(allCards, baseCardStyle.Width(MAP_CARD_WIDTH).Render(cardContent))
	}

	printGrid(allCards, cardsPerRow)
}

func GetChampCard(c Champ, index int) string { // Hozzáadtuk az indexet
	wrStr := formatPercent(c.Winrate, true)
	stats := fmt.Sprintf("Global WR: %s | PR: %.1f%%", wrStr, c.Pickrate)

	// Counter Picks felirat
	counterLabel := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("241")).
		MarginTop(1).
		Underline(true).
		Render("Counter Picks:")

	rows := [][]string{}
	limit := int(math.Min(float64(len(c.Counters)), 3))
	for i := 0; i < limit; i++ {
		cnt := c.Counters[i]
		rows = append(rows, []string{cnt.Name, formatPercent(cnt.LoseRate, false)})
	}

	t := table.New().
		Border(lipgloss.NormalBorder()).
		Headers("CHAMP", fmt.Sprintf("WR vs %s", strings.ToUpper(c.Name))).
		Rows(rows...).
		Render()

	// A fejléc, benne a sorszámmal (pl. #1 AATROX)
	headerText := fmt.Sprintf("#%d %s", index, strings.ToUpper(c.Name))

	content := fmt.Sprintf("%s\n%s\n%s\n%s",
		lipgloss.NewStyle().
			Bold(true).
			Background(lipgloss.Color("63")).  // Kék háttér
			Foreground(lipgloss.Color("255")). // Fehér szöveg
			Padding(0, 1).
			Render(headerText), // Itt jelenik meg a sorszám és a név
		stats,
		counterLabel,
		t)

	return baseCardStyle.Width(CHAMP_CARD_WIDTH).Render(content)
}

func RenderAllChamps(champs []Champ) {
	var allCards []string
	cardsPerRow := getCardsPerRow(CHAMP_CARD_WIDTH)

	for i, c := range champs {
		// Átadjuk a sorszámot (i+1)
		allCards = append(allCards, GetChampCard(c, i+1))
	}

	printGrid(allCards, cardsPerRow)
}

// Sorokba tördeli a kártyákat
func printGrid(cards []string, cardsPerRow int) {
	for i := 0; i < len(cards); i += cardsPerRow {
		end := i + cardsPerRow
		if end > len(cards) {
			end = len(cards)
		}
		fmt.Println(lipgloss.JoinHorizontal(lipgloss.Top, cards[i:end]...))
	}
}
