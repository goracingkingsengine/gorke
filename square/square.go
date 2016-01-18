package square

import "fmt"

const BOARD_WIDTH         = 8
const BOARD_WIDTHL        = BOARD_WIDTH-1

const BOARD_HEIGHT        = BOARD_WIDTH
const BOARD_HEIGHTL       = BOARD_HEIGHT-1

const BOARD_SIZE          = BOARD_WIDTH * BOARD_HEIGHT
const BOARD_SIZEL         = BOARD_SIZE-1

const NO_SQUARE           = BOARD_SIZE

const BOARD_SHIFT         = 3

const FILE_MASK           = BOARD_WIDTHL
const RANK_MASK           = FILE_MASK << BOARD_SHIFT

type TSquare byte
type TFile byte
type TRank byte

func FileOf(sq TSquare) TFile {
	return TFile(byte(sq)&FILE_MASK)
}

func RankOf(sq TSquare) TRank {
	return TRank((byte(sq)&RANK_MASK)>>BOARD_SHIFT)
}

func FileOk(f int) bool {
	return (f>=0) && (f<BOARD_WIDTH)
}

func RankOk(r int) bool {
	return (r>=0) && (r<BOARD_HEIGHT)
}

func FileRankOk(f int, r int) bool {
	return FileOk(f) && RankOk(r)
}

func AlgebFileToFile(af byte) (TFile, bool) {
	if (af<'a') || (af>('a'+BOARD_WIDTHL)) {
		return 0, true
	}
	return TFile(af-'a'), false
}

func FileToAlgebFile(f TFile) byte {
	return byte('a'+byte(f))
}

func AlgebRankToRank(ar byte) (TRank, bool) {
	if (ar<'1') || (ar>('1'+BOARD_HEIGHTL)) {
		return 0, true
	}
	return TRank(BOARD_HEIGHTL-(byte(ar)-'1')), false
}

func RankToAlgebRank(r TRank) byte {
	return byte('1'+(BOARD_HEIGHTL-byte(r)))
}

func FromFileRank(f TFile, r TRank) TSquare {
	return TSquare(byte(f)|(byte(r)<<BOARD_SHIFT))
}

func FromAlgeb(algeb string) TSquare {
	if len(algeb)<2 {
		return NO_SQUARE
	}
	f,err:=AlgebFileToFile(algeb[0])
	if(err) {
		return NO_SQUARE
	}
	r,err:=AlgebRankToRank(algeb[1])
	if(err) {
		return NO_SQUARE
	}
	return FromFileRank(f,r)
}

func ToAlgeb(sq TSquare) string {
	if (sq<0) || (sq>BOARD_SIZEL) {
		return "-"
	}
	return fmt.Sprintf("%c%c",FileToAlgebFile(FileOf(sq)),RankToAlgebRank(RankOf(sq)))
}