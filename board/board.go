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

var ALL_PIECE_TYPES=[]piece.TPieceType{piece.KING,piece.QUEEN,piece.ROOK,piece.BISHOP,piece.KNIGHT}

var MoveTable [MOVE_TABLE_MAX_SIZE]TMoveDescriptor

var MoveTablePtrs=make(map[TMoveTableKey]int)

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

type TBoard struct {
	pos TPosition
	sq square.TSquare
	p piece.TPiece
}

func (b *TBoard) SetFromFen(fen string) bool {
	return b.pos.SetFromFen(fen)
}

func (b *TBoard) ToPrintable() string {
	return b.pos.ToPrintable()
}