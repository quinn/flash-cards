package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quinn/flash-cards/flash"
	"github.com/quinn/flash-cards/pages"
	"github.com/quinn/flash-cards/ui"
)

func main() {
	// Initialize random with a source based on time
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize Echo
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Serve static files
	e.Static("/static", "static")

	// Routes
	e.GET("/", func(c echo.Context) error {
		return handleIndex(c, r)
	})

	// HTMX route for card interactions
	e.GET("/card", func(c echo.Context) error {
		return handleCardHtmx(c, r)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Start server
	e.Logger.Fatal(e.Start(":" + port))
}

// handleIndex serves the initial page load
func handleIndex(c echo.Context, r *rand.Rand) error {
	// Generate a random hiragana for initial page load
	hiragana, romaji := getRandomHiragana(r)
	component := pages.Index(hiragana, romaji, false)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// handleCardHtmx handles HTMX interactions for card flips
func handleCardHtmx(c echo.Context, r *rand.Rand) error {
	// Extract values from the HTMX request
	hiragana := c.QueryParam("hiragana")
	romaji := c.QueryParam("romaji")
	showRomaji := c.QueryParam("showRomaji") == "true"

	if showRomaji {
		// If we were showing romaji, get a new card and don't show romaji
		hiragana, romaji = getRandomHiragana(r)
		component := ui.Card(hiragana, romaji, false)
		return component.Render(c.Request().Context(), c.Response().Writer)
	} else {
		// If we weren't showing romaji, keep the same card and show romaji
		component := ui.Card(hiragana, romaji, true)
		return component.Render(c.Request().Context(), c.Response().Writer)
	}
}

func getRandomHiragana(r *rand.Rand) (string, string) {
	// Get the map from flash package
	hiraganaMap := flash.GetHiraganaToRomajiMap()

	// Convert map to slice of key-value pairs for random selection
	hiraganaSlice := make([]struct {
		Hiragana string
		Romaji   string
	}, 0, len(hiraganaMap))

	for k, v := range hiraganaMap {
		hiraganaSlice = append(hiraganaSlice, struct {
			Hiragana string
			Romaji   string
		}{k, v})
	}

	// Select random entry
	randomIndex := r.Intn(len(hiraganaSlice))
	return hiraganaSlice[randomIndex].Hiragana, hiraganaSlice[randomIndex].Romaji
}
