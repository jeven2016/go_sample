package read_ini

import "testing"

func TestReadIniFile(t *testing.T) {

	t.Run("read file", func(t *testing.T) {
		ReadIniFile()
	})

	t.Run("test Generate", func(t *testing.T) {
		println("= Generate============================")
		Generate()

	})
}
