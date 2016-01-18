package board

import(
	"fmt"
	"github.com/goracingkingsengine/gorke/piece"
	"github.com/goracingkingsengine/gorke/square"
)

const TURN_INDEX        = 64

const WHITE             = piece.WHITE
const BLACK             = piece.BLACK

type TPosition [65]byte
type TTurn piece.TColor

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