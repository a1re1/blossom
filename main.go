package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

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

	runShellScript := func(relPath string, args ...string) {
		ex, exErr := os.Executable()
		if exErr != nil {
			panic(exErr)
		}
		exPath := filepath.Dir(ex)

		if strings.Contains(exPath, "go-build") {
			_, filename, _, _ := runtime.Caller(0)
			exPath = path.Dir(filename)
		}

		var cmdArgs = append([]string{exPath + relPath}, args...)
		fmt.Println("args: ", cmdArgs)
		_, err := exec.Command("/bin/sh", cmdArgs...).Output()
		if err != nil {
			fmt.Printf("error %s", err)
		}
	}

	runWithLoading := func(relPath string, args ...string) {
		runnable := func() {
			runShellScript(relPath, args...)
		}
		_ = spinner.New().Title("Generating modules...").Accessible(accessible).Action(runnable).Run()
	}

	if language == "java" {
		var modules []string
		var newProject bool
		huh.NewConfirm().
			Title("Action:").
			Affirmative("New project").
			Negative("Existing project").
			Value(&newProject).
			Run()

		var groupId string
		var artifactId string

		if newProject {
			huh.NewForm(
				huh.NewGroup(
					huh.NewInput().Title("groupId").Placeholder("com.tylerwhitehurst").Value(&groupId),
					huh.NewInput().Title("artifactId").Value(&artifactId),
				),
			).
				Run()
			if len(groupId) == 0 {
				groupId = "com.tylerwhitehurst"
			}
			runWithLoading("/templates/java/pom-root.sh", groupId, artifactId)
		}

		huh.NewMultiSelect[string]().
			Title("Choose java modules").
			Options(
				huh.NewOption("Models", "models"),
				huh.NewOption("Sql", "sql"),
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

		for _, module := range modules {
			switch module {
			case "shared-immutables":
				var groupId string
				var artifactId string
				huh.NewInput().
					Title("groupId").
					Value(&groupId).
					Placeholder("com.tylerwhitehurst").
					Run()
				if len(groupId) == 0 {
					groupId = "com.tylerwhitehurst"
				}
				huh.NewInput().
					Title("artifactId").
					Value(&artifactId).
					Run()
				runWithLoading("/templates/java/pom-root.sh", groupId, artifactId)
				continue
			case "sql":
				continue
			}
		}
	}

	{
		var str = "Success 🥳"
		fmt.Println("")
		fmt.Println(
			lipgloss.NewStyle().
				Width(32).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("63")).
				Padding(1, 2).
				Render(str),
		)
	}
}
