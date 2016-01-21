package main

import(
	"fmt"
	"os"
	"bufio"
	"strings"
	"strconv"
	"time"
	"text/scanner"
	"github.com/goracingkingsengine/gorke/piece"
	"github.com/goracingkingsengine/gorke/square"
	"github.com/goracingkingsengine/gorke/board"
)

var b board.TBoard
var s scanner.Scanner
var commandline string
var node board.TNode

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

func Print() {
	fmt.Printf("%s",node.ToPrintable())
}

var stop=false
var ready=false

func Analyze() {
	ready=false
	for stop==false {
		for i:=0; (i<500) && (!stop); i++ {
			node.AddNode()
		}
		node.MiniMaxOut()
		Print()
		fmt.Printf("\ntotal nodes %d\n",len(board.NodeManager.Nodes))
	}
	ready=true
}

func main() {

	fmt.Printf("\nGorke - Go Racing Kings Chess Variant Engine\n")

	board.Init()

	Reset()

	node=b.CreateNode()

	game:=[]board.TNode{node}

	Print()

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
				fmt.Print("im - init move gen\n")
				fmt.Print("ns - next sq\n")
				fmt.Print("np - next pseudo legal move\n")
				fmt.Print("r - reset\n")
				fmt.Print("p - print\n")
				fmt.Print("a - analyze\n")
				fmt.Print("s - stop\n")
				fmt.Print("i - minimax out\n")
				fmt.Print("m i - make ith node move\n")
				fmt.Print("d - del move\n")
				fmt.Print("x - exit\n")
			}

			if command=="m" {
				tok=s.Scan()

				if tok!=scanner.EOF {
					i,err:=strconv.Atoi(s.TokenText())
					if err==nil {
						b.MakeMove(node.Moves[i])
						node=b.CreateNode()
						game=append(game,node)
						command="p"
					}
				}
			}

			if command=="d" {
				if len(game)>1 {
					game=game[0:len(game)-1]
				}
				node=game[len(game)-1]
				command="p"
			}

			if command=="r" {
				Reset()
				node=b.CreateNode()
				game=[]board.TNode{node}
				command="p"
			}

			if command=="a" {
				stop=false
				go Analyze()
			}

			if command=="s" {
				stop=true
				for ready==false {
					time.Sleep(100 * time.Millisecond)
				}
			}

			if command=="i" {
				node.MiniMaxOut()
			}

			if command=="p" {
				Print()
			}

			if command=="im" {
				b.InitMoveGen()
				fmt.Printf("%s\n",b.ReportMoveGen())
			}

			if command=="ns" {
				b.CurrentSq++
				b.NextSq()
				fmt.Printf("%s\n",b.ReportMoveGen())
			}

			if command=="np" {
				res:=b.NextPseudoLegalMove()
				fmt.Printf("res %v - %s\n",res,b.ReportMoveGen())
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