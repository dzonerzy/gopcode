package gopcode

import (
	"fmt"
	"strconv"
)

type AddrSpaceFlags int

const (
	BigEndian AddrSpaceFlags = 1 << iota
	Heritaged
	DoesDeadcode
	ProgramSpecific
	ReverseJustification
	FormalStackSpace
	Overlay
	OverlayBase
	Truncated
	HasPhysical
	IsOtherSpace
	HasNearPointers
)

type TranslateFlags int

const (
	BbTerminating TranslateFlags = 1 << iota
)

type OpCode int

const (
	CPUI_IMARK OpCode = iota

	CPUI_COPY  OpCode = 1
	CPUI_LOAD  OpCode = 2
	CPUI_STORE OpCode = 3 ///< Store at a pointer into a specified address space

	CPUI_BRANCH    OpCode = 4 ///< Always branch
	CPUI_CBRANCH   OpCode = 5 ///< Conditional branch
	CPUI_BRANCHIND OpCode = 6 ///< Indirect branch (jumptable)

	CPUI_CALL      OpCode = 7  ///< Call to an absolute address
	CPUI_CALLIND   OpCode = 8  ///< Call through an indirect address
	CPUI_CALLOTHER OpCode = 9  ///< User-defined operation
	CPUI_RETURN    OpCode = 10 ///< Return from subroutine

	// Integer/bit operations

	CPUI_INT_EQUAL      OpCode = 11 ///< Integer comparison, equality (==)
	CPUI_INT_NOTEQUAL   OpCode = 12 ///< Integer comparison, in-equality (!=)
	CPUI_INT_SLESS      OpCode = 13 ///< Integer comparison, signed less-than (<)
	CPUI_INT_SLESSEQUAL OpCode = 14 ///< Integer comparison, signed less-than-or-equal (<=)
	CPUI_INT_LESS       OpCode = 15 ///< Integer comparison, unsigned less-than (<)
	// This also indicates a borrow on unsigned substraction
	CPUI_INT_LESSEQUAL OpCode = 16 ///< Integer comparison, unsigned less-than-or-equal (<=)
	CPUI_INT_ZEXT      OpCode = 17 ///< Zero extension
	CPUI_INT_SEXT      OpCode = 18 ///< Sign extension
	CPUI_INT_ADD       OpCode = 19 ///< Addition, signed or unsigned (+)
	CPUI_INT_SUB       OpCode = 20 ///< Subtraction, signed or unsigned (-)
	CPUI_INT_CARRY     OpCode = 21 ///< Test for unsigned carry
	CPUI_INT_SCARRY    OpCode = 22 ///< Test for signed carry
	CPUI_INT_SBORROW   OpCode = 23 ///< Test for signed borrow
	CPUI_INT_2COMP     OpCode = 24 ///< Twos complement
	CPUI_INT_NEGATE    OpCode = 25 ///< Logical/bitwise negation (~)
	CPUI_INT_XOR       OpCode = 26 ///< Logical/bitwise exclusive-or (^)
	CPUI_INT_AND       OpCode = 27 ///< Logical/bitwise and (&)
	CPUI_INT_OR        OpCode = 28 ///< Logical/bitwise or (|)
	CPUI_INT_LEFT      OpCode = 29 ///< Left shift (<<)
	CPUI_INT_RIGHT     OpCode = 30 ///< Right shift, logical (>>)
	CPUI_INT_SRIGHT    OpCode = 31 ///< Right shift, arithmetic (>>)
	CPUI_INT_MULT      OpCode = 32 ///< Integer multiplication, signed and unsigned (*)
	CPUI_INT_DIV       OpCode = 33 ///< Integer division, unsigned (/)
	CPUI_INT_SDIV      OpCode = 34 ///< Integer division, signed (/)
	CPUI_INT_REM       OpCode = 35 ///< Remainder/modulo, unsigned (%)
	CPUI_INT_SREM      OpCode = 36 ///< Remainder/modulo, signed (%)

	CPUI_BOOL_NEGATE OpCode = 37 ///< Boolean negate (!)
	CPUI_BOOL_XOR    OpCode = 38 ///< Boolean exclusive-or (^^)
	CPUI_BOOL_AND    OpCode = 39 ///< Boolean and (&&)
	CPUI_BOOL_OR     OpCode = 40 ///< Boolean or (||)

	// Floating point operations

	CPUI_FLOAT_EQUAL     OpCode = 41 ///< Floating-point comparison, equality (==)
	CPUI_FLOAT_NOTEQUAL  OpCode = 42 ///< Floating-point comparison, in-equality (!=)
	CPUI_FLOAT_LESS      OpCode = 43 ///< Floating-point comparison, less-than (<)
	CPUI_FLOAT_LESSEQUAL OpCode = 44 ///< Floating-point comparison, less-than-or-equal (<=)
	// Slot 45 is currently unused
	CPUI_FLOAT_NAN OpCode = 46 ///< Not-a-number test (NaN)

	CPUI_FLOAT_ADD  OpCode = 47 ///< Floating-point addition (+)
	CPUI_FLOAT_DIV  OpCode = 48 ///< Floating-point division (/)
	CPUI_FLOAT_MULT OpCode = 49 ///< Floating-point multiplication (*)
	CPUI_FLOAT_SUB  OpCode = 50 ///< Floating-point subtraction (-)
	CPUI_FLOAT_NEG  OpCode = 51 ///< Floating-point negation (-)
	CPUI_FLOAT_ABS  OpCode = 52 ///< Floating-point absolute value (abs)
	CPUI_FLOAT_SQRT OpCode = 53 ///< Floating-point square root (sqrt)

	CPUI_FLOAT_INT2FLOAT   OpCode = 54 ///< Convert an integer to a floating-point
	CPUI_FLOAT_FLOAT2FLOAT OpCode = 55 ///< Convert between different floating-point sizes
	CPUI_FLOAT_TRUNC       OpCode = 56 ///< Round towards zero
	CPUI_FLOAT_CEIL        OpCode = 57 ///< Round towards +infinity
	CPUI_FLOAT_FLOOR       OpCode = 58 ///< Round towards -infinity
	CPUI_FLOAT_ROUND       OpCode = 59 ///< Round towards nearest

	// Internal opcodes for simplification. Not
	// typically generated in a direct translation.

	// Data-flow operations
	CPUI_MULTIEQUAL OpCode = 60 ///< Phi-node operator
	CPUI_INDIRECT   OpCode = 61 ///< Copy with an indirect effect
	CPUI_PIECE      OpCode = 62 ///< Concatenate
	CPUI_SUBPIECE   OpCode = 63 ///< Truncate

	CPUI_CAST      OpCode = 64 ///< Cast from one data-type to another
	CPUI_PTRADD    OpCode = 65 ///< Index into an array ([])
	CPUI_PTRSUB    OpCode = 66 ///< Drill down to a sub-field  (->)
	CPUI_SEGMENTOP OpCode = 67 ///< Look-up a \e segmented address
	CPUI_CPOOLREF  OpCode = 68 ///< Recover a value from the \e constant \e pool
	CPUI_NEW       OpCode = 69 ///< Allocate a new object (new)
	CPUI_INSERT    OpCode = 70 ///< Insert a bit-range
	CPUI_EXTRACT   OpCode = 71 ///< Extract a bit-range
	CPUI_POPCOUNT  OpCode = 72 ///< Count the 1-bits
	CPUI_LZCOUNT   OpCode = 73 ///< Count the leading 0-bits

	CPUI_MAX OpCode = 74
)

func (o OpCode) String() string {
	if name, ok := opCodeNames[o]; ok {
		return name
	}

	return fmt.Sprintf("unknown opcode %d", o)
}

var (
	opCodeNames = map[OpCode]string{
		CPUI_IMARK:             "IMARK",
		CPUI_COPY:              "COPY",
		CPUI_LOAD:              "LOAD",
		CPUI_STORE:             "STORE",
		CPUI_BRANCH:            "BRANCH",
		CPUI_CBRANCH:           "CBRANCH",
		CPUI_BRANCHIND:         "BRANCHIND",
		CPUI_CALL:              "CALL",
		CPUI_CALLIND:           "CALLIND",
		CPUI_CALLOTHER:         "CALLOTHER",
		CPUI_RETURN:            "RETURN",
		CPUI_INT_EQUAL:         "INT_EQUAL",
		CPUI_INT_NOTEQUAL:      "INT_NOTEQUAL",
		CPUI_INT_SLESS:         "INT_SLESS",
		CPUI_INT_SLESSEQUAL:    "INT_SLESSEQUAL",
		CPUI_INT_LESS:          "INT_LESS",
		CPUI_INT_LESSEQUAL:     "INT_LESSEQUAL",
		CPUI_INT_ZEXT:          "INT_ZEXT",
		CPUI_INT_SEXT:          "INT_SEXT",
		CPUI_INT_ADD:           "INT_ADD",
		CPUI_INT_SUB:           "INT_SUB",
		CPUI_INT_CARRY:         "INT_CARRY",
		CPUI_INT_SCARRY:        "INT_SCARRY",
		CPUI_INT_SBORROW:       "INT_SBORROW",
		CPUI_INT_2COMP:         "INT_2COMP",
		CPUI_INT_NEGATE:        "INT_NEGATE",
		CPUI_INT_XOR:           "INT_XOR",
		CPUI_INT_AND:           "INT_AND",
		CPUI_INT_OR:            "INT_OR",
		CPUI_INT_LEFT:          "INT_LEFT",
		CPUI_INT_RIGHT:         "INT_RIGHT",
		CPUI_INT_SRIGHT:        "INT_SRIGHT",
		CPUI_INT_MULT:          "INT_MULT",
		CPUI_INT_DIV:           "INT_DIV",
		CPUI_INT_SDIV:          "INT_SDIV",
		CPUI_INT_REM:           "INT_REM",
		CPUI_INT_SREM:          "INT_SREM",
		CPUI_BOOL_NEGATE:       "BOOL_NEGATE",
		CPUI_BOOL_XOR:          "BOOL_XOR",
		CPUI_BOOL_AND:          "BOOL_AND",
		CPUI_BOOL_OR:           "BOOL_OR",
		CPUI_FLOAT_EQUAL:       "FLOAT_EQUAL",
		CPUI_FLOAT_NOTEQUAL:    "FLOAT_NOTEQUAL",
		CPUI_FLOAT_LESS:        "FLOAT_LESS",
		CPUI_FLOAT_LESSEQUAL:   "FLOAT_LESSEQUAL",
		CPUI_FLOAT_NAN:         "FLOAT_NAN",
		CPUI_FLOAT_ADD:         "FLOAT_ADD",
		CPUI_FLOAT_DIV:         "FLOAT_DIV",
		CPUI_FLOAT_MULT:        "FLOAT_MULT",
		CPUI_FLOAT_SUB:         "FLOAT_SUB",
		CPUI_FLOAT_NEG:         "FLOAT_NEG",
		CPUI_FLOAT_ABS:         "FLOAT_ABS",
		CPUI_FLOAT_SQRT:        "FLOAT_SQRT",
		CPUI_FLOAT_INT2FLOAT:   "FLOAT_INT2FLOAT",
		CPUI_FLOAT_FLOAT2FLOAT: "FLOAT_FLOAT2FLOAT",
		CPUI_FLOAT_TRUNC:       "FLOAT_TRUNC",
		CPUI_FLOAT_CEIL:        "FLOAT_CEIL",
		CPUI_FLOAT_FLOOR:       "FLOAT_FLOOR",
		CPUI_FLOAT_ROUND:       "FLOAT_ROUND",
		CPUI_MULTIEQUAL:        "MULTIEQUAL",
		CPUI_INDIRECT:          "INDIRECT",
		CPUI_PIECE:             "PIECE",
		CPUI_SUBPIECE:          "SUBPIECE",
		CPUI_CAST:              "CAST",
		CPUI_PTRADD:            "PTRADD",
		CPUI_PTRSUB:            "PTRSUB",
		CPUI_SEGMENTOP:         "SEGMENTOP",
		CPUI_CPOOLREF:          "CPOOLREF",
		CPUI_NEW:               "NEW",
		CPUI_INSERT:            "INSERT",
		CPUI_EXTRACT:           "EXTRACT",
		CPUI_POPCOUNT:          "POPCOUNT",
		CPUI_LZCOUNT:           "LZCOUNT",
	}
)

func NewContext(LanguageID string) (*Context, error) {
	for _, al := range ArchLanguages {
		if al.LanguageID == LanguageID {
			ctx := pcode_context_create(al.Sla)

			for _, set := range al.ProcessorSpecs.ContextData.CtxSet.Set {
				v, _ := strconv.ParseUint(set.Val, 10, 32)
				ctx.SetVariableDefault(set.Name, uint32(v))
			}

			// populate registers
			_ = ctx.GetAllRegisters()

			return ctx, nil
		}
	}

	return nil, fmt.Errorf("language %s not found", LanguageID)
}
