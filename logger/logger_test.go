package logger

import (
	"testing"
)

func TestLoggerColor(t *testing.T) {
	//fmt.Printf("\033[1;37;41m%s\033[0m\n", "Red.")
	//d := color.New(color.FgHiYellow)
	//_, err := d.Printf("hello world")
	//if err != nil {
	//	t.Fatal(err)
	//}

	Fatal("hello world")
	select {}

}
