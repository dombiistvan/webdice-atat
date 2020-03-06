package main

import (
	b "github.com/dombiistvan/webdice-atat/base"
)

func main() {
	var sm b.SiteManager

	var exampleUrl string = "https://example.com"

	sm.Init(b.PC, 3)

	defer sm.Cancel()

	sm.GoToPath(exampleUrl, 10)
	sm.CreateScreenShot("example.com.png", 0)
	sm.WaitVisible(b.Header(0).ByContains("text()", "Example Domain", 0).String(), 0)

	sm.CreateScreenShot("done.png", 0)

	//in examples.go
	examples()
}
