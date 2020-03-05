# webdice-atat
A chromedp based interface to make easier to test websites from code

With this interface package, you can use chromedp (at this point) commonest features I implemented in, extended with some xpath interface.

for example, you can:
  - visit an url (GoToPath)
  - fill a form with data (FillFields)
  - click an element (ClickElement)
  - press a key down (KeyDown/raw version: KeyRawDown) 
  - release a key up (KeyUp)
  - key char event (KeyChar)
  - wait an element to be visible (WaitVisible)
  - wait an element not to be visible (WaitNotVisible)
  - make screenshot (CreateScreenShot)
  - wait an element to be enabled (WaitEnabled)
  - wait an element to be ready (WaitReady)
  - wait an element to be selected (WaitSelected)
  - wait an element not to be presented (WaitNotPresent)
  - just wait (Wait)
  
and all of these actions with own timeout

In the future we plan to extends these features with implementing more chromedp action

Check out main.go and examples.go for working examples with the chromedp interface and also the xpath interface usage!

You can start your own end-user test with this interface, by initializing ```SiteManager``` by the following method:
 - first you need to import this package, you can do via the rows ```import watat "github.com/dombiistvan/webdice-atat"``` (you can skip import this if you want to define your own devices for testing, but we have 2 predefined device will use in this) and ```import watatb "github.com/dombiistvan/webdice-atat/base"```
 - the next step is to define the variable with type ```SiteManager``` as ```var sm watatb.SiteManager```
 - the following is to initialize sitemanager with ```sm.Init()```. Init requires two parameters: 
   - device: There are two predefinied devices you can use:
      - ```watat.PC``` (Chrome UserAgent, fullHD resolution, scale 1.7778, name "pc", Landscape false, Mobile false, and Touch false options)
      - ```watat.SonyXPeriaXZPremium``` (Android UserAgent, 770x1560 resolution, 0.49 scale, name "Sony Xpereia XZ Premium", Landscape false, Mobile true, and Touch true options)
      - (if you dont want to use these devices, or you want to set up your own device, you can do it by using chromedp ```device.Info{}``` struct, it contains the mentioned options and you can pass that to the SiteManager Init function.. and in this case, you don't need both import, only the one ends with /base)
   - timeout: this is an int representation of the timeout seconds 
 - after Initialized the SiteManager struct, set up the Cancel with ```defer sm.Cancel()```, and now you can start with the concrete testing of your site. ```sm.GoToPath("https://www.example.com",10)```, ```sm.CreateScreenshot("example.com.png")``` etc.
