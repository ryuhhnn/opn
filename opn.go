// Copyright 2018 Ryan Coleman

// Permission is hereby granted, free of charge, to any person obtaining a copy of this
// software and associated documentation files (the "Software"), to deal in the Software
// without restriction, including without limitation the rights to use, copy, modify, merge,
// publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons
// to whom the Software is furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all copies
// or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
// OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/urfave/cli"
)

// Returns the user's home directory
// Note: I'm adding a trailing slash here out of convinience
// for when I use this function later on. This is more
// a personal preference than anything.
func userHomeDir() string {
	if runtime.GOOS == "Windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}

		return home + "/"
	}

	return os.Getenv("HOME") + "/"
}

// Checks to ensure the .opnrc file exists
func checkOpnExists(dir string) {
	if _, err := os.Stat(dir + ".opnrc"); os.IsNotExist(err) {
		f, err := os.Create(dir + ".opnrc")
		if err != nil {
			log.Fatal(err)
		}

		defer f.Close()
	}
}

// Opens the .opnrc file and returns a slice where each element is an alias
func openRcFile(dir string) map[string]string {
	content, err := ioutil.ReadFile(dir + ".opnrc")
	if err != nil {
		log.Fatal(err)
	}

	// Convert each line of the .opnrc file into a slice element
	aliases := strings.Split(string(content), "\n")

	// Strip all the whitespace in each slice element
	for i, val := range aliases {
		val = strings.TrimSpace(val)

		// Remove any empty values
		if val == "" {
			aliases = append(aliases[:i], aliases[1+i:]...)
		}
	}

	aliasMap := make(map[string]string)

	for _, val := range aliases {
		alias := strings.Split(val, "=")
		aliasMap[alias[0]] = alias[1]
	}

	return aliasMap
}

func addNewAlias(dir string, alias string, path string) error {
	f, err := os.OpenFile(dir+".opnrc", os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}

	defer f.Close()

	if _, err = f.WriteString(alias + "=" + path + "\n"); err != nil {
		return err
	}

	return nil
}

func main() {
	dir := userHomeDir()

	checkOpnExists(dir)

	app := cli.NewApp()
	app.Name = "opn"
	app.Usage = "quickly open any app using whatever name you see fit."
	app.Version = "0.1.0"
	app.Action = func(c *cli.Context) error {
		alias := c.Args().First()

		if alias == "" {
			fmt.Println("No app was specified. Need help? Run: opn help to see a list of available commands")
			return nil
		}

		aliases := openRcFile(dir)

		cmd := exec.Command("open", aliases[alias])
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}

		return nil
	}
	app.Commands = []cli.Command{
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "add a new app alias to opn",
			Action: func(c *cli.Context) error {
				alias := c.Args().First()
				path := c.Args().Get(1)

				if alias == "" || path == "" {
					fmt.Println("Please ensure you specify both an alias and a path.")
					return nil
				}

				cmd := addNewAlias(dir, alias, path)
				if cmd != nil {
					log.Fatal(cmd)
				}

				return nil
			},
		},
		{
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "lists all the aliases that are saved",
			Action: func(c *cli.Context) error {
				aliases := openRcFile(dir)

				fmt.Println("Saved aliases:")
				for name := range aliases {
					fmt.Println("\t" + name)
				}

				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
