package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/quinn/flash-cards/flash"
	"github.com/quinn/flash-cards/pages"
	"github.com/quinn/flash-cards/spec"
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
	pair := getRandomHiragana(r)
	choices := getRomajiChoicesForHiragana(r, pair.Hiragana)
	component := pages.Index(pair, choices)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

// handleCardHtmx handles HTMX interactions for card flips
func handleCardHtmx(c echo.Context, r *rand.Rand) error {
	var qs spec.QS
	if err := c.Bind(&qs); err != nil {
		return err
	}

	pair := spec.HiraganaRomajiPair{
		Hiragana: qs.Hiragana,
		Romaji:   qs.Answer,
	}

	component := ui.Card(pair, qs.Choices, qs.Answer)
	return component.Render(c.Request().Context(), c.Response().Writer)
}

func getRandomHiragana(r *rand.Rand) spec.HiraganaRomajiPair {
	// Get the map from flash package
	hiraganaMap := flash.GetHiraganaToRomajiMap()

	// Convert map to slice of key-value pairs for random selection
	hiraganaSlice := make([]spec.HiraganaRomajiPair, 0, len(hiraganaMap))

	for k, v := range hiraganaMap {
		hiraganaSlice = append(hiraganaSlice, spec.HiraganaRomajiPair{Hiragana: k, Romaji: v})
	}

	// Select random entry
	randomIndex := r.Intn(len(hiraganaSlice))
	return hiraganaSlice[randomIndex]
}

func getRomajiChoicesForHiragana(r *rand.Rand, hiragana string) []string {
	// Get the map from flash package
	hiraganaMap := flash.GetHiraganaToRomajiMap()

	// Get the romaji for the hiragana
	romaji := hiraganaMap[hiragana]

	// Get 4 other random romaji
	otherRomaji := make([]string, 4)
	for i := 0; i < 4; i++ {
		otherRomaji[i] = hiraganaMap[getRandomHiragana(r).Hiragana]
	}

	// Add the correct romaji to the slice
	romajiSlice := append(otherRomaji, romaji)

	// Shuffle the slice
	r.Shuffle(len(romajiSlice), func(i, j int) {
		romajiSlice[i], romajiSlice[j] = romajiSlice[j], romajiSlice[i]
	})

	return romajiSlice
}
