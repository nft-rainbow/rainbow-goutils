package empty

import "io"

type EmptyEncryptor struct {
}

func (a *EmptyEncryptor) Encrypt(input io.Reader, output io.Writer, key []byte) error {
	_, err := io.Copy(output, input)
	return err
}

func (a *EmptyEncryptor) Decrypt(input io.Reader, output io.Writer, key []byte) error {
	_, err := io.Copy(output, input)
	return err
}
