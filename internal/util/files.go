package util

import "os"

// RemoveFileIfExists removes filename if exists, or does nothing if the file
// is not there. Returns an error, if it occurred during deletion.
func RemoveFileIfExists(filename string) error {
	_, err := os.Stat(filename)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		// no db file
		return nil
	}

	return os.Remove(filename)
}
