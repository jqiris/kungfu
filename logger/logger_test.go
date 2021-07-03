package logger

import (
	"fmt"
	"github.com/mattn/go-colorable"
	"testing"
)

func TestLoggerColor(t *testing.T) {
	//fmt.Printf("\033[1;37;41m%s\033[0m\n", "Red.")
	//d := color.New(color.FgHiYellow)
	//_, err := d.Printf("hello world")
	//if err != nil {
	//	t.Fatal(err)
	//}

	var out = colorable.NewColorableStdout()
	fmt.Fprintf(out, "\x1b[%dm", 31)
	fmt.Fprintf(out, "%*s | ", 10, "adbd")
	fmt.Fprintf(out, "\x1b[m")

}
