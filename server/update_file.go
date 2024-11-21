package main

import (
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	dirPath := "./static/static/head" // 指定目录的路径

	i := 0
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() { // 如果不是目录，重命名文件
			oldName := info.Name()
			newName := fmt.Sprintf("u%d.jpg", i)
			i++
			err := os.Rename(filepath.Join(dirPath, oldName), filepath.Join(dirPath, newName))
			if err != nil {
				return err
			}
			fmt.Printf("%s has been renamed to %s.\n", oldName, newName)
		}
		return nil
	})

	if err != nil {
		fmt.Println(err)
	}
}
