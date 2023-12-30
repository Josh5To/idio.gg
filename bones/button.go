package bones

import (
	"fmt"
	"html"
	"html/template"
	"strings"
)

type ButtonData struct {
	Content string
	Form    string
	Type    string
}

/*
Add the 'button' function, allowing your to create buttons like this:
{{button "testbutton" "tbut" "type=\"button\""}}

The various values are defined like:
{{button "button text content", "id for button", "customAttributeKey=\"value\""... (up to 20)}}
*/
func AddButtonFuncs(fm template.FuncMap) error {
	fm["button"] = createButton()

	return nil
}

/*
Returns an escaped string of valid HTML for a <button>
The data of the button
*/
func createButton() func(innerText ...string) template.HTML {
	return func(innerText ...string) template.HTML {
		argLength := len(innerText)
		switch {
		case argLength == 1:
			return template.HTML((fmt.Sprintf("<button>%s</button>", innerText[0])))
		case argLength == 2:
			return template.HTML((fmt.Sprintf("<button id=\"%s\">%s</button>",
				innerText[1], innerText[0])))
		case argLength > 2:
			//Cap at 20 custom attributes
			if argLength > 20 {
				innerText = innerText[:20]
			}
			//Setup our initial string with our button id
			buttonStr := fmt.Sprintf("<button id=\"%s\" ", innerText[1])
			for i, text := range innerText {
				if i > 1 {
					//Any strings passed after the first two are user-defined custom attributes
					buttonStr = fmt.Sprintf("%s %s ", buttonStr, text)
				}
			}
			buttonStr = fmt.Sprintf("%s>%s</button>", strings.TrimSpace(buttonStr), innerText[0])
			return template.HTML(buttonStr)
			// % q>%s</button>
		default:
			return template.HTML(html.EscapeString("<button></button>"))
		}
	}
}
