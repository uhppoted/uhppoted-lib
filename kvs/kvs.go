package kvs

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"sync"
	"time"
)

type KeyValueStore struct {
	name      string
	store     map[string]interface{}
	version   uint64
	stored    uint64
	guard     sync.Mutex
	writeLock sync.Mutex
	re        *regexp.Regexp
	f         func(string) (interface{}, error)
}

func NewKeyValueStore(name string, f func(string) (interface{}, error)) *KeyValueStore {
	return &KeyValueStore{
		name:      name,
		store:     map[string]interface{}{},
		version:   0,
		stored:    0,
		guard:     sync.Mutex{},
		writeLock: sync.Mutex{},
		re:        regexp.MustCompile(`^\s*(.*?)(?:\s{2,})(\S.*)\s*`),
		f:         f,
	}
}

func (kv *KeyValueStore) Get(key string) (interface{}, bool) {
	kv.guard.Lock()
	defer kv.guard.Unlock()

	value, ok := kv.store[key]

	return value, ok
}

func (kv *KeyValueStore) Put(key string, value interface{}) {
	kv.guard.Lock()
	defer kv.guard.Unlock()

	kv.store[key] = value

	c := map[string]interface{}{}
	for k, v := range kv.store {
		c[k] = v
	}
}

func (kv *KeyValueStore) Store(key string, value interface{}, filepath string, log *log.Logger) {
	kv.guard.Lock()
	defer kv.guard.Unlock()

	kv.store[key] = value
	kv.version += 1

	go kv.save(filepath, log)
}

func (kv *KeyValueStore) LoadFromFile(filepath string) error {
	if filepath == "" {
		return nil
	}

	f, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer f.Close()

	return kv.load(f)
}

// Ref. https://www.joeshaw.org/dont-defer-close-on-writable-files/
func (kv *KeyValueStore) save(file string, log *log.Logger) {
	// ... copy current store

	kv.guard.Lock()

	version := kv.version
	store := map[string]interface{}{}
	for k, v := range kv.store {
		store[k] = v
	}

	kv.guard.Unlock()

	// ... store copy to file

	kv.writeLock.Lock()
	defer kv.writeLock.Unlock()

	if file == "" {
		return
	}

	dir := filepath.Dir(file)
	filename := fmt.Sprintf("%s.%d", filepath.Base(file), kv.version)
	tmpfile := filepath.Join(dir, filename)

	os.MkdirAll(dir, os.ModeDir|os.ModePerm)

	f, err := os.Create(tmpfile)
	if err != nil {
		log.Printf("ERROR: %s - %v", kv.name, err)
		return
	}

	for key, value := range store {
		if _, err := fmt.Fprintf(f, "%-20s  %v\n", key, value); err != nil {
			log.Printf("ERROR: %s - %v", kv.name, err)
			f.Close()
			return
		}
	}

	if err := f.Sync(); err != nil {
		log.Printf("ERROR: %s - %v", kv.name, err)
		f.Close()
		return
	}

	if err := f.Close(); err != nil {
		log.Printf("ERROR: %s - %v", kv.name, err)
		return
	}

	if version > kv.stored {
		if err := os.Rename(tmpfile, file); err != nil {
			log.Printf("ERROR: %s - %v", kv.name, err)
		} else {
			kv.stored = version
		}
	} else {
		log.Printf("WARN: %s - out of date version discarded", kv.name)
		if err := os.Remove(tmpfile); err != nil {
			log.Printf("ERROR: %s - %v", kv.name, err)
		}
	}
}

func (kv *KeyValueStore) Save(w io.Writer) error {
	for key, value := range kv.store {
		if _, err := fmt.Fprintf(w, "%-20s  %v\n", key, value); err != nil {
			return err
		}
	}

	return nil
}

// NOTE: interim file watcher implementation pending fsnotify in Go 1.4
//       (https://github.com/fsnotify/fsnotify requires workarounds for
//        files updated atomically by renaming)
func (kv *KeyValueStore) Watch(filepath string, logger *log.Logger) {
	go func() {
		finfo, err := os.Stat(filepath)
		if err != nil {
			logger.Printf("ERROR Failed to get file information for '%s': %v", filepath, err)
			return
		}

		lastModified := finfo.ModTime()
		logged := false
		for {
			time.Sleep(2500 * time.Millisecond)
			finfo, err := os.Stat(filepath)
			if err != nil {
				if !logged {
					logger.Printf("ERROR Failed to get file information for '%s': %v", filepath, err)
					logged = true
				}

				continue
			}

			logged = false
			if finfo.ModTime() != lastModified {
				log.Printf("INFO  Reloading information from %s\n", filepath)

				err := kv.LoadFromFile(filepath)
				if err != nil {
					log.Printf("ERROR Failed to reload information from %s: %v", filepath, err)
					continue
				}

				log.Printf("WARN  Updated %s from %s", kv.name, filepath)
				lastModified = finfo.ModTime()
			}
		}
	}()
}

func (kv *KeyValueStore) load(r io.Reader) error {
	store := map[string]interface{}{}
	s := bufio.NewScanner(r)
	for s.Scan() {
		match := kv.re.FindStringSubmatch(s.Text())
		if len(match) == 3 {
			key := strings.TrimSpace(match[1])
			value := strings.TrimSpace(match[2])

			if v, err := kv.f(value); err != nil {
				return err
			} else {
				store[key] = v
			}
		}
	}

	if s.Err() != nil {
		return s.Err()
	}

	return kv.merge(store)
}

func (kv *KeyValueStore) merge(store map[string]interface{}) error {
	kv.guard.Lock()
	defer kv.guard.Unlock()

	if !reflect.DeepEqual(store, kv.store) {
		for k, v := range store {
			kv.store[k] = v
		}

		for k, _ := range kv.store {
			if _, ok := store[k]; !ok {
				delete(kv.store, k)
			}
		}
	}

	return nil
}
