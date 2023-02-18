package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/spf13/cobra"
)

var (
	URL, ReqMethod, JSON string
	ErrInvalidArgsLength = errors.New("invalid args length. must be 2 args [URL] [METHOD] [JSON IF POST]")
	ErrInvalidURL        = errors.New("invalid URL was gave in first arg")
	ErrInvalidReqMethod  = errors.New("invalid request method. only GET, POST are processed")
	ErrInvalidJsonData   = errors.New("invalid json data. can't do POST request")
)

var (
	avalabelMethods = []string{"GET", "POST"}
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&ReqMethod, "method", "m", "GET", "method which request will use")
	rootCmd.MarkFlagRequired("method")
}

var rootCmd = cobra.Command{
	Use:   "gohttp URL [-m GET | -m POST file.json]",
	Short: "gohttp is a simple http cmd provider",
	Long:  "gohttp is a simple http cmd provider. Avalabel methods is GET, POST.",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return ErrInvalidArgsLength
		} else {
			if is.URL.Validate(args[0]) != nil {
				return ErrInvalidURL
			}
			URL = args[0]

			for _, m := range avalabelMethods {
				if ReqMethod == m {
					if ReqMethod == "POST" {
						if len(args) == 2 {
							j, err := jsonHandler(args[1])
							if err != nil {
								return ErrInvalidJsonData
							}
							JSON = j
						}
						return ErrInvalidArgsLength
					}
					return nil
				}
			}
		}
		return ErrInvalidReqMethod
	},
	DisableSuggestions: true,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.ParseFlags(args)
		if err != nil {
			cmd.PrintErrln(err)
		}

		client := http.Client{
			Timeout: 30 * time.Second,
		}

		var reader io.Reader
		if ReqMethod == "POST" {
			reader = strings.NewReader(JSON)
		}

		req, err := http.NewRequest(ReqMethod, URL, reader)
		if err != nil {
			cmd.PrintErrln(err)
		}

		res, err := client.Do(req)
		if err != nil {
			cmd.Println(err)
		}

		buf := make([]byte, 0, 16)
		b := bytes.NewBuffer(make([]byte, 16))
		for {
			n, err := res.Body.Read(b.Bytes())
			buf = append(buf, b.Bytes()[:n]...)
			if err != nil {
				if err == io.EOF {
					break
				}
				cmd.PrintErrln(err)
				break
			}
		}

		metaData()
		b = bytes.NewBuffer(buf)
		_, err = b.WriteTo(os.Stdout)
		if err != nil {
			cmd.PrintErrln(err)
		}
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func jsonHandler(path string) (string, error) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buf := make([]byte, 0, 16)
	bytes := bytes.NewBuffer(make([]byte, 16))

	for {
		n, err := file.Read(bytes.Bytes())
		buf = append(buf, bytes.Bytes()[:n]...)
		if err != nil {
			if err == io.EOF {
				break
			}
			return "", err
		}
	}

	if !json.Valid(buf) {
		return "", err
	}

	return string(buf), nil
}

func line(s rune, n int) {
	for i := 1; i <= n; i++ {
		fmt.Print(string(s))
	}
	fmt.Println()
}

func metaData() {
	line('-', 32)
	t := time.Now()
	fmt.Println("TIME:", fmt.Sprintf("%d:%d:%d %d:%d:%d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second()))
	fmt.Println("URL:", URL)
	fmt.Println("METHOD:", ReqMethod)
	if JSON == "" {
		fmt.Println("JSON: -")
	} else {
		fmt.Println("JSON:", JSON)
	}
	line('-', 32)
}
