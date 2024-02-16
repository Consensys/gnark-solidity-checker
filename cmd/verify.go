/*
Copyright © 2023 Consensys
*/
package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"

	"github.com/spf13/cobra"
)

var (
	fProof          string // hex encoded proof to verify
	fPublicInputs   string // hex encoded public inputs to verify
	fNbPublicInputs int    // number of public inputs
	fGroth16        bool
	fPlonK          bool
	fNbCommitments  int // number of commitments
)

// verifyCmd represents the verify command
var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "runs a simulated geth backend; deploys the contract and runs the verification of the provided proof and public inputs",
	Run:   runVerify,
}

func runVerify(cmd *cobra.Command, args []string) {
	fmt.Println("verify called")
	fBindings := filepath.Join(fBaseDir, bindingFile)
	if err := fileExists(fBindings); err != nil {
		fmt.Println(fBindings + " does not exist: " + err.Error())
		os.Exit(1)
	}

	if fGroth16 {
		var template string
		if fNbCommitments > 0 {
			template = tmplGroth16Commitment
		} else {
			template = tmplGroth16
		}
		if err := generateMain(template, filepath.Join(fBaseDir, "main.go"), fProof, fPublicInputs, fNbPublicInputs, fNbCommitments); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else if fPlonK {
		if err := generateMain(tmplPlonK, filepath.Join(fBaseDir, "main.go"), fProof, fPublicInputs, fNbPublicInputs, fNbCommitments); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Println("please specify either --groth16 or --plonk")
		os.Exit(1)
	}

	// generate go.mod file
	if err := generateGoMod(filepath.Join(fBaseDir, "go.mod")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// call go mod tidy
	cmdGoModTidy := exec.Command("go", "mod", "tidy")
	cmdGoModTidy.Dir = fBaseDir
	fmt.Println("running go mod tidy: " + cmdGoModTidy.String())
	if out, err := cmdGoModTidy.CombinedOutput(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		fmt.Printf("output: %s\n", string(out))
		os.Exit(1)
	}

	// call go run
	cmdGoRun := exec.Command("go", "run", "main.go", bindingFile)
	cmdGoRun.Dir = fBaseDir
	fmt.Println("running go run: " + cmdGoRun.String())
	if out, err := cmdGoRun.CombinedOutput(); err != nil {
		fmt.Printf("error: %s\n", err.Error())
		fmt.Printf("output: %s\n", string(out))
		os.Exit(1)
	}

}

func generateGoMod(filename string) error {
	fmt.Println("generating " + filename)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	tmpl, err := template.New("").Parse(tmplGoMod)
	if err != nil {
		return err
	}

	// execute template
	return tmpl.Execute(f, nil)
}

func generateMain(tmplStr, filename, proof, publicInputs string, nbPublicInputs, fNbCommitments int) error {
	fmt.Println("generating " + filename)
	f, err := os.Create(filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer f.Close()

	helpers := template.FuncMap{
		"mul": func(a, b int) int {
			return a * b
		},
	}

	tmpl, err := template.New("").Funcs(helpers).Parse(tmplStr)
	if err != nil {
		return err
	}

	data := struct {
		Proof          string
		PublicInputs   string
		NbPublicInputs int
		NbCommitments  int
	}{
		Proof:          proof,
		PublicInputs:   publicInputs,
		NbPublicInputs: nbPublicInputs,
		NbCommitments:  fNbCommitments,
	}

	// execute template
	return tmpl.Execute(f, data)

}

func init() {
	rootCmd.AddCommand(verifyCmd)

	verifyCmd.Flags().StringVarP(&fProof, "proof", "p", "", "hex encoded proof to verify")
	verifyCmd.Flags().StringVar(&fPublicInputs, "public-inputs", "", "hex encoded public inputs to verify")
	verifyCmd.Flags().IntVarP(&fNbPublicInputs, "nb-public-inputs", "n", 0, "number of public inputs")

	verifyCmd.MarkFlagRequired("proof")
	verifyCmd.MarkFlagRequired("public-inputs")
	verifyCmd.MarkFlagRequired("nb-public-inputs")

	verifyCmd.Flags().BoolVar(&fGroth16, "groth16", false, "use groth16 verification")
	verifyCmd.Flags().BoolVar(&fPlonK, "plonk", false, "use plonk verification")

	verifyCmd.MarkFlagsMutuallyExclusive("groth16", "plonk")

	verifyCmd.Flags().IntVar(&fNbCommitments, "commitment", 0, "number of commitments in proof")
}
