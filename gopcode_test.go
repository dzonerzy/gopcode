package gopcode_test

import (
	"fmt"
	"testing"

	"github.com/dzonerzy/gopcode"
)

func TestContext(t *testing.T) {
	ctx, err := gopcode.NewContext("x86:le:32:default")
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Destroy()
}

func TestGetAllRegisters(t *testing.T) {
	ctx, err := gopcode.NewContext("x86:le:32:default")
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Destroy()

	regs := ctx.GetAllRegisters()
	if len(regs) == 0 {
		t.Fatal("no registers found")
	}

	name := ctx.GetRegisterName(regs[0].Node.Space, regs[0].Node.Offset, regs[0].Node.Size)

	if name != regs[0].Name {
		t.Fatalf("expected %s, got %s", regs[0].Name, name)
	}
}

func TestDisassemble(t *testing.T) {
	ctx, err := gopcode.NewContext("x86:le:32:default")
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Destroy()

	data := []byte{0x90, 0x90, 0xc3}

	disas, err := ctx.Disassemble(data, 0x1000, 10)
	if err != nil {
		t.Fatal(err)
	}

	defer disas.Destroy()

	expected := []string{"NOP", "NOP", "RET"}

	for i, instr := range disas.Instructions {
		if instr.Mnemonic != expected[i] {
			t.Fatalf("expected %s, got %s", expected[i], instr.Mnemonic)
		}
	}
}

func TestTranslate(t *testing.T) {
	ctx, err := gopcode.NewContext("x86:le:32:default")
	if err != nil {
		t.Fatal(err)
	}
	defer ctx.Destroy()

	data := []byte{
		0x55,       // push ebp
		0x8b, 0xec, // mov ebp, esp
		0x83, 0xec, 0x08, // sub esp, 0x8
		0x90, // nop
		0x90, // nop
		0xc9, // leave
		0xc3, // ret
	}

	trans, err := ctx.Translate(data, 0x1000, 1024, gopcode.BbTerminating)
	if err != nil {
		t.Fatal(err)
	}

	defer trans.Destroy()

	expected := []gopcode.OpCode{
		gopcode.CPUI_IMARK,
		gopcode.CPUI_COPY,
		gopcode.CPUI_INT_SUB,
		gopcode.CPUI_STORE,
		gopcode.CPUI_IMARK,
		gopcode.CPUI_COPY,
		gopcode.CPUI_IMARK,
		gopcode.CPUI_INT_LESS,
		gopcode.CPUI_INT_SBORROW,
		gopcode.CPUI_INT_SUB,
		gopcode.CPUI_INT_SLESS,
		gopcode.CPUI_INT_EQUAL,
		gopcode.CPUI_INT_AND,
		gopcode.CPUI_POPCOUNT,
		gopcode.CPUI_INT_AND,
		gopcode.CPUI_INT_EQUAL,
		gopcode.CPUI_IMARK,
		gopcode.CPUI_IMARK,
		gopcode.CPUI_IMARK,
		gopcode.CPUI_COPY,
		gopcode.CPUI_LOAD,
		gopcode.CPUI_INT_ADD,
		gopcode.CPUI_IMARK,
		gopcode.CPUI_LOAD,
		gopcode.CPUI_INT_ADD,
		gopcode.CPUI_RETURN,
	}

	for i, bb := range trans.Ops {
		if bb.Opcode != expected[i] {
			t.Fatalf("expected %s, got %s at %d", expected[i], bb.Opcode, i)
		}
	}
}

func TestListArchitectures(t *testing.T) {
	if len(gopcode.ArchLanguages) == 0 {
		t.Fatal("no architecture languages found")
	}

	for _, al := range gopcode.ArchLanguages {
		fmt.Printf("Language: %s - %s\n", al.Description, al.LanguageID)
	}
}

func BenchmarkTranslate(b *testing.B) {
	ctx, err := gopcode.NewContext("x86:le:32:default")
	if err != nil {
		b.Fatal(err)
	}
	defer ctx.Destroy()

	data := []byte{0x90, 0x90, 0xc3}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		t, err := ctx.Translate(data, 0x1000, 10, gopcode.BbTerminating)
		if err != nil {
			b.Fatal(err)
		}
		t.Destroy()
	}
}

func BenchmarkDisassemble(b *testing.B) {
	ctx, err := gopcode.NewContext("x86:le:32:default")
	if err != nil {
		b.Fatal(err)
	}
	defer ctx.Destroy()

	data := []byte{0x90, 0x90, 0xc3}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		d, err := ctx.Disassemble(data, 0x1000, 10)
		if err != nil {
			b.Fatal(err)
		}
		d.Destroy()
	}
}
