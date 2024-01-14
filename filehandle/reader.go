package filehandle

import (
	"os"
	"path"
	"sync"
	"sync/atomic"
)

func CalculateDirectoryAllFileSize(dirPath string) (int64, error) {
	dirEntryList, readErr := os.ReadDir(dirPath)
	if os.IsNotExist(readErr) {
		return 0, nil
	}
	if readErr != nil {
		return 0, readErr
	}
	if len(dirEntryList) <= 0 {
		return 0, nil
	}

	var filesSize int64
	for _, entry := range dirEntryList {
		if entry.IsDir() {
			continue
		}
		fileInfo, errFileInfo := entry.Info()
		if errFileInfo != nil {
			return 0, errFileInfo
		}
		filesSize += fileInfo.Size()
	}

	return filesSize, nil
}

func CalculateAllFileSizeThroughRecursion(dirPath string, replyFilesSize *int64) error {
	dirEntryList, readErr := os.ReadDir(dirPath)
	if os.IsNotExist(readErr) {
		return nil
	}
	if readErr != nil {
		return readErr
	}
	if len(dirEntryList) <= 0 {
		return nil
	}

	var subDirPaths []string
	for _, entry := range dirEntryList {
		if entry.IsDir() {
			subDirPaths = append(subDirPaths, path.Join(dirPath, entry.Name()))
			continue
		}
		fileInfo, errFileInfo := entry.Info()
		if errFileInfo != nil {
			return errFileInfo
		}
		atomic.AddInt64(replyFilesSize, fileInfo.Size())
	}
	if len(subDirPaths) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	wg.Add(len(subDirPaths))
	for _, subDirPath := range subDirPaths {
		go func(dirPath string) {
			defer wg.Done()
			CalculateAllFileSizeThroughRecursion(dirPath, replyFilesSize)
		}(subDirPath)
	}
	wg.Wait()

	return nil
}
