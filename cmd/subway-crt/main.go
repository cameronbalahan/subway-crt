package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-isatty"
	"github.com/muesli/termenv"
	"subway-crt/internal/config"
	"subway-crt/internal/mta"
	"subway-crt/internal/transit"
	"subway-crt/internal/tui"
	"subway-crt/internal/web"
)

func main() {
	// Enforce truecolor output for vibrant 24-bit Outrun neon aesthetics
	lipgloss.SetColorProfile(termenv.TrueColor)

	cfg := config.Load()

	// Multi-channel broadcast to prevent channel contention between TUI and Web
	tuiChan := make(chan transit.TransitSnapshot, 1)
	webChan := make(chan transit.TransitSnapshot, 1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
		os.Exit(0)
	}()

	// Isolated GTFS-RT Polling Worker
	client := mta.NewClient(cfg.HTTPTimeout)
	go transit.RunPollingWorker(ctx, client, cfg, tuiChan, webChan)

	// Bubble Tea TUI Model (Target: 1080p console 240x67 chars)
	rootModel := tui.NewRootModel(tuiChan)

	// Web Server for browser viewing & AI Studio port 3000
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	webServer := web.NewServer(port)

	// Render loop for web server
	go func() {
		webModel := tui.NewRootModel(nil)
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		var latestSnap transit.TransitSnapshot
		for {
			select {
			case <-ctx.Done():
				return
			case snap := <-webChan:
				latestSnap = snap
				webServer.UpdateFrame(webModel.ViewWithSnapshot(latestSnap), latestSnap)
			case t := <-ticker.C:
				m, _ := webModel.Update(tui.TickMsg(t))
				if updated, ok := m.(tui.RootModel); ok {
					webModel = updated
				}
				// Decrement local seconds
				for i := range latestSnap.ManhattanTrains {
					if latestSnap.ManhattanTrains[i].SecondsAway > 0 {
						latestSnap.ManhattanTrains[i].SecondsAway--
						latestSnap.ManhattanTrains[i].MinutesAway = latestSnap.ManhattanTrains[i].SecondsAway / 60
					}
				}
				for i := range latestSnap.CanarsieTrains {
					if latestSnap.CanarsieTrains[i].SecondsAway > 0 {
						latestSnap.CanarsieTrains[i].SecondsAway--
						latestSnap.CanarsieTrains[i].MinutesAway = latestSnap.CanarsieTrains[i].SecondsAway / 60
					}
				}
				webServer.UpdateFrame(webModel.ViewWithSnapshot(latestSnap), latestSnap)
			}
		}
	}()

	go func() {
		_ = webServer.Start()
	}()

	// Determine if running in an interactive terminal (laptop/Raspberry Pi) or non-interactive daemon
	isTerminal := isatty.IsTerminal(os.Stdin.Fd()) || isatty.IsTerminal(os.Stdout.Fd())
	if os.Getenv("HEADLESS") == "true" || !isTerminal {
		select {
		case <-ctx.Done():
		}
		return
	}

	// Native Terminal TUI mode (e.g. running on laptop, SSH, or /dev/tty1)
	p := tea.NewProgram(rootModel, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "CRT terminal error: %v\n", err)
		os.Exit(1)
	}
}
