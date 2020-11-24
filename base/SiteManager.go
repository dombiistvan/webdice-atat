package base

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"io/ioutil"
	"reflect"
	"time"
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
	SonyXPeriaXZPremium = device.Info{
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

type SiteManager struct {
	ctx          context.Context
	cancel       context.CancelFunc
	info         chromedp.Device
	errorHandler func(err error)
	timeoutSec   int64

	fixActions []chromedp.Action
}

func (sm *SiteManager) Init(d chromedp.Device, defTimeoutSec int64) {
	sm.info = d
	sm.ctx, sm.cancel = chromedp.NewContext(context.Background())
	sm.fixActions = append(sm.fixActions, chromedp.EmulateViewport(sm.info.Device().Width, sm.info.Device().Height))
	sm.timeoutSec = defTimeoutSec
}

func (sm *SiteManager) GoToPath(url string, timoutSec int64) {
	err := sm.DoTimeoutContext(timoutSec, chromedp.Navigate(url))

	if err != nil {
		panic(err)
	}
}

func (sm *SiteManager) CreateScreenShot(filename string, timeoutSec int64) {
	var p []byte

	err := sm.DoTimeoutContext(timeoutSec, chromedp.CaptureScreenshot(&p))

	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(filename, p, 0755)

	if err != nil {
		panic(err)
	}
}

func (sm *SiteManager) Cancel() {
	sm.cancel()
}

func (sm SiteManager) ById(path, tag, id string) string {
	return sm.ByAttribute(path, tag, "id", id)
}

func (sm SiteManager) ByAttribute(path, tag, attributeKey, attributeValue string) string {
	return fmt.Sprintf(`%s%s[%s="%s"]`, path, tag, attributeKey, attributeValue)
}

func (sm *SiteManager) FillFields(fields []map[string]interface{}, timeoutSec int64) {
	var actions []chromedp.Action

	for _, fd := range fields {
		identifier, kok := fd["identifier"]
		value, vok := fd["value"]
		options, qok := fd["options"]
		if kok && vok && !qok {
			actions = append(actions, chromedp.SendKeys(identifier, value.(string)))
		} else if kok && vok {
			if reflect.TypeOf(options) != reflect.TypeOf([]chromedp.QueryOption{}) {
				panic(errors.New("options must be instance of []chromedp.QueryOption"))
			}
			actions = append(actions, chromedp.SendKeys(identifier, value.(string), options.([]chromedp.QueryOption)...))
		}
	}

	err := sm.DoTimeoutContext(timeoutSec, actions...)
	sm.Error(err)
}

func (sm *SiteManager) FillField(identifier string, value string, timeoutSec int64, options ...chromedp.QueryOption) {
	var actions []chromedp.Action

	actions = append(actions, chromedp.SendKeys(identifier, value, options...))

	err := sm.DoTimeoutContext(timeoutSec, actions...)
	sm.Error(err)
}

func (sm *SiteManager) WaitEnabled(selector string, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.WaitEnabled(selector))
	sm.Error(err)
}

func (sm *SiteManager) WaitNotPresent(selector string, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.WaitNotPresent(selector))
	sm.Error(err)
}

func (sm *SiteManager) WaitNotVisible(selector string, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.WaitNotVisible(selector))
	sm.Error(err)
}

func (sm *SiteManager) WaitVisible(selector string, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.WaitVisible(selector))
	sm.Error(err)
}

func (sm *SiteManager) WaitSelected(selector string, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.WaitSelected(selector))
	sm.Error(err)
}

func (sm *SiteManager) WaitReady(selector string, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.WaitReady(selector))
	sm.Error(err)
}

/*func (sm *SiteManager) PressButton(selector string, timeoutSec int64 ,options... chromedp.QueryOption){
	err := sm.DoTimeoutContext(timeoutSec,chromedp.Click(selector, options...))
	if err != nil{
		panic(err)
	}
}*/

func (sm *SiteManager) ClickElement(selector string, timeoutSec int64, options ...chromedp.QueryOption) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.Click(selector, options...))
	sm.Error(err)
}

func (sm *SiteManager) Wait(secs int64) {
	err := sm.DoTimeoutContext(0, chromedp.Sleep(time.Second*time.Duration(secs)))
	sm.Error(err)
}

func (sm *SiteManager) CustomAction(action chromedp.ActionFunc, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, action)
	sm.Error(err)
}

func (sm *SiteManager) FocusElement(selector string, timeoutSec int64, options ...chromedp.QueryOption) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.Focus(selector, options...))
	sm.Error(err)
}

func (sm *SiteManager) ClearElement(selector string, timeoutSec int64, options ...chromedp.QueryOption) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.Clear(selector, options...))
	sm.Error(err)
}

func (sm *SiteManager) DoubleClickElement(selector string, timeoutSec int64, options ...chromedp.QueryOption) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.DoubleClick(selector, options...))
	sm.Error(err)
}

func (sm *SiteManager) InnerHtmlInto(selector string, timeoutSec int64, html *string) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.InnerHTML(selector, html))
	sm.Error(err)
}

func (sm *SiteManager) TextInto(selector string, timeoutSec int64, text *string) {
	err := sm.DoTimeoutContext(timeoutSec, chromedp.Text(selector, text))
	sm.Error(err)
}

func (sm *SiteManager) GetElementAttributeValue(selector string, attribute string, timeoutSec int64, options ...chromedp.QueryOption) (string, bool) {
	var attributeValue string
	var ok bool
	err := sm.DoTimeoutContext(timeoutSec, chromedp.AttributeValue(selector, attribute, &attributeValue, &ok, options...))
	sm.Error(err)

	return attributeValue, ok
}

func (sm *SiteManager) GetElementAttributes(selector string, timeoutSec int64, options ...chromedp.QueryOption) map[string]string {
	var attributes map[string]string
	err := sm.DoTimeoutContext(timeoutSec, chromedp.Attributes(selector, &attributes, options...))
	sm.Error(err)

	return attributes
}

func (sm *SiteManager) GetElementsAttributes(selector string, timeoutSec int64, options ...chromedp.QueryOption) []map[string]string {
	var attributes []map[string]string
	err := sm.DoTimeoutContext(timeoutSec, chromedp.AttributesAll(selector, &attributes, options...))
	sm.Error(err)

	return attributes
}

func (sm *SiteManager) KeyDown(key string, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, input.DispatchKeyEvent(input.KeyDown).WithKey(key))
	sm.Error(err)
}

func (sm *SiteManager) KeyRawDown(key string, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, input.DispatchKeyEvent(input.KeyRawDown).WithKey(key))
	sm.Error(err)
}

func (sm *SiteManager) KeyUp(key string, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, input.DispatchKeyEvent(input.KeyUp).WithKey(key))
	sm.Error(err)
}

func (sm *SiteManager) KeyChar(key string, timeoutSec int64) {
	err := sm.DoTimeoutContext(timeoutSec, input.DispatchKeyEvent(input.KeyChar).WithKey(key))
	sm.Error(err)
}

func (sm SiteManager) getActions(actions ...chromedp.Action) []chromedp.Action {
	return append(sm.fixActions, actions...)
}

func (sm SiteManager) GetDefaultTimeoutSecs() int64 {
	return sm.timeoutSec
}

func (sm SiteManager) GetTimeoutDurationSecs(secs int64) time.Duration {
	return time.Duration(secs) * time.Second
}

func (sm *SiteManager) DoTimeoutContext(timeoutSec int64, action ...chromedp.Action) error {
	if timeoutSec == 0 && sm.timeoutSec > 0 {
		timeoutSec = sm.timeoutSec
	}

	var doCtx context.Context

	if timeoutSec > 0 {
		doCtx, _ = context.WithTimeout(sm.ctx, sm.GetTimeoutDurationSecs(timeoutSec))
	} else {
		doCtx = sm.ctx
	}

	err := chromedp.Run(doCtx, sm.getActions(action...)...)

	sm.Error(err)

	if doCtx.Err() != nil {
		sm.Error(doCtx.Err())
	}

	return err
}

func (sm *SiteManager) AddErrorHandler(eh func(err error)) {
	sm.errorHandler = eh
}

func (sm SiteManager) Error(err error) {
	if err == nil {
		return
	}
	if sm.errorHandler != nil {
		sm.errorHandler(err)
		return
	}

	panic(err)
}
