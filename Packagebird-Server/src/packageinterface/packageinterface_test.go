package packageinterface

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPackagePath(t *testing.T) {
	t.Logf("Testing Package path")
	expect := "C:\\Users\\ElishaAguilera\\Documents\\packages\\Oreo-v0.tar.gz"
	actual := PackagePath("Oreo-v0")
	assert.Equal(t, expect, actual)
}

func TestUnbundlePackage(t *testing.T) {
	t.Logf("Testing UnbundlePackage()")
	PackageName := "Oreo-v0"
	err := UnbundlePackage(PackageName)
	if err != nil {
		t.Log(err)
	}
}

func TestRunBuildCommand(t *testing.T) {
	t.Logf("Testing RunBuildCommand()")
	PackageName := "Oreo"
	err := RunBuildCommand(PackageName, "c")
	if err != nil {
		t.Log(err)
		assert.Fail(t, "Failed to build package")
	}
}

func TestBuildPath(t *testing.T) {
	t.Logf("Testing BuildPath()")
	PackageName := "Oreo"
	expect := "C:\\Users\\ElishaAguilera\\Documents\\packages\\Oreo\\Oreo"
	actual := BuildPath(PackageName)
	assert.Equal(t, expect, actual)
}

func TestFileExist(t *testing.T) {
	t.Logf("Tetsing FileExist()")
	if FileExist("Oreo") {
		t.Logf("File found")
	} else {
		t.Logf("File expected, but not found")
	}
}

func TestCompressFile(t *testing.T) {
	t.Logf("Testing CompressFile()")
	err := CompressFile("Oreo")
	if err != nil {
		t.Log(err)
	}
}
