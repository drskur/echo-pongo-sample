package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func main() {
	// initialize echo
	e := echo.New()
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		CookieMaxAge: 365 * 24 * 60 * 60,
		TokenLookup:  "form:csrfmiddlewaretoken",
	}))

	// Renderer configuration
	renderglobals := map[string]interface{}{
		"url": func(name string, params ...interface{}) string {
			return e.Reverse(name, params...)
		},
	}
	renderer := NewPongoRenderer(Debug(true), SetGlobals(renderglobals))
	renderer.AddDirectory("templates/")
	renderer.UseContextProcessor(func(echoCtx echo.Context, pongoCtx pongo2.Context) {
		if csrf, ok := echoCtx.Get(middleware.DefaultCSRFConfig.ContextKey).(string); ok {
			pongoCtx["csrf_token"] = csrf
			pongoCtx["csrf_token_input"] = fmt.Sprintf("<input type=\"hidden\" name=\"csrfmiddlewaretoken\" value=\"%s\" />", csrf)
		}
	})

	e.Renderer = renderer

	h := func(c echo.Context) error {
		return c.Render(http.StatusOK, "test.html", map[string]interface{}{
			"homelink": e.Reverse("home"),
		} /* context */)
	}

	// define the route and provide a name for them
	// long form
	homeroute := e.GET("/", h)
	homeroute.Name = "home"
	// short form
	e.POST("/form/:parameter", h).Name = "form"

	log.Fatal(e.Start(":8000"))
}
