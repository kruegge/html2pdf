package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/raff/godet"
)

func main() {
	fmt.Println("XHTML2PDF with godet")

	// connect to Chrome instance
	remote, err := godet.Connect("localhost:9222", false)
	if err != nil {
		fmt.Println("cannot connect to Chrome instance:", err)
		return
	}

	// disconnect when done
	defer remote.Close()

	// get browser and protocol version
	version, _ := remote.Version()
	fmt.Println(version)

	// install some callbacks
	remote.CallbackEvent(godet.EventClosed, func(params godet.Params) {
		fmt.Println("RemoteDebugger connection terminated.")
	})

	// re-enable events when changing active tab
	_ = remote.AllEvents(true) // enable all events

	_, _ = remote.Navigate("https://www.ab-in-den-urlaub.de")
	_ = remote.SetVisibleSize(1920, 1080)

	time.Sleep(time.Second * 8)

	doc, err := remote.GetDocument()
	if err != nil {
		log.Fatal(err)
	}

	type Document struct {
		Root struct {
			Children []struct {
				Children []struct {
					LocalName string `json:"localName"`
					NodeID    int    `json:"nodeId"`
					NodeName  string `json:"nodeName"`
				} `json:"children,omitempty"`
			} `json:"children"`
		} `json:"root"`
	}

	var document Document
	jsona, _ := json.Marshal(doc)
	if err := json.Unmarshal(jsona, &document); err != nil {
		log.Fatal(err)
	}

	var nodeId int
	for _, v := range document.Root.Children {
		for _, e := range v.Children {
			if e.NodeName == "BODY" {
				nodeId = e.NodeID
			}
		}
	}

	box, err := remote.GetBoxModel(nodeId)
	if err != nil {
		log.Fatal(err)
	}

	type Box struct {
		Model struct {
			Height int `json:"height"`
			Width  int `json:"width"`
		} `json:"model"`
	}

	var data Box
	jsond, _ := json.Marshal(box)
	if err := json.Unmarshal(jsond, &data); err != nil {
		log.Fatal(err)
	}

	// take a screenshot
	remote.SetVisibleSize(1920, data.Model.Height+100)
	_ = remote.SaveScreenshot("screenshot.png", 0644, 0, true)

	// calculate correct dimmensions for PDF creation
	pdfDimensionWidth := 8.2
	pdfRatio := float64(data.Model.Width/72) / pdfDimensionWidth
	pdfDimensionHeight := float64(float64(data.Model.Height) / pdfRatio / 72)

	// or save page as PDF
	_ = remote.SavePDF("page.pdf", 0644,
		godet.Margins(0, 0, 0, 0),
		godet.PrintBackground(),
		godet.Dimensions(pdfDimensionWidth, pdfDimensionHeight),
		godet.Scale(0.4))
}
