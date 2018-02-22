package router

import (
	"html/template"
	"math/rand"
	"net/http"
	"time"

	"github.com/Yara-Rules/yara-endpoint/server/context"

	"github.com/go-macaron/csrf"
	"github.com/go-macaron/session"
	macaron "gopkg.in/macaron.v1"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func NewMacaron() *macaron.Macaron {
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(context.Contexter())
	m.Use(session.Sessioner())
	// m.Use(macaron.Context())
	m.Use(csrf.Csrfer(csrf.Options{
		// The global secret value used to generate Tokens. Default is a random string.
		Secret: randStringBytes(15),
		// HTTP header used to set and get token. Default is "X-CSRFToken".
		Header: "X-CSRF-Token",
		// Form value used to set and get token. Default is "_csrf".
		Form: "_csrf",
		// Cookie value used to set and get token. Default is "_csrf".
		Cookie: "_csrf",
		// Cookie path. Default is "/".
		CookiePath: "/",
		// Key used for getting the unique ID per user. Default is "uid".
		SessionKey: "csrf",
		// If true, send token via header. Default is false.
		SetHeader: true,
		// If true, send token via cookie. Default is false.
		SetCookie: false,
		// Set the Secure flag to true on the cookie. Default is false.
		Secure: true,
		// Disallow Origin appear in request header. Default is false.
		Origin: false,
		// The function called when Validate fails. Default is a simple error print.
		ErrorFunc: func(w http.ResponseWriter) {
			http.Error(w, "Invalid csrf token.", http.StatusBadRequest)
		},
	}))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs: []template.FuncMap{map[string]interface{}{
			"URLFor": m.URLFor,
		}},
	}))
	// m.Use(tplextender.CompoRender(tplextender.RenderOptions{
	//     Funcs: map[string]interface{}{
	//         "URLFor": m.URLFor,
	//     },
	// }))
	m.Use(macaron.Static("public", macaron.StaticOptions{
		Prefix:      "public",
		SkipLogging: false,
		Expires: func() string {
			return time.Now().Add(24 * 60 * time.Minute).UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
		},
	}))

	return m
}

func randStringBytes(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
