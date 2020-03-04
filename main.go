package main

import (
	"github.com/chromedp/chromedp/device"
	b "webdiceatat/base"
)

var (
	PC = device.Info{
		"pc",
		"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/78.0.3904.108 Safari/537.36",
		1920,
		1080,
		1.7778,
		false,
		false,
		false,
	}
	sonyXPeriaXZPremium = device.Info{
		Name:      "Sony Xpereia XZ Premium",
		UserAgent: "Mozilla/5.0 (Linux; Android 5.1.1; Nexus 5 Build/LMY48B; wv) AppleWebKit/537.36 (KHTML, like Gecko) Version/4.0 Chrome/43.0.2357.65 Mobile Safari/537.36",
		Width:     770,
		Height:    1560,
		Scale:     0.49,
		Landscape: false,
		Mobile:    true,
		Touch:     true,
	}
)

func main() {
	var sm b.SiteManager

	var exampleUrl string = "https://example.com"

	sm.Init(PC, 3)

	defer sm.Cancel()

	sm.GoToPath(exampleUrl, 10)
	sm.CreateScreenShot("example.com.png", 0)
	sm.WaitVisible(b.Header(0).ByContains("text()", "Example Domain", 0).String(), 0)

	sm.CreateScreenShot("done.png", 0)

	//in examples.go
	examples()
}
