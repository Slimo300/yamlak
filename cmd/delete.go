/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	yamlak "github.com/Slimo300/yamlak/pkg/yaml"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "yamlak delete deletes a node in yaml file",
	Long: `yamlak delete lets you find a node with a given path in yaml file and delete it from file\n\n
	Example usage: 
	
	yamlak delete spec.template.spec.containers[0].image file.yaml -- will delete "image" in file.yaml
	.`,

	Args: cobra.ExactArgs(2),

	Aliases: []string{"del"},

	RunE: func(cmd *cobra.Command, args []string) error {
		nodePath := args[0]
		filePath := args[1]

		originalFile, err := os.OpenFile(filePath, os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer originalFile.Close()

		tmpFile, err := os.OpenFile(TEMP_FILE, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
		if err != nil {
			return err
		}
		defer tmpFile.Close()

		dec := yaml.NewDecoder(originalFile)
		enc := yaml.NewEncoder(tmpFile)

		for {
			var doc interface{}
			if err := dec.Decode(&doc); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				return err
			}

			if !CheckConditions(doc, conditions) {
				continue
			}

			if err := yamlak.DeleteValueByQuery(doc, nodePath); err != nil && !errors.Is(err, yamlak.ErrValueNotFound) {
				return err
			}

			if err := enc.Encode(&doc); err != nil {
				return err
			}
		}

		if err := originalFile.Close(); err != nil {
			return err
		}
		if err := tmpFile.Close(); err != nil {
			return err
		}

		if err := os.Rename(TEMP_FILE, filePath); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	deleteCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(os.Stderr, "Usage:\n  \tyamlak delete <path_to_node> <file> [flags]\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		cmd.Flags().PrintDefaults()

		return nil
	})
	deleteCmd.Flags().StringSliceVar(&conditions, "condition", []string{}, "condition for target objects to fulfill")
}
