package tar_test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/plgd-dev/client-application/pkg/tar"
	"github.com/stretchr/testify/require"
)

func TestTar(t *testing.T) {
	data := bytes.NewBuffer(make([]byte, 0, 512))
	ex, err := os.Executable()
	require.NoError(t, err)
	err = tar.Tar(ex, data)
	require.NoError(t, err)
	err = tar.Untar(os.TempDir()+string(os.PathSeparator)+filepath.Base(ex), data)
	require.NoError(t, err)
}
