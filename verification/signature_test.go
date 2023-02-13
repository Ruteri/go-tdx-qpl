package verification

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/edgelesssys/go-tdx-qpl/verification/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

// RawQuoteBlob is a Base64-encoded version of an example quote generated on an Intel TDX development platform.
// This is a duplicate from the types package, but this test likely will be thrown away later anyway as it's mainly just here for prototyping.
const rawQuoteBlob = "BAACAIEAAAAAAAAAk5pyM/ecTKmUCg2zlX8GB5/OUj/OJupF09PbkG1RcaEAAAAAAwAFAAAAAAAAAAAAAAAAAC/SecFhZKk91b83PYNDKNRgCMK2k6+eu4ZbCLLO0yDJqJtIaan6tg++nQxaU2PGVgAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEAAAAADnAgYAAAAAALZeoAnkJOb3Yf3T18iWJDlFOzfs32LaBPe8XTJ2hruLr8il0kqcMc7mDkq6h8L3GwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAOGvdeYZJ0EOQrVLOfZoHPmwv7rlErFehw5MjZ1aXLOFVxsOHcL3C/nM7whWDworWCFf8fwMMUQsHwYaMXvkCUCxgsE9Q8bbLlsqV33em+6T1FKv091GxuEvmzA5EvMQsQAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAEhlbGxvIGZyb20gRWRnZWxlc3MgU3lzdGVtcyEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAADMEAAAYbPmffGRNtL5ViDWxe44+/k3th7PC6R186hE9iAfQQG6Mf45s2kKBHhCNJLxC+YMlyrm/FGWWa5SdRXVyhdki9DGtp/Gtnj07btzjqn+YZfht2Mp6Yi/SjGCyeT6esHHdPHZl9I+/HuyYncR0NmwjNd90PWsnCCM5B37x9yk5skGAEYQAAAFBQ8RA/8AAwAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAVAAAAAAAAAOcAAAAAAAAAhT4pjzt83iiwZJPQb7Ktb5VmqX/qbT3mYjayrxo1FQ8AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAANyeKnxvlI8XR040p/xD7QMPfBVj8bq932NAyC4OVKjFAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAgAEAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAALfu0NIBSRpqy0gLANRkDypIPqV0QxpsHMiF7hlJ4u/wAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAACUvwcuuQwxcA6DGLDisjrN+8+W7wNfIcn3jBSpBnjtMvTFo9YytDsa3oOhXc2cWu3vCHl+Ylv1nbbD2b4FHencIAAAAQIDBAUGBwgJCgsMDQ4PEBESExQVFhcYGRobHB0eHwUAXg4AAC0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQpNSUlFOERDQ0JKYWdBd0lCQWdJVVNqZGJmR0dsZzN6YXFQUGp3eXhFcnI3Qk9TRXdDZ1lJS29aSXpqMEVBd0l3CmNERWlNQ0FHQTFVRUF3d1pTVzUwWld3Z1UwZFlJRkJEU3lCUWJHRjBabTl5YlNCRFFURWFNQmdHQTFVRUNnd1IKU1c1MFpXd2dRMjl5Y0c5eVlYUnBiMjR4RkRBU0JnTlZCQWNNQzFOaGJuUmhJRU5zWVhKaE1Rc3dDUVlEVlFRSQpEQUpEUVRFTE1Ba0dBMVVFQmhNQ1ZWTXdIaGNOTWpNd01USTJNVEV3TlRJNVdoY05NekF3TVRJMk1URXdOVEk1CldqQndNU0l3SUFZRFZRUUREQmxKYm5SbGJDQlRSMWdnVUVOTElFTmxjblJwWm1sallYUmxNUm93R0FZRFZRUUsKREJGSmJuUmxiQ0JEYjNKd2IzSmhkR2x2YmpFVU1CSUdBMVVFQnd3TFUyRnVkR0VnUTJ4aGNtRXhDekFKQmdOVgpCQWdNQWtOQk1Rc3dDUVlEVlFRR0V3SlZVekJaTUJNR0J5cUdTTTQ5QWdFR0NDcUdTTTQ5QXdFSEEwSUFCR2NqCnlzNDhPNEZBanRySVd0cGFHTERMRFRwYjBVdUNFS1VWUXZBWTczWHdtZFRBeDNZbXp1clphckd6Nzl2T1RncDUKYm5VUGRMRW5haERyeDFETzhqZWpnZ01NTUlJRENEQWZCZ05WSFNNRUdEQVdnQlNWYjEzTnZSdmg2VUJKeWRUMApNODRCVnd2ZVZEQnJCZ05WSFI4RVpEQmlNR0NnWHFCY2hscG9kSFJ3Y3pvdkwyRndhUzUwY25WemRHVmtjMlZ5CmRtbGpaWE11YVc1MFpXd3VZMjl0TDNObmVDOWpaWEowYVdacFkyRjBhVzl1TDNZMEwzQmphMk55YkQ5allUMXcKYkdGMFptOXliU1psYm1OdlpHbHVaejFrWlhJd0hRWURWUjBPQkJZRUZKWFRNdHJoRHNUazZuQkV0SGl1ellyTwp6ZlM5TUE0R0ExVWREd0VCL3dRRUF3SUd3REFNQmdOVkhSTUJBZjhFQWpBQU1JSUNPUVlKS29aSWh2aE5BUTBCCkJJSUNLakNDQWlZd0hnWUtLb1pJaHZoTkFRMEJBUVFRZGdhS3JXR1hBMnViWmdPOTY3M2I0ekNDQVdNR0NpcUcKU0liNFRRRU5BUUl3Z2dGVE1CQUdDeXFHU0liNFRRRU5BUUlCQWdFRk1CQUdDeXFHU0liNFRRRU5BUUlDQWdFRgpNQkFHQ3lxR1NJYjRUUUVOQVFJREFnRU5NQkFHQ3lxR1NJYjRUUUVOQVFJRUFnRUNNQkFHQ3lxR1NJYjRUUUVOCkFRSUZBZ0VETUJBR0N5cUdTSWI0VFFFTkFRSUdBZ0VCTUJBR0N5cUdTSWI0VFFFTkFRSUhBZ0VBTUJBR0N5cUcKU0liNFRRRU5BUUlJQWdFRE1CQUdDeXFHU0liNFRRRU5BUUlKQWdFQU1CQUdDeXFHU0liNFRRRU5BUUlLQWdFQQpNQkFHQ3lxR1NJYjRUUUVOQVFJTEFnRUFNQkFHQ3lxR1NJYjRUUUVOQVFJTUFnRUFNQkFHQ3lxR1NJYjRUUUVOCkFRSU5BZ0VBTUJBR0N5cUdTSWI0VFFFTkFRSU9BZ0VBTUJBR0N5cUdTSWI0VFFFTkFRSVBBZ0VBTUJBR0N5cUcKU0liNFRRRU5BUUlRQWdFQU1CQUdDeXFHU0liNFRRRU5BUUlSQWdFTE1COEdDeXFHU0liNFRRRU5BUUlTQkJBRgpCUTBDQXdFQUF3QUFBQUFBQUFBQU1CQUdDaXFHU0liNFRRRU5BUU1FQWdBQU1CUUdDaXFHU0liNFRRRU5BUVFFCkJnQ0Fid1VBQURBUEJnb3Foa2lHK0UwQkRRRUZDZ0VCTUI0R0NpcUdTSWI0VFFFTkFRWUVFSDNOaWwxZitycHEKT0tBSmhkTjg3QXN3UkFZS0tvWklodmhOQVEwQkJ6QTJNQkFHQ3lxR1NJYjRUUUVOQVFjQkFRSC9NQkFHQ3lxRwpTSWI0VFFFTkFRY0NBUUVBTUJBR0N5cUdTSWI0VFFFTkFRY0RBUUgvTUFvR0NDcUdTTTQ5QkFNQ0EwZ0FNRVVDCklRREN6Ly9KNVV4bXViRjNoWVJlR0lyL1laNUlnT2dEVkZybUJ4dzFkMm5sR3dJZ0hWc2UybjRabnBOaXc2bTAKVWEyalBTWVRQWlRKWlB1K1Uwd1Y1d0syQXVBPQotLS0tLUVORCBDRVJUSUZJQ0FURS0tLS0tCi0tLS0tQkVHSU4gQ0VSVElGSUNBVEUtLS0tLQpNSUlDbGpDQ0FqMmdBd0lCQWdJVkFKVnZYYzI5RytIcFFFbkoxUFF6emdGWEM5NVVNQW9HQ0NxR1NNNDlCQU1DCk1HZ3hHakFZQmdOVkJBTU1FVWx1ZEdWc0lGTkhXQ0JTYjI5MElFTkJNUm93R0FZRFZRUUtEQkZKYm5SbGJDQkQKYjNKd2IzSmhkR2x2YmpFVU1CSUdBMVVFQnd3TFUyRnVkR0VnUTJ4aGNtRXhDekFKQmdOVkJBZ01Ba05CTVFzdwpDUVlEVlFRR0V3SlZVekFlRncweE9EQTFNakV4TURVd01UQmFGdzB6TXpBMU1qRXhNRFV3TVRCYU1IQXhJakFnCkJnTlZCQU1NR1VsdWRHVnNJRk5IV0NCUVEwc2dVR3hoZEdadmNtMGdRMEV4R2pBWUJnTlZCQW9NRVVsdWRHVnMKSUVOdmNuQnZjbUYwYVc5dU1SUXdFZ1lEVlFRSERBdFRZVzUwWVNCRGJHRnlZVEVMTUFrR0ExVUVDQXdDUTBFeApDekFKQmdOVkJBWVRBbFZUTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFTlNCLzd0MjFsWFNPCjJDdXpweHc3NGVKQjcyRXlER2dXNXJYQ3R4MnRWVExxNmhLazZ6K1VpUlpDbnFSN3BzT3ZncUZlU3hsbVRsSmwKZVRtaTJXWXozcU9CdXpDQnVEQWZCZ05WSFNNRUdEQVdnQlFpWlF6V1dwMDBpZk9EdEpWU3YxQWJPU2NHckRCUwpCZ05WSFI4RVN6QkpNRWVnUmFCRGhrRm9kSFJ3Y3pvdkwyTmxjblJwWm1sallYUmxjeTUwY25WemRHVmtjMlZ5CmRtbGpaWE11YVc1MFpXd3VZMjl0TDBsdWRHVnNVMGRZVW05dmRFTkJMbVJsY2pBZEJnTlZIUTRFRmdRVWxXOWQKemIwYjRlbEFTY25VOURQT0FWY0wzbFF3RGdZRFZSMFBBUUgvQkFRREFnRUdNQklHQTFVZEV3RUIvd1FJTUFZQgpBZjhDQVFBd0NnWUlLb1pJemowRUF3SURSd0F3UkFJZ1hzVmtpMHcraTZWWUdXM1VGLzIydWFYZTBZSkRqMVVlCm5BK1RqRDFhaTVjQ0lDWWIxU0FtRDV4a2ZUVnB2bzRVb3lpU1l4ckRXTG1VUjRDSTlOS3lmUE4rCi0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KLS0tLS1CRUdJTiBDRVJUSUZJQ0FURS0tLS0tCk1JSUNqekNDQWpTZ0F3SUJBZ0lVSW1VTTFscWROSW56ZzdTVlVyOVFHemtuQnF3d0NnWUlLb1pJemowRUF3SXcKYURFYU1CZ0dBMVVFQXd3UlNXNTBaV3dnVTBkWUlGSnZiM1FnUTBFeEdqQVlCZ05WQkFvTUVVbHVkR1ZzSUVOdgpjbkJ2Y21GMGFXOXVNUlF3RWdZRFZRUUhEQXRUWVc1MFlTQkRiR0Z5WVRFTE1Ba0dBMVVFQ0F3Q1EwRXhDekFKCkJnTlZCQVlUQWxWVE1CNFhEVEU0TURVeU1URXdORFV4TUZvWERUUTVNVEl6TVRJek5UazFPVm93YURFYU1CZ0cKQTFVRUF3d1JTVzUwWld3Z1UwZFlJRkp2YjNRZ1EwRXhHakFZQmdOVkJBb01FVWx1ZEdWc0lFTnZjbkJ2Y21GMAphVzl1TVJRd0VnWURWUVFIREF0VFlXNTBZU0JEYkdGeVlURUxNQWtHQTFVRUNBd0NRMEV4Q3pBSkJnTlZCQVlUCkFsVlRNRmt3RXdZSEtvWkl6ajBDQVFZSUtvWkl6ajBEQVFjRFFnQUVDNm5Fd01ESVlaT2ovaVBXc0N6YUVLaTcKMU9pT1NMUkZoV0dqYm5CVkpmVm5rWTR1M0lqa0RZWUwwTXhPNG1xc3lZamxCYWxUVll4RlAyc0pCSzV6bEtPQgp1ekNCdURBZkJnTlZIU01FR0RBV2dCUWlaUXpXV3AwMGlmT0R0SlZTdjFBYk9TY0dyREJTQmdOVkhSOEVTekJKCk1FZWdSYUJEaGtGb2RIUndjem92TDJObGNuUnBabWxqWVhSbGN5NTBjblZ6ZEdWa2MyVnlkbWxqWlhNdWFXNTAKWld3dVkyOXRMMGx1ZEdWc1UwZFlVbTl2ZEVOQkxtUmxjakFkQmdOVkhRNEVGZ1FVSW1VTTFscWROSW56ZzdTVgpVcjlRR3prbkJxd3dEZ1lEVlIwUEFRSC9CQVFEQWdFR01CSUdBMVVkRXdFQi93UUlNQVlCQWY4Q0FRRXdDZ1lJCktvWkl6ajBFQXdJRFNRQXdSZ0loQU9XLzVRa1IrUzlDaVNEY05vb3dMdVBSTHNXR2YvWWk3R1NYOTRCZ3dUd2cKQWlFQTRKMGxySG9NcytYbzVvL3NYNk85UVd4SFJBdlpVR09kUlE3Y3ZxUlhhcUk9Ci0tLS0tRU5EIENFUlRJRklDQVRFLS0tLS0KAA=="

/*
	This is a collection of verification snippets which we later can use to build the full verification.
	This is mainly done to understand how the crypto works.
*/

// 4.1.2.4.16
// Use given public key & signature over SGXQuote4Header + SGXReport2.
func TestQuoteSignatureVerificationBasic(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	rawQuote, err := base64.StdEncoding.DecodeString(rawQuoteBlob)
	require.NoError(err)

	parsedQuote, err := types.ParseQuote(rawQuote)
	require.NoError(err)

	signature := parsedQuote.Signature.Signature
	publicKey := parsedQuote.Signature.PublicKey // This key is called attestKey in Intel's code.

	headerBytes := parsedQuote.Header.Marshal()
	reportBytes := parsedQuote.Body.Marshal()
	toVerify := sha256.Sum256(append(headerBytes[:], reportBytes[:]...)) // Quote header + TDReport

	// It's crypto time!
	key := new(ecdsa.PublicKey)
	key.Curve = elliptic.P256()

	// Either construct the key manually...
	// https://github.com/intel/SGXDataCenterAttestationPrimitives/blob/c057b236790834cf7e547ebf90da91c53c7ed7f9/QuoteVerification/QVL/Src/AttestationLibrary/src/OpensslHelpers/KeyUtils.cpp#L63
	key.X = new(big.Int).SetBytes(publicKey[:32])
	key.Y = new(big.Int).SetBytes(publicKey[32:64])

	// Or use this one trick Go does not want you to know!
	// elliptic.Unmarshal expects the input to be *65* bytes for our curve.
	// We only have 64 bytes. So, what's the extra byte?
	// Well, apparently to look like valid ASN.1, you need to prepend a 0x04 (OCTET STRING).

	// key.X, key.Y = elliptic.Unmarshal(key, append([]byte{0x04}, publicKey[:]...))

	assert.NotNil(key.X)
	assert.NotNil(key.Y)

	// However, the ASN.1 trick does not seem to work for ecdsa.VerifyASN1.
	// The function seems to expect an ASN.1 SEQUENCE.
	// No idea what that looks like... but Intel does this: https://github.com/intel/SGXDataCenterAttestationPrimitives/blob/c057b236790834cf7e547ebf90da91c53c7ed7f9/QuoteVerification/QVL/Src/AttestationLibrary/src/OpensslHelpers/SignatureVerification.cpp#L76
	// So let's do the same here, too.
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])

	verified := ecdsa.Verify(key, toVerify[:], r, s)
	assert.True(verified)
}

// 4.1.2.4.13
// Then, the public key from above is verified/authenticated over the QEReportCertificationData.
// A hash over the AttestKey and the QEAuthData is added as report data to the QE EnclaveReport, which then is signed with the PCK (?).
func TestQEReportAttestKeyReportDataConcat(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	rawQuote, err := base64.StdEncoding.DecodeString(rawQuoteBlob)
	require.NoError(err)

	parsedQuote, err := types.ParseQuote(rawQuote)
	require.NoError(err)

	qeReport, ok := parsedQuote.Signature.CertificationData.Data.(types.QEReportCertificationData)
	require.True(ok)

	attestKeyData := parsedQuote.Signature.PublicKey
	qeAuthData := qeReport.QEAuthData.Data
	concat := append(attestKeyData[:], qeAuthData...)
	concatSHA256 := sha256.Sum256(concat)

	assert.Equal(concatSHA256[:], qeReport.EnclaveReport.ReportData[:32])
}

// 4.1.2.4.12
func TestQEReportSignatureVerification(t *testing.T) {
	assert := assert.New(t)
	require := require.New(t)

	rawQuote, err := base64.StdEncoding.DecodeString(rawQuoteBlob)
	require.NoError(err)

	parsedQuote, err := types.ParseQuote(rawQuote)
	require.NoError(err)

	qeReport, ok := parsedQuote.Signature.CertificationData.Data.(types.QEReportCertificationData)
	require.True(ok)
	pemChain, ok := qeReport.CertificationData.Data.([]byte)
	require.True(ok)

	// Parse certificate chain
	pckLeafPEM, rest := pem.Decode(pemChain)
	assert.NotEmpty(pckLeafPEM)
	assert.NotEmpty(rest)
	pckLeaf, err := x509.ParseCertificate(pckLeafPEM.Bytes)
	assert.NoError(err)

	pckIntermediatePEM, rest := pem.Decode(rest)
	assert.NotEmpty(pckIntermediatePEM)
	assert.NotEmpty(rest)

	//pckIntermediate, err := x509.ParseCertificate(pckIntermediatePEM.Bytes)
	//assert.NoError(err)

	rootCAPEM, rest := pem.Decode(rest)
	assert.NotEmpty(rootCAPEM)
	assert.Equal([]byte{0x0}, rest) // C terminated string with 0x0 byte

	//rootCA, err := x509.ParseCertificate(rootCAPEM.Bytes)
	//assert.NoError(err)

	enclaveReport := qeReport.EnclaveReport
	marshaledEnclaveReport := enclaveReport.Marshal()
	marshaledEnclaveReportHash := sha256.Sum256(marshaledEnclaveReport[:])

	pckLeafECDSAPublicKey := pckLeaf.PublicKey.(*ecdsa.PublicKey)
	signature := qeReport.Signature

	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])

	verified := ecdsa.Verify(pckLeafECDSAPublicKey, marshaledEnclaveReportHash[:], r, s)
	assert.True(verified)
}
