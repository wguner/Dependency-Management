package filesystem

import (
	"log"
	"os"
)

func CreateDefaultSubdirectories() error {
	subdirectories := []string{"projects", "packages", "builds", "tmp"}
	if err := CreateSubdirectories(subdirectories); err != nil {
		return err
	}
	return nil
}

func CreateSubdirectories(paths []string) error {
	for _, path := range paths {
		if err := CreateSubdirectory(path); err != nil {
			log.Print(err)
			return err
		}
	}
	return nil
}

func CreateSubdirectory(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0755)
		return nil
	} else if err == nil {
		return nil
	} else {
		log.Print(err)
		return err
	}
}
