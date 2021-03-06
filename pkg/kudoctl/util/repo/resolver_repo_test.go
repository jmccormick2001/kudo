package repo

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing findfirst for the following combo
// foo-2.0.0_1.0.1.tgz			# app version semver
// foo-3.0.0_1.0.1.tgz			# app version semver
// foo-bar_1.0.1.tgz			#  app version not semver
// foo-1.0.0.tgz,   			# no app version
// the rules are semver takes ordered precedence
// non semver takes order from position in index file (first is highest)
func TestIndexFile_FindFirstMatch(t *testing.T) {

	index := createTestIndexFile()

	// first operator version should be the latest app version
	pv, _ := index.FindFirstMatch("Foo", "", "1.0.1")
	assert.Equal(t, "3.0.0", pv.AppVersion)

	// first buzz operator version is the one without an app version
	pv, _ = index.FindFirstMatch("Buzz", "", "1.0.1")
	assert.Equal(t, "", pv.AppVersion)
	assert.Equal(t, "1.0.1", pv.OperatorVersion)

	//  buzz search without app or op version should have the latest buzz
	pv, _ = index.FindFirstMatch("Buzz", "", "")
	assert.Equal(t, "1.0.2", pv.OperatorVersion)

	// specific app version should have the app version
	pv, _ = index.FindFirstMatch("Foo", "2.0.0", "")
	assert.Equal(t, "2.0.0", pv.AppVersion)
	assert.Equal(t, "1.0.1", pv.OperatorVersion)

	// same search but specific for app and op versions
	pv, _ = index.FindFirstMatch("Foo", "2.0.0", "1.0.1")
	assert.Equal(t, "2.0.0", pv.AppVersion)
	assert.Equal(t, "1.0.1", pv.OperatorVersion)

}

func createTestIndexFile() *IndexFile {

	index := &IndexFile{}

	pv := createPackageVersion("Foo", "2.0.0", "1.0.1")
	addToIndex(index, pv)
	pv = createPackageVersion("Foo", "3.0.0", "1.0.1")
	addToIndex(index, pv)
	pv = createPackageVersion("Foo", "", "1.0.1")
	addToIndex(index, pv)
	pv = createPackageVersion("Foo", "bar", "1.0.1")
	addToIndex(index, pv)
	pv = createPackageVersion("Buzz", "", "1.0.1")
	addToIndex(index, pv)
	pv = createPackageVersion("Buzz", "", "1.0.2")
	addToIndex(index, pv)

	index.sortPackages()
	return index
}

func addToIndex(index *IndexFile, pv PackageVersion) {
	err := index.AddPackageVersion(&pv)
	if err != nil {
		fmt.Printf("err in test %v", err)
	}
}

func createPackageVersion(name, appVersion, operatorVersion string) PackageVersion {

	m := Metadata{
		Name:            name,
		OperatorVersion: operatorVersion,
		AppVersion:      appVersion,
	}
	return PackageVersion{
		Metadata: &m,
	}
}
