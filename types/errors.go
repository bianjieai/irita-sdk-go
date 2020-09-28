package types

import (
	"fmt"

	"github.com/pkg/errors"
)

const (
	// RootCodespace is the codespace for all errors defined in irita
	RootCodespace = "sdk"

	OK                Code = 0
	Internal          Code = 1
	TxDecode          Code = 2
	InvalidSequence   Code = 3
	Unauthorized      Code = 4
	InsufficientFunds Code = 5
	UnknownRequest    Code = 6
	InvalidAddress    Code = 7
	InvalidPubkey     Code = 8
	UnknownAddress    Code = 9
	InvalidCoins      Code = 10
	OutOfGas          Code = 11
	MemoTooLarge      Code = 12
	InsufficientFee   Code = 13
	TooManySignatures Code = 14
	NoSignatures      Code = 15
	ErrJsonMarshal    Code = 16
	ErrJsonUnmarshal  Code = 17
	InvalidRequest    Code = 18
	TxInMempoolCache  Code = 19
	MempoolIsFull     Code = 20
	TxTooLarge        Code = 21
)

var (
	// errUnknown = register(RootCodespace, 111222, "unknown error")
	errInvalid = register(RootCodespace, 999999, "sdk check error")
)

func init() {
	_ = register(RootCodespace, OK, "success")
	_ = register(RootCodespace, Internal, "internal")
	_ = register(RootCodespace, TxDecode, "tx parse error")
	_ = register(RootCodespace, InvalidSequence, "invalid sequence")
	_ = register(RootCodespace, Unauthorized, "unauthorized")
	_ = register(RootCodespace, InsufficientFunds, "insufficient funds")
	_ = register(RootCodespace, UnknownRequest, "unknown request")
	_ = register(RootCodespace, InvalidAddress, "invalid address")
	_ = register(RootCodespace, InvalidPubkey, "invalid pubkey")
	_ = register(RootCodespace, UnknownAddress, "unknown address")
	_ = register(RootCodespace, InvalidCoins, "invalid coins")
	_ = register(RootCodespace, OutOfGas, "out of gas")
	_ = register(RootCodespace, MemoTooLarge, "memo too large")
	_ = register(RootCodespace, InsufficientFee, "insufficient fee")
	_ = register(RootCodespace, TooManySignatures, "maximum number of signatures exceeded")
	_ = register(RootCodespace, NoSignatures, "no signatures supplied")
	_ = register(RootCodespace, ErrJsonMarshal, "failed to marshal JSON bytes")
	_ = register(RootCodespace, ErrJsonUnmarshal, "failed to unmarshal JSON bytes")
	_ = register(RootCodespace, InvalidRequest, "invalid request")
	_ = register(RootCodespace, TxInMempoolCache, "tx already in mempool")
	_ = register(RootCodespace, MempoolIsFull, "mempool is full")
	_ = register(RootCodespace, TxTooLarge, "tx too large")
}

type Code uint32
type CodeV017 uint32

// Error represents a root error.
//
// Weave framework is using root error to categorize issues. Each instance
// created during the runtime should wrap one of the declared root errors. This
// allows error tests and returning all errors to the client in a safe manner.
//
// All popular root errors are declared in this package. If an extension has to
// declare a custom root error, always use register function to ensure
// error code uniqueness.
type Error interface {
	Error() string
	Code() uint32
	Codespace() string
}

// GetError is used to covert irita error to sdk error
func GetError(codespace string, code uint32, log ...string) Error {
	codeV1, ok := v17CodeMap[code]
	if !ok {
		codeV1 = InvalidRequest
	}
	return sdkError{
		codespace: codespace,
		code:      uint32(codeV1),
		desc:      log[0],
	}
}

// Wrap extends given error with an additional information.
//
// If the wrapped error does not provide ABCICode method (ie. stdlib errors),
// it will be labeled as internal error.
//
// If err is nil, this returns nil, avoiding the need for an if statement when
// wrapping a error returned at the end of a function
func Wrap(err error) Error {
	if err == nil {
		return nil
	}

	return sdkError{
		codespace: errInvalid.Codespace(),
		code:      errInvalid.Code(),
		desc:      err.Error(),
	}
}

func WrapWithMessage(err error, format string, args ...interface{}) Error {
	desc := fmt.Sprintf(format, args...)
	return Wrap(errors.WithMessage(err, desc))
}

// Wrapf extends given error with an additional information.
//
// This function works like Wrap function with additional functionality of
// formatting the input as specified.
func Wrapf(format string, args ...interface{}) Error {
	desc := fmt.Sprintf(format, args...)
	return Wrap(errors.New(desc))
}

type sdkError struct {
	codespace string
	code      uint32
	desc      string
}

func (e sdkError) Error() string {
	return e.desc
}

func (e sdkError) Code() uint32 {
	return e.code
}

func (e sdkError) Codespace() string {
	return e.codespace
}

// register returns an error instance that should be used as the base for
// creating error instances during runtime.
//
// Popular root errors are declared in this package, but extensions may want to
// declare custom codes. This function ensures that no error code is used
// twice. Attempt to reuse an error code results in panic.
//
// Use this function only during a program startup phase.
func register(codespace string, code Code, description string) Error {
	err := sdkError{
		codespace: codespace,
		code:      uint32(code),
		desc:      description,
	}
	setUsed(err)

	return err
}

// usedCodes is keeping track of used codes to ensure their uniqueness. No two
// error instances should share the same (codespace, code) tuple.
var usedCodes = map[string]Error{}

func errorID(codespace string, code uint32) string {
	return fmt.Sprintf("%s:%d", codespace, code)
}

func setUsed(err Error) {
	usedCodes[errorID(err.Codespace(), err.Code())] = err
}

// will remove from irita v1.0
var v17CodeMap = map[uint32]Code{
	0:  OK,
	1:  Internal,
	2:  TxDecode,
	3:  InvalidSequence,
	4:  Unauthorized,
	5:  InsufficientFunds,
	6:  UnknownRequest,
	7:  InvalidAddress,
	8:  InvalidPubkey,
	9:  UnknownAddress,
	10: InvalidCoins,
	11: InvalidCoins,
	12: OutOfGas,
	13: MemoTooLarge,
	14: InsufficientFee,
	15: UnknownRequest,
	16: TooManySignatures,
	17: OutOfGas,
	18: InvalidRequest,
	19: InvalidRequest,
	20: InvalidRequest,
	21: TxTooLarge,
	22: InvalidRequest,
	23: InvalidRequest,
}

func CatchPanic(fn func(errMsg string)) {
	if err := recover(); err != nil {
		var msg string
		switch e := err.(type) {
		case error:
			msg = e.Error()
		case string:
			msg = e
		}
		fn(msg)
	}
}
