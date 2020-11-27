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
	cancel       []context.CancelFunc
	info         chromedp.Device
	errorHandler func(err error)
	timeoutSec   int64

	fixActions []chromedp.Action
}

func (sm *SiteManager) Init(d chromedp.Device, defTimeoutSec int64, headless bool, ignoreCertErrors bool) {
	sm.info = d

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", headless),
		chromedp.Flag("ignore-certificate-errors", ignoreCertErrors),
	)

	neaCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	sm.cancel = append(sm.cancel, cancel)

	// create a timeout
	taskCtx, ttCancel := context.WithTimeout(neaCtx, time.Duration(defTimeoutSec)*time.Second)
	sm.cancel = append(sm.cancel, ttCancel)

	sm.ctx = taskCtx

	sm.fixActions = append(sm.fixActions, chromedp.EmulateViewport(sm.info.Device().Width, sm.info.Device().Height))
	sm.timeoutSec = defTimeoutSec
}

func (sm *SiteManager) GoToPath(url string, timoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timoutSec, false, chromedp.Navigate(url))

	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) CreateScreenShot(filename string, timeoutSec int64, handleError bool) error {
	var p []byte

	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.CaptureScreenshot(&p))

	sm.Error(err, handleError)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, p, 0755)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) Cancel() {
	for i := len(sm.cancel) - 1; i >= 0; i++ {
		sm.cancel[i]()
	}
}

func (sm SiteManager) ByID(path, tag, id string) string {
	return sm.ByAttribute(path, tag, "id", id)
}

func (sm SiteManager) ByAttribute(path, tag, attributeKey, attributeValue string) string {
	return fmt.Sprintf(`%s%s[%s="%s"]`, path, tag, attributeKey, attributeValue)
}

func (sm *SiteManager) FillFields(fields []map[string]interface{}, timeoutSec int64, handleError bool) error {
	var actions []chromedp.Action

	for _, fd := range fields {
		identifier, kok := fd["identifier"]
		value, vok := fd["value"]
		options, qok := fd["options"]
		if kok && vok && !qok {
			actions = append(actions, chromedp.SendKeys(identifier, value.(string)))
		} else if kok && vok {
			if reflect.TypeOf(options) != reflect.TypeOf([]chromedp.QueryOption{}) {
				sm.Error(errors.New("options must be instance of []chromedp.QueryOption"), handleError)
			}
			actions = append(actions, chromedp.SendKeys(identifier, value.(string), options.([]chromedp.QueryOption)...))
		}
	}

	err := sm.DoTimeoutContext(timeoutSec, false, actions...)

	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) FillField(identifier string, value string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	var actions []chromedp.Action

	actions = append(actions, chromedp.SendKeys(identifier, value, options...))

	err := sm.DoTimeoutContext(timeoutSec, false, actions...)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) ScrollTo(identifier string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.ScrollIntoView(identifier))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitEnabled(selector string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.WaitEnabled(selector))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitNotPresent(selector string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.WaitNotPresent(selector))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitNotVisible(selector string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.WaitNotVisible(selector))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitVisible(selector string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.WaitVisible(selector))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitSelected(selector string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.WaitSelected(selector))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitReady(selector string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.WaitReady(selector))
	sm.Error(err, handleError)

	return err
}

/*func (sm *SiteManager) PressButton(selector string, timeoutSec int64 ,options... chromedp.QueryOption){
	err := sm.DoTimeoutContext(timeoutSec,chromedp.Click(selector, options...))
	if err != nil{
		panic(err)
	}
}*/

func (sm *SiteManager) ClickElement(selector string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.Click(selector, options...))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) Wait(secs int64, handleError bool) error {
	err := sm.DoTimeoutContext(0, false, chromedp.Sleep(time.Second*time.Duration(secs)))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) CustomAction(action chromedp.ActionFunc, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) FocusElement(selector string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.Focus(selector, options...))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) ClearElement(selector string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.Clear(selector, options...))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) DoubleClickElement(selector string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.DoubleClick(selector, options...))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) InnerHTMLInto(selector string, timeoutSec int64, html *string, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.InnerHTML(selector, html))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) OuterHTMLInto(selector string, timeoutSec int64, html *string, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.OuterHTML(selector, html))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) TextInto(selector string, timeoutSec int64, text *string, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.Text(selector, text))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) GetElementAttributeValue(selector string, attribute string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) (string, bool, error) {
	var attributeValue string
	var ok bool
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.AttributeValue(selector, attribute, &attributeValue, &ok, options...))
	sm.Error(err, handleError)

	return attributeValue, ok, err
}

func (sm *SiteManager) GetElementAttributes(selector string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) (map[string]string, error) {
	var attributes map[string]string
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.Attributes(selector, &attributes, options...))
	sm.Error(err, handleError)

	return attributes, err
}

func (sm *SiteManager) GetElementsAttributes(selector string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) ([]map[string]string, error) {
	var attributes []map[string]string
	err := sm.DoTimeoutContext(timeoutSec, false, chromedp.AttributesAll(selector, &attributes, options...))
	sm.Error(err, handleError)

	return attributes, err
}

func (sm *SiteManager) KeyDown(key string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, input.DispatchKeyEvent(input.KeyDown).WithKey(key))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) KeyRawDown(key string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, input.DispatchKeyEvent(input.KeyRawDown).WithKey(key))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) KeyUp(key string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, input.DispatchKeyEvent(input.KeyUp).WithKey(key))
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) KeyChar(key string, timeoutSec int64, handleError bool) error {
	err := sm.DoTimeoutContext(timeoutSec, false, input.DispatchKeyEvent(input.KeyChar).WithKey(key))
	sm.Error(err, handleError)

	return err
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

func (sm *SiteManager) DoTimeoutContext(timeoutSec int64, handleError bool, action ...chromedp.Action) error {
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

	sm.Error(err, handleError)

	if doCtx.Err() != nil {
		sm.Error(doCtx.Err(), handleError)
	}

	return err
}

func (sm *SiteManager) AddErrorHandler(eh func(err error)) {
	sm.errorHandler = eh
}

func (sm SiteManager) Error(err error, handleError bool) {
	if !handleError || err == nil {
		return
	}

	if sm.errorHandler != nil {
		sm.errorHandler(err)
		return
	}

	panic(err)
}
