package sbbs

import "os"

// A helpful utility function that creates a file and logs the file's path.
func CreateFile(name string) (*os.File, error) {
	LogQuietInfo("Creating and Opening File: '%s'", name)
	return os.Create(name)
}

// A helpful utility function that opens a file and logs the file's path.
func Open(name string) (*os.File, error) {
	LogQuietInfo("Opening File: '%s'", name)
	return os.Open(name)
}

// A helpful utility function that creates but does not open a file and logs the
// file's path.
func Touch(name string) error {
	LogQuietInfo("Creating File: '%s'", name)
	f, err := os.Create(name)
	defer f.Close()
	return err
}

// A helpful utility function that creates the supplied directory as well as all
// necessary parent directories.
func Mkdir(path string) error {
	LogQuietInfo("Creating Dir(s): '%s'", path)
	return os.MkdirAll(path, 0755)
}

// A helpful utility function that changes the programs current working
// directory and logs the old and new current working directories.
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
