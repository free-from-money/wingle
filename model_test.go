package wingle

import "testing"

func TestPwCipher(t *testing.T) {
	const (
		token = `m8vvUi16CjN7sIOAsZr7EW2PDnzMDcMe`
		psd   = `ODc0MzliY2MyMjcyMmNkZWMwYjJkZWNmNGNjODMwZGVhYjdmODgyNjY0MDI4YmY4YzdhNjkyMmZhOTdmMjljNw==`
	)

	t.Log(PwCipher("admin", "12345678", token))
	t.Log(psd)
}
