package addressutil

import "testing"

func Test_decode(t *testing.T) {

	t.Log(IsBech32Address("0x355de3D14d343458115B424ee7C3CbdB0068E04A"))
	t.Log(IsHexAddress("0x355de3D14d343458115B424ee7C3CbdB0068E04A"))

	t.Log(IsBech32Address("lat1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2ydaxq4"))
	t.Log(IsHexAddress("lat1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2ydaxq4"))

	t.Log(IsBech32Address("atp1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2amt7l6"))
	t.Log(IsHexAddress("atp1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2amt7l6"))

	t.Log(ConvertEthAddress("0x355de3D14d343458115B424ee7C3CbdB0068E04A").Hex())
	t.Log(ConvertEthAddress("lat1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2ydaxq4").Hex())
	t.Log(ConvertEthAddress("atp1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2amt7l6").Hex())
}
