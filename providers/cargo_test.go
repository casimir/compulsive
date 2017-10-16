package providers

import (
	"reflect"
	"testing"
)

func TestLoadManifest(t *testing.T) {
	manifestContent := []byte(`[v1]
"clippy 0.0.123 (registry+https://github.com/rust-lang/crates.io-index)" = ["cargo-clippy.exe"]
"rustfmt 0.8.3 (registry+https://github.com/rust-lang/crates.io-index)" = ["cargo-fmt.exe", "rustfmt.exe"]
"ufind 0.3.0 (path+file:///C:/Users/casimir/dev/src/github.com/casimir/ufind)" = ["ufind.exe"]
`)
	expected := []cargoManifestEntry{
		{
			name:     "clippy",
			version:  "0.0.123",
			uri:      "registry+https://github.com/rust-lang/crates.io-index",
			binaries: []string{"cargo-clippy.exe"},
		},
		{
			name:     "rustfmt",
			version:  "0.8.3",
			uri:      "registry+https://github.com/rust-lang/crates.io-index",
			binaries: []string{"cargo-fmt.exe", "rustfmt.exe"},
		},
		{
			name:     "ufind",
			version:  "0.3.0",
			uri:      "path+file:///C:/Users/casimir/dev/src/github.com/casimir/ufind",
			binaries: []string{"ufind.exe"},
		},
	}
	got := unmarshalManifest(manifestContent)
	if !reflect.DeepEqual(expected, got) {
		t.Fail()
	}
}
