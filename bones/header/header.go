package header

import (
	"bytes"
	"html/template"

	"github.com/rs/zerolog/log"
)

const headerTemplate = `{{define "header"}}
<header class="{{.ClassName}}"></header>
{{end}}`

type Header struct {
	ClassName     string
	LogoImageUrl  string
	LogoImageLink string
}

/*
 */
func AddHeaderFuncs(fm template.FuncMap) error {
	fm["header"] = createHeader()

	return nil
}

/*
Returns an escaped string of valid HTML for a <button>
The data of the button
*/
func createHeader() func(h Header) template.HTML {
	return func(h Header) template.HTML {
		temp, err := template.New("header").Parse(headerTemplate)
		if err != nil {
			log.Info().
				Err(err).
				Msgf("Error parsing header: %v\n", err)
			return template.HTML("")
		}
		log.Debug().Msgf("new header template created with name: %s", temp.Name())

		sb := new(bytes.Buffer)

		if err := temp.Execute(sb, h); err != nil {
			log.Info().
				Err(err).
				Msgf("Error executing header template: %v\n", err)
			return template.HTML("")
		}

		log.Printf("header data: %v\n", h)

		log.Printf("executed header: (string): %s\n", sb.String())

		return template.HTML(sb.String())
	}
}
