package util

type Code int

const (
	Success = iota
	ErrRegister
	ErrWriteDB
	ErrDataBinding
	ErrReadDB
	ErrUnmarshalObj
	ErrStartGame
	ErrNameOccupied
	ErrGameNotEnd
	ErrGameNotExist
	ErrGameOccupied
	ErrDuplicateOperation
	ErrRequestFormat
	ErrInvalidPlay
	ErrGameDataMissing
)