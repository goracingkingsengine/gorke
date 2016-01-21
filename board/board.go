package board

import(
	"fmt"
	"math"
	"github.com/goracingkingsengine/gorke/piece"
	"github.com/goracingkingsengine/gorke/square"
)

const WHITE             = piece.WHITE
const BLACK             = piece.BLACK

//////////////////////////////////////////

const TURN_INDEX        = 64
const DEPTH_INDEX       = 65

type TPosition [square.BOARD_SIZE+2]byte

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
}

type TNode struct {
	Pos TPosition
	Moves []TMove
}

var ALL_PIECE_TYPES=[]piece.TPieceType{piece.KING,piece.QUEEN,piece.ROOK,piece.BISHOP,piece.KNIGHT}

var MoveTable [MOVE_TABLE_MAX_SIZE]TMoveDescriptor

var MoveTablePtrs=make(map[TMoveTableKey]int)

func (b *TBoard) CreateNode() TNode {
	var node TNode
	node.Pos=b.Pos
	b.InitMoveGen()
	node.Moves=[]TMove{}
	for b.NextLegalMove() {
		node.Moves=append(node.Moves,b.CurrentMove)
	}
	return node
}

func (b *TBoard) SetSq(sq square.TSquare,p piece.TPiece) {
	b.Pos[byte(sq)]=byte(p)
}

func (b *TBoard) MakeMove(m TMove) {
	b.SetSq(m.From,piece.NO_PIECE)
	b.SetSq(m.To,m.Piece)
	b.SetTurn(InvTurnOf(b.GetTurn()))
	b.SetDepth(b.GetDepth()+1)
}

func (b *TBoard) UnMakeMove(m TMove) {
	b.SetSq(m.From,m.Piece)
	b.SetSq(m.To,m.CapPiece)
	b.SetTurn(InvTurnOf(b.GetTurn()))
	b.SetDepth(b.GetDepth()-1)
}

func (n *TNode) ToPrintable() string {
	var buff="----- Node: -----\n\n"
	buff+=fmt.Sprintf("%s\nMoves (%d):\n\n",n.Pos.ToPrintable(),len(n.Moves))
	for i := range n.Moves {
		buff+=fmt.Sprintf("%3d %s\n",i,n.Moves[i].ToPrintable())
	}
	return fmt.Sprintf("%s\n----- End Node -----",buff)
}

func (b *TBoard) GetMoveTablePtr(sq square.TSquare, p piece.TPiece) int {
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
			b.CurrentPtr=b.GetMoveTablePtr(b.CurrentSq,b.CurrentPiece)
			return true
		}
		b.CurrentSq++
	}
	return false
}

func (m *TMove) ToAlgeb() string {
	return fmt.Sprintf("%s%s",square.ToAlgeb(m.From),square.ToAlgeb(m.To))
}

func (m *TMove) ToPrintable() string {
	return fmt.Sprintf("%s (%c) %d",m.ToAlgeb(),piece.ToFenChar(m.CapPiece),m.Eval)
}

func (b *TBoard) ReportMoveGen() string {
	return fmt.Sprintf("current sq %d : current piece %c : current ptr %d current move %s",
		b.CurrentSq,piece.ToFenChar(b.CurrentPiece),b.CurrentPtr,b.CurrentMove.ToAlgeb())
}

func (b *TBoard) InitMoveGen() {
	b.CurrentSq=0
	b.NextSq()
}

func (b *TBoard) NextLegalMove() bool {
	return b.NextPseudoLegalMove()
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
				ptr++
			}
		}

		ok=(ptr<square.BOARD_SIZE)
	}

	if (ptr<square.BOARD_SIZE) || (feni>=(l-1)) {
		return false
	}

	pos[TURN_INDEX]=WHITE

	feni++

	if fen[feni]=='b' {
		pos[TURN_INDEX]=piece.BLACK
	}

	pos[DEPTH_INDEX]=0

	return true
}

func TurnToChar(t TTurn) byte {
	if t==BLACK {
		return 'b'
	}
	return 'w'
}

func (pos *TPosition) GetTurn() TTurn {
	return TTurn(pos[TURN_INDEX])
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

	buff+=fmt.Sprintf("\nturn : %c , depth : %d\n",TurnToChar(pos.GetTurn()),pos.GetDepth())

	return buff
}

func (b *TBoard) SetFromFen(fen string) bool {
	return b.Pos.SetFromFen(fen)
}

func (b *TBoard) ToPrintable() string {
	return b.Pos.ToPrintable()
}