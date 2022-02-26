package helper

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMakeSliceFromDir(t *testing.T) {
	currentPath := "./"
	filePaths, err := MakeSliceFromDir(currentPath)
	require.NoError(t, err)
	require.Equal(t, "Quit", filePaths[0])
	require.Equal(t, "../", filePaths[1])

}
