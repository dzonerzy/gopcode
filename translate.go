package gopcode

// #include <pcode.h>
import (
	"C"
)
import (
	"fmt"
	"unsafe"
)

type PcodeOp struct {
	Output *VarNode
	Inputs []*VarNode
	Opcode OpCode
}

type PcodeTranslation struct {
	_formatter prettyPrinter
	_trans     *C.PcodeTranslationC
	Ops        []PcodeOp
}

func (p *PcodeTranslation) Destroy() {
	C.pcode_translation_free(p._trans)
}

func (p *PcodeTranslation) Format(pco PcodeOp) string {
	return p._formatter.formatPcodeOp(pco)
}

// PcodeContext *ctx, const char *bytes, unsigned int num_bytes, uint64_t base_address, unsigned int max_instructions, uint32_t flags)
func pcode_translate(ctx *Context, dat []byte, baseAddress uint64, maxInstructions uint32, flags TranslateFlags) (*PcodeTranslation, error) {
	data := unsafe.Pointer(&dat[0])
	var trans *C.PcodeTranslationC = C.pcode_translate(ctx._ctx, (*C.char)(data), C.uint(len(dat)), C.ulonglong(baseAddress), C.uint(maxInstructions), C.uint(flags))

	if trans == nil {
		return nil, fmt.Errorf("translation failed")
	}

	var pcodetrans = &PcodeTranslation{
		_formatter: newPcodeFormatter(),
	}

	pcodetrans._trans = trans

	for i := 0; i < int(trans.num_ops); i++ {
		op := (*C.PcodeOpC)(unsafe.Pointer(uintptr(unsafe.Pointer(trans.ops)) + uintptr(i)*unsafe.Sizeof(C.PcodeOpC{})))

		var pcodeop PcodeOp
		pcodeop.Opcode = OpCode(op.opcode)
		if op.output != nil {
			pcodeop.Output = &VarNode{
				Offset: uint64(op.output.offset),
				Size:   int32(op.output.size),
				Space: &AddrSpace{
					Name:               C.GoString(op.output.space.name),
					Index:              uint32(op.output.space.index),
					AddressSize:        uint32(op.output.space.address_size),
					WordSize:           uint32(op.output.space.word_size),
					Flags:              AddrSpaceFlags(op.output.space.flags),
					Highest:            uint64(op.output.space.highest),
					PointerLowerBound:  uint64(op.output.space.pointer_lower_bound),
					PointerUpperBound:  uint64(op.output.space.pointer_upper_bound),
					NativeAddrSpacePtr: op.output.space.n_space,
				},
			}
		}

		for j := 0; j < int(op.num_inputs); j++ {

			inp := (*C.VarnodeDataC)(unsafe.Pointer(uintptr(unsafe.Pointer(op.inputs)) + uintptr(j)*unsafe.Sizeof(C.VarnodeDataC{})))

			pcodeop.Inputs = append(pcodeop.Inputs, &VarNode{
				Offset: uint64(inp.offset),
				Size:   int32(inp.size),
				Space: &AddrSpace{
					Name:               C.GoString(inp.space.name),
					Index:              uint32(inp.space.index),
					AddressSize:        uint32(inp.space.address_size),
					WordSize:           uint32(inp.space.word_size),
					Flags:              AddrSpaceFlags(inp.space.flags),
					Highest:            uint64(inp.space.highest),
					PointerLowerBound:  uint64(inp.space.pointer_lower_bound),
					PointerUpperBound:  uint64(inp.space.pointer_upper_bound),
					NativeAddrSpacePtr: inp.space.n_space,
				},
			})

		}

		pcodetrans.Ops = append(pcodetrans.Ops, pcodeop)

	}

	return pcodetrans, nil
}
