package api

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// POST /api/gpg
func apiGPGAddKey(c *gin.Context) {
	var b struct {
		Keyserver   string
		GpgKeyID    string
		GpgKeyArmor string
		Keyring     string
	}

	if !c.Bind(&b) {
		return
	}

	var err error
	args := []string{"--no-default-keyring"}
	keyring := "trustedkeys.gpg"
	if len(b.Keyring) > 0 {
		keyring = b.Keyring
	}
	args = append(args, "--keyring", keyring)
	if len(b.Keyserver) > 0 {
		args = append(args, "--keyserver", b.Keyserver)
	}
	if len(b.GpgKeyArmor) > 0 {
		var tempdir string
		tempdir, err = ioutil.TempDir(os.TempDir(), "aptly")
		if err != nil {
			c.Fail(400, err)
			return
		}
		defer os.RemoveAll(tempdir)

		keypath := filepath.Join(tempdir, "key")
		keyfile, err := os.Create(keypath)
		if err != nil {
			c.Fail(400, err)
			return
		}
		if _, err = keyfile.WriteString(b.GpgKeyArmor); err != nil {
			c.Fail(400, err)
		}
		args = append(args, "--import", keypath)

	}
	if len(b.GpgKeyID) > 0 {
		args = append(args, "--recv", b.GpgKeyID)
	}

	cmd := exec.Command("gpg", args...)
	cmd.Stdout = os.Stdout
	if err = cmd.Run(); err != nil {
		c.Fail(400, err)
		return
	}

	c.JSON(200, gin.H{})
}
