/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	yamlak "github.com/Slimo300/yamlak/pkg/yaml"
)

const TEMP_FILE = "/tmp/yamlak.yaml"

var forceCreateFlag bool

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "yamlak set sets a value in yaml file",
	Long: `yamlak set lets you find value of a given node in yaml file\n\n
	Example usage: 
	
	yamlak set spec.template.spec.containers[0].image nginx file.yaml -- will change value in file.yaml to "nginx"
	.`,

	Args: cobra.ExactArgs(3),

	RunE: func(cmd *cobra.Command, args []string) error {
		nodePath := args[0]
		value := args[1]
		filePath := args[2]

		originalFile, err := os.OpenFile(filePath, os.O_RDWR, 0644)
		if err != nil {
			return err
		}
		defer originalFile.Close()

		buf := &bytes.Buffer{}
		dec := yaml.NewDecoder(originalFile)
		enc := yaml.NewEncoder(buf)

		for {
			var doc interface{}
			if err := dec.Decode(&doc); errors.Is(err, io.EOF) {
				break
			} else if err != nil {
				return err
			}

			if CheckConditions(doc, conditions) {
				if forceCreateFlag {
					if err := yamlak.CreateValueByQuery(doc, nodePath, value); err != nil {
						return err
					}
				} else {
					if err := yamlak.SetValueByQuery(doc, nodePath, value); err != nil {
						return err
					}
				}
			}

			if err := enc.Encode(&doc); err != nil {
				return err
			}
		}

		if err := outputResult(originalFile, buf, inPlaceFlag); err != nil {
			return err
		}

		return nil
	},
}

func init() {
	setCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(os.Stderr, "Usage:\n  \tyamlak set <path_to_node> <value> <file> [flags]\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		cmd.Flags().PrintDefaults()

		return nil
	})
	setCmd.Flags().StringSliceVar(&conditions, "condition", []string{}, "condition for target objects to fulfill")
	setCmd.Flags().BoolVarP(&forceCreateFlag, "force", "f", false, "use to force creation of node path in file")
	setCmd.Flags().BoolVarP(&inPlaceFlag, "in-place", "i", false, "use to modify given file and not output to stdout")
}
