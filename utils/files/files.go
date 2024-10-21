package files

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// GetPackageRuntimePath - retrieves the complete path of the package being executed
// Result: {directory} string - the result is a string representing the file path of the package being executed
// Result: {err} error - returns a error if the could not retrieve the file path
func GetPackageRuntimePath() (string, error) {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}
	return dir, nil
}

// GetCallerFileRuntimePath - retrieves the complete path of the file being executed
// Result: {filePath} string - the result is a string representing the file path of the package being executed
// Result: {err} error - returns a error if the could not retrieve the runtime caller
func GetCallerFileRuntimePath() (string, error) {
	// gets the runtime path from the env package
	_, f, _, success := runtime.Caller(0)
	if !success {
		return "", errors.New("error retrieving runtime path")
	}

	return f, nil
}

// CombineGivenPathAndGivenFilePath - retrieves a combination for both given paths, one for the base path and other for the file
// Param: {basePath} string - a existent base path for the package
// Param: {filePath} string - a complete file path for a existing file
// Result: {directory} string - the result is a string representing the complete file path resulting from the combination for both given paths, one for the package and other for the .env
func CombineGivenPathAndGivenFilePath(basePath, filePath string) string {
	dir := filepath.Join(filepath.Dir(basePath), filePath)

	return dir
}

// ReadFile - retrieves a byte array for a read file from the given complete path
// Param: {fullPath} string - a complete file path for a existing file
// Result: {file} []byte - the result is a byte array for the read file
// Result: {err} error - returns a error if the could not read/find a file for the given complete file path
func ReadFile(fullPath string) ([]byte, error) {
	result, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// ReadFileWithoutBasePath - retrieves a byte array for a read file from the given partial path
// Param: {fullPath} string - a complete file path for a existing file
// Result: {file} []byte - the result is a byte array for the read file
// Result: {err} error - returns a error if the could not read/find a file for the given partial file path
func ReadFileWithoutBasePath(fullPath string) ([]byte, error) {
	path := fullPath
	if strings.EqualFold(os.Getenv("ENVIROMENT"), "development") {
		basePath, _ := GetCallerFileRuntimePath()
		path = CombineGivenPathAndGivenFilePath(basePath, fullPath)
	}

	result, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// WriteFile - create a file and save it with the given complete path
// Param: {fullPath} string - a complete file path for a existing file
// Param: {file} []byte - the file to be saved
// Result: {err} error - returns a error if the could not rsave the file for the given complete file path
func WriteFile(fullPath string, file []byte) error {
	err := os.WriteFile(fullPath, file, 0644)
	if err != nil {
		return err
	}

	return nil
}
