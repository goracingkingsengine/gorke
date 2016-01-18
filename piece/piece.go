package piece

const IS_SLIDING  = 1<<6
const IS_STRAIGHT = 1<<5
const IS_DIAGONAL = 1<<4
const IS_SINGLE   = 1<<3
const IS_JUMPING  = 1<<2
const IS_WHITE    = 1<<1
const IS_BLACK    = 1<<0

const IS_PIECE    = IS_WHITE|IS_BLACK

const TYPE_MASK   = IS_SLIDING|IS_STRAIGHT|IS_DIAGONAL|IS_SINGLE|IS_JUMPING

const WHITE       = IS_WHITE

const COLOR_MASK  = IS_PIECE

const BLACK       = IS_BLACK

const KING        = IS_STRAIGHT|IS_DIAGONAL|IS_SINGLE
const QUEEN       = IS_STRAIGHT|IS_DIAGONAL|IS_SLIDING
const ROOK        = IS_STRAIGHT|IS_SLIDING
const BISHOP      = IS_DIAGONAL|IS_SLIDING
const KNIGHT      = IS_JUMPING|IS_SINGLE

const NO_PIECE    = 0
const NO_COLOR    = 0

type TPiece byte
type TPieceType byte
type TColor byte

func TypeOf(p TPiece) TPieceType {
	return TPieceType(p&TYPE_MASK)
}

func ColorOf(p TPiece) TColor {
	return TColor(p&COLOR_MASK)
}

func FromTypeAndColor(t TPieceType, c TColor) TPiece {
	return TPiece(byte(t)|byte(c))
}

func ToFenChar(p TPiece) byte {
	var fenchar byte

	switch TypeOf(p) {
		case KING : fenchar='k'
		case QUEEN : fenchar='q'
		case ROOK : fenchar='r'
		case BISHOP : fenchar='b'
		case KNIGHT : fenchar='n'
		default : return ' '
	}

	if ColorOf(p)==WHITE {
		fenchar-='a'-'A'
	}

	return fenchar
}

func FromFenChar(fenchar byte) TPiece {
	var c TColor=WHITE

	if(fenchar>'a') {
		c=BLACK
		fenchar-='a'-'A'
	}

	var t TPieceType

	switch fenchar {
		case 'K' : t=KING
		case 'Q' : t=QUEEN
		case 'R' : t=ROOK
		case 'B' : t=BISHOP
		case 'N' : t=KNIGHT
		default : return NO_PIECE
	}

	return FromTypeAndColor(t,c)
}