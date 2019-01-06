package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode/utf8"
)

var foxs = []string{
	`        \
         \ /\   /\
          //\\_//\\     ____
          \_     _/    /   /
           / * * \    /^^^]
           \_\O/_/    [   ]
            /   \_    [   /
            \     \_  /  /
            [ [ /  \/ _/
           _[ [ \  /_/
	`,
	`                     /
                    /
___                /
"._'-.         (\-.
   '-.';.--.___/ _'>
     ''"( )    , )
         \\----\-\
         ""    """
	`,
	`      |
      \	
       \ |\__/|
        /     \
       /_.~ ~,_\
          \@/
	`,
	`     \
      \ ,-.      .-,
        |-.\ __ /.-|
        \  '    '  /
        / _     _  \
        | _'q  p _ |
        '._=/  \=_.'
          {'\()/'}'\
          {      }  \
          |{    }    \
          \ '--'   .- \
          |-      /    \
          | | | | |     ;
          | | |.;.,..__ |
        .-"";'         '|
       /    |           /
       '-../____,..---''
	`,
}

func buildDialogue(lines []string, maxwidth int) string {
	var borders, limits []string
	count := len(lines)

	borders = []string{"/", "\\", "\\", "/", "|", "<", ">"}

	top := " " + strings.Repeat("_", maxwidth+2)
	bottom := " " + strings.Repeat("-", maxwidth+2)

	limits = append(limits, top)
	if count == 1 {
		s := fmt.Sprintf("%s %s %s", borders[5], lines[0], borders[6])
		limits = append(limits, s)
	} else {
		s := fmt.Sprintf(`%s %s %s`, borders[0], lines[0], borders[1])
		limits = append(limits, s)
		i := 1
		for ; i < count-1; i++ {
			s = fmt.Sprintf(`%s %s %s`, borders[4], lines[i], borders[4])
			limits = append(limits, s)
		}
		s = fmt.Sprintf(`%s %s %s`, borders[2], lines[i], borders[3])
		limits = append(limits, s)
	}

	limits = append(limits, bottom)
	return strings.Join(limits, "\n")
}

func tabsToSpaces(lines []string) []string {
	var linesWithoutTabs []string
	for _, l := range lines {
		l = strings.Replace(l, "\t", "    ", -1)
		linesWithoutTabs = append(linesWithoutTabs, l)
	}
	return linesWithoutTabs
}

func calculateMaxWidth(lines []string) int {
	maxWidth := 0
	for _, l := range lines {
		len := utf8.RuneCountInString(l)
		if len > maxWidth {
			maxWidth = len
		}
	}

	return maxWidth
}

func normalizeStringsLength(lines []string, maxwidth int) []string {
	var linesNormalized []string
	for _, l := range lines {
		s := l + strings.Repeat(" ", maxwidth-utf8.RuneCountInString(l))
		linesNormalized = append(linesNormalized, s)
	}
	return linesNormalized
}

func rgb(i int) (int, int, int) {
	var f = 0.1
	return int(math.Sin(f*float64(i)+0)*127 + 128),
		int(math.Sin(f*float64(i)+2*math.Pi/3)*127 + 128),
		int(math.Sin(f*float64(i)+4*math.Pi/3)*127 + 128)
}

func printRGB(output string) {
	for j := 0; j < len(output); j++ {
		r, g, b := rgb(j)
		fmt.Printf("\033[38;2;%d;%d;%dm%c\033[0m", r, g, b, output[j])
	}
	fmt.Println()
}

func main() {
	info, _ := os.Stdin.Stat()

	if info.Mode()&os.ModeNamedPipe == 0 {
		log.Fatal("No se esta recibiendo entrada por pipes")
	}

	var lines []string

	reader := bufio.NewReader(os.Stdin)

	for {
		line, _, err := reader.ReadLine()
		if err != nil && err == io.EOF {
			break
		}
		lines = append(lines, string(line))
	}

	// Las tabulaciones se convierten es espacio
	lines = tabsToSpaces(lines)

	// Se obtiene el tamaño de la linea mas grande
	maxwidth := calculateMaxWidth(lines)

	// Se normaliza el espaciado a la derecha de las lineas en base la linea mas grande para evitar deformaciones del cuadro de dialogo
	messages := normalizeStringsLength(lines, maxwidth)

	// Se dibuja el cuadro de dialogo
	dialogue := buildDialogue(messages, maxwidth)

	printRGB(dialogue)

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)
	printRGB(foxs[r1.Intn(len(foxs))])
}
