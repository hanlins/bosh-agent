package httpblobprovider

import (
	"fmt"
	"net/http"
	"os"

	boshcrypto "github.com/cloudfoundry/bosh-utils/crypto"
	boshsys "github.com/cloudfoundry/bosh-utils/system"
)

var DefaultCryptoAlgorithms = []boshcrypto.Algorithm{boshcrypto.DigestAlgorithmSHA1, boshcrypto.DigestAlgorithmSHA512}

type HTTPBlobImpl struct {
	fs               boshsys.FileSystem
	createAlgorithms []boshcrypto.Algorithm
}

func NewHTTPBlobImpl(fs boshsys.FileSystem) HTTPBlobImpl {
	return HTTPBlobImpl{
		fs: fs,
	}
}

func (h HTTPBlobImpl) WithDefaultAlgorithms() HTTPBlobImpl {
	h.createAlgorithms = DefaultCryptoAlgorithms
	return h
}

func (h HTTPBlobImpl) WithAlgorithms(a []boshcrypto.Algorithm) HTTPBlobImpl {
	h.createAlgorithms = a
	return h
}

func (h HTTPBlobImpl) Upload(signedURL, filepath string) (boshcrypto.MultipleDigest, error) {
	digest, err := boshcrypto.NewMultipleDigestFromPath(filepath, h.fs, h.createAlgorithms)
	if err != nil {
		return boshcrypto.MultipleDigest{}, err
	}

	// Do not close the file in the happy path because the client.Do will handle that.
	file, err := h.fs.OpenFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		return boshcrypto.MultipleDigest{}, err
	}

	stat, err := h.fs.Stat(filepath)
	if err != nil {
		defer file.Close()
		return boshcrypto.MultipleDigest{}, err
	}

	req, err := http.NewRequest("PUT", signedURL, file)
	if err != nil {
		defer file.Close()
		return boshcrypto.MultipleDigest{}, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Expect", "100-continue")
	req.ContentLength = stat.Size()

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return boshcrypto.MultipleDigest{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return boshcrypto.MultipleDigest{}, fmt.Errorf("Error executing PUT to %s for %s, response was %+v", signedURL, file, resp)
	}

	return digest, nil
}

func (h HTTPBlobImpl) Get(signedURL string) (string, error) {
	return "", nil
}