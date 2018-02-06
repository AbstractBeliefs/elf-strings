# elf-strings
elf-strings will programmatically read an ELF binary's string sections within a given binary. This is meant to be much like the `strings` UNIX utility, however is purpose built for ELF binaries. 

This means that you can get suitable information about the strings within the binary, such as the section they reside in, the offset in the section, etc.. This utility also has the functionality to 'demangle' C++ symbols, iterate linked libraries and print basic information about the ELF.

# Output
![alt text](https://i.imgur.com/plIdQCF.png "example of demangled strings")

# Arguments
```
  -binary string
        the path to the ELF you wish to parse
  -demangle
        demangle C++ symbols into their original source indentifiers, prettify found C++ symbols (optional)
  -hex
        output the strings as a hexadecimal literal (optional)
  -libs
        show the linked libraries in the binary (optional)
  -max uint
        the maximum amount of strings that you wish to be output (optional)
  -offset
        show the offset of the string in the section (default, recommended) (default true)
  -output-file string
        the path of the output file that you want to output to (optional)
  -output-format string
        the format you want to output as (optional, plain/json/xml) (default "plain")
```

An example grabbing the strings from the binary `hello`.
`./elfstrings --binary=hello --demangle`
