package gopcode

// #include <pcode.h>
import (
	"C"
)
import (
	"fmt"
	"unsafe"
)

type DisassemblyInstruction struct {
	Mnemonic string
	Body     string
	Address  uint64
	Length   uint64
}

type PcodeDisassembly struct {
	_disas       *C.PcodeDisassemblyC
	Instructions []DisassemblyInstruction
}

func (p *PcodeDisassembly) Destroy() {
	C.pcode_disassembly_free(p._disas)
}

// PcodeDisassemblyC *pcode_disassemble(PcodeContext *ctx, const char *bytes, unsigned int num_bytes, uint64_t address, unsigned int max_instructions);
func pcode_disassemble(ctx *Context, dat []byte, baseAddress uint64, maxInstructions uint32) (*PcodeDisassembly, error) {
	data := unsafe.Pointer(&dat[0])
	var disas *C.PcodeDisassemblyC = C.pcode_disassemble(ctx._ctx, (*C.char)(data), C.uint(len(dat)), C.ulonglong(baseAddress), C.uint(maxInstructions))

	if disas == nil {
		return nil, fmt.Errorf("disassembly failed")
	}

	var pcodeDis = &PcodeDisassembly{}
	pcodeDis._disas = disas

	for i := 0; i < int(disas.num_instructions); i++ {
		instr := (*C.DisassemblyInstructionC)(unsafe.Pointer(uintptr(unsafe.Pointer(disas.instructions)) + uintptr(i)*unsafe.Sizeof(C.DisassemblyInstructionC{})))

		var disInstr DisassemblyInstruction

		disInstr.Mnemonic = C.GoString(instr.mnemonic)
		disInstr.Body = C.GoString(instr.body)
		disInstr.Address = uint64(instr.address)
		disInstr.Length = uint64(instr.length)

		pcodeDis.Instructions = append(pcodeDis.Instructions, disInstr)
	}

	return pcodeDis, nil
}
