package main

import (
	"github.com/chromedp/chromedp/kb"
	b "webdiceatat/base"
)

func examples() {
	var sm b.SiteManager

	sm.Init(PC, 3)

	defer sm.Cancel()

	sm.GoToPath("https://www.w3schools.com/html/html_forms.asp", 10)
	sm.CreateScreenShot("w3school.load.png", 0)
	sm.WaitVisible(b.Input(0).ByAttribute("id", "fname", 0).String(), 0)

	sm.FillFields([]map[string]interface{}{
		{
			`identifier`: b.Input(0).ByAttribute("id", "fname", 0).String(),
			`value`:      "Firstname",
		},
		{
			`identifier`: b.Input(0).ByAttribute("id", "lname", 0).String(),
			`value`:      "Lastname",
		},
	}, 3)

	sm.CreateScreenShot("w3school.filled.png", 0)

	sm.ClickElement(b.Input(0).ByAttribute("type", "submit", 0).ByAttribute("value", "Submit", 0).String(), 0)
	sm.CreateScreenShot("w3school.clicked.png", 0)

	sm.KeyDown(kb.Enter, 3)

	//sm.WaitVisible(b.Button(0).ByAttribute("value", "Kattints ide", 0).String(), 0)
	//sm.ClickElement(b.Button(0).ByAttribute("value", "Kattints ide", 0).String(), 0)
}
