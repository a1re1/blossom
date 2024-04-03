package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/huh/spinner"
	"github.com/charmbracelet/lipgloss"
)

func main() {
	accessible, _ := strconv.ParseBool(os.Getenv("ACCESSIBLE"))

	var language string

	huh.NewSelect[string]().
		Title("Choose language").
		Options(
			huh.NewOption("Java", "java"),
			huh.NewOption("JavaScript", "javascript"),
			huh.NewOption("TypeScript", "typescript"),
		).
		Value(&language).
		Validate(func(str string) error {
			if str == "javascript" || str == "typescript" {
				return errors.New("There are currently no templates for " + str + ", please check back in future updates. :)")
			}
			return nil
		}).
		WithAccessible(accessible).
		Run()

	var modules []string

	if language == "java" {
		huh.NewMultiSelect[string]().
			Title("Choose java modules").
			Options(
				huh.NewOption("Pom Root", "pom-root"),
				huh.NewOption("SharedImmutables", "SharedImmutables"),
				huh.NewOption("Sqlite", "Sqlite"),
			).
			Value(&modules).
			Validate(func(args []string) error {
				if len(args) == 0 {
					return errors.New("Must choose at least 1 module.")
				}
				return nil
			}).
			WithAccessible(accessible).
			Run()
	}

	loadTest := func() {
		time.Sleep(2 * time.Second)
	}

	_ = spinner.New().Title("Generating modules...").Accessible(accessible).Action(loadTest).Run()

	// Print order summary.
	{
		var sb strings.Builder

		fmt.Fprintf(&sb, "Created new %s modules: %s", language, modules)

		fmt.Println(
			lipgloss.NewStyle().
				Width(40).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2).
				Render(sb.String()),
		)
	}
}
