package base

import (
	"context"
	"errors"
	"fmt"
	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/device"
	"io/ioutil"
	"log"
	"os"
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

	activeGroup  string
	groupActions map[string][]chromedp.Action

	fixActions []chromedp.Action
}

func (sm *SiteManager) Init(d chromedp.Device, defTimeoutSec int64, headless bool, ignoreCertErrors bool) {
	sm.info = d
	sm.groupActions = make(map[string][]chromedp.Action)

	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.NoDefaultBrowserCheck,
		chromedp.Flag("headless", headless),
		chromedp.Flag("ignore-certificate-errors", ignoreCertErrors),
	)

	neaCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	sm.cancel = append(sm.cancel, cancel)

	// also set up a custom logger
	ctx, tCancel := chromedp.NewContext(neaCtx, chromedp.WithLogf(log.Printf))
	sm.cancel = append(sm.cancel, tCancel)

	if defTimeoutSec > 0 {
		// create a timeout
		taskCtx, ttCancel := context.WithTimeout(ctx, time.Duration(defTimeoutSec)*time.Second)
		sm.cancel = append(sm.cancel, ttCancel)
		sm.ctx = taskCtx
	} else {
		sm.ctx = ctx
	}

	sm.fixActions = append(sm.fixActions, chromedp.EmulateViewport(sm.info.Device().Width, sm.info.Device().Height))
	sm.timeoutSec = defTimeoutSec
}

func (sm *SiteManager) Group(group string) {
	sm.activeGroup = group
}

// if group is empty, it will check for the actual group and run actions from it
func (sm *SiteManager) GroupProcess(group string, timeoutSecs int64, handleError bool) error {
	if group == "" && sm.activeGroup == "" {
		return errors.New("could not run actions connected to empty group")
	}

	if group == "" {
		group = sm.activeGroup
	}

	actions, gok := sm.groupActions[group]

	if !gok || len(actions) == 0 {
		return errors.New(`your group does not exists, or there is no related action`)
	}

	err := sm.DoTimeoutContext(timeoutSecs, false, sm.groupActions[group]...)

	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) GoToPath(url string, timoutSec int64, handleError bool) error {
	action := chromedp.Navigate(url)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}
	err := sm.DoTimeoutContext(timoutSec, false, action)

	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) CaptureScreenshotInto(contentInto *[]byte, timeoutSec int64, handleError bool) error {
	action := chromedp.CaptureScreenshot(contentInto)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}
	err := sm.DoTimeoutContext(timeoutSec, false, action)

	sm.Error(err, handleError)
	if err != nil {
		return err
	}

	return err
}

func (sm *SiteManager) CreateScreenShot(filename string, timeoutSec int64, handleError bool) error {
	action := chromedp.ActionFunc(func(ctx context.Context) error {
		var p []byte
		err := sm.CaptureScreenshotInto(&p, timeoutSec, handleError)
		sm.Error(err, handleError)
		err = ioutil.WriteFile(filename, p, os.ModePerm)
		sm.Error(err, handleError)
		return err
	})
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}
	err := sm.DoTimeoutContext(timeoutSec, false, action)

	sm.Error(err, handleError)
	if err != nil {
		return err
	}

	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) Cancel() {
	for i := len(sm.cancel) - 1; i >= 0; i-- {
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

	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], actions...)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, actions...)

	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) FillField(identifier string, value string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	var actions []chromedp.Action

	actions = append(actions, chromedp.SendKeys(identifier, value, options...))

	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], actions...)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, actions...)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) ScrollTo(identifier string, timeoutSec int64, handleError bool) error {
	action := chromedp.ScrollIntoView(identifier)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitEnabled(selector string, timeoutSec int64, handleError bool) error {
	action := chromedp.WaitEnabled(selector)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitNotPresent(selector string, timeoutSec int64, handleError bool) error {
	action := chromedp.WaitNotPresent(selector)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitNotVisible(selector string, timeoutSec int64, handleError bool) error {
	action := chromedp.WaitNotVisible(selector)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitVisible(selector string, timeoutSec int64, handleError bool) error {
	action := chromedp.WaitVisible(selector)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitSelected(selector string, timeoutSec int64, handleError bool) error {
	action := chromedp.WaitSelected(selector)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) WaitReady(selector string, timeoutSec int64, handleError bool) error {
	action := chromedp.WaitReady(selector)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
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
	action := chromedp.Click(selector, options...)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) Wait(secs int64, handleError bool) error {
	action := chromedp.Sleep(time.Second * time.Duration(secs))
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(0, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) CustomAction(action chromedp.ActionFunc, timeoutSec int64, handleError bool) error {
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) FocusElement(selector string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	action := chromedp.Focus(selector, options...)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) ClearElement(selector string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	action := chromedp.Clear(selector, options...)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) DoubleClickElement(selector string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	action := chromedp.DoubleClick(selector, options...)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) InnerHTMLInto(selector string, timeoutSec int64, html *string, handleError bool) error {
	action := chromedp.InnerHTML(selector, html)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) OuterHTMLInto(selector string, timeoutSec int64, html *string, handleError bool) error {
	action := chromedp.OuterHTML(selector, html)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) TextInto(selector string, timeoutSec int64, text *string, handleError bool) error {
	action := chromedp.Text(selector, text)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) GetElementAttributeValue(selector string, attribute string, into *string, ok *bool, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	action := chromedp.AttributeValue(selector, attribute, into, ok, options...)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) GetElementAttributes(selector string, into *map[string]string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	action := chromedp.Attributes(selector, into, options...)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) GetElementsAttributes(selector string, into *[]map[string]string, timeoutSec int64, handleError bool, options ...chromedp.QueryOption) error {
	action := chromedp.AttributesAll(selector, into, options...)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) KeyDown(key string, timeoutSec int64, handleError bool) error {
	action := input.DispatchKeyEvent(input.KeyDown).WithKey(key)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) KeyRawDown(key string, timeoutSec int64, handleError bool) error {
	action := input.DispatchKeyEvent(input.KeyRawDown).WithKey(key)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) KeyUp(key string, timeoutSec int64, handleError bool) error {
	action := input.DispatchKeyEvent(input.KeyUp).WithKey(key)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
	sm.Error(err, handleError)

	return err
}

func (sm *SiteManager) KeyChar(key string, timeoutSec int64, handleError bool) error {
	action := input.DispatchKeyEvent(input.KeyChar).WithKey(key)
	if sm.activeGroup != "" {
		sm.groupActions[sm.activeGroup] = append(sm.groupActions[sm.activeGroup], action)
		return nil
	}

	err := sm.DoTimeoutContext(timeoutSec, false, action)
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
