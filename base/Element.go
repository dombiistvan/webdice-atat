package base

import (
	"fmt"
	"strings"
)

const containsKey, equalKey, attributeKey string = "contains", "equal", "attribute"

type Element struct {
	child  *Element
	parent *Element
	path   string
	tag    string
	tagPos int
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
	},
	{
	"equal" :
		[
			{
				"option": "text() || @attr || . || anything gaven",
				"value": "Valami",
				"position:1
			}
		],
		...,
	},
	"attribute" :
		[
			{
				"attribute": "id",
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
	var selector string

	selector = fmt.Sprintf("%s%s%s", e.Path(), e.Tag(), e.Filters())

	if e.child != nil {
		selector += "/" + e.child.String()
	}

	return selector
}

func (e *Element) Path() string {
	if e.parent == nil && e.path == "" {
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

func Radio(tagPos int, attrPos int) *Element {
	var e Element
	return e.ByTag("input", tagPos).ByAttribute("type", "radio", attrPos)
}

func Checkbox(tagPos int, attrPos int) *Element {
	var e Element
	return e.ByTag("input", tagPos).ByAttribute("type", "checkbox", attrPos)
}

func InpButton(tagPos int, attrPos int) *Element {
	var e Element
	return e.ByTag("input", tagPos).ByAttribute("type", "button", attrPos)
}

func InpSubmit(tagPos int, attrPos int) *Element {
	var e Element
	return e.ByTag("input", tagPos).ByAttribute("type", "submit", attrPos)
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

func Header5(tagPos int) *Element {
	var e Element
	return e.ByTag("h5", tagPos)
}

func Header6(tagPos int) *Element {
	var e Element
	return e.ByTag("h6", tagPos)
}

func Anchor(tagPos int) *Element {
	var e Element
	return e.ByTag("a", tagPos)
}

func Body(tagPos int) *Element {
	var e Element
	return e.ByTag("body", tagPos)
}

func Canvas(tagPos int) *Element {
	var e Element
	return e.ByTag("canvas", tagPos)
}

func Caption(tagPos int) *Element {
	var e Element
	return e.ByTag("caption", tagPos)
}

func Center(tagPos int) *Element {
	var e Element
	return e.ByTag("center", tagPos)
}

func Code(tagPos int) *Element {
	var e Element
	return e.ByTag("code", tagPos)
}

func Colgroup(tagPos int) *Element {
	var e Element
	return e.ByTag("colgroup", tagPos)
}

func Data(tagPos int) *Element {
	var e Element
	return e.ByTag("data", tagPos)
}

func DataList(tagPos int) *Element {
	var e Element
	return e.ByTag("datalist", tagPos)
}

func Dd(tagPos int) *Element {
	var e Element
	return e.ByTag("dd", tagPos)
}

func Dt(tagPos int) *Element {
	var e Element
	return e.ByTag("dt", tagPos)
}

func Dl(tagPos int) *Element {
	var e Element
	return e.ByTag("dl", tagPos)
}

func Details(tagPos int) *Element {
	var e Element
	return e.ByTag("details", tagPos)
}

func Dialog(tagPos int) *Element {
	var e Element
	return e.ByTag("Dialog", tagPos)
}

func Dir(tagPos int) *Element {
	var e Element
	return e.ByTag("dir", tagPos)
}

func Em(tagPos int) *Element {
	var e Element
	return e.ByTag("em", tagPos)
}

func Embed(tagPos int) *Element {
	var e Element
	return e.ByTag("embed", tagPos)
}

func Fieldset(tagPos int) *Element {
	var e Element
	return e.ByTag("fieldset", tagPos)
}

func Figure(tagPos int) *Element {
	var e Element
	return e.ByTag("figure", tagPos)
}

func Footer(tagPos int) *Element {
	var e Element
	return e.ByTag("footer", tagPos)
}

func Form(tagPos int) *Element {
	var e Element
	return e.ByTag("form", tagPos)
}

func Frame(tagPos int) *Element {
	var e Element
	return e.ByTag("frame", tagPos)
}

func Frameset(tagPos int) *Element {
	var e Element
	return e.ByTag("frameset", tagPos)
}

func Head(tagPos int) *Element {
	var e Element
	return e.ByTag("head", tagPos)
}

func Hr(tagPos int) *Element {
	var e Element
	return e.ByTag("hr", tagPos)
}

func Br(tagPos int) *Element {
	var e Element
	return e.ByTag("br", tagPos)
}

func Html(tagPos int) *Element {
	var e Element
	return e.ByTag("html", tagPos)
}

func Link(tagPos int) *Element {
	var e Element
	return e.ByTag("link", tagPos)
}

func Main(tagPos int) *Element {
	var e Element
	return e.ByTag("main", tagPos)
}

func Map(tagPos int) *Element {
	var e Element
	return e.ByTag("map", tagPos)
}

func Nav(tagPos int) *Element {
	var e Element
	return e.ByTag("nav", tagPos)
}

func Object(tagPos int) *Element {
	var e Element
	return e.ByTag("object", tagPos)
}

func Optgroup(tagPos int) *Element {
	var e Element
	return e.ByTag("optgroup", tagPos)
}

func Option(tagPos int) *Element {
	var e Element
	return e.ByTag("option", tagPos)
}

func Picture(tagPos int) *Element {
	var e Element
	return e.ByTag("picture", tagPos)
}

func Pre(tagPos int) *Element {
	var e Element
	return e.ByTag("pre", tagPos)
}

func Quote(tagPos int) *Element {
	var e Element
	return e.ByTag("q", tagPos)
}

func Script(tagPos int) *Element {
	var e Element
	return e.ByTag("script", tagPos)
}

func Section(tagPos int) *Element {
	var e Element
	return e.ByTag("section", tagPos)
}

func Select(tagPos int) *Element {
	var e Element
	return e.ByTag("select", tagPos)
}

func Small(tagPos int) *Element {
	var e Element
	return e.ByTag("small", tagPos)
}

func Source(tagPos int) *Element {
	var e Element
	return e.ByTag("source", tagPos)
}

func Span(tagPos int) *Element {
	var e Element
	return e.ByTag("span", tagPos)
}

func Strong(tagPos int) *Element {
	var e Element
	return e.ByTag("strong", tagPos)
}

func Style(tagPos int) *Element {
	var e Element
	return e.ByTag("stly", tagPos)
}

func Svg(tagPos int) *Element {
	var e Element
	return e.ByTag("svg", tagPos)
}

func Table(tagPos int) *Element {
	var e Element
	return e.ByTag("table", tagPos)
}

func Thead(tagPos int) *Element {
	var e Element
	return e.ByTag("thead", tagPos)
}

func Tbody(tagPos int) *Element {
	var e Element
	return e.ByTag("tbody", tagPos)
}

func Tfoot(tagPos int) *Element {
	var e Element
	return e.ByTag("tfoot", tagPos)
}

func Tr(tagPos int) *Element {
	var e Element
	return e.ByTag("tr", tagPos)
}

func Td(tagPos int) *Element {
	var e Element
	return e.ByTag("td", tagPos)
}

func Th(tagPos int) *Element {
	var e Element
	return e.ByTag("th", tagPos)
}

func Title(tagPos int) *Element {
	var e Element
	return e.ByTag("title", tagPos)
}

func UnorderedList(tagPos int) *Element {
	return Ul(tagPos)
}

func Ul(tagPos int) *Element {
	var e Element
	return e.ByTag("ul", tagPos)
}

func OrderedList(tagPos int) *Element {
	return Ol(tagPos)
}

func Ol(tagPos int) *Element {
	var e Element
	return e.ByTag("ol", tagPos)
}

func ListItem(tagPos int) *Element {
	return Li(tagPos)
}

func Li(tagPos int) *Element {
	var e Element
	return e.ByTag("li", tagPos)
}
