package main

import (
	"fmt"
	"github.com/chromedp/chromedp/kb"
	b "github.com/dombiistvan/webdice-atat/base"
)

func examples() {
	var sm b.SiteManager

	sm.Init(b.PC, 5)

	defer sm.Cancel()

	sm.GoToPath("https://www.w3schools.com/html/html_forms.asp", 10)
	sm.CreateScreenShot("w3school.load.png", 0)
	sm.WaitVisible(b.Input(0).ByPath("//").ByAttribute("id", "fname", 0).String(), 0)

	sm.FillFields([]map[string]interface{}{
		{
			`identifier`: b.Input(0).ByPath("//").ByAttribute("id", "fname", 0).String(),
			`value`:      "Firstname",
		},
		{
			`identifier`: b.Input(0).ByPath("//").ByAttribute("id", "lname", 0).String(),
			`value`:      "Lastname",
		},
	}, 3)

	sm.CreateScreenShot("w3school.filled.png", 0)

	hrefAttribute, ok := sm.GetElementAttributeValue(b.Html(0).AddChild(b.Body(0).AddChild(b.Div(0).ByEqual("@class", "w3-container top", 0).AddChild(b.Anchor(0).ByContains("@class", "w3schools-logo", 0)))).String(), "href", 0)
	fmt.Println(ok, hrefAttribute)

	elementAttributes := sm.GetElementAttributes(b.Div(0).ByPath("//").ByAttribute("class", "w3-right w3-hide-small w3-wide toptext", 0).String(), 0)
	fmt.Println(elementAttributes)

	multipleElementAttributes := sm.GetElementsAttributes(b.HtmlTag("a", 0).ByContains("@class", "w3-btn", 0).String(), 0)
	fmt.Println(multipleElementAttributes)

	sm.ClickElement(b.Input(0).ByPath("//").ByAttribute("type", "submit", 0).ByAttribute("value", "Submit", 0).String(), 0)
	sm.CreateScreenShot("w3school.clicked.png", 0)

	sm.KeyDown(kb.Enter, 3)

	//sm.WaitVisible(b.Button(0).ByAttribute("value", "Kattints ide", 0).String(), 0)
	//sm.ClickElement(b.Button(0).ByAttribute("value", "Kattints ide", 0).String(), 0)
}
