package gopcode

// #include <pcode.h>
import (
	"C"
)
import (
	"fmt"
	"sync"
	"unsafe"
)

var varNodePool = sync.Pool{
	New: func() interface{} { return &VarNode{} },
}

var addrSpacePool = sync.Pool{
	New: func() interface{} { return &AddrSpace{} },
}

// sync.Map for caching AddrSpace instances by *C.char name
var addrSpaceCache = sync.Map{}

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

// getOrCreateAddrSpace retrieves an AddrSpace from cache or creates a new one if it doesn't exist
func getOrCreateAddrSpace(space *C.AddrSpaceC) *AddrSpace {
	// Use the space name pointer as the unique key for caching
	namePtr := space.name

	// Try to load AddrSpace from the cache
	if cached, ok := addrSpaceCache.Load(namePtr); ok {
		return cached.(*AddrSpace)
	}

	// If not found in cache, get a new AddrSpace from the pool
	addrSpace := addrSpacePool.Get().(*AddrSpace)

	// Initialize the AddrSpace fields only once for this unique name
	addrSpace.Name = C.GoString(namePtr)
	addrSpace.Index = uint32(space.index)
	addrSpace.AddressSize = uint32(space.address_size)
	addrSpace.WordSize = uint32(space.word_size)
	addrSpace.Flags = AddrSpaceFlags(space.flags)
	addrSpace.Highest = uint64(space.highest)
	addrSpace.PointerLowerBound = uint64(space.pointer_lower_bound)
	addrSpace.PointerUpperBound = uint64(space.pointer_upper_bound)
	addrSpace.NativeAddrSpacePtr = space.n_space

	// Store in cache for future reuse
	addrSpaceCache.Store(namePtr, addrSpace)

	return addrSpace
}

func (p *PcodeTranslation) Destroy() {
	for _, op := range p.Ops {
		if op.Output != nil {
			addrSpacePool.Put(op.Output.Space) // Return AddrSpace to pool
			varNodePool.Put(op.Output)         // Return VarNode to pool
		}
		for _, input := range op.Inputs {
			addrSpacePool.Put(input.Space) // Return AddrSpace to pool
			varNodePool.Put(input)         // Return VarNode to pool
		}
	}
	p.Ops = nil // Clear Ops slice

	C.pcode_translation_free(p._trans) // Free C-side resources if applicable
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

	pcodetrans := &PcodeTranslation{
		_formatter: DefaultPcodeFormatter,
		_trans:     trans,
		Ops:        make([]PcodeOp, 0, int(trans.num_ops)), // Pre-allocate Ops slice based on num_ops
	}

	for i := 0; i < int(trans.num_ops); i++ {
		op := (*C.PcodeOpC)(unsafe.Pointer(uintptr(unsafe.Pointer(trans.ops)) + uintptr(i)*unsafe.Sizeof(C.PcodeOpC{})))

		// Initialize PcodeOp with pooled VarNode and AddrSpace
		pcodeop := PcodeOp{
			Opcode: OpCode(op.opcode),
			Inputs: make([]*VarNode, 0, int(op.num_inputs)), // Pre-allocate Inputs slice
		}

		if op.output != nil {
			// Get VarNode from the pool
			varNode := varNodePool.Get().(*VarNode)

			// Retrieve AddrSpace from cache or create a new one
			varNode.Space = getOrCreateAddrSpace(op.output.space)
			varNode.Offset = uint64(op.output.offset)
			varNode.Size = int32(op.output.size)
			pcodeop.Output = varNode
		}

		for j := 0; j < int(op.num_inputs); j++ {
			inp := (*C.VarnodeDataC)(unsafe.Pointer(uintptr(unsafe.Pointer(op.inputs)) + uintptr(j)*unsafe.Sizeof(C.VarnodeDataC{})))

			// Get a VarNode from pool for each input
			inputNode := varNodePool.Get().(*VarNode)

			// Retrieve AddrSpace from cache or create a new one
			inputNode.Space = getOrCreateAddrSpace(inp.space)
			inputNode.Offset = uint64(inp.offset)
			inputNode.Size = int32(inp.size)
			pcodeop.Inputs = append(pcodeop.Inputs, inputNode)
		}

		pcodetrans.Ops = append(pcodetrans.Ops, pcodeop)
	}

	return pcodetrans, nil
}
