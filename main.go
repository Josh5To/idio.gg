package main

import (
	"flag"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/josh5to/idio.gg/tt-fn/auth"
	"github.com/josh5to/idio.gg/tt-fn/bones"
	"github.com/josh5to/idio.gg/tt-fn/bones/header"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const VERSION = "v0.1"

type DocumentData struct {
	Head   bones.Head
	Header header.Header
}

type Page struct {
	Data         DocumentData
	PageTemplate *template.Template
}

var (
	tikTokRedirectURL string
)

func init() {
	flag.StringVar(&tikTokRedirectURL, "tik-tok-rurl", "http://localhost:8180/oauth/tiktok/validate", "Provide the \"Redirect URI\" for TikTok integration.")
}

func main() {
	flag.Parse()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	page, err := createHomepage()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create homepage")
	}

	privacyPage, err := createPrivacyPolicyPage()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create privacy page")
	}

	tosPage, err := createTermsOfServicePage()
	if err != nil {
		log.Fatal().Err(err).Msg("unable to create terms of service page")
	}

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if err := page.PageTemplate.ExecuteTemplate(w, "base", page.Data); err != nil {
			fmt.Printf("unable to execute: %v\n", err)
		}
	})

	//Setup handler function for TikTok oauth entrypoint
	http.HandleFunc("/oauth/tiktok", auth.TikTokAuthHandler(tikTokRedirectURL))

	//Setup handler function for TikTok oauth redirect callback
	http.HandleFunc("/oauth/tiktok/validate", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Err(err).Msg("error reading body of url in callback")
		}
		log.Info().Msg(string(body))
	})

	//Privacy policy
	http.HandleFunc("/privacy", func(w http.ResponseWriter, r *http.Request) {
		if err := privacyPage.PageTemplate.ExecuteTemplate(w, "base", privacyPage.Data); err != nil {
			fmt.Printf("unable to execute: %v\n", err)
		}
	})

	//Terms of Service page
	http.HandleFunc("/tos", func(w http.ResponseWriter, r *http.Request) {
		if err := tosPage.PageTemplate.ExecuteTemplate(w, "base", tosPage.Data); err != nil {
			fmt.Printf("unable to execute: %v\n", err)
		}
	})

	log.Info().Msg("Listening on :8180...")
	err = http.ListenAndServe(":8180", nil)
	if err != nil {
		log.Fatal().Err(err)
	}

}

func createHead() *bones.Head {
	return &bones.Head{
		Meta: bones.Meta{
			Charset:  "utf-8",
			Viewport: "width=device-width, initial-scale=1",
		},
		Title:      "ttfn",
		Stylesheet: []string{"../static/stylesheets/main.css"},
	}
}

func createHomepage() (*Page, error) {
	htmldoc := DocumentData{
		Head:   *createHead(),
		Header: header.Header{ClassName: "top-bar"},
	}

	//Add stylesheet
	htmldoc.Head.Stylesheet = append(htmldoc.Head.Stylesheet, "../static/home/pagestyle.css")

	docBase := filepath.Join("./templates", "base.tmpl")
	docHead := filepath.Join("./templates", "head.tmpl")
	docBody := filepath.Join("./static/home", "body.tmpl")

	ourFuncMap := template.FuncMap{
		"printVersion": func() template.HTML {
			return template.HTML(fmt.Sprintf("<p>%s</p>", VERSION))
		},
	}

	if err := bones.AddButtonFuncs(ourFuncMap); err != nil {
		return nil, err
	}
	if err := header.AddHeaderFuncs(ourFuncMap); err != nil {
		return nil, err
	}

	page, err := template.New("main-page").Funcs(ourFuncMap).ParseFiles(docBase, docHead, docBody)
	if err != nil {
		return nil, err
	}

	return &Page{
		Data:         htmldoc,
		PageTemplate: page,
	}, nil
}

func createPrivacyPolicyPage() (*Page, error) {
	htmldoc := DocumentData{
		Head: bones.Head{
			Meta: bones.Meta{
				Charset:  "utf-8",
				Viewport: "width=device-width, initial-scale=1",
			},
			Title:      "ttfn - Privacy Policy",
			Stylesheet: []string{"../static/stylesheets/main.css"},
		},
	}

	docBase := filepath.Join("./templates", "base.tmpl")
	docHead := filepath.Join("./templates", "head.tmpl")
	docBody := filepath.Join("./static/privacy", "body.tmpl")

	ourFuncMap := template.FuncMap{}

	if err := bones.AddButtonFuncs(ourFuncMap); err != nil {
		return nil, err
	}

	if err := header.AddHeaderFuncs(ourFuncMap); err != nil {
		return nil, err
	}

	page, err := template.New("privacy-policy-page").Funcs(ourFuncMap).ParseFiles(docBase, docHead, docBody)
	if err != nil {
		return nil, err
	}

	return &Page{
		Data:         htmldoc,
		PageTemplate: page,
	}, nil
}

func createTermsOfServicePage() (*Page, error) {
	htmldoc := DocumentData{
		Head: bones.Head{
			Meta: bones.Meta{
				Charset:  "utf-8",
				Viewport: "width=device-width, initial-scale=1",
			},
			Title:      "ttfn - Terms of Service",
			Stylesheet: []string{"../static/stylesheets/main.css"},
		},
	}

	docBase := filepath.Join("./templates", "base.tmpl")
	docHead := filepath.Join("./templates", "head.tmpl")
	docBody := filepath.Join("./static/tos", "body.tmpl")

	ourFuncMap := template.FuncMap{}

	if err := bones.AddButtonFuncs(ourFuncMap); err != nil {
		return nil, err
	}

	if err := header.AddHeaderFuncs(ourFuncMap); err != nil {
		return nil, err
	}

	page, err := template.New("tos-page").Funcs(ourFuncMap).ParseFiles(docBase, docHead, docBody)
	if err != nil {
		return nil, err
	}

	return &Page{
		Data:         htmldoc,
		PageTemplate: page,
	}, nil
}
