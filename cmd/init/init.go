/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package initialize

import (
	"embed"
	"html/template"
	"io/fs"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/fabianrwx/micro/internal/application/core"
	"github.com/spf13/cobra"
)

//go:embed templates/*
var templates embed.FS

var (
	name string
)

var (
	swaggerPath = "docs/swagger"
)

// initCmd represents the init command
var InitCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		projectName := strings.ToLower(name)

		// check to see if project already exist
		if _, err := os.Stat(projectName); err == nil {
			slog.Error("Failed to create project", slog.String("error", "project already exist"))
			return
		}

		// make dir if not exist
		if err := os.Mkdir(projectName, os.ModePerm); err != nil {
			slog.Error("Failed to create directory", slog.String("error", "project already exist"))
			return
		}

		// make dir if not exist
		if err := os.MkdirAll(filepath.Join(projectName, swaggerPath), os.ModePerm); err != nil {
			slog.Error("Failed to create directory", slog.String("error", "project already exist"))
			return
		}

		// cd into dir and run go mod init
		if err := os.Chdir(projectName); err != nil {
			slog.Error("Failed to change directory", slog.String("error", err.Error()))
			return
		}

		// copy files to project
		if err := copyFiles(projectName); err != nil {
			slog.Error("Failed to copy files", slog.String("error", err.Error()))
			return
		}

	},
}

func init() {
	// add flag for project name
	InitCmd.Flags().StringVarP(&name, "name", "n", "", "name of project")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func copyFiles(projectDir string) error {
	project, err := core.NewProject(name)
	if err != nil {
		slog.Error("Failed to create project", slog.String("error", err.Error()))
		return err
	}

	return fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !d.IsDir() {
			// Remove the "templates/" prefix from the path
			relativePath := strings.TrimPrefix(path, "templates/")
			destination := projectDir

			// Parse the template
			tmpl, err := template.ParseFS(templates, path)
			if err != nil {
				slog.Error("Failed to parse template", slog.String("file", path), slog.String("error", err.Error()))
				return err
			}

			// Create the destination file
			f, err := os.Create(relativePath)
			if err != nil {
				slog.Error("Failed to create file", slog.String("file", destination), slog.String("error", err.Error()))
				return err
			}

			// Execute the template with the project data
			if err := tmpl.Execute(f, project); err != nil {
				slog.Error("Failed to execute template", slog.String("file", destination), slog.String("error", err.Error()))
				f.Close() // Explicitly close file on error
				return err
			}

			// Close the file
			if err := f.Close(); err != nil {
				slog.Error("Failed to close file", slog.String("file", destination), slog.String("error", err.Error()))
				return err
			}
		}

		return nil
	})
}
