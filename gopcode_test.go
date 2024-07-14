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

	regs := ctx.GetAllRegisters()
	name := ctx.GetRegisterName(regs[0].Node.Space, regs[0].Node.Offset, regs[0].Node.Size)

	if name != regs[0].Name {
		t.Fatalf("expected %s, got %s", regs[0].Name, name)
	}

	var data = []byte{0x90, 0x90, 0xc3}

	// disassemble
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

	// translate
	trans, err := ctx.Translate(data, 0x1000, 10, gopcode.BbTerminating)
	if err != nil {
		t.Fatal(err)
	}
	defer trans.Destroy()

	for _, op := range trans.Ops {
		fmt.Printf("Opcode: %s\n", op.Opcode.String())
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
