package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const backupFolder = "img-backup"

func main() {
	run()
}

func run() {
	searchDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	hashes := make([]string, 0)
	first := false
	err = filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if f.IsDir() || !isImage(path) || isBackupDir(path) {
			return nil
		}
		hashed := hash(path)
		if !exists(hashes, hashed) {
			hashes = append(hashes, hashed)
		} else {
			moveToBackup(path, f.Name(), first)
			first = false
		}
		return err
	})

	if err != nil {
		panic(err)
	}

	return
}

func exists(hashes []string, hash string) bool {
	for _, h := range hashes {
		if h == hash {
			return true
		}
	}

	return false
}

func hash(path string) string {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		panic(err)
	}

	return string(h.Sum(nil))
}

func isImage(path string) bool {
	imgFormats := []string{
		".jpeg", ".jpg", ".png", ".tif",
		".tiff", ".bmp", ".raw", ".gif",
	}
	for _, f := range imgFormats {
		if strings.HasSuffix(path, f) {
			return true
		}
	}
	return false
}

func isBackupDir(path string) bool {
	if strings.Contains(path, backupFolder) {
		return true
	}
	return false
}

func moveToBackup(oldPath string, filename string, first bool) {
	if first {
		os.Mkdir(backupFolder, 0700)
	}
	currDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	backupPath := currDir + "/" + backupFolder + "/" + getRandomStr() + filename
	fmt.Println("Moving " + oldPath + " to backup folder.")

	err = os.Rename(oldPath, backupPath)
	if err != nil {
		panic(err)
	}
}

func getRandomStr() string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, 6)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	rand.Seed(time.Now().UnixNano())

	return string(b) + "-"
}
