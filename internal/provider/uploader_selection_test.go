package provider

import "testing"

// selectUploader exercises the same precedence New() uses, so the test doesn't
// depend on unexported provider internals beyond the factory function.
func TestUploaderSelectionPrecedence(t *testing.T) {
	p := New("test")().(*hcloudimageProvider)

	t.Run("fake forced by env", func(t *testing.T) {
		t.Setenv("HCLOUDIMAGE_FAKE", "1")
		u, err := p.newUploader(providerConfig{Token: "real-token"})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := u.(*FakeUploader); !ok {
			t.Errorf("HCLOUDIMAGE_FAKE=1 should force fake, got %T", u)
		}
	})

	t.Run("token selects real", func(t *testing.T) {
		t.Setenv("HCLOUDIMAGE_FAKE", "")
		u, err := p.newUploader(providerConfig{Token: "real-token"})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := u.(*hcloudUploader); !ok {
			t.Errorf("token should select real uploader, got %T", u)
		}
	})

	t.Run("no token falls back to fake", func(t *testing.T) {
		t.Setenv("HCLOUDIMAGE_FAKE", "")
		u, err := p.newUploader(providerConfig{})
		if err != nil {
			t.Fatal(err)
		}
		if _, ok := u.(*FakeUploader); !ok {
			t.Errorf("no token should fall back to fake, got %T", u)
		}
	})
}
