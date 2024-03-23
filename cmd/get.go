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

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "yamlak get lets you find a value in yaml file",

	Long: `yamlak get lets you find value of a given node in yaml file\n\n
	Example usage: 
	
	yamlak get spec.template.spec.containers[0].image file.yaml -- will return value of image in Kubernetes deployment configuration
	`,

	Args: cobra.ExactArgs(2),

	RunE: func(cmd *cobra.Command, args []string) error {
		nodePath := args[0]
		filePath := args[1]

		f, err := os.Open(filePath)
		if err != nil {
			return err
		}
		defer f.Close()

		dec := yaml.NewDecoder(f)

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

			val, err := yamlak.GetValueByQuery(doc, nodePath)
			if err != nil {
				return err
			}

			fmt.Println(val)
		}

		return nil
	},
}

func init() {
	getCmd.SetUsageFunc(func(cmd *cobra.Command) error {
		fmt.Fprintf(os.Stderr, "Usage:\n  \tyamlak get <path_to_node> <file> [flags]\n")
		fmt.Fprintf(os.Stderr, "Flags:\n")
		cmd.Flags().PrintDefaults()

		return nil
	})
	getCmd.Flags().StringSliceVar(&conditions, "condition", []string{}, "condition for target objects to fulfill")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
