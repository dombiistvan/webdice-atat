# webdice-atat
A chromedp based interface to make easier to test websites from code

##### The following awesome packages are used in this package  
- https://github.com/chromedp/cdproto/input
- https://github.com/chromedp/chromedp
- https://github.com/chromedp/chromedp/device
##### thanks for the contributors, and look at these packages if you want to extend mine, or you have a similar problem you want to solve

With this interface package, you can use chromedp (at this point) commonest features I implemented in, extended with some XPATH interface.

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

Check out main.go and examples.go for working examples with the chromedp interface and also the XPATH interface usage!

You can start your own end-user test with this interface, by initializing ```SiteManager``` by the following method:
 - first you need to import this package, you can do via the row ```import watat "github.com/dombiistvan/webdice-atat/base"```
 - the next step is to define the variable with type ```SiteManager``` as ```var sm watatb.SiteManager```
 - the following step is to initialize sitemanager with ```sm.Init(...)``` - it requires two parameters: 
   - device: There are two predefinied devices you can use:
      - ```watat.PC``` (Chrome UserAgent, fullHD resolution, scale 1.7778, name "pc", Landscape false, Mobile false, and Touch false options)
      - ```watat.SonyXPeriaXZPremium``` (Android UserAgent, 770x1560 resolution, 0.49 scale, name "Sony Xpereia XZ Premium", Landscape false, Mobile true, and Touch true options)
      - (if you dont want to use these devices, or you want to set up your own device, you can do it by using chromedp ```device.Info{}``` struct, it contains the mentioned options and you can just pass that to the SiteManager Init function)
   - timeout: this is an int representation of the timeout seconds 
 - after Initialized the SiteManager struct, set up the Cancel with ```defer sm.Cancel()```, and now you can start with the concrete testing of your site, for example ```sm.GoToPath("https://www.example.com",10)```, ```sm.CreateScreenshot("example.com.png")``` etc.
 
 The selectors when you want to wait/click for, or do any other action with the elements, are XPATH selectors. You can eighter type the selector by hand, as a string, when you pass it, or you can use my predefined functions, and call the ```.String()``` or ```string(element)```.
 There are predefined Tag caller functions: ```Button(); Input(); Textarea(); Submit(); Header(); Header2(); Header3(); Header4(); Div(); Paragraph()```, but you can use ```HtmlTag()``` function too to use other html tags what are not defined with an own function yet. The predefined functions has only one parameter, the position what is the html tag position under the path you gave.
 For example ```string(Header(2))``` or ```Header(2).String()``` both will return a string with XPATH selector to //h2[2]. Let check an example if you need an html tag I didn't predefined with an own function yet, fe.: article. Get the article on the 3th position. The XPATH selector generation would be the following: ```string(HtmlTag("article",3))``` or ```HtmlTag("article",3).String()``` is the same.
 
 Now get a longer example with more filter.
 
 Lets just have an XPATH selector like: 
  - //html/body/div[contains(@class,"something")][@someAttribute="groupAttribute"][3] 
    - The function will be ```Div(0).ByPath("//html/body").ByContains("@class","something",0).ByAttribute("someAttribute","groupAttribute",3)```
    The position (3) is comes in the ByAttribute, because it is after filtering the tag. If you would like to have a reverse filtering, for example ..../div[3]... you should put the position to the ```ByTag()``` function.
    
 There are 3 kind of filtering functions right now, ```ByContains(subject, value, position)```, ```ByEqual(subject, value, position)``` and ```ByAttribute(attribute, value, position)```
  - ByContains make a string like htmlTag[contains(subject,value)][position only if greater than 0] -> ```.../div[contains(@id,"-list-product")][1]``` 
  - ByEqual make a string like htmlTag[subject="value"][position only if greater than 0] -> ```.../div[text()="First Product")]``` 
  - ByAttribute make a string like htmlTag[@attribute="value"][position only if greater than 0] -> ```.../div[@class="list-product")][1]``` 
  
  #### NESTED XPATH SELECTORS ARE SUPPORTED ALREADY
  
  You can now making nested XPATH selectors by calling on the parent the ```AddChild(child)``` function, and passing the children element. For example, if you want to have a selector like
  //html/div[@class="product-list-item-container"]/div[contains(@class,"list-item")][@id="first-product"][1] you can make it by the following code: 
```
    base.HtmlTag("html",0).AddChild(
        base.Div(0).ByAttribute("class","product-list-item-container",0).AddChild(
            base.Div(0).ByContains("@class","list-item",0).ByEqual("@id","first-product",1)
        )
    )
```
  
  
  
