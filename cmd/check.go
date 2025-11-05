package cmd

import (
	"errors"
	"fmt"
	"sync"

	"github.com/vcircosta/go-livecoding/internal/checker"
	"github.com/vcircosta/go-livecoding/internal/config"
	"github.com/vcircosta/go-livecoding/internal/reporter"

	"github.com/spf13/cobra"
)

var (
	inputFilePath  string
	outputFilePath string
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Vérifie l'accessibilité d'une liste d'URLs.",
	Long:  `La commande 'check' parcourt une liste prédéfinie d'URLs et affiche leur statut d'accessibilité en utilisant des goroutines pour la concurrence.`,
	Run: func(cmd *cobra.Command, args []string) {

		if inputFilePath == "" {
			fmt.Println("Erreur: le chemin du fichier d'entrée (--input) est obligatoire.")
			return
		}

		targets, err := config.LoadTargetsFromFile(inputFilePath)
		if err != nil {
			fmt.Printf("Erreur lors du chargement des URLs: %v\n", err)
			return
		}

		if len(targets) == 0 {
			fmt.Println("Aucune URL à vérifier dans le fichier d'entrée.")
			return
		}

		var wg sync.WaitGroup
		resultsChan := make(chan checker.CheckResult, len(targets))

		wg.Add(len(targets))
		for _, target := range targets {
			go func(t config.InputTarget) {
				defer wg.Done()
				result := checker.CheckURL(t)
				resultsChan <- result // Envoyer le résultat au canal
			}(target)
		}
		wg.Wait()
		close(resultsChan)

		var finalReport []checker.ReportEntry
		for res := range resultsChan { // Récupérer tous les résultats du canal
			reportEntry := checker.ConvertToReportEntry(res)
			finalReport = append(finalReport, reportEntry)

			// Affichage immédiat comme avant
			if res.Err != nil {
				var unreachable *checker.UnreachableError
				if errors.As(res.Err, &unreachable) {
					fmt.Printf("KO %s (%s) est inaccessible : %v\n", res.InputTarget.Name, unreachable.URL, unreachable.Err)
				} else {
					fmt.Printf("KO %s (%s) : erreur - %v\n", res.InputTarget.Name, res.InputTarget.URL, res.Err)
				}
			} else {
				fmt.Printf("OK %s (%s) : OK - %s\n", res.InputTarget.Name, res.InputTarget.URL, res.Status)
			}
		}

		if outputFilePath != "" {
			err := reporter.ExportResultsToJsonFile(outputFilePath, finalReport)
			if err != nil {
				fmt.Printf("Erreur lors de l'exportation des résultats: %v\n", err)
			} else {
				fmt.Printf("✅ Résultats exportés vers %s\n", outputFilePath)
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&inputFilePath, "input", "i", "", "Chemin vers le fichier d'entrée")
	checkCmd.Flags().StringVarP(&outputFilePath, "output", "o", "", "Chemin vers le fichier de sortie")

	checkCmd.MarkFlagRequired("input")
}
