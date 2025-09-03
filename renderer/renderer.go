package renderer

import (
	"fmt"
	"github.com/fatih/color"
)

func Render(s string) {
	fmt.Println(s)
}

func RenderGreen(s string) {
	color.Green(s)
}

func RenderRed(s string) {
	color.Red(s)
}

func RenderYellow(s string) {
	color.Yellow(s)
}

func RenderBlue(s string) {
	color.Blue(s)
}