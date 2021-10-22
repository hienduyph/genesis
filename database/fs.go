package database

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func getDatabaseDirPath(dataDir string) string {
	return filepath.Join(dataDir, "database")
}

func getGenesisJSONPathFile(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "genesis.json")
}

func getBlocksDBFilePath(dataDir string) string {
	return filepath.Join(getDatabaseDirPath(dataDir), "block.db")
}

func initDataIfNotExists(dataDir string) error {
	if fileExist(getGenesisJSONPathFile(dataDir)) {
		return nil
	}
	fmt.Printf("Initializing system at `%s`\n", dataDir)
	if err := os.MkdirAll(getDatabaseDirPath(dataDir), os.ModePerm); err != nil {
		return fmt.Errorf("create db dir failed: %w", err)
	}
	if err := writeGenesisToDisk(getGenesisJSONPathFile(dataDir)); err != nil {
		return fmt.Errorf("write genesis file failed: %w", err)
	}
	if err := writeEmptyBlocksDbToDisk(getBlocksDBFilePath(dataDir)); err != nil {
		return fmt.Errorf("write block.db file failed: %w", err)
	}
	return nil
}

func fileExist(filePah string) bool {
	_, err := os.Stat(filePah)
	return !(err != nil && os.IsNotExist(err))
}

func dirExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func writeEmptyBlocksDbToDisk(path string) error {
	return ioutil.WriteFile(path, []byte(""), os.ModePerm)
}
