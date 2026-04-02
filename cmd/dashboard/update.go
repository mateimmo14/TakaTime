package main

import (
	"github.com/Rtarun3606k/TakaTime/internal/db"
	"github.com/Rtarun3606k/TakaTime/internal/types"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

type dataLoadedMsg struct {
	updatedModel Model
	err          error
	FromCache    bool
}

func fetchData(uri string) tea.Cmd {
	return func() tea.Msg {
		// Initialize SQLite connection
		sqliteDB, err := db.InitSQLite()
		if err == nil {
			defer sqliteDB.Close()

			// Check Cache First
			cachedData, err := db.GetDashboardCache(sqliteDB)
			if err == nil && cachedData != nil {
				// Cache HIT
				tempModel := Model{
					LanguageListStats: cachedData.Languages,
					ProjectListStats:  cachedData.Projects,
					OsListStats:       cachedData.OS,
					editorListStats:   cachedData.Editors,
					TimeStats:         cachedData.TimeStats,
				}
				return dataLoadedMsg{updatedModel: tempModel, err: nil, FromCache: true}
			}
		}

		// Cache MISS! Fetch fresh from MongoDB
		tempModel := Model{}
		filledModel, _, err := tempModel.GetData(uri)
		if err != nil {
			return dataLoadedMsg{updatedModel: tempModel, err: err}
		}

		// Save to SQLite Cache for next time
		if sqliteDB != nil {
			db.SaveDashboardCache(sqliteDB, types.CacheData{
				Languages: filledModel.LanguageListStats,
				Projects:  filledModel.ProjectListStats,
				OS:        filledModel.OsListStats,
				Editors:   filledModel.editorListStats,
				TimeStats: filledModel.TimeStats,
			})
		}

		return dataLoadedMsg{
			updatedModel: filledModel,
			err:          nil,
			FromCache:    false,
		}
	}
}

// ----------------------------
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmd  tea.Cmd
		cmds []tea.Cmd
	)

	switch msg := msg.(type) {

	//spinner

	case spinner.TickMsg:
		var spinCmd tea.Cmd
		m.Spinner, spinCmd = m.Spinner.Update(msg)
		cmds = append(cmds, spinCmd)

		// If we are still loading, return immediately so the UI repaints the new frame
		if m.Loading {
			return m, tea.Batch(cmds...)
		}

		return m, cmd

	case tea.WindowSizeMsg:
		// Save the raw dimensions
		m.Width = msg.Width
		m.Height = msg.Height

		// We know our Header is roughly 3 lines, and Footer is 3 lines.
		// So the vertical margin we need to subtract from the viewport is 6.
		headerHeight := 6
		footerHeight := 3
		verticalMargin := headerHeight + footerHeight

		if !m.Ready {
			// 1. Initialize the viewport on the first render
			m.Viewport = viewport.New(msg.Width, msg.Height-verticalMargin)
			m.Viewport.YPosition = headerHeight                  // Start it below the header
			m.Viewport.SetContent(m.generateScrollableContent()) // Helper func we will write
			m.Ready = true
		} else {
			// 2. If already ready, just resize it
			m.Viewport.Width = msg.Width
			m.Viewport.Height = msg.Height - verticalMargin
			m.Viewport.SetContent(m.generateScrollableContent())
		}
		return m, nil

	case dataLoadedMsg:
		m.Loading = false
		m.CacheData = msg.FromCache

		//  ASSIGN THE DATA HERE!
		if msg.err == nil {
			m.LanguageListStats = msg.updatedModel.LanguageListStats
			m.ProjectListStats = msg.updatedModel.ProjectListStats
			m.OsListStats = msg.updatedModel.OsListStats
			m.editorListStats = msg.updatedModel.editorListStats
			m.TimeStats = msg.updatedModel.TimeStats
		}

		// 3. Update the viewport content now that we have data!
		if m.Ready {
			m.Viewport.SetContent(m.generateScrollableContent())
		}
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		case "r":
			m.Loading = true
			return m, fetchData(m.MongoURI)

		case "m":
			// Toggle the boolean (True becomes False, False becomes True)
			m.ViewMore = !m.ViewMore

			// Recalculate the viewport content with the new size!
			if m.Ready {
				m.Viewport.SetContent(m.generateScrollableContent())
			}
			return m, nil

		}
	}

	// 4. Pass ANY unhandled messages (like scrolling keys) to the viewport!
	if m.Ready {
		m.Viewport, cmd = m.Viewport.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}
