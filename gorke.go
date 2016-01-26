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

var token string

func NextToken() bool {
	if s.Scan()!=scanner.EOF {
		token=s.TokenText()
		return true
	} else {
		token=""
		return false
	}
}

func GetRest() string {
	var p=s.Pos().Offset+1
	if (p>=len(commandline)) {
		return ""
	}
	return commandline[p:]
}

func Log(what string) {
	f, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
	    panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(what); err != nil {
	    panic(err)
	}
}

var enginerunning=false

func StopEngine() {
	g.AbortAnalysis()
	for g.Ready==false {
		time.Sleep(100 * time.Millisecond)
	}
	enginerunning=false
	g.SendBestMove()
}

func Printu(ucistr string) {
	os.Stdout.Write([]byte(ucistr))
}

func main() {

	f,err:=os.Create("log.txt")
	if err!=nil {
		panic(err)
	} else {
		f.Close()
	}

	board.Init()

	Printu("Gorke Racing Kings Chess Variant Engine\n\n")

	g.Reset()

	var command string=""

	reader := bufio.NewReader(os.Stdin)

	for command!="x" {

		//fmt.Print("\n> ")

		commandline, _ = reader.ReadString('\n')

		Log(commandline)

		commandline=strings.TrimSpace(commandline)

		s.Init(strings.NewReader(commandline))
		s.Mode=scanner.ScanIdents|scanner.ScanInts

		if NextToken() {

			command=s.TokenText()

			if command=="l" {
				fmt.Print("l - list commands\n")
				fmt.Print("f - set from fen\n")
				fmt.Print("p - print\n")
				fmt.Print("r - reset\n")
				fmt.Print("a - analyze\n")
				fmt.Print("b - alphabeta\n")
				fmt.Print("q - toggle quiescence search\n")
				fmt.Print("s - stop\n")
				fmt.Print("m [algeb] - make algeb move\n")
				fmt.Print("d - del move\n")
				fmt.Print("dd - del all moves\n")
				fmt.Print("x - exit\n")
			}

			if command=="quit" {
				return
			}

			if command=="isready" {
				Printu("readyok\n")
			}

			if command=="uci" {
				Printu("id name Gorke\n")
				Printu("id author golang\n")
				Printu("\n")
				Printu("option name MultiPV type spin default 1 min 1 max 500\n")
				Printu("uciok\n")
			}

			if command=="q" {
				board.DoQuiescence=!board.DoQuiescence
				fmt.Printf("doquiescence %v\n",board.DoQuiescence)
			}

			if command=="p" {
				g.Print()
			}

			if command=="m" {
				if NextToken() {
					g.MakeAlgebMove(token)
					g.Print()
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

			if (command=="go") {
				if g.Multipv>1 {
					command="a"
				} else {
					command="b"
				}
			}

			if (command=="a") {
				if enginerunning {
					StopEngine()
				}
				enginerunning=true
				go g.Analyze()
			}

			if (command=="b") {
				if enginerunning {
					StopEngine()
				}
				enginerunning=true
				go g.AlphaBeta()
			}

			if (command=="s") || (command=="stop") {
				StopEngine()
			}

			if command=="f" {
				g.SetFromFen(GetRest())
			}

			if command=="position" {
				// 0 - position
				// 1 - fen
				// 2 - startpos
				// 3 - moves
				// 4 - move1
				// 5 - move2
				// 6 - ...

				// or...

				// 0 - position
				// 1 - fen
				// 2 - posfen
				// 3 - turnfen
				// 4 - castlefen
				// 5 - epfen
				// 6 - halfmovefen
				// 7 - fullmovefen
				// 8 - moves
				// 9 - move1
				// 10 - move2
				// 11 - ...

				if enginerunning {
					StopEngine()
				}
				if NextToken() {
					if(token=="fen") {						
						if NextToken() {
							fen:=game.START_FEN
							parts:=strings.Split(commandline," ")
							movesat:=3
							if(token!="startpos") {								
								if len(parts)>=8 {
									fen=fmt.Sprintf("%s %s %s %s %s %s",parts[2],parts[3],parts[4],parts[5],parts[6],parts[7])
									movesat=8
								}
							}
							g.SetFromFen(fen)
							if len(parts)>movesat {
								if parts[movesat]=="moves" {
									for i:=movesat+1; i<len(parts); i++ {
										g.MakeAlgebMove(parts[i])
									}
									g.B.SetDepth(0)
									g.Node.Depth=0
								}
							}
						}						
					}
				}
				if enginerunning {
					go g.Analyze()
				}
			}

			if command=="setoption" {
				if NextToken() {
					if(token=="name") {
						if NextToken() {
							name:=token
							if NextToken() {
								if(token=="value") {
									if NextToken() {
										value:=token

										if name=="MultiPV" {
											i,err:=strconv.Atoi(value)
											if err==nil {
												g.Multipv=i
											}
										}
									}
								}
							}
						}
					}
				}
			}
			
		}

	}

}