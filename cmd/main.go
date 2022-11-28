package main

import (
	"fmt"
	"log"
	"math"

	"github.com/spf13/cobra"

	"github.com/landru29/fraction"
)

func main() {
	var (
		n    fraction.Number
		deep int
	)

	cmd := &cobra.Command{
		Use:   "fraction",
		Short: "Convert a floating number in fraction",
		Long:  "Try to find the rationalfraction of a floating number",
		RunE: func(cmd *cobra.Command, args []string) error {
			if n.IsZero() {
				fmt.Println("0")
				return nil
			}

			numerator, denominator, exact := n.Render(deep)

			globalNumerator := numerator + n.Numerator(denominator)

			fmt.Printf(
				"%s %d / %d %s\n",
				n.Sign(),
				globalNumerator,
				denominator,
				map[bool]string{
					true:  "!",
					false: fmt.Sprintf("(%f)", math.Abs(n.Raw())-float64(globalNumerator)/float64(denominator)),
				}[exact],
			)

			return nil
		},
	}

	cmd.PersistentFlags().VarP(&n, "float", "f", "floating number")
	cmd.PersistentFlags().IntVarP(&deep, "deep", "d", 1000000, "deep computing")

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
