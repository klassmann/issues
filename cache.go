package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"time"
)

// Here we are using GOB serialization in cache for efficiency

// CacheFile is used per query name
// and saved per Configuration.Cache duration
type CacheFile struct {
	ID      string // ID is the query name in configuration, it is used in the name of file
	Result  IssueQueryResult
	Updated time.Time
}

func (c *CacheFile) printCache() {
	fmt.Printf("ID: %v\n", c.ID)
	fmt.Printf("Result: %v\n", c.Result)
	fmt.Printf("Updated: %v\n", c.Updated)
}

const cacheFilenamePattern string = "%s.cache"

func getCacheFilename(id string) string {
	return path.Join(os.TempDir(), fmt.Sprintf(cacheFilenamePattern, id))
}

func loadCache(id string) (*CacheFile, error) {

	filename := getCacheFilename(id)
	if !fileExists(filename) {
		return nil, fmt.Errorf("cache loading: cache file for %s does not exist", id)
	}

	f, err := os.Open(filename)
	defer f.Close()

	if err != nil {
		return nil, fmt.Errorf("open cache file: %v", err)
	}

	reader := bufio.NewReader(f)
	dec := gob.NewDecoder(reader)

	var cache CacheFile
	err = dec.Decode(&cache)
	// cache.printCache()

	if err != nil {
		return nil, fmt.Errorf("reading cache file: %v", err)
	}
	return &cache, nil
}

func saveCache(id string, result *IssueQueryResult) error {
	now := time.Now()
	cache := CacheFile{ID: id, Result: *result, Updated: now}
	filename := getCacheFilename(id)
	f, err := os.Create(filename)
	defer f.Close()

	if err != nil {
		return err
	}

	writer := bufio.NewWriter(f)
	enc := gob.NewEncoder(writer)
	err = enc.Encode(&cache)
	writer.Flush()

	return err
}
