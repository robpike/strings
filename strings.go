// Copyright 2014 The rspace Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Strings is a more capable, UTF-8 aware version of the standard strings utility.
package main // import "robpike.io/cmd/strings"

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var (
	min    = flag.Int("min", 6, "minimum length of UTF-8 strings printed, in runes")
	max    = flag.Int("max", 256, "maximum length of UTF-8 strings printed, in runes")
	ascii  = flag.Bool("ascii", false, "restrict strings to ASCII")
	offset = flag.Bool("offset", false, "show file name and offset of start of each string")
)

var stdout *bufio.Writer

func main() {
	log.SetFlags(0)
	log.SetPrefix("strings")
	stdout = bufio.NewWriter(os.Stdout)
	defer stdout.Flush()

	flag.Parse()
	if *max < *min {
		*max = *min
	}

	if flag.NArg() == 0 {
		do("<stdin>", os.Stdin)
	} else {
		for _, arg := range flag.Args() {
			fd, err := os.Open(arg)
			if err != nil {
				log.Print(err)
				continue
			}
			do(arg, fd)
			stdout.Flush()
			fd.Close()
		}
	}
}

func do(name string, file *os.File) {
	in := bufio.NewReader(file)
	str := make([]rune, 0, *max)
	filePos := int64(0)
	print := func() {
		if len(str) >= *min {
			s := string(str)
			if *offset {
				fmt.Printf("%s:#%d:\t%s\n", name, filePos-int64(len(s)), s)
			} else {
				fmt.Println(s)
			}
		}
		str = str[0:0]
	}
	for {
		var (
			r   rune
			wid int
			err error
		)
		// One string per loop.
		for ; ; filePos += int64(wid) {
			r, wid, err = in.ReadRune()
			if err != nil {
				if err != io.EOF {
					log.Print(err)
				}
				return
			}
			if !strconv.IsPrint(r) || *ascii && r >= 0xFF {
				print()
				continue
			}
			// It's printable. Keep it.
			if len(str) >= *max {
				print()
			}
			str = append(str, r)
		}
	}
}
