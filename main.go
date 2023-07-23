package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

type Quote struct {
	Content string `json:"content"`
	Author string `json:"author"`

}

func getQuote(client *http.Client) (Quote) {
	res, err := client.Get("https://api.quotable.io/quotes/random")
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	var quotes []Quote 
	decodingErr := json.Unmarshal([]byte(body), &quotes); if decodingErr != nil {
		fmt.Println("error could not unmarshall")
		fmt.Println(decodingErr)
	}

	return quotes[len(quotes)-1]
}

func main() {
	client := http.Client{Timeout: time.Second + 10}
	basicTheApp := app.New()
	window := basicTheApp.NewWindow("Random Quote")
	window.Resize(fyne.NewSize(800, 500))

	quoteLabel := widget.NewLabel("")
	quoteLabel.Alignment = fyne.TextAlignCenter
	quoteLabel.Wrapping = fyne.TextWrapWord

	authorLabel := widget.NewLabel("")
	authorLabel.Alignment = fyne.TextAlignCenter

	quoteButton := widget.NewButton("Get Random Quote", func() {
		quote := getQuote(&client)
		quoteLabel.SetText(quote.Content)
		authorLabel.SetText(quote.Author)
	})
	quoteButton.Resize(fyne.NewSize(100, 100))
	hStack := container.NewHBox(
		layout.NewSpacer(),
		quoteButton,
		layout.NewSpacer(),
	)
	vStack := container.NewVBox(
		hStack,
		layout.NewSpacer(),
		quoteLabel,
		authorLabel,
		layout.NewSpacer(),
		
	)
	window.SetContent(container.NewPadded(vStack))
	window.ShowAndRun()
}
