package gopcode

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"os"
	"strings"
)

type Set struct {
	Name string `xml:"name,attr"`
	Val  string `xml:"val,attr"`
}

type ContextSet struct {
	Space string `xml:"space,attr"`
	Set   []Set  `xml:"set"`
}

type TrackedSet struct {
	Space string `xml:"space,attr"`
	Set   []Set  `xml:"set"`
}

type ContextData struct {
	CtxSet   ContextSet `xml:"context_set"`
	TrackSet TrackedSet `xml:"tracked_set"`
}

type ProcessorSpec struct {
	ContextData ContextData `xml:"context_data"`
}

type ArchitectureLanguage struct {
	Description    string
	LanguageID     string
	ProcessorSpecs ProcessorSpec
	Sla            []byte
}

var (
	ArchLanguages []ArchitectureLanguage
)

type languageDef struct {
	Processor   string `xml:"processor,attr"`
	Endian      string `xml:"endian,attr"`
	Size        string `xml:"size,attr"`
	Variant     string `xml:"variant,attr"`
	Version     string `xml:"version,attr"`
	SLAFile     string `xml:"slafile,attr"`
	PSpec       string `xml:"processorspec,attr"`
	ManualIdx   string `xml:"manualindexfile,attr"`
	ID          string `xml:"id,attr"`
	Description string `xml:"description"`
}

type archLanguages struct {
	Langs []languageDef `xml:"language"`
}

func init() {
	archs, err := ProcessorsFS.ReadDir("processors")
	if err != nil {
		panic("could not read processors directory")
	}

	for _, arch := range archs {
		processArchitecture(arch.Name())
	}
}

func processArchitecture(archName string) {
	files, err := ProcessorsFS.ReadDir(fmt.Sprintf("processors/%s/data/languages", archName))
	if err != nil {
		panic(fmt.Sprintf("could not read processors/%s/data/languages", archName))
	}

	ldefs := filterLdefFiles(files)
	for _, ldef := range ldefs {
		processLdefFile(archName, ldef)
	}
}

func filterLdefFiles(files []os.DirEntry) []string {
	var ldefs []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".ldefs") {
			ldefs = append(ldefs, file.Name())
		}
	}
	return ldefs
}

func processLdefFile(archName, ldef string) {
	langs, err := ProcessorsFS.ReadFile(fmt.Sprintf("processors/%s/data/languages/%s", archName, ldef))
	if err != nil {
		panic(fmt.Sprintf("could not read %s.ldefs", archName))
	}

	if bytes.Contains(langs, []byte("version=\"1.1\"")) {
		langs = bytes.ReplaceAll(langs, []byte("version=\"1.1\""), []byte("version=\"1.0\""))
	}

	var l archLanguages
	if err := xml.Unmarshal(langs, &l); err != nil {
		panic(fmt.Sprintf("could not unmarshal %s.ldefs: %v", archName, err))
	}

	for _, lang := range l.Langs {
		processLanguage(archName, lang)
	}
}

func processLanguage(archName string, lang languageDef) {
	var al ArchitectureLanguage
	al.Description = lang.Description
	al.LanguageID = strings.ToLower(lang.ID)

	pspec, err := ProcessorsFS.ReadFile(fmt.Sprintf("processors/%s/data/languages/%s", archName, lang.PSpec))
	if err != nil {
		panic(fmt.Sprintf("could not read %s", lang.PSpec))
	}

	var ps ProcessorSpec
	if err := xml.Unmarshal(pspec, &ps); err != nil {
		panic(fmt.Sprintf("could not unmarshal %s", lang.PSpec))
	}

	sla, err := ProcessorsFS.ReadFile(fmt.Sprintf("processors/%s/data/languages/%s", archName, lang.SLAFile))
	if err != nil {
		panic(fmt.Sprintf("could not read %s", lang.SLAFile))
	}

	al.ProcessorSpecs = ps
	al.Sla = sla
	ArchLanguages = append(ArchLanguages, al)
}
