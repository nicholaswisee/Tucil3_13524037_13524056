package main

import (
	"fmt"
	"os"

	"github.com/nicholaswisee/Tucil3_13524037_13524056/core/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <path-to-input.txt>")
		fmt.Println("Example: go run main.go ../test/input/input1.txt")
		os.Exit(1)
	}

	path := os.Args[1]
	m, err := parser.ParseFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Parsed map: %dx%d\n", m.Height, m.Width)
	fmt.Printf("Start: %+v\n", m.StartPos)
	fmt.Printf("Goal:  %+v\n", m.GoalPos)
	fmt.Printf("Numbers: %d (", m.TotalNumbers)
	for i := 0; i < m.TotalNumbers; i++ {
		if i > 0 {
			fmt.Print(", ")
		}
		fmt.Printf("%d→%+v", i, m.NumberPos[i])
	}
	fmt.Println(")")
	fmt.Println("Grid:")
	for _, row := range m.Grid {
		fmt.Println(string(row))
	}
	fmt.Println("Costs:")
	for _, row := range m.Costs {
		for j, c := range row {
			if j > 0 {
				fmt.Print(" ")
			}
			fmt.Printf("%3d", c)
		}
		fmt.Println()
	}
}
