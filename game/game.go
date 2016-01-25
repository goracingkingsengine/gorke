package game

import(
	"fmt"
	"os"
	"time"
	"github.com/goracingkingsengine/gorke/board"
)

//////////////////////////////////////////////////////

const START_FEN         = "8/8/8/8/8/8/krbnNBRK/qrbnNBRQ w - - 0 1"

//////////////////////////////////////////////////////

type TGame struct {
	B board.TBoard
	Moves board.TMoveList
	Node board.TNode
	Stop bool
	Ready bool
	Multipv int
}

func (g *TGame) ToPrintable() string {
	return fmt.Sprintf("%s\n%s\n",g.Node.ToPrintable(),g.Moves.ToPrintable())
}

func (g *TGame) Print() {
	//fmt.Printf("%s",g.ToPrintable())
}

func (g *TGame) Init() {
	g.Moves=[]board.TMove{}
	board.InitNodeManager()
	board.EvalDepth=1
	g.Node=g.B.CreateNode()
	g.Print()
}

func (g *TGame) Reset() {
	g.B.SetFromFen(START_FEN)
	g.Init()
	g.Multipv=1
}

func (g *TGame) SetFromFen(fen string) bool {
	success:=g.B.SetFromFen(fen)
	if success {
		g.Init()
		return true
	} else {
		g.Reset()
		return false
	}
}

func (g *TGame) MakeMove(i int) {
	if (i<0) || (i>len(g.Node.Moves)) {
		return
	}
	m:=g.Node.Moves[i]
	g.B.MakeMove(m)
	g.Moves=append(g.Moves,m)
	g.Node=g.B.CreateNode()
	g.Print()
}

func (g *TGame) MakeAlgebMove(algeb string) bool {
	for i:=0; i<len(g.Node.Moves); i++ {
		currentalgeb:=g.Node.Moves[i].ToAlgeb()
		if currentalgeb==algeb {
			g.MakeMove(i)
			return true
		}
	}
	return false
}

func (g *TGame) DelMove() {
	g.DelMoveInner()
	g.Print()
}

func (g *TGame) DelMoveInner() bool {
	l:=len(g.Moves)
	if l>0 {
		g.B.UnMakeMove(g.Moves[l-1])
		g.Node=g.B.CreateNode()
		g.Moves=g.Moves[0:l-1]
		return true
	}
	return false
}

func (g *TGame) DelAllMoves() {
	for g.DelMoveInner() {
	}
	g.Print()
}

//////////////////////////////////////////////////////

func (g *TGame) ClearAbortAnalysis() {
	g.Stop=false
	board.AbortMiniMax=false	
}

func (g *TGame) AbortAnalysis() {
	g.Stop=true
	board.AbortMiniMax=true
}

func (g *TGame) Analyze() {

	startingTime := time.Now().UTC()

	g.Ready=false

	depth:=1

	board.Nodes=0

	board.EvalDepth=1

	g.Init()

	board.EvalDepth=2

	for g.Stop==false {

		for k:=0; k<depth; k++ {

			nodes:=board.Nodes

			//fmt.Printf("\n%s\ndepth %d nodes %d\n",g.ToPrintable(),depth,nodes)
			
			for i:=0; i<g.Multipv; i++ {
				if len(g.Node.Moves)>i {
					g.B.MakeMove(g.Node.Moves[i])
					currentTime := time.Now().UTC()
					durationMilliSeconds:=currentTime.Sub(startingTime).Nanoseconds()/1e6
					durationSeconds:=currentTime.Sub(startingTime).Nanoseconds()/1e9
					nps:=int64(0)
					if(durationSeconds>0) {
						nps=int64(nodes)/durationSeconds
					}
					line:=fmt.Sprintf("info depth %d time %d nps %d multipv %d nodes %d score cp %d pv %s %s\n",
						depth,durationMilliSeconds,nps,i+1,nodes,g.Node.Moves[i].Eval,g.Node.Moves[i].ToAlgeb(),g.B.GetLine())
					os.Stdout.Write([]byte(line)[0:])
					g.B.UnMakeMove(g.Node.Moves[i])
				}
			}

			for j:=0; (j<depth) && (!g.Stop); j++ {
				//fmt.Printf("adding node ... ")
				//res:=
				g.Node.AddNode(depth)
				//fmt.Printf("res %v\n",res)
			}

			g.Node.MiniMaxOut(depth)

		}

		depth++
	}

	board.EvalDepth=1

	g.Ready=true

	g.ClearAbortAnalysis()	
}

func (g *TGame) SendBestMove() {
	sendbestmove:=fmt.Sprintf("bestmove %s\n",g.Node.Moves[0].ToAlgeb())
	os.Stdout.Write([]byte(sendbestmove))
}

//////////////////////////////////////////////////////