package base

import (
	"fmt"
	"strings"
)

const containsKey, equalKey, attributeKey string = "contains", "equal", "attribute"

type Element struct {
	child    *Element
	parent   *Element
	fromRoot bool
	path     string
	tag      string
	tagPos   int
	/**
	{
	"contains" :
		[
			{
				"target": "@valami",
				"value": "Valami",
				"position:1
			}
		],
		...,
	}
	*/
	filters map[string][]map[string]interface{}
}

func (e *Element) AddChild(child *Element) *Element {
	e.child = child
	e.child.parent = e

	return e
}

func (e *Element) init() {

	if e.filters == nil {
		e.filters = make(map[string][]map[string]interface{})
		e.filters[containsKey] = []map[string]interface{}{}
		e.filters[equalKey] = []map[string]interface{}{}
		e.filters[attributeKey] = []map[string]interface{}{}
	}
}

func (e *Element) ByPath(path string) *Element {
	e.path = strings.TrimSpace(path)

	return e
}

func (e *Element) ByTag(tag string, tagPos int) *Element {
	e.tag = strings.TrimSpace(tag)
	e.tagPos = tagPos

	return e
}

func (e *Element) ByAttribute(attribute string, value string, filterPos int) *Element {
	e.init()
	e.filters[attributeKey] = append(e.filters[attributeKey], map[string]interface{}{
		"attribute": attribute,
		"value":     value,
		"position":  filterPos,
	})

	return e
}

func (e *Element) ByContains(option string, value string, filterPos int) *Element {
	e.init()
	e.filters[containsKey] = append(e.filters[containsKey], map[string]interface{}{
		"option":   option,
		"value":    value,
		"position": filterPos,
	})

	return e
}

func (e *Element) ByEqual(option string, value string, filterPos int) *Element {
	e.init()
	e.filters[equalKey] = append(e.filters[equalKey], map[string]interface{}{
		"option":   option,
		"value":    value,
		"position": filterPos,
	})

	return e
}

func (e Element) String() string {
	var selector string = fmt.Sprintf("%s%s%s", e.Path(), e.Tag(), e.Filters())

	if e.child != nil {
		return selector + e.child.String()
	}
	return selector
}

func (e *Element) Path() string {
	if e.parent == nil && e.path == "" {
		e.path = "//"
		e.fromRoot = true
	} else if e.parent != nil && e.path == "" {
		e.path = "/"
	}

	return e.path
}

func (e Element) Tag() string {
	if e.tag == "" {
		e.tag = "*"
	}

	return e.tag + e.TagPos()
}

func (e Element) TagPos() string {
	if e.tagPos > 0 {
		return fmt.Sprintf("[%v]", e.tagPos)
	}

	return ""
}

func (e Element) Filters() string {
	var joinFilters []string

	for _, af := range e.filters[attributeKey] {
		var posision string
		if af["position"].(int) > 0 {
			posision = fmt.Sprintf("[%v]", af["position"].(int))
		}
		joinFilters = append(joinFilters, fmt.Sprintf(`[@%s="%s"]%s`, af["attribute"].(string), af["value"].(string), posision))
	}

	for _, cf := range e.filters[containsKey] {
		var posision string
		if cf["position"].(int) > 0 {
			posision = fmt.Sprintf("[%v]", cf["position"].(int))
		}
		joinFilters = append(joinFilters, fmt.Sprintf(`[contains(%s,"%s")]%s`, cf["option"].(string), cf["value"].(string), posision))
	}

	for _, ef := range e.filters[equalKey] {
		var posision string
		if ef["position"].(int) > 0 {
			posision = fmt.Sprintf("[%v]", ef["position"].(int))
		}
		joinFilters = append(joinFilters, fmt.Sprintf(`[%s="%s"]%s`, ef["option"].(string), ef["value"].(string), posision))
	}

	return strings.Join(joinFilters, "")
}

func HtmlTag(tag string, tagPos int) *Element {
	var e Element
	return e.ByTag(tag, tagPos)
}

func Submit(tagPos int) *Element {
	var e Element
	return e.ByTag("submit", tagPos)
}

func Button(tagPos int) *Element {
	var e Element

	return e.ByTag("button", tagPos)
}

func Input(tagPos int) *Element {
	var e Element
	return e.ByTag("input", tagPos)
}

func Textarea(tagPos int) *Element {
	var e Element
	return e.ByTag("textarea", tagPos)
}

func Paragraph(tagPos int) *Element {
	var e Element
	return e.ByTag("p", tagPos)
}

func Div(tagPos int) *Element {
	var e Element
	return e.ByTag("div", tagPos)
}

func Header(tagPos int) *Element {
	var e Element
	return e.ByTag("h1", tagPos)
}

func Header2(tagPos int) *Element {
	var e Element
	return e.ByTag("h2", tagPos)
}

func Header3(tagPos int) *Element {
	var e Element
	return e.ByTag("h3", tagPos)
}

func Header4(tagPos int) *Element {
	var e Element
	return e.ByTag("h4", tagPos)
}
