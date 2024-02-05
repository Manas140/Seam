package main

import (
	"fmt"
	"os"
)

func help() {
	fmt.Printf(`%s %s: A Smart Stitch Tool.

  %s
    -h: show this help
    -i input_directory
    -o output_directory (defaults to ./out)
    -r rough_height (defaults to 5000)
    -f output_format [jpeg|jpg|png] (defaults to jpeg)
    -q jpeg_quality (defaults to 100)
    
  %s
    -a: use absolute height (defaults to false)

  %s
    -n neighbor_count (defaults to 5)
    -s skip_step (defaults to 5)
    -t threshold (defaults to 50)

  %s
    https://github.com/manas140/seam/issues
  %s`, col("Seam",
		34),
		col("v0.0.2", 2),
		col("General:", 32),
		col("Toggle:", 36),
		col("Advanced:", 35),
		col("Report Issues:", 31),
		"\n")
	os.Exit(0)
}
