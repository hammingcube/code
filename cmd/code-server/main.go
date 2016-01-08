package main

import (
	"encoding/json"
	"errors"
	"flag"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/maddyonline/code"
	"io/ioutil"
	"net/http"
)

func decodeJsonPayload(r *http.Request, v interface{}) error {
	content, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return err
	}
	if len(content) == 0 {
		return errors.New("Empty Payload")
	}
	err = json.Unmarshal(content, v)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var pathToRunner string
	var debug bool
	flag.StringVar(&pathToRunner, "runner", ".", "path to runner binary")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()
	if debug {
		log.SetLevel(log.DEBUG)
	}

	runner := code.NewRunner(pathToRunner)
	e := echo.New()
	e.Use(mw.Logger())
	e.Get("/", func(c *echo.Context) error {
		return c.String(http.StatusOK, "All good")
	})
	e.Post("/build", func(c *echo.Context) error {
		input := &code.Input{}
		err := decodeJsonPayload(c.Request(), input)
		if err != nil {
			return err
		}
		output, err := runner.Run(input)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, output)
	})
	log.Info("Listening on 3014")
	e.Run(":3014")
}
