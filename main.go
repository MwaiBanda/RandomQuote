package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Quote struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

func getQuote(client *http.Client, channel chan<- Quote, waitGroup *sync.WaitGroup) {
	res, err := client.Get("https://api.quotable.io/quotes/random")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	defer waitGroup.Done()
	body, _ := io.ReadAll(res.Body)
	var quotes []Quote
	decodingErr := json.Unmarshal(body, &quotes)
	if decodingErr != nil {
		fmt.Println("error could not unmarshall")
		fmt.Println(decodingErr)
	}

	channel <- quotes[len(quotes)-1]
}

func getQuotes(client *http.Client, onQuoteAdded func(quote Quote)) {
	startTime := time.Now()
	quoteChannel := make(chan Quote)
	var waitGroup sync.WaitGroup

	for i := 0; i < 100; i++ {
		waitGroup.Add(1)
		go getQuote(client, quoteChannel, &waitGroup)
	}

	go func() {
		waitGroup.Wait()
		close(quoteChannel)
	}()

	for quote := range quoteChannel {
		fmt.Println(quote)
		onQuoteAdded(quote)
	}

	fmt.Printf("Processed 100 Requests in %s", time.Since(startTime))
}

func main() {
	client := http.Client{Timeout: time.Second * 10}
	basicTheApp := app.New()
	window := basicTheApp.NewWindow("Random Quote")
	window.Resize(fyne.NewSize(1920, 1080))

	var quotes []Quote

	quotesList := widget.NewList(
		func() int {
			return len(quotes)
		},
		func() fyne.CanvasObject {
			return widget.NewLabelWithStyle("template", fyne.TextAlignLeading, fyne.TextStyle{Bold: true, Symbol: true})
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(quotes[i].Content + " - " + quotes[i].Author)
			o.(*widget.Label).Wrapping = fyne.TextTruncate
		})
	quoteButton := widget.NewButtonWithIcon("Get Random Quotes", theme.MediaReplayIcon(), func() {
		quotes = []Quote{}
		getQuotes(&client, func(quote Quote) {
			quotes = append(quotes, quote)
			quotesList.Refresh()
			time.Sleep(time.Millisecond * 15)
		})
	})

	quoteButton.Resize(fyne.NewSize(100, 100))

	hStack := container.NewHBox(
		quoteButton,
		layout.NewSpacer(),
	)

	window.SetContent(container.NewPadded(container.NewBorder(hStack, nil, nil, nil, container.NewMax(quotesList))))
	go func() {
		window.Show()
		getQuotes(&client, func(quote Quote) {
			quotes = append(quotes, quote)
			quotesList.Refresh()
			time.Sleep(time.Millisecond * 15)
		})
	}()
	basicTheApp.Run()
}
