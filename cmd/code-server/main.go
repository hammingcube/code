package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"github.com/labstack/gommon/log"
	"github.com/maddyonline/code"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
)

func prefixed(s string) string {
	return fmt.Sprintf("CODE_PROJECT_%s", s)
}

var (
	PORT   = flag.String("port", "", fmt.Sprintf("port to start server on (or use, %s env variable)", prefixed("PORT")))
	ROOT   = flag.String("root", "", fmt.Sprintf("root directory to serve content from (or use, %s env variable)", prefixed("ROOT")))
	RUNNER = flag.String("runner", "", fmt.Sprintf("path to runner binary (or use, %s env variable)", prefixed("RUNNER")))
)

func initialize() {
	flag.Parse()
	assignString(PORT, *PORT, os.Getenv(prefixed("PORT")), "3014")
	assignString(ROOT, *ROOT, os.Getenv(prefixed("ROOT")), defaultRootDir())
	assignString(RUNNER, *RUNNER, os.Getenv(prefixed("RUNNER")), "./runner")
	log.Info("Using PORT=%s, ROOT=%s, RUNNER=%s", *PORT, *ROOT, *RUNNER)
}

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
	flag.Parse()
	initialize()

	runner := code.NewRunner(*RUNNER)
	e := echo.New()
	e.Use(mw.Logger())
	e.Static("/", filepath.Join(*ROOT, "static"))
	e.Index(filepath.Join(*ROOT, "static/index.html"))
	e.Post("/run", func(c *echo.Context) error {
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
	e.Run(fmt.Sprintf(":%s", *PORT))
}

func assignString(v *string, args ...string) {
	for _, arg := range args {
		if arg != "" {
			*v = arg
			break
		}
	}
}

func defaultRootDir() string {
	// Configure opts.StaticFilesRoot
	defaultDir := "."
	if GOPATH := os.Getenv("GOPATH"); GOPATH != "" {
		srcDir, err := filepath.Abs(filepath.Join(GOPATH, "src/github.com/maddyonline/code/cmd/code-server"))
		if err == nil {
			defaultDir = srcDir
		}
	}
	return defaultDir
}
