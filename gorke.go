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
	"github.com/goracingkingsengine/gorke/board"
)

var b board.TBoard
var s scanner.Scanner
var commandline string

func GetRest() string {
	var p=s.Pos().Offset+1
	if (p>=len(commandline)) {
		return ""
	}
	return commandline[p:]
}

func Reset() {
	b.SetFromFen("8/8/8/8/8/8/krbnNBRK/qrbnNBRQ w - - 0 1")
}

func main() {

	fmt.Printf("Gorke - Go Racing Kings Chess Variant Engine")

	board.InitMoveTable()

	Reset()

	var command string=""

	reader := bufio.NewReader(os.Stdin)

	for command!="x" {

		fmt.Print("\n> ")

		commandline, _ = reader.ReadString('\n')

		s.Init(strings.NewReader(commandline))
		s.Mode=scanner.ScanIdents|scanner.ScanInts
		var tok rune

		tok = s.Scan()

		if tok!=scanner.EOF {
			command=s.TokenText()

			if command=="l" {
				fmt.Print("l - list commands\n")
				fmt.Print("ftop - fenchar to piece\n")
				fmt.Print("atos - algeb to square\n")
				fmt.Print("stoa - square to algeb\n")
				fmt.Print("mt sq p - move table\n")
				fmt.Print("f - set from fen\n")
				fmt.Print("r - reset\n")
				fmt.Print("p - print\n")
				fmt.Print("x - exit\n")
			}

			if command=="p" {
				fmt.Print(b.ToPrintable())
			}

			if command=="r" {
				Reset()
			}

			if command=="f" {
				b.SetFromFen(GetRest())
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

			if command=="mt" {
				tok=s.Scan()

				if tok!=scanner.EOF {
					sq,err:=strconv.Atoi(s.TokenText())

					if err==nil {
						tok=s.Scan()

						if tok!=scanner.EOF {
							p:=piece.FromFenChar(s.TokenText()[0])

							var key board.TMoveTableKey
							key.Sq=square.TSquare(sq)
							key.P=p

							ptr:=board.MoveTablePtrs[key]

							for !board.MoveTable[ptr].EndPiece {
								md:=board.MoveTable[ptr]
								fmt.Printf("ptr %d to %d next %d\n",ptr,md.To,md.NextVector)
								ptr++
							}
						}
					}
				}
			}
		}

	}

}