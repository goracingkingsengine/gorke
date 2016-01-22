package main

import(
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"time"
	"text/scanner"
	"github.com/goracingkingsengine/gorke/board"
	"github.com/goracingkingsengine/gorke/game"
)

var g=new(game.TGame)

var s scanner.Scanner
var commandline string

func GetRest() string {
	var p=s.Pos().Offset+1
	if (p>=len(commandline)) {
		return ""
	}
	return commandline[p:]
}

func main() {

	board.Init()

	fmt.Printf("\nGorke - Go Racing Kings Chess Variant Engine\n")

	g.Reset()

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
				fmt.Print("f - set from fen\n")
				fmt.Print("r - reset\n")
				fmt.Print("a - analyze\n")
				fmt.Print("s - stop\n")
				fmt.Print("m [i] - make ith node move\n")
				fmt.Print("d - del move\n")
				fmt.Print("dd - del all moves\n")
				fmt.Print("x - exit\n")
			}

			if command=="m" {
				tok=s.Scan()

				if tok!=scanner.EOF {
					i,err:=strconv.Atoi(s.TokenText())
					if err==nil {
						g.MakeMove(i-1)
					}
				}
			}

			if command=="d" {
				g.DelMove()
			}

			if command=="dd" {
				g.DelAllMoves()
			}

			if command=="r" {
				g.Reset()
			}

			if command=="a" {
				g.ClearAbortAnalysis()
				go g.Analyze()
			}

			if command=="s" {
				g.AbortAnalysis()
				for g.Ready==false {
					time.Sleep(100 * time.Millisecond)
				}
			}

			if command=="f" {
				g.SetFromFen(GetRest())
			}
			
		}

	}

}