package piece

const IS_PIECE    = 1<<6
const IS_SLIDING  = 1<<5
const IS_STRAIGHT = 1<<4
const IS_DIAGONAL = 1<<3
const IS_SINGLE   = 1<<2
const IS_JUMPING  = 1<<1

const TYPE_MASK   = IS_PIECE|IS_SLIDING|IS_STRAIGHT|IS_DIAGONAL|IS_SINGLE|IS_JUMPING

const WHITE       = 1<<0

const COLOR_MASK  = WHITE

const BLACK       = 0

const KING        = IS_PIECE|IS_STRAIGHT|IS_DIAGONAL|IS_SINGLE
const QUEEN       = IS_PIECE|IS_STRAIGHT|IS_DIAGONAL|IS_SLIDING
const ROOK        = IS_PIECE|IS_STRAIGHT|IS_SLIDING
const BISHOP      = IS_PIECE|IS_DIAGONAL|IS_SLIDING
const KNIGHT      = IS_PIECE|IS_JUMPING|IS_SINGLE

const NO_PIECE    = 0

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