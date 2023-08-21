package main

import (
	"fmt"
	"image"
	"os"
	"testing"
)

func TestListDominantColors(t *testing.T) {
	reader, err := os.Open("images/test.jpg")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	defer reader.Close()

	image, _, err := image.Decode(reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s", err)
	}

	colors := ListDominantColors(image, 5)
	fmt.Println(colors.Theme())
	for _, color := range colors {
		fmt.Printf("%s %s %s %.2f%%\n", color.Theme, color.Hex, color.Name, color.Proportion)
	}
}
