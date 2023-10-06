package main

import (
	"bufio"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"sync"
)

// WalkFiles sends _normal_ files on the files channel. if done is closed, fs.SkipAll is returned.
func WalkFiles(done chan struct{}, files chan<- string) fs.WalkDirFunc {
	return func(path string, d fs.DirEntry, err error) error {
		// fail early if received error
		if err != nil {
			return err
		}

		select {
		case <-done:
			return fs.SkipAll
		default:
		}

		info, err := d.Info()
		if err != nil {
			return err
		}

		// not a normal file
		if info.Mode()&fs.ModeType != 0 {
			return nil
		}
		select {
		case <-done:
			return fs.SkipAll
		case files <- path:
		}

		return nil
	}
}

type FileHash struct {
	absPath string
	hash    string
	err     error
}

// FileHashWorker receives files from the `files` channel, reads and computes the file's sha256
// hash, and send the result on the hashes channel.
func FileHashWorker(wg *sync.WaitGroup, files <-chan string, hashes chan<- FileHash) {
	defer wg.Done()

	buffReader := new(bufio.Reader)
	hasher := sha256.New()
	for filePath := range files {
		fp, err := os.Open(filePath)
		if err != nil {
			hashes <- FileHash{absPath: filePath, err: err}
			continue
		}

		buffReader.Reset(fp)
		hasher.Reset()

		_, err = io.Copy(hasher, buffReader)
		if err != nil {
			hashes <- FileHash{absPath: filePath, err: err}
			continue
		}

		err = fp.Close()
		if err != nil {
			hashes <- FileHash{absPath: filePath, err: err}
			continue
		}

		hash := hasher.Sum(nil)
		hashes <- FileHash{absPath: filePath, hash: fmt.Sprintf("%x", hash)}
	}
}
