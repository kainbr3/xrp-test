package files

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCases_Files_Unit(t *testing.T) {
	tests := []struct {
		name     string
		testFunc func(*testing.T)
	}{
		{"Success retrieving package runtime path", testGetPackageRuntimePath},
		{"Success retrieving caller file runtime path", testGetCallerFileRuntimePath},
		{"Success retrieving combined runtime path with given path", testCombineRuntimePathWithGivenPath},
		{"Success reading data from an existing file", testReadFile},
		{"Failure reading data from non-existent path or file", testFailReadFile},
		{"Success reading data from a given path", testReadFileWithoutBasePath},
		{"Failure reading data from an non-existent given path", testFailReadFileWithoutBasePath},
	}

	for _, tt := range tests {
		tt := tt // capture range variable
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() //adds parallel execution for each test case
			tt.testFunc(t)
		})
	}
}

func testGetPackageRuntimePath(t *testing.T) {
	t.Log("TestGetPackageRuntimePath - Testing a success clause for retreiving runtime path")
	executionRuntimeDir, err := GetPackageRuntimePath()
	assert.NoError(t, err)
	assert.NotEmpty(t, executionRuntimeDir)
}

func testGetCallerFileRuntimePath(t *testing.T) {
	t.Log("TestGetCallerFileRuntimePath - Testing a success clause for retreiving runtime file path")
	result, err := GetCallerFileRuntimePath()
	callerFileName := strings.Split(result, "/files/")[1]
	assert.NoError(t, err)
	assert.NotEmpty(t, result)
	assert.Equal(t, callerFileName, "files.go")
}

func testCombineRuntimePathWithGivenPath(t *testing.T) {
	t.Log("TestCombineRuntimePathWithGivenPath - Testing a success clause for retreiving the complete path for test-file.json")
	basePath, _ := GetCallerFileRuntimePath()
	result := CombineGivenPathAndGivenFilePath(basePath, "/json/test-file.json")
	fileName := strings.Split(result, "/json/")[1]
	assert.NotEmpty(t, result)
	assert.Equal(t, fileName, "test-file.json")
}

func testReadFile(t *testing.T) {
	t.Log("TestReadFile - Testing a success clause for reading a file from a given partial file path")
	basePath, _ := GetCallerFileRuntimePath()
	fullpath := CombineGivenPathAndGivenFilePath(basePath, "/json/test-file.json")
	result, err := ReadFile(fullpath)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	jsonResult := string(result)
	assert.Contains(t, jsonResult, "\"TestProp\"")
	assert.Contains(t, jsonResult, "\"TestValue\"")
}

func testFailReadFile(t *testing.T) {
	t.Log("TestFailReadFile - Testing a failure clause for reading a file from a given partial file path")
	basePath, _ := GetCallerFileRuntimePath()
	fullpath := CombineGivenPathAndGivenFilePath(basePath, "/../invalid_path/test-file.json")
	result, err := ReadFile(fullpath)
	errorMessage := strings.Split(err.Error(), ": ")[1]
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, errorMessage, "no such file or directory")
}

func testReadFileWithoutBasePath(t *testing.T) {
	t.Log("TestReadFileWithoutBasePath - Testing a success clause for reading a file from given complete path")
	os.Setenv("ENVIROMENT", "development")
	result, err := ReadFileWithoutBasePath("/json/test-file.json")
	assert.NoError(t, err)
	assert.NotNil(t, result)

	jsonResult := string(result)
	assert.Contains(t, jsonResult, "\"TestProp\"")
	assert.Contains(t, jsonResult, "\"TestValue\"")
}

func testFailReadFileWithoutBasePath(t *testing.T) {
	t.Log("TestReadFileWithoutBasePath - Testing a failure clause for reading a file from given complete path")
	result, err := ReadFileWithoutBasePath("/../invalid_path/test-file.json")
	errorMessage := strings.Split(err.Error(), ": ")[1]
	assert.Nil(t, result)
	assert.Error(t, err)
	assert.Equal(t, errorMessage, "no such file or directory")
}
