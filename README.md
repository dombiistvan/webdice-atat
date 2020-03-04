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
