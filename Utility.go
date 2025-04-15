package sbbs

import (
	"os"
)

// A utility function that creates a file and logs the file's path.
func CreateFile(name string) (*os.File, error) {
	LogQuietInfo("Creating and Opening File: '%s'", name)
	return os.Create(name)
}

// A utility function that opens a file and logs the file's path.
func Open(name string) (*os.File, error) {
	LogQuietInfo("Opening File: '%s'", name)
	return os.Open(name)
}

// A utility function that creates but does not open a file and logs the file's
// path.
func Touch(name string) error {
	LogQuietInfo("Creating File: '%s'", name)
	f, err := os.Create(name)
	defer f.Close()
	return err
}

// A utility function that creates the supplied directory as well as all
// necessary parent directories.
func Mkdir(path string) error {
	LogQuietInfo("Creating Dir(s): '%s'", path)
	return os.MkdirAll(path, 0755)
}

// A utility function that changes the programs current working directory and
// logs the old and new current working directories.
func Cd(dir string) error {
	old, err := os.Getwd()
	if err != nil {
		return err
	}
	LogQuietInfo("Previous Cwd: '%s'", old)

	err = os.Chdir(dir)
	if err != nil {
		return err
	}
	LogQuietInfo("New Cwd: '%s'", dir)
	return nil
}

// A utility function that changes the supplied env variable to the supplied
// value, returning a closure that can be used to set the env variable back to
// it's original value. If the supplied env variable did not exist before
// calling this function then the returned closure will remove the env variable
// instead of reseting it to it's original value.
func TmpEnvVarSet(name string, val string) (reset func() error, err error) {
	oldVal, ok := os.LookupEnv(name)
	reset = func() error {
		if ok {
			LogQuietInfo("Restoring '%s' env var to '%s'", name, oldVal)
			return os.Setenv(name, oldVal)
		} else {
			LogQuietInfo("Deleting '%s' env var sice it did not exist before", name)
			return os.Unsetenv(name)
		}
	}
	LogQuietInfo("Temporarily setting '%s' env var to '%s'", name, val)
	err = os.Setenv(name, val)
	return
}
