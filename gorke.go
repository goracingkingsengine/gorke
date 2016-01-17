package main

import(
	"fmt"
	"os"
	"bufio"
	"github.com/goracingkingsengine/gorke/piece"
)

func main() {

	fmt.Printf("Gorke - Go Racing Kings Chess Variant Engine")

	var command string=""

	reader := bufio.NewReader(os.Stdin)

	for command!="x" {

		fmt.Print("\n> ")

		text, _ := reader.ReadString('\n')
		
		command=text[0:1]

		if command=="l" {
			fmt.Print("l - list commands\n")
			fmt.Print("f - fenchar to piece\n")
			fmt.Print("x - exit\n")
		}

		if command=="f" {
			var fenchar=text[2]

			var p=piece.FromFenChar(fenchar)

			fmt.Printf("fenchar %c piece code %d type %d color %d",fenchar,p,piece.TypeOf(p),piece.ColorOf(p))
		}

	}

}