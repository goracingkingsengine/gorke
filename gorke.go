package main

import(
	"fmt"
	"os"
	"bufio"
	"strings"
	"text/scanner"
	"github.com/goracingkingsengine/gorke/piece"
)

func main() {

	fmt.Printf("Gorke - Go Racing Kings Chess Variant Engine")

	var command string=""

	reader := bufio.NewReader(os.Stdin)

	for command!="x" {

		fmt.Print("\n> ")

		commandline, _ := reader.ReadString('\n')

		var s scanner.Scanner
		s.Init(strings.NewReader(commandline))
		s.Mode=scanner.ScanIdents|scanner.ScanInts
		var tok rune

		tok = s.Scan()

		if tok!=scanner.EOF {
			command=s.TokenText()

			if command=="list" {
				fmt.Print("list - list commands\n")
				fmt.Print("ftop - fenchar to piece\n")
				fmt.Print("x - exit\n")
			}

			if command=="ftop" {
				tok=s.Scan()

				var fenchar byte=' '

				if tok!=scanner.EOF {
					fenchar=s.TokenText()[0]
				}

				var p=piece.FromFenChar(fenchar)

				fmt.Printf("piece code %d type %d color %d",p,piece.TypeOf(p),piece.ColorOf(p))
			}
		}

	}

}