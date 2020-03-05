# webdice-atat
A chromedp based interface to make easier testing websites from code

##### The following awesome packages are used for this codebase  
- https://github.com/chromedp/cdproto/input
- https://github.com/chromedp/chromedp
- https://github.com/chromedp/chromedp/device
##### thanks for the contributors, and look at these go libaries if you want to extend mine, or you have a similar problem you want to solve

With this interface package, you can use chromedp's few features implemented in, extended with some XPATH interface.

for example, you can:
  - visit an url (GoToPath)
  - fill a form with data (FillFields)
  - click on an element (ClickElement)
  - press a key down (KeyDown/raw version: KeyRawDown) 
  - release a key up (KeyUp)
  - key char event (KeyChar)
  - wait for an element to be visible (WaitVisible)
  - wait for an element not to be visible (WaitNotVisible)
  - make a screenshot (CreateScreenShot)
  - wait for an element to be enabled (WaitEnabled)
  - wait for an element to be ready (WaitReady)
  - wait for an element to be selected (WaitSelected)
  - wait for an element not to be presented (WaitNotPresent)
  - just wait (Wait)
  
and all of these actions with own timeout

In the future we plan to extend these features by implementing more chromedp action

Check out main.go and examples.go for working examples with the chromedp interface and also the XPATH interface usage!

You can start your own end-user test with this interface, by initializing ```SiteManager``` by the following method:
 - first you need to import this package, you can do it via the row ```import watat "github.com/dombiistvan/webdice-atat/base"```
 - the next step is to define the variable with type ```SiteManager``` as ```var sm watat.SiteManager```
 - the following step is to initialize sitemanager with ```sm.Init(...)``` - it requires two parameters: 
   - device: There are two predefinied devices you can use:
      - ```watat.PC``` (Chrome UserAgent, fullHD resolution, scale 1.7778, name "pc", Landscape false, Mobile false, and Touch false options)
      - ```watat.SonyXPeriaXZPremium``` (Android UserAgent, 770x1560 resolution, 0.49 scale, name "Sony Xpereia XZ Premium", Landscape false, Mobile true, and Touch true options)
      - (if you dont want to use these devices, or you want to set up your own device, you can do it by using chromedp ```device.Info{}``` struct, it contains the mentioned options and you can just pass that to the SiteManager Init function)
   - timeout: this is an integer value for set the timeoutin seconds
 - after Initialized the SiteManager, set up destroying the session after all, by ```defer sm.Cancel()```, and now you can start with the concrete testing of your site...
 
 The selectors - that you want to wait for/click on, or do any other action with, - are XPATH selectors. You can either type the selector by hand as a string when you pass it, or you can use the predefined functions, and call the ```.String()``` or cast the struct to string: ```string(element)```.
 There are predefined HTML TAG caller functions: ```Button(); Input(); Textarea(); Submit(); Header(); Header2(); Header3(); Header4(); Div(); Paragraph()```, but you can use ```HtmlTag()``` function too, to use other html tags, that are not defined with an individual function (yet). These functions have only one parameter, the position, what is the html tag's position under the path you gave.
 For example ```string(Header(2))``` or ```Header(2).String()``` both will return a string with XPATH selector to //h2[2]. Let's check an example, if you need an html tag We didn't predefined with a function yet: ```<article>```. Get the article on the 3th position. The XPATH selector generation would be the following: ```string(HtmlTag("article",3))``` or ```HtmlTag("article",3).String()``` is the same.
 
 Now let's see a longer example with more filter in the xpath.
 
 Let's just have an XPATH selector like: 
  - //html/body/div[contains(@class,"something")][@someAttribute="groupAttribute"][3] 
    - The function will be ```watat.HtmlTag("html",0).AddChild(watat.HtmlTag("body",0).AddChild(Div(0).ByContains("@class","something",0).ByAttribute("someAttribute","groupAttribute",3)))```
    The position (3) comes in the ByAttribute method, because it is after filtering the tag. If you would like to have a reverse filtering, for example ..../div[3]... you should put the position to the ```ByTag()``` function.
    
 There are 3 type of filtering function right now, ```ByContains(subject, value, position)```, ```ByEqual(subject, value, position)``` and ```ByAttribute(attribute, value, position)```
  - ByContains makes a string like htmlTag[contains(subject,value)][position only if greater than 0] -> ```.../div[contains(@id,"-list-product")][1]``` 
  - ByEqual make a string like htmlTag[subject="value"][position only if greater than 0] -> ```.../div[text()="First Product")]``` 
  - ByAttribute make a string like htmlTag[@attribute="value"][position only if greater than 0] -> ```.../div[@class="list-product")][1]``` 
  
  #### NESTED XPATH SELECTORS ARE SUPPORTED ALREADY
  
  You can now create nested XPATH selectors by calling on the parent the ```AddChild(child)``` function, and passing the children element. For example, if you want to have a selector like
  //html/div[@class="product-list-item-container"]/div[contains(@class,"list-item")][@id="first-product"][1] you can make it by the following code: 
```
    base.HtmlTag("html",0).AddChild(
        base.Div(0).ByAttribute("class","product-list-item-container",0).AddChild(
            base.Div(0).ByContains("@class","list-item",0).ByEqual("@id","first-product",1)
        )
    )
```

..and you have to cast to string or call the .String() method on only the outest element, it will render the entire path until it has no more nested children.
  
  
