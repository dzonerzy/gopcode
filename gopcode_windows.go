//go:build windows
// +build windows

package gopcode

/*
#cgo CFLAGS: -I${SRCDIR}/include
#cgo LDFLAGS: -L${SRCDIR}/lib/windows/x64 -lpcode -lstdc++
*/
// #include <windows.h>
// #include <pcode.h>
import (
	"C"
)
import (
	"unsafe"
)

type AddrSpace struct {
	NativeAddrSpacePtr *C.NativeAddrSpace
	Name               string
	RegisterName       string
	Flags              AddrSpaceFlags
	Highest            uint64
	PointerLowerBound  uint64
	PointerUpperBound  uint64
	Index              uint32
	AddressSize        uint32
	WordSize           uint32
}

type VarNode struct {
	Space  *AddrSpace
	Offset uint64
	Size   int32
}

func (v *VarNode) GetRegisterName() string {
	name := C.pcode_varcode_get_register_name(v.Space.NativeAddrSpacePtr, C.ulonglong(v.Offset), C.int32_t(v.Size))
	return C.GoString(name)
}

type Register struct {
	Node *VarNode
	Name string
}

type Context struct {
	_ctx       *C.PcodeContext
	LanguageID string
	_registers []*Register
}

func (c *Context) Destroy() {
	C.pcode_context_free(c._ctx)
}

func (c *Context) SetVariableDefault(name string, value uint32) {
	cname := C.CString(name)
	C.pcode_context_set_variable_default(c._ctx, cname, C.uint32_t(value))
	C.free(unsafe.Pointer(cname))
}

func (c *Context) GetAllRegisters() []*Register {
	if c._registers != nil {
		return c._registers
	}

	var regs []*Register
	reglist := C.pcode_context_get_all_registers(c._ctx)
	for i := 0; i < int(reglist.count); i++ {
		// get reg[i] considering it is a pointer to a C struct
		reg := (*C.RegisterInfoC)(unsafe.Pointer(uintptr(unsafe.Pointer(reglist.registers)) + uintptr(i)*unsafe.Sizeof(C.RegisterInfoC{})))
		regs = append(regs, &Register{
			Name: C.GoString(reg.name),
			Node: &VarNode{
				Space: &AddrSpace{
					Name:               C.GoString(reg.varnode.space.name),
					Index:              uint32(reg.varnode.space.index),
					AddressSize:        uint32(reg.varnode.space.address_size),
					WordSize:           uint32(reg.varnode.space.word_size),
					Flags:              AddrSpaceFlags(reg.varnode.space.flags),
					Highest:            uint64(reg.varnode.space.highest),
					PointerLowerBound:  uint64(reg.varnode.space.pointer_lower_bound),
					PointerUpperBound:  uint64(reg.varnode.space.pointer_upper_bound),
					NativeAddrSpacePtr: reg.varnode.space.n_space,
				},
				Offset: uint64(reg.varnode.offset),
				Size:   int32(reg.varnode.size),
			},
		})
	}

	if c._registers == nil {
		c._registers = regs
	}

	return regs
}

func (c *Context) GetRegisterName(space *AddrSpace, offset uint64, size int32) string {
	cname := C.pcode_context_get_register_name(c._ctx, space.NativeAddrSpacePtr, C.ulonglong(offset), C.int32_t(size))
	return C.GoString(cname)
}

func (c *Context) Disassemble(data []byte, baseAddress uint64, maxInstructions uint32) (*PcodeDisassembly, error) {
	return pcode_disassemble(c, data, baseAddress, maxInstructions)
}

func (c *Context) Translate(data []byte, baseAddress uint64, maxInstructions uint32, flags TranslateFlags) (*PcodeTranslation, error) {
	return pcode_translate(c, data, baseAddress, maxInstructions, flags)
}

func pcode_context_create(sla []byte) *Context {
	// convert sla to C.uchar *
	var csla = C.CString(string(sla))
	ctx := C.pcode_context_create((*C.uchar)(unsafe.Pointer(&sla[0])), C.size_t(len(sla)))
	C.free(unsafe.Pointer(csla))
	return &Context{_ctx: ctx}
}
