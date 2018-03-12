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

func NewMacaron(version string) *macaron.Macaron {
	m := macaron.New()
	m.Use(macaron.Logger())
	m.Use(macaron.Recovery())
	m.Use(context.Contexter())
	m.Use(session.Sessioner())
	m.Use(csrf.Csrfer(csrf.Options{
		Secret:     randStringBytes(15),
		Header:     "X-CSRF-Token",
		Form:       "_csrf",
		Cookie:     "_csrf",
		CookiePath: "/",
		SessionKey: "csrf",
		SetHeader:  true,
		SetCookie:  false,
		Secure:     true,
		Origin:     false,
		ErrorFunc: func(w http.ResponseWriter) {
			http.Error(w, "Invalid csrf token.", http.StatusBadRequest)
		},
	}))
	m.Use(macaron.Renderer(macaron.RenderOptions{
		Funcs: []template.FuncMap{map[string]interface{}{
			"URLFor":  m.URLFor,
			"Version": func() string { return version },
		}},
	}))
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
