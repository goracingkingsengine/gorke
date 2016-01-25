package board

import(
	"fmt"
	"math"
	"math/rand"
	"time"
	"sort"
	"sync"
	"github.com/goracingkingsengine/gorke/piece"
	"github.com/goracingkingsengine/gorke/square"
)

const INDEX_OF_WHITE    = 1
const INDEX_OF_BLACK    = 0

const DRAW_SCORE        = 0

const MATE_LIMIT        = 9000
const MATE_SCORE        = 10000
const INFINITE_SCORE    = 15000
const INVALID_SCORE     = 20000

//////////////////////////////////////////

const TURN_INDEX        = 64
const DEPTH_INDEX       = 65
const KINGSPOS_INDEX    = 66

var PIECE_VALUES=map [piece.TPieceType]int{
	piece.NO_PIECE : 0,
	piece.KING : 0,
	piece.QUEEN : 700,
	piece.ROOK : 500,
	piece.BISHOP : 300,
	piece.KNIGHT : 300}

const KING_ADVANCE_VALUE = 400

type TPosition [square.BOARD_SIZE+4]byte

//////////////////////////////////////////

type TTurn piece.TColor

const MOVE_TABLE_MAX_SIZE = 20000

type TMoveTableKey struct {
	Sq square.TSquare
	P piece.TPiece
}

type TMoveDescriptor struct {
	To square.TSquare
	NextVector int
	EndPiece bool
}

type TMove struct {
	From square.TSquare
	To square.TSquare
	Piece piece.TPiece
	CapPiece piece.TPiece
	Eval int
}

type TBoard struct {
	Pos TPosition
	CurrentSq square.TSquare
	CurrentPiece piece.TPiece
	CurrentPtr int
	CurrentMove TMove
	KingPos [2]byte
	Material [2]int
}

type TMoveList []TMove

type TNode struct {
	B TBoard
	Moves TMoveList
	Visited int
}

type TNodeManager struct {
	Nodes map[TPosition]TNode
}

var ALL_PIECE_TYPES=[]piece.TPieceType{piece.KING,piece.QUEEN,piece.ROOK,piece.BISHOP,piece.KNIGHT}

var MoveTable [MOVE_TABLE_MAX_SIZE]TMoveDescriptor

var MoveTablePtrs=make(map[TMoveTableKey]int)

var NodeManager TNodeManager

var AbortMiniMax=false

var EvalDepth=1

var Nodes=0

var r=rand.New(rand.NewSource(time.Now().UnixNano()))

var randStartTime=time.Now().UTC()

////////////////////////////////////////

func (b *TBoard) CreateNode() TNode {
	if NodeManager.DoesNodeExist(b.Pos) {
		//fmt.Printf("node already exists\n")
		return NodeManager.GetNode(b.Pos)
	}
	var node TNode
	node.B=*b
	b.InitMoveGen()
	node.Moves=[]TMove{}
	node.Visited=1
	for b.NextLegalMove() {
		node.Moves=append(node.Moves,b.CurrentMove)
	}
	node.Eval()
	NodeManager.AddNode(b.Pos,node)
	return node
}

func (b *TBoard) SetSq(sq square.TSquare,p piece.TPiece) {
	b.Pos[byte(sq)]=byte(p)
}

func (b *TBoard) IsSquareColInCheck(sq square.TSquare, c piece.TColor) bool {
	ksq:=b.Pos.GetKingPos(c)
	for _, pt := range ALL_PIECE_TYPES {
		test_piece:=piece.FromTypeAndColor(pt,piece.InvColorOf(c))
		ptr:=GetMoveTablePtr(ksq,test_piece)
		for !MoveTable[ptr].EndPiece {
			md:=MoveTable[ptr]
			p:=b.PieceAtSq(md.To)
			if p==test_piece {
				return true
			}
			if piece.ColorOf(p)!=piece.NO_COLOR {
				ptr=md.NextVector
			} else {
				ptr++
			}
		}
	}
	return false
}

func (b *TBoard) IsSqInCheck(sq square.TSquare) bool {
	return b.IsSquareColInCheck(sq,b.GetColorOfTurn())
}

func (b *TBoard) IsColInCheck(c piece.TColor) bool {
	return b.IsSquareColInCheck(b.Pos.GetKingPos(c),c)
}

func (b *TBoard) IsInCheck() bool {
	return b.IsColInCheck(b.GetColorOfTurn())
}

func (b *TBoard) IsOppInCheck() bool {
	return b.IsColInCheck(piece.InvColorOf(b.GetColorOfTurn()))
}

func (b *TBoard) MakeMove(m TMove) {

	b.SetSq(m.From,piece.NO_PIECE)

	b.SetSq(m.To,m.Piece)

	if(m.Piece==(piece.WHITE|piece.KING)) {
		b.Pos[KINGSPOS_INDEX+1]=byte(m.To)
	}

	if(m.Piece==(piece.BLACK|piece.KING)) {
		b.Pos[KINGSPOS_INDEX]=byte(m.To)
	}

	b.Material[IndexOfColor(b.GetColorOfInvTurn())]-=PIECE_VALUES[piece.TypeOf(m.CapPiece)]

	b.SetDepth(b.GetDepth()+1)

	b.SetTurn(InvTurnOf(b.GetTurn()))

}

func (b *TBoard) UnMakeMove(m TMove) {

	b.SetSq(m.From,m.Piece)

	if(m.Piece==(piece.WHITE|piece.KING)) {
		b.Pos[KINGSPOS_INDEX+1]=byte(m.From)
	}

	if(m.Piece==(piece.BLACK|piece.KING)) {
		b.Pos[KINGSPOS_INDEX]=byte(m.From)
	}

	b.SetSq(m.To,m.CapPiece)

	b.Material[IndexOfColor(b.GetColorOfTurn())]+=PIECE_VALUES[piece.TypeOf(m.CapPiece)]

	b.SetDepth(b.GetDepth()-1)

	b.SetTurn(InvTurnOf(b.GetTurn()))

}

func GetLineRecursive(b TBoard, line string) string {
	if NodeManager.DoesNodeExist(b.Pos) {
		n:=NodeManager.GetNode(b.Pos)
		if len(n.Moves)<=0 {
			return line
		}
		m:=n.Moves[0]
		b.MakeMove(m)
		return GetLineRecursive(b, line+m.ToAlgeb()+" ")
	} else {
		return line
	}
}

func (b TBoard) GetLine() string {
	return GetLineRecursive(b, "")
}

func (n *TNode) ToPrintable() string {
	var buff="\n----- Node: -----\n\n"
	buff+=fmt.Sprintf("%s\nMoves (%d):\n\n",n.B.ToPrintable(),len(n.Moves))
	b:=n.B
	for i,m := range n.Moves {
		b.MakeMove(m)
		line:=b.GetLine()
		buff+=fmt.Sprintf("%3d %s pv %s\n",i+1,m.ToPrintable(),line)
		b.UnMakeMove(m)
	}
	return fmt.Sprintf("%s\n----- End Node ( visited %d ) -----\n",buff,n.Visited)
}

func (b *TBoard) GetKingPos(c piece.TColor) square.TSquare {
	return b.Pos.GetKingPos(c)
}

// heuristics

func (b *TBoard) EvalCol(c piece.TColor) int {
	eval:=b.Material[IndexOfColor(c)]
	ksq:=b.GetKingPos(c)
	ksqr:=int(square.BOARD_HEIGHTL-square.RankOf(ksq))
	eval+=ksqr*KING_ADVANCE_VALUE
	return eval
}

var mutex = &sync.Mutex{}

func (b *TBoard) Eval() int {
	//mutex.Lock()
	e:=b.EvalCol(piece.WHITE)-b.EvalCol(piece.BLACK)//+r.Intn(25)-50
	//mutex.Unlock()
	if b.GetTurn()==piece.BLACK {
		return e
	}
	return -e
}

func GetMoveTablePtr(sq square.TSquare, p piece.TPiece) int {
	return MoveTablePtrs[TMoveTableKey{sq,p}]
}

func (b *TBoard) PieceAtSq(sq square.TSquare) piece.TPiece {
	return piece.TPiece(b.Pos[byte(sq)])
}

func (b *TBoard) ColorOfSq(sq square.TSquare) piece.TColor {
	return piece.ColorOf(piece.TPiece(b.Pos[byte(sq)]))
}

func (b *TBoard) NextSq() bool {
	for b.CurrentSq<square.BOARD_SIZE {
		if b.ColorOfSq(b.CurrentSq)==b.GetColorOfTurn() {
			b.CurrentPiece=b.PieceAtSq(b.CurrentSq)
			b.CurrentPtr=GetMoveTablePtr(b.CurrentSq,b.CurrentPiece)
			return true
		}
		b.CurrentSq++
	}
	return false
}

func (m *TMove) ToAlgeb() string {
	return fmt.Sprintf("%s%s",square.ToAlgeb(m.From),square.ToAlgeb(m.To))
}

func SignedEval(eval int) string {
	if eval<=0 {
		return fmt.Sprintf("%d",eval)
	}
	return fmt.Sprintf("+%d",eval)
}

func (m *TMove) ToPrintable() string {
	return fmt.Sprintf("%s %s",m.ToAlgeb(),SignedEval(m.Eval))
}

func (b *TBoard) ReportMoveGen() string {
	return fmt.Sprintf("current sq %d : current piece %c : current ptr %d current move %s",
		b.CurrentSq,piece.ToFenChar(b.CurrentPiece),b.CurrentPtr,b.CurrentMove.ToAlgeb())
}

func (b *TBoard) InitMoveGen() {
	b.CurrentSq=0
	b.NextSq()
}

func (b *TBoard) KingOnBaseRank(c piece.TColor) bool {
	return (square.RankOf(b.GetKingPos(c))==0)
}

func (b *TBoard) IsWhiteTurn() bool {
	return (b.GetTurn()==piece.WHITE)
}

func (b *TBoard) IsBlackTurn() bool {
	return (b.GetTurn()==piece.BLACK)
}

func (b *TBoard) WhiteKingOnBaseRank() bool {
	return b.KingOnBaseRank(piece.WHITE)
}

func (b *TBoard) BlackKingOnBaseRank() bool {
	return b.KingOnBaseRank(piece.BLACK)
}

func (b *TBoard) TerminalEval() int {
	wb:=b.WhiteKingOnBaseRank()
	bb:=b.BlackKingOnBaseRank()

	if wb && bb {
		return DRAW_SCORE
	}

	if bb {
		return -MATE_SCORE
	}

	if wb {
		return -MATE_SCORE
	}

	return DRAW_SCORE
}

func (b *TBoard) NextLegalMove() bool {
	if b.BlackKingOnBaseRank() {
		return false
	}
	wb:=b.WhiteKingOnBaseRank()
	if wb && b.IsWhiteTurn() {
		return false
	}
	for b.NextPseudoLegalMove() {
		b.MakeMove(b.CurrentMove)
		incheck:=b.IsInCheck()
		oppincheck:=b.IsOppInCheck()
		blackmated:=(wb && (!b.BlackKingOnBaseRank()))
		b.UnMakeMove(b.CurrentMove)
		ok:=(!incheck)&&(!oppincheck)&&(!blackmated)
		if !ok {
			//fmt.Printf("move thrown out for check %s\n",b.CurrentMove.ToAlgeb())
		}
		if ok {
			return true
		}
	}
	return false
}

func (b *TBoard) NextPseudoLegalMove() bool {
	for b.CurrentSq<square.BOARD_SIZE {
		md:=MoveTable[b.CurrentPtr]
		if md.EndPiece {
			b.CurrentSq++
			b.NextSq()
		} else {
			cp:=b.PieceAtSq(md.To)
			c:=piece.ColorOf(cp)
			if c==b.GetColorOfTurn() {
				// own piece
				b.CurrentPtr=md.NextVector
			} else if c==piece.NO_COLOR {
				// empty
				b.CurrentMove=TMove{b.CurrentSq,md.To,b.CurrentPiece,cp,0}
				b.CurrentPtr++
				return true
			} else {
				// capture
				b.CurrentMove=TMove{b.CurrentSq,md.To,b.CurrentPiece,cp,0}
				b.CurrentPtr=md.NextVector
				return true
			}
		}
	}
	return false
}

func Init() {
	InitMoveTable()
	InitNodeManager()
}

func InitNodeManager() {
	NodeManager.Nodes=make(map [TPosition]TNode)
}

func (nm *TNodeManager) DoesNodeExist(pos TPosition) bool {
	_,err := nm.Nodes[pos]
	return err
}

func (nm *TNodeManager) GetNode(pos TPosition) TNode {
	n,_ := nm.Nodes[pos]
	n.Visited++
	nm.Nodes[pos]=n
	return n
}

func (nm *TNodeManager) AddNode(pos TPosition, n TNode) {
	nm.Nodes[pos]=n
}

func InitMoveTable() {
	ptr:=0
	for _,pt := range ALL_PIECE_TYPES {
		for sq:=0; sq<square.BOARD_SIZE; sq++ {
			
			MoveTablePtrs[TMoveTableKey{square.TSquare(sq), piece.FromTypeAndColor(piece.TPieceType(pt),piece.BLACK)}]=ptr
			MoveTablePtrs[TMoveTableKey{square.TSquare(sq), piece.FromTypeAndColor(piece.TPieceType(pt),piece.WHITE)}]=ptr

			for df:=-2; df<=2; df++ {
				for dr:=-2; dr<=2; dr++ {
					vector_ok:=false
					dfabs:=math.Abs(float64(df))
					drabs:=math.Abs(float64(dr))
					prodabs:=dfabs*drabs
					sumabs:=dfabs+drabs
					if (pt&piece.IS_JUMPING)!=0 {
						vector_ok=vector_ok||(prodabs==2)
					}
					if (pt&piece.IS_STRAIGHT)!=0 {
						vector_ok=vector_ok||(sumabs==1)
					}
					if (pt&piece.IS_DIAGONAL)!=0 {
						vector_ok=vector_ok||(prodabs==1)
					}
					if vector_ok {
						ok:=true
						f:=int(square.FileOf(square.TSquare(sq)))
						r:=int(square.RankOf(square.TSquare(sq)))
						vector_start:=ptr
						for ok {
							f+=df
							r+=dr
							if square.FileRankOk(f,r) {
								tsq:=square.FromFileRank(square.TFile(f),square.TRank(r))
								MoveTable[ptr].To=tsq
								MoveTable[ptr].EndPiece=false
								ptr++
							} else {
								ok=false
							}
							if (pt&piece.IS_SLIDING)==0 {
								ok=false
							}
						}						
						for vector_next_ptr:=vector_start; vector_next_ptr<ptr; vector_next_ptr++ {
							MoveTable[vector_next_ptr].NextVector=ptr
						}
					}
				}
			}

			MoveTable[ptr].EndPiece=true
			ptr++
		}
	}
}

func InvTurnOf(t TTurn) TTurn {
	return TTurn(piece.InvColorOf(piece.TColor(t)))
}

func InvColorOfTurn(t TTurn) piece.TColor {
	return piece.TColor(InvTurnOf(t))
}

func InvTurnOfColor(c piece.TColor) TTurn {
	return TTurn(piece.InvColorOf(c))
}

func IndexOfColor(c piece.TColor) byte {
	return byte(c)>>1
}

func (pos *TPosition) SetFromFen(fen string) bool {
	var l=len(fen)

	if l<=0 {
		return false
	}

	var ok=true

	var ptr=0

	var feni=0

	for ; (feni<l) && ok ; feni++ {
		c:=fen[feni]

		if (c>='1') && (c<='8') {
			for fill:=0; fill<int(c-'0'); fill++ {
				if(ptr<square.BOARD_SIZE) {
					pos[ptr]=piece.NO_PIECE
					ptr++
				}
			}
		} else if(ptr<square.BOARD_SIZE) {
			p:=piece.FromFenChar(c)
			if p!=piece.NO_PIECE {
				pos[ptr]=byte(p)
				if piece.TypeOf(p)==piece.KING {
					pos[KINGSPOS_INDEX+IndexOfColor(piece.ColorOf(p))]=byte(ptr)
				}
				ptr++
			}
		}

		ok=(ptr<square.BOARD_SIZE)
	}

	if (ptr<square.BOARD_SIZE) || (feni>=(l-1)) {
		return false
	}

	pos[TURN_INDEX]=piece.WHITE

	feni++

	if fen[feni]=='b' {
		pos[TURN_INDEX]=piece.BLACK
	}

	pos[DEPTH_INDEX]=0

	return true
}

func TurnToChar(t TTurn) byte {
	if t==piece.BLACK {
		return 'b'
	}
	return 'w'
}

func (pos *TPosition) GetTurn() TTurn {
	return TTurn(pos[TURN_INDEX])
}

func (pos *TPosition) GetInvTurn() TTurn {
	return InvTurnOf(TTurn(pos[TURN_INDEX]))
}

func (pos *TPosition) GetDepth() int {
	return int(pos[DEPTH_INDEX])
}

func (pos *TPosition) SetTurn(t TTurn) {
	pos[TURN_INDEX]=byte(t)
}

func (pos *TPosition) SetDepth(d int) {
	pos[DEPTH_INDEX]=byte(d)
}

func (b *TBoard) GetTurn() TTurn {
	return b.Pos.GetTurn()
}

func (b *TBoard) GetInvTurn() TTurn {
	return b.Pos.GetInvTurn()
}

func (b *TBoard) GetDepth() int {
	return b.Pos.GetDepth()
}

func (b *TBoard) SetTurn(t TTurn) {
	b.Pos.SetTurn(t)
}

func (b *TBoard) SetDepth(d int) {
	b.Pos.SetDepth(d)
}

func (b *TBoard) GetColorOfTurn() piece.TColor {
	return piece.TColor(b.Pos.GetTurn())
}

func (b *TBoard) GetColorOfInvTurn() piece.TColor {
	return piece.TColor(b.Pos.GetInvTurn())
}

func (pos *TPosition) GetKingPos(c piece.TColor) square.TSquare {
	return square.TSquare(pos[KINGSPOS_INDEX+IndexOfColor(c)])
}

func (pos *TPosition) ToPrintable() string {
	var buff=""

	for i:=0; i<square.BOARD_SIZE; i++ {
		var fenchar=piece.ToFenChar(piece.TPiece(pos[i]))
		if fenchar==' ' {
			fenchar='.'
		}
		buff+=string(fenchar)
		if ((i+1) % square.BOARD_WIDTH) == 0 {
			buff+="\n"
		}
	}

	buff+=fmt.Sprintf("\nturn : %c , depth : %d , wkpos : %s , bkpos : %s\n",
		TurnToChar(pos.GetTurn()),pos.GetDepth(),
		square.ToAlgeb(pos.GetKingPos(piece.WHITE)),
		square.ToAlgeb(pos.GetKingPos(piece.BLACK)))

	return buff
}

func (b *TBoard) CalcMaterial() {
	b.Material[INDEX_OF_WHITE]=0
	b.Material[INDEX_OF_BLACK]=0
	for sq:=0; sq<square.BOARD_SIZE; sq++ {
		p:=b.PieceAtSq(square.TSquare(sq))
		c:=piece.ColorOf(p)
		t:=piece.TypeOf(p)
		b.Material[IndexOfColor(c)]+=PIECE_VALUES[t]
	}
}

func (b *TBoard) SetFromFen(fen string) bool {
	ok:=b.Pos.SetFromFen(fen)
	if !ok {
		return false
	}
	b.CalcMaterial()
	return true
}

func (b *TBoard) ToPrintable() string {
	return fmt.Sprintf("%s\nmaterial w : %d, b : %d | eval %d / w : %d, b : %d\n",
		b.Pos.ToPrintable(),b.Material[INDEX_OF_WHITE],b.Material[INDEX_OF_BLACK],
		b.Eval(),b.EvalCol(piece.WHITE),b.EvalCol(piece.BLACK))
}

var AlphaBetaEvals [1000]int

func (n *TNode) Eval() {
	for i:=0; i<len(n.Moves); i++ {
		AlphaBetaEvals[i]=INVALID_SCORE
	}
	for i, m:=range n.Moves {
		n.B.MakeMove(m)
		go n.AlphaBeta(i,n.B,0,EvalDepth,-INFINITE_SCORE,INFINITE_SCORE)
		//n.Moves[i].Eval=n.AlphaBeta(i,n.B,0,EvalDepth,-INFINITE_SCORE,INFINITE_SCORE)
		n.B.UnMakeMove(m)
	}
	notready:=true
	for notready && (!AbortMiniMax) {
		notready=false
		for i:=0; i<len(n.Moves); i++ {
			if AlphaBetaEvals[i]==INVALID_SCORE {
				notready=true
			}
		}
		time.Sleep(10 * time.Millisecond)
	}
	if AbortMiniMax {
		return
	}
	for i:=0; i<len(n.Moves); i++ {
		n.Moves[i].Eval=AlphaBetaEvals[i]+r.Intn(10)-20
	}
	n.Sort()
}

func (n *TNode) Len() int {
	return len(n.Moves)
}

func (n *TNode) Less(i,j int) bool {
	return n.Moves[i].Eval>n.Moves[j].Eval
}

func (n *TNode) Swap(i,j int) {
	n.Moves[i] , n.Moves[j] = n.Moves[j] , n.Moves[i]
}

func (n *TNode) Sort() {
	sort.Sort(n)
}

func (ownern *TNode) AddNodeRecursive(b TBoard, depth int, max_depth int, line TMoveList) bool {
	n:=b.CreateNode()
	l:=len(n.Moves)
	if l<=0 {
		return false
	}
	if depth>=max_depth {
		return false
	}
	var selm TMove
	ok:=false
	li:=l-1
	for i:=0; (i<l) && (!ok); i++ {
		selm=n.Moves[i]
		if i==li {
			ok=true
		} else {
			ok=r.Intn(100)>60
			if(math.Abs(float64(selm.Eval))>MATE_LIMIT) {
				ok=false
			}
		}
	}
	if ok {
		b.MakeMove(selm)
		if NodeManager.DoesNodeExist(b.Pos) {
			return n.AddNodeRecursive(b, depth+1, max_depth, append(line,selm))
		} else {
			//fmt.Printf("\nadding %s\n",append(line,selm).ToPrintable())
			b.CreateNode()			
			return true
		}
	}
	return false
}

func (n *TNode) AddNode(max_depth int) bool {
	b:=n.B
	success:=n.AddNodeRecursive(b, 0, max_depth, []TMove{})
	if !success {
		//fmt.Printf("\nadding node failed\n")
	}
	return success
}

func (ml TMoveList) ToPrintable() string {
	buff:=""
	for i,m := range ml {
		buff+=fmt.Sprintf("%s",m.ToAlgeb())
		if i<(len(ml)-1) {
			buff+=" "
		}
	}
	return fmt.Sprintf("line [ %s ]",buff)
}

func (ownern *TNode) MinimaxOutRecursive(b TBoard, depth int, max_depth int, positions []TPosition) int {
	if AbortMiniMax {
		return INVALID_SCORE
	}
	if !NodeManager.DoesNodeExist(b.Pos) {
		return INVALID_SCORE
	}
	if(depth>max_depth) {
		return INVALID_SCORE
	}
	n:=NodeManager.GetNode(b.Pos)
	max:=(-INFINITE_SCORE)
	if len(n.Moves)==0 {
		max=b.TerminalEval()
	} else {
		for i,m := range n.Moves {
		b.MakeMove(m)
		found:=false
		for i:=0; ((i<len(positions)) && (!found)); i++ {
			found=(positions[i]==b.Pos)
		}
		geteval:=0
		if !found {			
			geteval=ownern.MinimaxOutRecursive(b, depth+1, max_depth, append(positions,b.Pos))
		}
		b.UnMakeMove(m)
		eval:=m.Eval
		if(geteval!=INVALID_SCORE) {
			eval=geteval
			if math.Abs(float64(geteval))>MATE_LIMIT {
				if(geteval<0) {
					eval++
				} else {
					eval--
				}
			}
		}
		if(eval>max) {
			max=eval
		}
		n.Moves[i].Eval=eval
		}
		n.Sort()
	}
	return -max
}

func (n *TNode) MiniMaxOut(max_depth int) {
	b:=n.B
	n.MinimaxOutRecursive(b, 0, max_depth, []TPosition{})
}

func (n *TNode) AlphaBeta(store int,b TBoard, depth int, max_depth int, alpha int, beta int) int {

	if (depth>max_depth) || (AbortMiniMax) {
		return INVALID_SCORE
	}

	Nodes++
	
	b.InitMoveGen()

	hasmove:=b.NextLegalMove()

	if hasmove {
		for hasmove {
			b.MakeMove(b.CurrentMove)
			eval:=n.AlphaBeta(store,b,depth+1, max_depth, -beta, -alpha)
			if eval==INVALID_SCORE {
				eval=b.Eval()
			}
			b.UnMakeMove(b.CurrentMove)
			if math.Abs(float64(eval))>MATE_LIMIT {
				if(eval<0) {
					eval++
				} else {
					eval--
				}
			}
			if eval>alpha {
				alpha=eval
			}
			if alpha>beta {
				if depth==0 {
					AlphaBetaEvals[store]=-alpha
				}
				return -alpha
			}
			hasmove=b.NextLegalMove()
		}
		if depth==0 {
			AlphaBetaEvals[store]=-alpha
		}
		return -alpha
	} else {
		t:=-b.TerminalEval()
		if depth==0 {
			AlphaBetaEvals[store]=t
		}
		return t
	}
}