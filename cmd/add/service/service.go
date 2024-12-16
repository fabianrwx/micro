/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package service

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/fabianrwx/micro/internal/application/core"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

const (
	baseDir = "./"
	rootFS  = "templates"
)

//go:embed templates/*
var templates embed.FS

var (
	name string
)

// serviceCmd represents the service command
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		serviceName := strings.ToLower(name)

		// check to see if service already exist
		if _, err := os.Stat(filepath.Join(baseDir, serviceName)); err == nil {
			slog.Error("Failed to create service", slog.String("error", "service already exist"))
			return
		}

		// check to see if current dir has a mod file
		if _, err := os.Stat("Taskfile.yml"); err != nil {
			slog.Error("Failed to create service", slog.String("error", "expected to find a Taskfile.yml, are you in the project root?"))
			return
		}

		service, err := core.NewService(serviceName)
		if err != nil {
			slog.Error("Failed to create service", slog.String("error", err.Error()))
			return
		}

		// create folder structure and run all files through go template engine
		if err := copyAndParseTemplates(templates, rootFS, serviceName, baseDir, service); err != nil {
			slog.Error("Failed to copy and parse templates", slog.String("error", err.Error()))
			return
		}

		tf, err := parseTaskFile("Taskfile.yml")
		if err != nil {
			slog.Error("Failed to parse template file", slog.String("error", err.Error()))
			return
		}

		updated, err := updateTaskfile(tf, serviceName)
		if err != nil {
			slog.Error("Failed to update taskfile", slog.String("error", err.Error()))
			return
		}

		b, err := yaml.Marshal(updated)
		if err != nil {
			slog.Error("Failed to marshal taskfile", slog.String("error", err.Error()))
			return
		}

		if err := os.WriteFile("Taskfile.yml", b, 0644); err != nil {
			slog.Error("Failed to write taskfile", slog.String("error", err.Error()))
			return
		}

		if err := goModGen(); err != nil {
			slog.Error("Failed to generate go mod", slog.String("error", err.Error()))
			return
		}

		if err := protoGen(); err != nil {
			slog.Error("Failed to generate proto", slog.String("error", err.Error()))
			return
		}

		if err := goModTidy(); err != nil {
			slog.Error("Failed to run go mod tidy", slog.String("error", err.Error()))
			return
		}

	},
}

func makedirs(dirs map[string]struct{}) error {
	for dir := range dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return err
		}
	}
	return nil
}

func init() {

	ServiceCmd.Flags().StringVarP(&name, "name", "n", "", "name of service")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func goModGen() error {
	// run go generate
	cmd := exec.Command("go", "mod", "init")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go generate: %w", err)
	}

	return nil
}

func goModTidy() error {
	// run go generate
	cmd := exec.Command("go", "mod", "tidy")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go generate: %w", err)
	}

	return nil
}

func protoGen() error {
	// run go generate
	cmd := exec.Command("task", "proto")

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to run go generate: %w", err)
	}

	return nil
}

func parseTaskFile(filename string) (*core.TaskFile, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var config core.TaskFile
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return nil, fmt.Errorf("failed to decode yaml: %w", err)
	}

	return &config, nil
}

func updateTaskfile(tf *core.TaskFile, service string) (*core.TaskFile, error) {
	// Check if the "proto" task exists
	if task, exists := tf.Tasks["proto"]; exists {
		// If "proto" task exists, append the new command
		newCmd := fmt.Sprintf("protoc --proto_path=%s/proto --go_out=%s/pb --go_opt=paths=source_relative --go-grpc_out=%s/pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=%s/pb --grpc-gateway_opt=paths=source_relative --openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=api %s/proto/*.proto", service, service, service, service, service)

		// Append the new command to the existing cmds list
		task.Cmds = append(task.Cmds, newCmd)

		// Update the "proto" task in the map
		tf.Tasks["proto"] = task
	} else {
		// If "proto" task doesn't exist, create a new one with the first command
		tf.Tasks["proto"] = core.TaskCommands{
			Cmds: []string{
				fmt.Sprintf("protoc --proto_path=%s/proto --go_out=%s/pb --go_opt=paths=source_relative --go-grpc_out=%s/pb --go-grpc_opt=paths=source_relative --grpc-gateway_out=%s/pb --grpc-gateway_opt=paths=source_relative --openapiv2_out=docs/swagger --openapiv2_opt=allow_merge=true,merge_file_name=api %s/proto/*.proto", service, service, service, service, service),
			},
		}
	}

	return tf, nil
}

// copyAndParseTemplates replicates the folder structure from the embedded filesystem
// and processes files using Go's template package before writing them to the destination directory.
func copyAndParseTemplates(embedFS fs.FS, srcDir, destBaseDir, serviceName string, data interface{}) error {

	// Destination directory based on service name
	destDir := filepath.Join(destBaseDir, serviceName)

	// Walk through the embedded filesystem
	err := fs.WalkDir(embedFS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error walking path %s: %w", path, err)
		}

		// Calculate relative path to replicate the structure
		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return fmt.Errorf("error calculating relative path for %s: %w", path, err)
		}

		// Full destination path
		destPath := filepath.Join(destDir, relPath)

		// If it's a directory, create it
		if d.IsDir() {
			if err := os.MkdirAll(destPath, os.ModePerm); err != nil {
				return fmt.Errorf("error creating directory %s: %w", destPath, err)
			}
			return nil
		}

		// If it's a file, process it as a template
		return processTemplateFile(embedFS, path, destPath, data)
	})

	return err
}

// processTemplateFile reads a file from the embedded filesystem,
// processes it as a Go template, and writes the output to the destination path.
func processTemplateFile(embedFS fs.FS, srcPath, destPath string, data interface{}) error {
	// Define the custom template function
	funcMap := template.FuncMap{
		"title": capitalize,
	}
	// Open the source file from the embedded filesystem
	srcFile, err := embedFS.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open source file %s: %w", srcPath, err)
	}
	defer srcFile.Close()

	// Read the file content into a string
	content, err := io.ReadAll(srcFile)
	if err != nil {
		return fmt.Errorf("failed to read source file %s: %w", srcPath, err)
	}

	// Parse the content as a Go template
	tmpl, err := template.New(filepath.Base(srcPath)).Funcs(funcMap).Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse template %s: %w", srcPath, err)
	}

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
	}
	defer destFile.Close()

	// Execute the template with the provided data
	if err := tmpl.Execute(destFile, data); err != nil {
		return fmt.Errorf("failed to execute template for %s: %w", destPath, err)
	}

	fmt.Printf("Processed template: %s -> %s\n", srcPath, destPath)
	return nil
}

func capitalize(value core.ServiceName) string {
	// Convert the value to a string
	str := value.GetServiceName()

	// Capitalize the first letter
	if len(str) == 0 {
		return str
	}

	capitalized := strings.ToUpper(string(str[0])) + str[1:]

	slog.Info("Capitalized value", slog.String("value", capitalized))

	return capitalized
}
