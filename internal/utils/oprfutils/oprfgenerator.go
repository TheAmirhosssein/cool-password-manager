package oprfutils

import (
	"crypto"
	"fmt"
	"os"
	"path/filepath"

	"github.com/bytemare/ksf"
	"github.com/bytemare/opaque"
)

func GenerateAndSaveKeys(dir string) error {
	// Resolve absolute path
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return fmt.Errorf("resolve absolute path: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(absDir, 0700); err != nil {
		return err
	}

	files := []string{
		"oprf.bin",
		"private.bin",
		"public.bin",
	}

	// If all exist → do nothing (OPAQUE safety)
	allExist := true
	for _, name := range files {
		if _, err := os.Stat(filepath.Join(absDir, name)); err != nil {
			if os.IsNotExist(err) {
				allExist = false
				break
			}
			return err
		}
	}

	if allExist {
		return nil
	}

	conf := &opaque.Configuration{
		OPRF: opaque.P256Sha256,
		AKE:  opaque.P256Sha256,

		Hash: crypto.SHA256,
		KDF:  crypto.SHA256,
		MAC:  crypto.SHA256,

		KSF: ksf.Argon2id,
	}

	// Generate keys (RAW BYTES)
	oprfSeed := conf.GenerateOPRFSeed()
	privateKey, publicKey := conf.KeyGen()

	data := map[string][]byte{
		"oprf_seed.bin":      oprfSeed,
		"server_private.bin": privateKey,
		"server_public.bin":  publicKey,
	}

	for name, content := range data {
		path := filepath.Join(absDir, name)

		// Never overwrite existing keys
		if _, err := os.Stat(path); err == nil {
			continue
		}

		if err := os.WriteFile(path, content, 0600); err != nil {
			return fmt.Errorf("write %s: %w", name, err)
		}
	}

	return nil
}
