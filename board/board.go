package board

import(
	"fmt"
	"math"
	"github.com/goracingkingsengine/gorke/piece"
	"github.com/goracingkingsengine/gorke/square"
)

const TURN_INDEX        = 64

const WHITE             = piece.WHITE
const BLACK             = piece.BLACK

type TPosition [65]byte
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
}

type TBoard struct {
	pos TPosition
	CurrentSq square.TSquare
	CurrentPiece piece.TPiece
	CurrentPtr int
	CurrentMove TMove
}

var ALL_PIECE_TYPES=[]piece.TPieceType{piece.KING,piece.QUEEN,piece.ROOK,piece.BISHOP,piece.KNIGHT}

var MoveTable [MOVE_TABLE_MAX_SIZE]TMoveDescriptor

var MoveTablePtrs=make(map[TMoveTableKey]int)

func (b *TBoard) GetMoveTablePtr(sq square.TSquare, p piece.TPiece) int {
	return MoveTablePtrs[TMoveTableKey{sq,p}]
}

func (b *TBoard) PieceAtSq(sq square.TSquare) piece.TPiece {
	return piece.TPiece(b.pos[byte(sq)])
}

func (b *TBoard) ColorOfSq(sq square.TSquare) piece.TColor {
	return piece.ColorOf(piece.TPiece(b.pos[byte(sq)]))
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

func (b *TBoard) ReportMoveGen() string {
	return fmt.Sprintf("current sq %d : current piece %c : current ptr %d current move %s",
		b.CurrentSq,piece.ToFenChar(b.CurrentPiece),b.CurrentPtr,b.CurrentMove.ToAlgeb())
}

func (b *TBoard) InitMoveGen() {
	b.CurrentSq=0
	b.NextSq()
}

func (b *TBoard) NextPseudoLegalMove() bool {
	for b.CurrentSq<square.BOARD_SIZE {
		md:=MoveTable[b.CurrentPtr]
		if md.EndPiece {
			b.CurrentSq++
			b.NextSq()
		} else {
			c:=b.ColorOfSq(md.To)
			if c==b.GetColorOfTurn() {
				// own piece
				b.CurrentPtr=md.NextVector
			} else if c==piece.NO_COLOR {
				// empty
				b.CurrentMove=TMove{b.CurrentSq,md.To}
				b.CurrentPtr++
				return true
			} else {
				// capture
				b.CurrentMove=TMove{b.CurrentSq,md.To}
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

func (b *TBoard) GetTurn() TTurn {
	return b.pos.GetTurn()
}

func (b *TBoard) GetColorOfTurn() piece.TColor {
	return piece.TColor(b.pos.GetTurn())
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

	buff+=fmt.Sprintf("\nturn : %c\n",TurnToChar(pos.GetTurn()))

	return buff
}

func (b *TBoard) SetFromFen(fen string) bool {
	return b.pos.SetFromFen(fen)
}

func (b *TBoard) ToPrintable() string {
	return b.pos.ToPrintable()
}