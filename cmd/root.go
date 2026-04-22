package cmd

import (
	"fmt"
	"os"

	"github.com/adetola/convert/converter"
	"github.com/spf13/cobra"
)

var (
	listFlag  bool
	forceFlag bool
)

var rootCmd = &cobra.Command{
	Use:   "cv <input> <output>",
	Short: "Convert files between formats",
	Long:  "A fast CLI tool that converts files between image, data, and document formats.",
	Args: func(cmd *cobra.Command, args []string) error {
		if listFlag {
			return nil
		}
		if len(args) != 2 {
			return fmt.Errorf("requires exactly 2 arguments: input and output file paths")
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if listFlag {
			fmt.Println("Supported conversions:")
			for _, pair := range converter.ListSupported() {
				fmt.Printf("  %s\n", pair)
			}
			return nil
		}

		input := args[0]
		output := args[1]

		if _, err := os.Stat(input); os.IsNotExist(err) {
			return fmt.Errorf("input file not found: %s", input)
		}

		if !forceFlag {
			if _, err := os.Stat(output); err == nil {
				return fmt.Errorf("output file already exists: %s. Use --force to overwrite", output)
			}
		}

		if err := converter.Convert(input, output); err != nil {
			return err
		}

		fmt.Printf("Converted %s → %s\n", input, output)
		return nil
	},
}

func init() {
	rootCmd.Flags().BoolVar(&listFlag, "list", false, "List all supported format conversions")
	rootCmd.Flags().BoolVar(&forceFlag, "force", false, "Overwrite output file if it already exists")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
