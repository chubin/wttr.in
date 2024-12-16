package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/chubin/vt10x"
	"github.com/fogleman/gg"
)

func StringSliceToRuneSlice(s string) [][]rune {
	strings := strings.Split(s, "\n")
	result := make([][]rune, len(strings))

	i := 0
	for _, str := range strings {
		if len(str) == 0 {
			continue
		}
		result[i] = []rune(str)
		i++
	}

	return result
}

func maxRowLength(rows [][]rune) int {
	maxLen := 0
	for _, row := range rows {
		if len(row) > maxLen {
			maxLen = len(row)
		}
	}
	return maxLen
}

func GeneratePng() {
	runes := StringSliceToRuneSlice(`
Weather report: Hochstadt an der Aisch, Germany

     \  /       Partly cloudy
   _ /"".-.     +5(2) °C
     \_(   ).   ↗ 9 km/h
     /(___(__)  10 km
                0.0 mm
                        ┌─────────────┐
┌───────────────────────┤  Sat 11 Nov ├───────────────────────┐
│             Noon      └──────┬──────┘      Night            │
├──────────────────────────────┼──────────────────────────────┤
│  _'/"".-.     Patchy rain po…│  _'/"".-.     Patchy rain po…│
│   ,\_(   ).   +6(3) °C       │   ,\_(   ).   +5(2) °C       │
│    /(___(__)  → 22-29 km/h   │    /(___(__)  ↗ 14-20 km/h   │
│      ‘ ‘ ‘ ‘  10 km          │      ‘ ‘ ‘ ‘  10 km          │
│     ‘ ‘ ‘ ‘   0.1 mm | 86%   │     ‘ ‘ ‘ ‘   0.0 mm | 89%   │
└──────────────────────────────┴──────────────────────────────┘
                        ┌─────────────┐
┌───────────────────────┤  Sun 12 Nov ├───────────────────────┐
│             Noon      └──────┬──────┘      Night            │
├──────────────────────────────┼──────────────────────────────┤
│    \  /       Partly cloudy  │      .-.      Light drizzle  │
│  _ /"".-.     +8(7) °C       │     (   ).    +5(2) °C       │
│    \_(   ).   ↑ 7-8 km/h     │    (___(__)   ↑ 13-18 km/h   │
│    /(___(__)  10 km          │     ‘ ‘ ‘ ‘   2 km           │
│               0.0 mm | 0%    │    ‘ ‘ ‘ ‘    0.3 mm | 76%   │
└──────────────────────────────┴──────────────────────────────┘
`)

	// Dimensions of each rune in pixels
	runeWidth := 8
	runeHeight := 14

	// Compute the width and height of the final image
	imageWidth := runeWidth * maxRowLength(runes)
	imageHeight := runeHeight * len(runes)

	// Create a new context with the computed dimensions
	dc := gg.NewContext(imageWidth, imageHeight)

	// fontPath := "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf"
	// fontPath := "/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc"
	fontPath := "/usr/share/fonts/truetype/lexi/LexiGulim.ttf"

	err := dc.LoadFontFace(fontPath, 13)
	if err != nil {
		log.Fatal(err)
	}

	// Loop through each rune in the array and draw it on the context
	for i, row := range runes {
		for j, char := range row {
			// Compute the x and y coordinates for drawing the current rune
			x := float64(j*runeWidth + runeWidth/2)
			y := float64(i*runeHeight + runeHeight/2)

			// Set the appropriate color for the current rune
			if char == '#' {
				dc.SetRGB(0, 0, 0) // Black
			} else if char == '@' {
				dc.SetRGB(1, 0, 0) // Red
			} else {
				dc.SetRGB(1, 1, 1) // White
			}

			character := string(char)
			// if char == ' ' {
			// 	character = fmt.Sprint(j % 10)
			// }
			dc.DrawRectangle(x, y, x+float64(runeWidth), y+float64(runeHeight))
			dc.Fill()

			// Draw a rectangle with the rune's dimensions and color
			dc.DrawString(character, x, y) // Draw the character centered on the canvas
			// dc.DrawStringAnchored(character, x, y, 0.5, 0.5) // Draw the character centered on the canvas
		}
	}

	// Save the image to a PNG file
	err = dc.SavePNG("output.png")
	if err != nil {
		fmt.Println("Error saving PNG:", err)
		return
	}

	fmt.Println("PNG generated successfully")
}

func GeneratePngFromANSI(input []byte, outputFile string) error {
	// Dimensions of each rune in pixels
	runeWidth := 8
	runeHeight := 14
	fontSize := 13.0
	// fontPath := "/usr/share/fonts/truetype/dejavu/DejaVuSansMono.ttf"
	fontPath := "/usr/share/fonts/truetype/wqy/wqy-zenhei.ttc"

	imageCols := 80
	imageRows := 25

	// Compute the width and height of the final image
	imageWidth := runeWidth * imageCols
	imageHeight := runeHeight * imageRows

	// Create terminal and feed it with input.
	term := vt10x.New(vt10x.WithSize(imageCols, imageRows))
	_, err := term.Write([]byte("\033[20h"))
	if err != nil {
		return fmt.Errorf("virtual terminal write error: %w", err)
	}

	_, err = term.Write(input)
	if err != nil {
		return fmt.Errorf("virtual terminal write error: %w", err)
	}

	// Create a new context with the computed dimensions
	dc := gg.NewContext(imageWidth, imageHeight)

	err = dc.LoadFontFace(fontPath, fontSize) // Set font size to 96
	if err != nil {
		return fmt.Errorf("error loading font: %w", err)
	}

	// Loop through each rune in the array and draw it on the context
	for i := 0; i < imageRows; i++ {
		for j := 0; j < imageCols; j++ {
			// Compute the x and y coordinates for drawing the current rune
			x := float64(j * runeWidth)
			y := float64(i * runeHeight)

			cell := term.Cell(j, i)
			character := string(cell.Char)

			dc.DrawRectangle(x, y, float64(runeWidth), float64(runeHeight))
			bg := colorANSItoRGB(cell.BG)
			dc.SetRGB(bg[0], bg[1], bg[2])
			dc.Fill()

			fg := colorANSItoRGB(cell.FG)
			dc.SetRGB(fg[0], fg[1], fg[2])

			// Draw a rectangle with the rune's dimensions and color
			dc.DrawString(character, x, y+float64(runeHeight)-3) // Draw the character centered on the canvas
			// dc.DrawStringAnchored(character, x, y, 0.5, 0.5) // Draw the character centered on the canvas
		}
	}

	// Save the image to a PNG file
	err = dc.SavePNG(outputFile)
	if err != nil {
		return fmt.Errorf("error saving png: %w", err)
	}

	return nil
}

func colorANSItoRGB(colorANSI vt10x.Color) [3]float64 {
	defaultBG := vt10x.Color(0)
	defaultFG := vt10x.Color(8)

	if colorANSI == vt10x.DefaultFG {
		colorANSI = defaultFG
	}
	if colorANSI == vt10x.DefaultBG {
		colorANSI = defaultBG
	}

	if colorANSI > 255 {
		return [3]float64{127, 127, 127}
	}
	return ansiColorsDB[colorANSI]
}

func main() {
	data, err := os.ReadFile("zh-text.txt")
	if err != nil {
		log.Fatalln(err)
	}

	err = GeneratePngFromANSI(data, "output.png")
	if err != nil {
		log.Fatalln(err)
	}
}
