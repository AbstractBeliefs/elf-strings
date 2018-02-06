package main

import (
	"flag"
	"fmt"
	"log"

	humanize "github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

var (
	demangleOpt = flag.Bool("demangle", false, "demangle C++ symbols into their original source identifiers, prettify found C++ symbols (optional)")
	hexOpt      = flag.Bool("hex", false, "output the strings as a hexadecimal literal (optional)")
	offsetOpt   = flag.Bool("offset", true, "show the offset of the string in the section (default, recommended)")
	binaryOpt   = flag.String("binary", "", "the path to the ELF you wish to parse")
	formatOpt   = flag.String("output-format", "plain", "the format you want to output as (optional, plain/json/xml)")
	outputOpt   = flag.String("output-file", "", "the path of the output file that you want to output to (optional)")
	maxOpt      = flag.Uint64("max", 0, "the maximum amount of strings that you wish to be output (optional)")
	libOpt      = flag.Bool("libs", false, "show the linked libraries in the binary (optional)")
	infoOpt     = flag.Bool("no-info", false, "don't show any information about the binary")
)

// ReadSection is the main logic here
// it combines all of the modules, etc.
func ReadSection(reader *ElfReader, section *string) {
	var err error
	var writer *OutWriter
	var count uint64

	sect := reader.ReaderParseSection(*section)

	defer reader.Close()

	if *outputOpt != "" {
		writer, err = NewOutWriter(*outputOpt, OutParseTypeStr(*formatOpt))
		if err != nil {
			log.Fatal(err.Error())
		}
	}

	if sect != nil {
		nodes := reader.ReaderParseStrings(sect)
		for k, v := range nodes {
			if *maxOpt != 0 {
				if count == *maxOpt {
					break
				}
			}

			str := string(v)
			if *demangleOpt {
				demangled, err := UtilDemangle(&str)
				if err == nil {
					str = demangled
				}
			}

			if *hexOpt {
				str = UtilConvHex(str)
			}

			if *offsetOpt {
				fmt.Printf("[%s%s]: %s\n",
					color.BlueString(*section),
					color.GreenString("+%#x", k),
					str)
			} else {
				fmt.Println(str)
			}

			if writer != nil {
				writer.WriteResult(str, k)
			}

			count++
		}
	}
}

// ReadBasic will read the basic information
// about the ELF
func ReadBasic(reader *ElfReader) {
	stat, err := reader.File.Stat()
	if err != nil {
		return
	}

	size := humanize.Bytes(uint64(stat.Size()))

	fmt.Printf(
		"[+] Size: %s\n"+
			"[+] Arch: %s\n"+
			"[+] Entry point: %#x\n"+
			"[+] Class: %s\n"+
			"[+] Byte order: %s\n",
		size,
		UtilConvertMachine(reader.ExecReader.Machine),
		reader.ExecReader.Entry,
		reader.ExecReader.Class.String(),
		reader.ExecReader.ByteOrder.String(),
	)

	if *libOpt {
		fmt.Println("[+] Libraries:")
		libs, err := reader.ExecReader.ImportedLibraries()
		if err == nil {
			for _, lib := range libs {
				fmt.Printf("\t [!] %s\n", lib)
			}
		}
	}

	fmt.Println("------------------------")
}

// main is the entrypoint for this program
func main() {
	flag.Parse()

	if *binaryOpt == "" {
		flag.PrintDefaults()
		return
	}

	r, err := NewELFReader(binaryOpt)
	if err != nil {
		log.Fatal(err.Error())
	}

	ReadBasic(r)

	sections := []string{".dynstr", ".rodata", ".rdata",
		".strtab", ".comment", ".note",
		".stab", ".stabstr"}

	for _, section := range sections {
		ReadSection(r, &section)
	}
}
