package cmd

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var (
	isParallel bool
	nTaken     int
	rootCmd    = &cobra.Command{
		Use:   "argg",
		Short: "Go implementation of xargs",
		Long:  "Less complete and less performant version of xargs but done in 1h and by me :)",
		Run: func(_ *cobra.Command, args []string) {
			pipedArgs := readPipedArgs(os.Stdin)
			finalArgs := splitArgsByN(pipedArgs, nTaken)
			var executor func([][]string, func([]string))
			if isParallel {
				executor = execCommandInParallel
			} else {
				executor = execCommandSequentially
			}
			executor(finalArgs, execCommand(args))
		},
	}
)

func init() {
	rootCmd.PersistentFlags().BoolVarP(&isParallel, "parallel", "P", false, "Executes commands parallelly")
	rootCmd.PersistentFlags().IntVarP(&nTaken, "number", "n", math.MaxInt, "Take n from stdin")
}

func execCommandSequentially(args [][]string, f func(arg []string)) {
	for _, arg := range args {
		f(arg)
	}
}

func execCommandInParallel(args [][]string, f func(arg []string)) {
	var wg sync.WaitGroup
	wg.Add(len(args))
	for _, arg := range args {
		go func(arg []string) {
			f(arg)
			wg.Done()
		}(arg)
	}
	wg.Wait()
}

func execCommand(args []string) func([]string) {
	return func(arg []string) {
		allArgs := mergeArgs(args, arg)
		cmd := exec.Command(args[0], allArgs...)
		cmd.Stderr = os.Stderr
		cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed executing: %v with %v\n", args[0], err)
		}
	}
}

func splitArgsByN(args []string, n int) [][]string {
	var result [][]string
	for len(args) > 0 {
		var chunk []string
		if len(args) >= n {
			chunk = args[:n]
			args = args[n:]
		} else {
			chunk = args
			args = nil
		}
		result = append(result, chunk)
	}
	return result
}

func mergeArgs(args, pipedArgs []string) []string {
	var allArgs []string
	if len(args) > 1 {
		allArgs = append(args[1:], pipedArgs...)
	} else {
		allArgs = append([]string{}, pipedArgs...)
	}
	return allArgs
}

func readPipedArgs(r io.Reader) []string {
	var args []string
	sc := bufio.NewScanner(r)
	for sc.Scan() {
		args = append(args, strings.Split(sc.Text(), " ")...)
	}
	return args
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
