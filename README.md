## Sample of using django-like templating engine (pongo2) with echo framework (golang)

This sample uses [pongo2 templating engine](https://github.com/flosch/pongo2) with [Labstack Echo](https://github.com/labstack/echo) web framework.

It demonstrates the use of named routes in the templates (see `test.html` and specifically the `action` attribute) and reversing the routes in the program itself (see the only handler method in the `main.go`)
