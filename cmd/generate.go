/*
Copyright © 2023 Consensys
*/
package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

const (
	bindingFile = "gnark_solidity.go"
)

var (
	fSolFile string // solidity file to generate from
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "generates wraps solc and abigen to generate the solidity contract and go bindings of a given .sol contract",
	Run:   runGenerate,
}

func runGenerate(cmd *cobra.Command, args []string) {
	if fSolFile == "" {
		fmt.Println("please specify --solidity")
		os.Exit(1)
	}
	fmt.Println("generate called")
	fSolFile = filepath.Join(fBaseDir, fSolFile)
	if err := fileExists(fSolFile); err != nil {
		fmt.Println(fSolFile + " does not exist: " + err.Error())
		os.Exit(1)
	}

	// call solc
	cmdSolc := exec.Command("solc", "--via-ir", "--evm-version", "paris", "--combined-json", "abi,bin", fSolFile, "-o", fBaseDir, "--overwrite")
	fmt.Println("running solc: " + cmdSolc.String())
	if out, err := cmdSolc.CombinedOutput(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		fmt.Printf("output: %s\n", string(out))

		os.Exit(1)
	}

	// call abigen
	cmdAbigen := exec.Command("abigen", "--combined-json", filepath.Join(fBaseDir, "combined.json"), "--pkg", "main", "--out", filepath.Join(fBaseDir, bindingFile))
	fmt.Println("running abigen: " + cmdAbigen.String())
	if out, err := cmdAbigen.CombinedOutput(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		fmt.Printf("output: %s\n", string(out))
		os.Exit(1)
	}

}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().StringVar(&fSolFile, "solidity", "", "path to the solidity file to generate from")
	generateCmd.MarkFlagRequired("solidity")
}

func fileExists(path string) error {
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("file %s does not exist", path)
	}
	return nil
}
