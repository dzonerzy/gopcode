package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"regexp"

	"github.com/dzonerzy/gopcode"
)

var (
	lid       = flag.String("lid", "x86:le:32:default", "Language ID")
	translate = flag.Bool("translate", false, "Translate")
	disasm    = flag.Bool("disasm", false, "Disassemble")
	data      = flag.String("data", "90 90 c3", "Data to disassemble/translate")
)

func isValidHexString(s string) bool {
	// valid hex string is 90 90 c3
	return regexp.MustCompile(`^([0-9a-fA-F]{2}\s?)+$`).MatchString(s)
}

func hexStringToBytes(s string) []byte {
	s = regexp.MustCompile(`\s`).ReplaceAllString(s, "")
	b, err := hex.DecodeString(s)
	if err != nil {
		log.Fatalf("failed to decode hex string: %v", err)
	}

	return b
}

func main() {
	flag.Parse()

	ctx, err := gopcode.NewContext(*lid)
	if err != nil {
		log.Fatalf("failed to create context: %v", err)
	}
	defer ctx.Destroy()

	// check if both disasm and translate are false
	if !*disasm && !*translate {
		log.Fatalf("disasm and translate cannot be both false")
	}

	// check if both are true
	if *disasm && *translate {
		log.Fatalf("only one of disasm and translate can be chosen")
	}

	// check if data is populated
	if *data == "" {
		log.Fatalf("data cannot be empty")
	}

	// check if is valid hex string
	if !isValidHexString(*data) {
		log.Fatalf("data is not a valid hex string")
	}

	// convert hex string to bytes
	b := hexStringToBytes(*data)

	if *disasm {
		disas, err := ctx.Disassemble(b, 0x401000, 1024)
		if err != nil {
			log.Fatalf("failed to disassemble: %v", err)
		}
		defer disas.Destroy()

		for _, instr := range disas.Instructions {
			fmt.Printf("0x%x: %s %s\n", instr.Address, instr.Mnemonic, instr.Body)
		}
	} else if *translate {
		trans, err := ctx.Translate(b, 0x401000, 1024, 0)
		if err != nil {
			log.Fatalf("failed to translate: %v", err)
		}
		defer trans.Destroy()

		for _, op := range trans.Ops {
			fmt.Println(trans.Format(op))
		}
	}

}
