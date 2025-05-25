package files

import (
	"encoding/gob"
	"errors"
	"example/tlgbot/storage"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
)

const defaultPerm = 0774



type Storage struct {
	basePath string
}

func New(basePath string) *Storage {
	return &Storage{basePath: basePath}
}

func (s Storage) Save(page *storage.Page) (err error) {
	fPath := filepath.Join(s.basePath, page.UserName)
	if err := os.MkdirAll(fPath, defaultPerm); err != nil {
		return fmt.Errorf("save file fail: %w", err)
	}

	fName, err := fileName(page)
	if err != nil {
		return fmt.Errorf("save file fail: %w", err)
	}

	fPath = filepath.Join(fPath, fName)

	file, err := os.Create(fPath)
	if err != nil {
		return fmt.Errorf("create file fail: %w", err)
	}

	defer func() { _ = file.Close() }()

	if err := gob.NewEncoder(file).Encode(page); err != nil {
		return fmt.Errorf("encode file fail: %w", err)
	}
	return nil
}

func (s Storage) PickRandom(userName string) (page *storage.Page, err error) {
	path := filepath.Join(s.basePath, userName)

	haveFolder := folderExist(path)
	if !haveFolder {
		return nil, storage.ErrNoSavedPages
	} 

	files, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("read dir fail: %w", err)
	}

	if len(files) == 0 {
		return nil, storage.ErrNoSavedPages
	}

	// 0-9
	n := rand.Intn(len(files))

	file := files[n]

	return s.decodePage(filepath.Join(path, file.Name()))
}

func folderExist(path string) bool {
	info, err := os.Stat(path)
    if os.IsNotExist(err) {
        return false
    }
    return info.IsDir()
}

func (s Storage) Remove(p *storage.Page) error {
	fileName, err := fileName(p)
	if err != nil {
		return fmt.Errorf("get file name fail: %w ", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("remove page fail: %w with path %s", err, path)
	}

	return nil
}

func (s Storage) IsExist(p *storage.Page) (bool, error) {
	fileName, err := fileName(p)
	if err != nil {
		return false, fmt.Errorf("get file name fail: %w", err)
	}

	path := filepath.Join(s.basePath, p.UserName, fileName)

	switch _, err := os.Stat(path); {
	case errors.Is(err, os.ErrNotExist):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("is exist page fail: %w with path %s", err, path)
	}
	return true, nil
}

func (s Storage) decodePage(filePath string) (*storage.Page, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open page fail: %w", err)
	}

	defer func() { _ = f.Close() }()

	var p storage.Page

	if err := gob.NewDecoder(f).Decode(&p); err != nil {
		return nil, fmt.Errorf("decode page fail: %w", err)
	}

	return &p, nil
}

func fileName(p *storage.Page) (string, error) {
	return p.Hash()
}
