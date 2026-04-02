package main

import (
	"fmt"
	"strings"

	"github.com/Rtarun3606k/TakaTime/internal/Styles"
	utils "github.com/Rtarun3606k/TakaTime/internal/Utils"
	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/charmbracelet/lipgloss"
)

// MakeLink wraps text in an OSC 8 ANSI escape sequence to make it clickable
func MakeLink(url, text string) string {
	return fmt.Sprintf("\x1b]8;;%s\x1b\\%s\x1b]8;;\x1b\\", url, text)
}

// Add 'showAll bool' to the end of the parameters
func buildStatsList(title string, stats []types.ListStats, styles Styles.AppStyles, width int, showAll bool) string {
	if len(stats) == 0 {
		return ""
	}

	var b strings.Builder
	// b.WriteString(styles.SubText.Render(fmt.Sprintf("--- %s ---", title)) + "\n")
	// 1. Create a prominent, centered style dynamically based on the box width
	titleStyle := lipgloss.NewStyle().
		Bold(true).             // Make it thick
		Width(width).           // Stretch it to the exact width of the box
		Align(lipgloss.Center). // Perfectly center the text inside that width
		MarginBottom(1)         // Add a blank line below it for breathing room

	// 2. Format the string (Uppercase makes it feel "bigger" and more important)
	formattedTitle := fmt.Sprintf("━ %s ━", strings.ToUpper(title))

	// 3. Render and write it to the builder
	b.WriteString(titleStyle.Render(formattedTitle) + "\n")
	// 1. Determine which stats to show
	limit := 10

	displayStats := stats
	hiddenCount := 0

	// 2. Safely check and slice using the exact same variable
	if !showAll && len(stats) > limit {
		displayStats = stats[:limit]
		hiddenCount = len(stats) - limit
	}

	barWidth := 15

	// 2. Loop through the SLICED array
	for _, stat := range displayStats {
		label := styles.ListLabel.Render(utils.SafeTruncateString(stat.Label, 10))
		value := styles.ListValue.Render(stat.Value)
		percentStr := styles.ListPercent.Render(fmt.Sprintf("%8.1f%%", stat.Percent*100))

		filledCount := int(stat.Percent * float64(barWidth))
		if filledCount > barWidth {
			filledCount = barWidth
		}

		filledBar := styles.Color2.Render(strings.Repeat("█", filledCount))
		emptyBar := styles.SubText.Render(strings.Repeat("░", barWidth-filledCount))
		visualBar := filledBar + emptyBar

		b.WriteString(fmt.Sprintf("%s | %s | %s %s\n\n", label, value, visualBar, percentStr))
	}

	// 3. Add the "See More" text if we are hiding things
	if hiddenCount > 0 {
		indicator := fmt.Sprintf("... and %d more (press 'm')", hiddenCount)
		// Right-align or center the indicator for a clean look
		b.WriteString(styles.SubText.Render(indicator) + "\n")
	}

	return styles.Box.Width(width).Render(b.String())
}

// buildTimeGrid creates a row of 4 horizontal cards for your summary stats
func buildTimeGrid(data types.TimeGridStruct, styles Styles.AppStyles, width int) string {
	// 1. Determine Layout & Calculate Width based on Breakpoints
	var cardWidth int
	var columns int

	if width >= 80 {
		columns = 4
		cardWidth = (width / 4) - 2 // 4 cards across
	} else if width >= 45 {
		columns = 2
		cardWidth = (width / 2) - 2 // 2 cards across
	} else {
		columns = 1
		cardWidth = width - 2 // 1 card across (takes full width)
	}

	// 2. Helper function to build a single card
	buildCard := func(title, value string) string {
		titleBlock := styles.StatCardTitle.Render(title)
		valueBlock := styles.StatCardValue.Render(value)
		content := lipgloss.JoinVertical(lipgloss.Center, titleBlock, valueBlock)

		return styles.StatCard.Width(cardWidth).Render(content)
	}

	// 3. Build the 4 individual cards
	yesterday := buildCard("Yesterday", data.Yestarday)
	week := buildCard("7 Days", data.Week)
	month := buildCard("30 Days", data.Month)
	allTime := buildCard("All Time", data.AllTime)

	// 4. Render the layout based on the calculated columns
	switch columns {
	case 4:
		// Wide screen: All in one row
		return lipgloss.JoinHorizontal(lipgloss.Top, yesterday, "  ", week, "  ", month, "  ", allTime)

	case 2:
		// Medium screen: 2x2 Grid
		row1 := lipgloss.JoinHorizontal(lipgloss.Top, yesterday, "  ", week)
		row2 := lipgloss.JoinHorizontal(lipgloss.Top, month, "  ", allTime)
		// Join the two rows vertically with a blank line between them
		return lipgloss.JoinVertical(lipgloss.Left, row1, "\n", row2)

	default:
		// Small screen: Stacked vertically
		return lipgloss.JoinVertical(lipgloss.Left, yesterday, "\n", week, "\n", month, "\n", allTime)
	}
}

func (m Model) generateScrollableContent() string {
	if m.Loading {
		loadingText := fmt.Sprintf("%s %s",
			m.Spinner.View(),
			m.AppStyles.SubText.Render("Fetching your coding stats..."),
		)
		return lipgloss.Place(
			m.Viewport.Width, m.Viewport.Height,
			lipgloss.Center, lipgloss.Center,
			loadingText,
		)
	}

	var b strings.Builder
	contentWidth := m.Viewport.Width - 4
	if contentWidth < 40 {
		contentWidth = 40
	}

	if m.CacheData {
		b.WriteString(m.AppStyles.SubText.Render("⚡ Loaded from local cache") + "\n\n")
	}

	b.WriteString(buildTimeGrid(m.TimeStats, m.AppStyles, contentWidth))
	b.WriteString("\n\n")

	halfWidth := (contentWidth / 2) - 1
	lanuagesBlock := buildStatsList("Languages", m.LanguageListStats, m.AppStyles, halfWidth, m.ViewMore)
	projectsBlock := buildStatsList("Projects", m.ProjectListStats, m.AppStyles, halfWidth, m.ViewMore)

	flexRowLanguageAndProjects := lipgloss.JoinHorizontal(lipgloss.Top, lanuagesBlock, " ", projectsBlock)
	b.WriteString(flexRowLanguageAndProjects + "\n")
	osBox := buildStatsList("Operating Systems", m.OsListStats, m.AppStyles, halfWidth, m.ViewMore)
	editorBox := buildStatsList("Editors", m.editorListStats, m.AppStyles, halfWidth, m.ViewMore)

	flexRow := lipgloss.JoinHorizontal(lipgloss.Top, osBox, "  ", editorBox)
	b.WriteString(flexRow + "\n")

	// Center the whole block
	return lipgloss.NewStyle().
		Width(m.Viewport.Width).
		Align(lipgloss.Center).
		Render(b.String())
}

// 2. The Main View
func (m Model) View() string {
	if !m.Ready {
		return "Initializing..."
	}

	// --- HEADER ---
	logo := m.AppStyles.Title.Render("TakaTime ")
	links := m.AppStyles.SubText.Render("Home (H) | About (A)")

	gapWidth := m.Width - lipgloss.Width(logo) - lipgloss.Width(links) - 4
	if gapWidth < 0 {
		gapWidth = 0
	}
	header := lipgloss.JoinHorizontal(lipgloss.Center, logo, strings.Repeat(" ", gapWidth), links)

	// 👉 Use MarginTop(1) and MarginBottom(1)
	header = lipgloss.NewStyle().MarginTop(1).MarginBottom(1).Render(header)

	// --- FOOTER ---
	// Add the scroll percentage to the footer!
	scrollPercent := fmt.Sprintf("Scroll: %3.f%%", m.Viewport.ScrollPercent()*100)
	helpText := m.AppStyles.SubText.Render("q: quit • r: refresh • j/k/↑/↓: scroll • " + scrollPercent)

	// 1. Create the clickable ANSI strings
	githubLink := MakeLink("https://github.com/Rtarun3606k/TakaTime", "GitHub: Rtarun3606k")
	discordLink := MakeLink("https://discord.gg/YOUR_DISCORD_LINK", "Discord: TakaTime") // Replace with actual link!

	// 2. Combine them with your pipe separator
	socialsText := fmt.Sprintf("%s | %s", githubLink, discordLink)

	// 3. Render the combined string with your AppStyles
	socials := m.AppStyles.Text.Render(socialsText)

	// Combine into the final footer block
	footerContent := lipgloss.JoinVertical(lipgloss.Center, helpText, socials)
	footer := lipgloss.NewStyle().
		Width(m.Width).
		Align(lipgloss.Center).
		MarginTop(1).
		Render(footerContent)

	// --- COMBINE EVERYTHING ---
	// Notice we use m.Viewport.View() for the middle section!
	return lipgloss.JoinVertical(lipgloss.Top,
		header,
		m.Viewport.View(),
		footer,
	)
}
