package game

import(
	"fmt"
	"github.com/goracingkingsengine/gorke/board"
)

const START_FEN         = "8/8/8/8/8/8/krbnNBRK/qrbnNBRQ w - - 0 1"

//////////////////////////////////////////////////////

type TGame struct {
	B board.TBoard
	Moves board.TMoveList
	Node board.TNode
	Stop bool
	Ready bool
}

func (g *TGame) ToPrintable() string {
	return fmt.Sprintf("%s\n%s\n",g.Node.ToPrintable(),g.Moves.ToPrintable())
}

func (g *TGame) Print() {
	fmt.Printf("%s",g.ToPrintable())
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

	for g.Stop==false {
		tnodes:=len(board.NodeManager.Nodes)
		addnodes:=500
		if tnodes>100 {
			addnodes=5000
		}
		if tnodes>20000 {
			addnodes=10000
		}
		if tnodes>50000 {
			addnodes=25000
		}
		for i:=0; (i<addnodes) && (!g.Stop); i++ {
			g.Node.AddNode()
		}
		g.Node.MiniMaxOut()

		fmt.Printf("\n%s\ntotal nodes %d\n",g.ToPrintable(),len(board.NodeManager.Nodes))
	}

	g.Ready=true
	
}

//////////////////////////////////////////////////////