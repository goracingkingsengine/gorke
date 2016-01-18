package main

import(
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"text/scanner"
	"github.com/goracingkingsengine/gorke/piece"
	"github.com/goracingkingsengine/gorke/square"
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
				fmt.Print("atos - algeb to square\n")
				fmt.Print("stoa - square to algeb\n")
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

			if command=="atos" {
				tok=s.Scan()

				if tok!=scanner.EOF {
					var s=square.FromAlgeb(s.TokenText())

					fmt.Printf("square %d",s)
				}
			}

			if command=="stoa" {
				tok=s.Scan()

				if tok!=scanner.EOF {
					sq,err:=strconv.Atoi(s.TokenText())

					if err==nil {
						fmt.Printf("algeb %s",square.ToAlgeb(square.TSquare(sq)))
					}
				}
			}
		}

	}

}