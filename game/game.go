package game

import(
	"fmt"
	"os"
	"github.com/goracingkingsengine/gorke/board"
)

//////////////////////////////////////////////////////

const UCI=true

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
	if !UCI {
		fmt.Printf("%s",g.ToPrintable())
	}
}

func (g *TGame) Init() {
	g.Moves=[]board.TMove{}
	board.InitNodeManager()
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
	g.Node.AbortMiniMax=false
}

func (g *TGame) AbortAnalysis() {
	g.Stop=true
	g.Node.AbortMiniMax=true
}

func (g *TGame) Analyze() {

	g.Ready=false

	depth:=1

	for g.Stop==false {

		for i:=0; (i<5000) && (!g.Stop); i++ {
				g.Node.AddNode(depth)
		}

		g.Node.MiniMaxOut(depth)

		nodes:=len(board.NodeManager.Nodes)

		if !UCI {
			fmt.Printf("\n%s\ndepth %d nodes %d\n",g.ToPrintable(),depth,nodes)
		} else {
			for i:=0; i<g.Multipv; i++ {
				if len(g.Node.Moves)>i {
					g.B.MakeMove(g.Node.Moves[i])
					line:=fmt.Sprintf("info depth %d multipv %d nodes %d score cp %d pv %s %s\n",
						depth,i+1,nodes,g.Node.Moves[i].Eval,g.Node.Moves[i].ToAlgeb(),g.B.GetLine())
					os.Stdout.Write([]byte(line))
					g.B.UnMakeMove(g.Node.Moves[i])
				}
			}
		}

		depth++
	}

	g.Ready=true
	
}

func (g *TGame) SendBestMove() {
	sendbestmove:=fmt.Sprintf("bestmove %s\n",g.Node.Moves[0].ToAlgeb())
	os.Stdout.Write([]byte(sendbestmove))
}

//////////////////////////////////////////////////////