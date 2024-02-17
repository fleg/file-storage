package app

import "os"

func CheckFolders(storagePath string) error {
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return err
	}

	return nil
}
