package gopcode

import (
	"fmt"
	"strings"
)

var (
	DefaultPcodeFormatter = newPcodeFormatter()
)

type prettyPrinter interface {
	formatPcodeOp(code PcodeOp) string
	formatVarNode(code VarNode) string
	formatOpCode(code OpCode) string
}

type pcodePrettyPrinter struct {
	opcodeHandlers map[OpCode]prettyPrinter
}

func (pp pcodePrettyPrinter) formatPcodeOp(pco PcodeOp) string {
	var formatter prettyPrinter

	if handler, ok := pp.opcodeHandlers[pco.Opcode]; ok {
		formatter = handler
	} else {
		formatter = pcodeDefaultPrettyPrinter{}
	}

	var formatted string

	if pco.Output != nil {
		formatted += fmt.Sprintf("%s = ", formatter.formatVarNode(*pco.Output))
	}

	formatted += formatter.formatPcodeOp(pco)

	return formatted
}

func (pp pcodePrettyPrinter) formatVarNode(vn VarNode) string {
	return "<not implemented>"
}

func (pp pcodePrettyPrinter) formatOpCode(op OpCode) string {
	return "<not implemented>"
}

type pcodeDefaultPrettyPrinter struct{}

func (pp pcodeDefaultPrettyPrinter) formatPcodeOp(pco PcodeOp) string {
	var formatted_inputs []string

	for _, input := range pco.Inputs {
		formatted_inputs = append(formatted_inputs, pp.formatVarNode(*input))
	}

	return fmt.Sprintf("%s %s", pp.formatOpCode(pco.Opcode), strings.Join(formatted_inputs, ", "))
}

func (pp pcodeDefaultPrettyPrinter) formatVarNode(vn VarNode) string {
	if vn.Space.Name == "const" {
		return fmt.Sprintf("0x%x", vn.Offset)
	} else if vn.Space.Name == "register" {
		return vn.GetRegisterName()
	}

	return fmt.Sprintf("%s[%x:%d]", vn.Space.Name, vn.Offset, vn.Size)
}

func (pp pcodeDefaultPrettyPrinter) formatOpCode(op OpCode) string {
	return op.String()
}

type pcodePrettyUnary struct {
	pcodeDefaultPrettyPrinter
	operator string
}

func (pp pcodePrettyUnary) formatPcodeOp(pco PcodeOp) string {
	return fmt.Sprintf("%s%s", pp.operator, pp.formatVarNode(*pco.Inputs[0]))
}

type pcodePrettyBinary struct {
	pcodeDefaultPrettyPrinter
	operator string
}

func (pp pcodePrettyBinary) formatPcodeOp(pco PcodeOp) string {
	return fmt.Sprintf("%s %s %s", pp.formatVarNode(*pco.Inputs[0]), pp.operator, pp.formatVarNode(*pco.Inputs[1]))
}

type pcodePrettyFunction struct {
	pcodeDefaultPrettyPrinter
	operator string
}

func (pp pcodePrettyFunction) formatPcodeOp(pco PcodeOp) string {
	var formatted_inputs []string
	for _, input := range pco.Inputs {
		formatted_inputs = append(formatted_inputs, pp.formatVarNode(*input))
	}
	return fmt.Sprintf("%s(%s)", pp.operator, strings.Join(formatted_inputs, ", "))
}

type pcodePrettySpecial struct {
	pcodeDefaultPrettyPrinter
	operator string
}

func (pp pcodePrettySpecial) format_BRANCH(pco PcodeOp) string {
	return fmt.Sprintf("goto %s", pp.formatVarNode(*pco.Inputs[0]))
}

func (pp pcodePrettySpecial) format_BRANCHIND(pco PcodeOp) string {
	return fmt.Sprintf("goto [%s]", pp.formatVarNode(*pco.Inputs[0]))
}

func (pp pcodePrettySpecial) format_CALL(pco PcodeOp) string {
	return fmt.Sprintf("call %s", pp.formatVarNode(*pco.Inputs[0]))
}

func (pp pcodePrettySpecial) format_CALLIND(pco PcodeOp) string {
	return fmt.Sprintf("call [%s]", pp.formatVarNode(*pco.Inputs[0]))
}

func (pp pcodePrettySpecial) format_CBRANCH(pco PcodeOp) string {
	return fmt.Sprintf("if (%s) goto %s", pp.formatVarNode(*pco.Inputs[1]), pp.formatVarNode(*pco.Inputs[0]))
}

func (pp pcodePrettySpecial) format_LOAD(pco PcodeOp) string {
	return fmt.Sprintf("*[%s]%s", pco.Inputs[0].GetSpaceFromConst().Name, pp.formatVarNode(*pco.Inputs[1]))
}

func (pp pcodePrettySpecial) format_RETURN(pco PcodeOp) string {
	return fmt.Sprintf("return %s", pp.formatVarNode(*pco.Inputs[0]))
}

func (pp pcodePrettySpecial) format_STORE(pco PcodeOp) string {
	return fmt.Sprintf("*[%s]%s = %s", pco.Inputs[0].GetSpaceFromConst().Name, pp.formatVarNode(*pco.Inputs[1]), pp.formatVarNode(*pco.Inputs[2]))
}

func (pp pcodePrettySpecial) formatPcodeOp(pco PcodeOp) string {
	switch pco.Opcode {
	case CPUI_BRANCH:
		return pp.format_BRANCH(pco)
	case CPUI_BRANCHIND:
		return pp.format_BRANCHIND(pco)
	case CPUI_CALL:
		return pp.format_CALL(pco)
	case CPUI_CALLIND:
		return pp.format_CALLIND(pco)
	case CPUI_CBRANCH:
		return pp.format_CBRANCH(pco)
	case CPUI_LOAD:
		return pp.format_LOAD(pco)
	case CPUI_RETURN:
		return pp.format_RETURN(pco)
	case CPUI_STORE:
		return pp.format_STORE(pco)
	default:
		return pp.pcodeDefaultPrettyPrinter.formatPcodeOp(pco)
	}
}

func newPcodeFormatter() prettyPrinter {
	return pcodePrettyPrinter{
		opcodeHandlers: map[OpCode]prettyPrinter{
			CPUI_BOOL_AND:          pcodePrettyBinary{operator: "&&"},
			CPUI_BOOL_NEGATE:       pcodePrettyUnary{operator: "!"},
			CPUI_BOOL_OR:           pcodePrettyBinary{operator: "||"},
			CPUI_BOOL_XOR:          pcodePrettyBinary{operator: "^^"},
			CPUI_BRANCH:            pcodePrettySpecial{},
			CPUI_BRANCHIND:         pcodePrettySpecial{},
			CPUI_CALL:              pcodePrettySpecial{},
			CPUI_CALLIND:           pcodePrettySpecial{},
			CPUI_CBRANCH:           pcodePrettySpecial{},
			CPUI_COPY:              pcodePrettyUnary{operator: ""},
			CPUI_CPOOLREF:          pcodePrettyFunction{operator: "cpool"},
			CPUI_FLOAT_ABS:         pcodePrettyFunction{operator: "abs"},
			CPUI_FLOAT_ADD:         pcodePrettyBinary{operator: "f+"},
			CPUI_FLOAT_CEIL:        pcodePrettyFunction{operator: "ceil"},
			CPUI_FLOAT_DIV:         pcodePrettyBinary{operator: "f/"},
			CPUI_FLOAT_EQUAL:       pcodePrettyBinary{operator: "f=="},
			CPUI_FLOAT_FLOAT2FLOAT: pcodePrettyFunction{operator: "float2float"},
			CPUI_FLOAT_FLOOR:       pcodePrettyFunction{operator: "floor"},
			CPUI_FLOAT_INT2FLOAT:   pcodePrettyFunction{operator: "int2float"},
			CPUI_FLOAT_LESS:        pcodePrettyBinary{operator: "f<"},
			CPUI_FLOAT_LESSEQUAL:   pcodePrettyBinary{operator: "f<="},
			CPUI_FLOAT_MULT:        pcodePrettyBinary{operator: "f*"},
			CPUI_FLOAT_NAN:         pcodePrettyFunction{operator: "nan"},
			CPUI_FLOAT_NEG:         pcodePrettyUnary{operator: "f-"},
			CPUI_FLOAT_NOTEQUAL:    pcodePrettyBinary{operator: "f!="},
			CPUI_FLOAT_ROUND:       pcodePrettyFunction{operator: "round"},
			CPUI_FLOAT_SQRT:        pcodePrettyFunction{operator: "sqrt"},
			CPUI_FLOAT_SUB:         pcodePrettyBinary{operator: "f-"},
			CPUI_FLOAT_TRUNC:       pcodePrettyFunction{operator: "trunc"},
			CPUI_INT_2COMP:         pcodePrettyUnary{operator: "-"},
			CPUI_INT_ADD:           pcodePrettyBinary{operator: "+"},
			CPUI_INT_AND:           pcodePrettyBinary{operator: "&"},
			CPUI_INT_CARRY:         pcodePrettyFunction{operator: "carry"},
			CPUI_INT_DIV:           pcodePrettyBinary{operator: "/"},
			CPUI_INT_EQUAL:         pcodePrettyBinary{operator: "=="},
			CPUI_INT_LEFT:          pcodePrettyBinary{operator: "<<"},
			CPUI_INT_LESS:          pcodePrettyBinary{operator: "<"},
			CPUI_INT_LESSEQUAL:     pcodePrettyBinary{operator: "<="},
			CPUI_INT_MULT:          pcodePrettyBinary{operator: "*"},
			CPUI_INT_NEGATE:        pcodePrettyUnary{operator: "~"},
			CPUI_INT_NOTEQUAL:      pcodePrettyBinary{operator: "!="},
			CPUI_INT_OR:            pcodePrettyBinary{operator: "|"},
			CPUI_INT_REM:           pcodePrettyBinary{operator: "%"},
			CPUI_INT_RIGHT:         pcodePrettyBinary{operator: ">>"},
			CPUI_INT_SBORROW:       pcodePrettyFunction{operator: "sborrow"},
			CPUI_INT_SCARRY:        pcodePrettyFunction{operator: "scarry"},
			CPUI_INT_SDIV:          pcodePrettyBinary{operator: "s/"},
			CPUI_INT_SEXT:          pcodePrettyFunction{operator: "sext"},
			CPUI_INT_SLESS:         pcodePrettyBinary{operator: "s<"},
			CPUI_INT_SLESSEQUAL:    pcodePrettyBinary{operator: "s<="},
			CPUI_INT_SREM:          pcodePrettyBinary{operator: "s%"},
			CPUI_INT_SRIGHT:        pcodePrettyBinary{operator: "s>>"},
			CPUI_INT_SUB:           pcodePrettyBinary{operator: "-"},
			CPUI_INT_XOR:           pcodePrettyBinary{operator: "^"},
			CPUI_INT_ZEXT:          pcodePrettyFunction{operator: "zext"},
			CPUI_LOAD:              pcodePrettySpecial{},
			CPUI_NEW:               pcodePrettyFunction{operator: "newobject"},
			CPUI_POPCOUNT:          pcodePrettyFunction{operator: "popcount"},
			CPUI_LZCOUNT:           pcodePrettyFunction{operator: "lzcount"},
			CPUI_RETURN:            pcodePrettySpecial{},
			CPUI_STORE:             pcodePrettySpecial{},
		},
	}
}
