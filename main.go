package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"

	"github.com/jedib0t/go-pretty/v6/list"
)

type config struct {
	nWorkers int
}

type configFunc func(*config)

func defaultConfig() *config {
	return &config{
		nWorkers: runtime.NumCPU(),
	}
}

func (cfg *config) NWorkers() int {
	return cfg.nWorkers
}

func (cfg *config) MustValidate() {
	if cfg.nWorkers < 1 {
		panic("nWorkers must be greater than 0")
	}
}

func WithNWorkers(nWorkers int) configFunc {
	return func(cfg *config) {
		cfg.nWorkers = nWorkers
	}
}

// WalkFilesAndBuildTree takes and walks some root directory hashing (sha256) each file on the walk
// and adding the file and its hash to a Merkle tree. The Merkle tree is built and returned.
func WalkFilesAndBuildTree(root string, done chan struct{}, cfgFns ...configFunc) *HashTree {
	cfg := defaultConfig()
	for _, fn := range cfgFns {
		fn(cfg)
	}
	cfg.MustValidate()

	var wg sync.WaitGroup

	files := make(chan string)
	hashes := make(chan FileHash)

	// walk file tree and put files on `files` channel
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := filepath.WalkDir(root, WalkFiles(done, files))
		close(files)

		if err != nil {
			// TODO: improve error handling
			log.Fatal(err)
		}
	}()

	var workerWg sync.WaitGroup
	nWorkers := cfg.NWorkers()
	workerWg.Add(nWorkers)
	for i := 0; i < nWorkers; i++ {
		go FileHashWorker(&workerWg, files, hashes)
	}
	wg.Add(1)
	go func() {
		defer wg.Done()
		workerWg.Wait()
		close(hashes)
	}()

	tree := NewHashTree(root, DefaultHasher{})

	for hash := range hashes {
		// TODO: improve error handling
		if hash.err != nil {
			log.Fatal(hash.err)
		}
		tree.Add(hash.absPath, hash.hash)
	}

	// NOTE: this _should_ not be necessary b.c. when the `hashes` channel is closed it should be
	// the last running goroutine. waiting to catch any potentially hanging goroutines
	wg.Wait()

	tree.BuildTree()
	return &tree
}

func shutdownHandler(sig chan os.Signal, done chan struct{}) {
	<-sig
	close(done)
	// TODO: improve how shutdowns are handled
	os.Exit(1)
}

func main() {
	nWorkers := flag.Int("n", runtime.NumCPU(), "Maximum number of workers. Defaults to number of cpus.")
	format := flag.String("fmt", "flat", "Output format. Options are `flat`, `tree`.")
	flag.Parse()

	// TODO: support stdin
	path := "."
	if len(flag.Args()) > 1 {
		path = flag.Args()[0]
	}

	// channel that is used to signal an early exit.
	// this channel is _only_ ever closed.
	// each stage of the pipeline tries to read from this channel and returns early upon reading.
	done := make(chan struct{})

	// setup _early_ shutdown handling
	sig := make(chan os.Signal)
	go shutdownHandler(sig, done)
	signal.Notify(sig, os.Interrupt)

	merkelTree := WalkFilesAndBuildTree(path, done, WithNWorkers(*nWorkers))

	switch *format {
	case "flat":
		fmt.Println(merkelTree)
	case "tree":
		l := list.List{}
		l.SetStyle(list.StyleConnectedRounded)
		BuildTreeView(*merkelTree, &l)
		fmt.Println(l.Render())
	default:
		fmt.Println(merkelTree)
	}
}
