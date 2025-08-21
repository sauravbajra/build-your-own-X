package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"unicode"
)

// Counts holds the counting results for a given input.
type Counts struct {
	lines int
	words int
	bytes int
	chars int
}

func main() {
	// Define command-line flags.
	// -l for lines, -w for words, -c for bytes, -m for characters.
	linesFlag := flag.Bool("l", false, "count lines")
	wordsFlag := flag.Bool("w", false, "count words")
	bytesFlag := flag.Bool("c", false, "count bytes")
	charsFlag := flag.Bool("m", false, "count characters")
	longestLine := flag.Bool("L", false, "length of line containing most byte")	

	// Parse the flags from the command line.
	flag.Parse()

	// If no flags are specified, the default behavior is to count lines, words, and bytes.
	// We check if all flags are false (their default value).
	noFlags := !*linesFlag && !*wordsFlag && !*bytesFlag && !*charsFlag && !*longestLine
	if noFlags {
		*linesFlag = true
		*wordsFlag = true
		*bytesFlag = true
	}

	// Get the list of files from the command-line arguments.
	files := flag.Args()

	// If no files are provided, read from standard input.
	if len(files) == 0 {
		counts, err := count(os.Stdin)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading from stdin: %v\n", err)
			os.Exit(1)
		}
		printCounts(counts, "", *linesFlag, *wordsFlag, *bytesFlag, *charsFlag)
		return
	}

	// Process each file provided as an argument.
	var totalCounts Counts
	for _, filename := range files {
		// Open the file for reading.
		file, err := os.Open(filename)
		if err != nil {
fmt.Fprintf(os.Stderr, "error opening file %s: %v\n", filename, err)
			continue // Skip to the next file on error.
		}

		// The count function does the actual work.
		counts, err := count(file)
		file.Close() // Close the file after reading.
		if err != nil {
			fmt.Fprintf(os.Stderr, "error reading file %s: %v\n", filename, err)
			continue
		}

		// Print the counts for the current file.
		printCounts(counts, filename, *linesFlag, *wordsFlag, *bytesFlag, *charsFlag)

		// Add the counts of the current file to the total.
		totalCounts.lines += counts.lines
		totalCounts.words += counts.words
		totalCounts.bytes += counts.bytes
		totalCounts.chars += counts.chars
	}

	// If more than one file was processed, print the total counts.
	if len(files) > 1 {
		printCounts(totalCounts, "total", *linesFlag, *wordsFlag, *bytesFlag, *charsFlag)
	}
}

// count reads from an io.Reader and returns the line, word, byte, and character counts.
// This function performs all counting in a single pass for efficiency.
func count(r io.Reader) (Counts, error) {
	var counts Counts
	// A bufio.Reader is used for efficient reading.
	reader := bufio.NewReader(r)
	inWord := false

	for {
		// ReadRune reads a single UTF-8 encoded Unicode character (rune) and returns
		// the rune, its size in bytes, and an error.
		rune, size, err := reader.ReadRune()
		if err != nil {
			// io.EOF signals the end of the input.
			if err == io.EOF {
				break
			}
			return Counts{}, err
		}

		// Increment byte and character counts.
		counts.bytes += size
		counts.chars++

		// Increment line count on newline characters.
		if rune == '\n' {
			counts.lines++
		}

		// Word counting logic: A word is a sequence of non-whitespace characters.
		// We are in a word if the current character is not a space.
		// We count a new word when we transition from not being in a word to being in one.
		if unicode.IsSpace(rune) {
			inWord = false
		} else if !inWord {
			counts.words++
			inWord = true
		}
	}
	return counts, nil
}

// printCounts formats and prints the counts based on the active flags.
func printCounts(counts Counts, filename string, showLines, showWords, showBytes, showChars bool) {
	// Conditionally print each count field with padding for alignment.
	// The order follows the standard `wc` output: lines, words, chars, bytes.
	if showLines {
		fmt.Printf("%8d", counts.lines)
	}
	if showWords {
		fmt.Printf("%8d", counts.words)
	}
	// Note: `wc` shows characters if -m is specified, otherwise it shows bytes.
	// If both are specified, it shows both.
	if showChars {
		fmt.Printf("%8d", counts.chars)
	}
	if showBytes && !showChars { // Only show bytes if chars are not requested, unless both are flagged
		fmt.Printf("%8d", counts.bytes)
	}
	if showBytes && showChars {
		fmt.Printf("%8d", counts.bytes)
	}


	// Print the filename if one is provided.
	if filename != "" {
		fmt.Printf(" %s\n", filename)
	} else {
		// If reading from stdin, just print a newline.
		fmt.Println()
	}
}
