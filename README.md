## GoPCode
GoPCode is a library which provides a simple set of easy to use API to interact with ghidra PCode.

## Installation

To install GoPCode, simply run:
```bash
go get github.com/dzonerzy/gopcode
```

## Usage

Translation is the process of converting raw bytes into PCode instructions. The translation process is done by providing the raw bytes, the address of the first byte, the maximum number of instructions to translate, and the flags to use during translation. The flags are used to control the translation process, such as whether to stop at the first branch instruction or to translate the entire block of instructions.

```go
// create a new Context by providing the language ID
ctx, err := gopcode.NewContext("x86:LE:32:default")
if err != nil {
    panic(err)
}
// always remember to destroy object when done
defer ctx.Destroy()

// translate example Translate(data, address, max_instructions, flags)
pcode, err := ctx.Translate([]byte{0x55, 0x89, 0xe5}, 0x401000, 1024, 0)
if err != nil {
    panic(err)
}
defer pcode.Destroy()

// iterate over the translated opcodes
for _, op := range pcode.Ops {
    fmt.Printf("Opcode: %s\n", op.Opcode.String())
}
```

Disassembly is the process of converting PCode instructions into human-readable assembly instructions. The disassembly process is done by providing the PCode instructions and the address of the first byte. The disassembly process will return a list of assembly instructions.

```go
// create a new Context by providing the language ID
ctx, err := gopcode.NewContext("x86:LE:32:default")
if err != nil {
    panic(err)
}
// always remember to destroy object when done
defer ctx.Destroy()

// disassemble example Disassemble(data, address, max_instructions)
disas, err := ctx.Disassemble([]byte{0x55, 0x89, 0xe5}, 0x401000, 1024)
if err != nil {
    panic(err)
}
defer disas.Destroy()

// iterate over the disassembled instructions
for _, instr := range disas.Instructions {
    fmt.Printf("0x%x: %s %s\n", instr.Address, instr.Mnemonic, instr.Body)
}
```

## Supported architectures

GoPCode is based on the [pcode_c](https://github.com/dzonerzy/pcode_c) repository, so far GoPCode include precompiled binaries for the following architectures:

- windows x86 (64-bit)
- linux x86 (64-bit)

building for diffrent architectures is not tested and may not work.