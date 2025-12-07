package usecases

import (
	"testing"
	"time"
)

const privateKeyPem = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQDt61qIkyM/CIw0
v34a1qchCyOBYTuZRTXo6SvPQ9YpA2HoCRCcKfp/N337lv11xAGZUB2YnzE0LwKs
1UcIuyOgyYDbjRkEJdTkV4NFQ2Px/1EjtupZ+lm3SS/dfcYbj9W3VSspApPQ/XkP
adruVR8XxFuu/dJAj3ZIVTCpsfn9cD5VEAKmQFnVmNQZuhgHFlZ20Hy4IJpV/Qky
rmFQIqeRKqQ/DNUkOAu+Rwy/wbBfH+gtpyovA8i3B/GJZTKvbcpjTuKTnsj3N3Nd
naj5wYN6cuWFKbJp6aXAvrTbhtiGh4tnZtyKhjyAcgvayeV5N9tKYby0yPGR1PW6
kRIzM0YfAgMBAAECggEABPhKntm7/cAW9a8eWj8rpJQP/M7kKNJ6StA8GwtGuPqa
G/e8ghaaZffpyMyhpMkgY2x6AcspgvaMbsHRxwvptZ0f9PYglKaZqN9vHY5H0zFL
J5zVjmWdZCfCOTU8Yy0BAOBlk2i7X707vyet7BaZHKz8YU5qCvE0PlSRPKo8F6Ar
TM3yTr8AWp+DLcEtIFdRTudi80bdoENXh7R9jS1zRIQzSYmtCNF8MRy3TyZ/LHGi
2/S7JLe7rYjTEaMroM1QLxPbDPgcHCe0Ve2mcqtV+kQV+yOsV2sVqSyV/Jt2IPat
32fWVMzmUJvSFAMStf+6sf51G6cxI20eJSFlXwR5cQKBgQD9A0xaIq63EhDU5oZP
olrbyfc0vCAurHXJAd14p+Vzj5DgVJg7BC39LtcmprGep+tDe48stZFFb0W7zc/Z
L1VtyyGa2yQEBNLdmTP6F7tA1GXsPzXYM2vN3f540TAKpvyZnRGxKiR6fn+phhhq
+8ar6Suyl1u9gCFb33tAKOG2BwKBgQDwum/cvmpY+EojPZuWvgjGXi4EumZaP1HO
3HnoKITtU8rLobhgEDv7v7MjpSkQ+o93tJnvXZuUD2Dbs0+MvjB8r1tx0/Awa8oI
dy902pXA8lV7kZB3df28KWPT2LNq3zRlJvNrNePjAttBQ3YQosJmOMp/J6fFz9SX
wJ4MTmopKQKBgARiukAVudGSjpgiJtHajpigt5hCaoxkkOYbEiu1PVTzeB9rV/gt
6l4pIbGZ0hpd7sYMrj6oJwx9EUhgGOo619A/ZSW6BrXLH5yXuz7qimRlSh7+OYC1
43h+EJsnhR2qJ1bCUjwv7tHwv2XA3Ut9ccQpFojR9tUiE3H0Pb6u9rqhAoGBALWZ
8Al3HINBy6wKLfXqJnR/V/f5Jn2uhuinKtAYwS7Ip5Q2zACsPpQMaffaAMDuRIzp
kbchxtxLPaZ//uMOF0X4g+O7HtdoeWEpiIN+4rpMFnDBv1pfiKsKDmUidTeKatxk
Jf4bCW+YGA+D9O1X24+CCEEkiUyRHK/ef1yJS00BAoGAANHq9G9aJOHWXZk+M0uE
DKnNQVkV4t3lvm95IWbIyIcF4UekW2+S5OizSB+nG/j6sC7F5v/J9vaf8Erf9eab
ApMsjVLuIv4in6Bl30fpmj30RRC7/AXxtjC49HgQfko7PKREfeDoyVYTFahHL23L
De8WmllT+uQU0qciOHGcnOI=
-----END PRIVATE KEY-----`

func Test_generateGetTokenSignature(t *testing.T) {

	type args struct {
		privateKey string
		xTimestamp string
		clientId   string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test with valid private_key",
			args: args{
				privateKey: privateKeyPem,
				xTimestamp: time.Now().UTC().Format(time.RFC3339),
				clientId:   "BRN-0230-1696820291289",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateGetTokenSignature(tt.args.privateKey, tt.args.xTimestamp, tt.args.clientId)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateGetTokenSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			t.Logf("Signature: %s", got)
		})
	}
}

//func Test_generateRequestSignature(t *testing.T) {
//	type args struct {
//		privateKey  string
//		httpMethod  string
//		requestPath string
//		xTimestamp  string
//		accessToken string
//		jsonBody    []byte
//	}
//
//	tests := []struct {
//		name    string
//		args    args
//		want    string
//		wantErr bool
//	}{
//		{
//			name: "Test with valid private_key",
//			args: args{
//				privateKey:  privateKeyPem,
//				httpMethod:  "POST",
//				requestPath: "/snap/v1.1/emoney/bank-account-inquiry",
//				accessToken: "eyJhbGciOiJSUzI1NiJ9.eyJleHAiOjE3NjM5OTc4MjcsImlzcyI6IkRPS1UiLCJjbGllbnRJZCI6IkJSTi0wMjMwLTE2OTY4MjAyOTEyODkifQ.gvxwrKe6CBqjDvYoZmB3U_NMUrYhyjYGquuX8ssxpTvG0UcsuOW9BNMsj1FzCj3oLVGL75n3vu-z4M1lKEhOnP6-_Sz_XJeCT0UDTGJRKarpowaWaEdQcLQotM0EbGj0qAukOgKRe64_eNtxzOa-WtFeNJQEkJHG9T26jsjcqbSbyhVyjWasI-z6CNPJhutFZ9D3rTX70tDqYGjD7roeLb1U8S1yVicfqELRjH_loP22JJMLDPnzO0xf88X2FajsNV7y1h87MXjkhsUspOi2_ZPInxc2PY04P7X40iXMn4Gh77wykUqj2U2ZKYewP6cn48b8uLrTBWw0Nq2kJyO19g",
//				xTimestamp:  time.Now().UTC().Format(time.RFC3339),
//				jsonBody:    []byte(`{"partnerReferenceNo":"hsjkans284b2he54","customerNumber":"628115678890","amount":{"value":"200000.00","currency":"IDR"},"beneficiaryAccountNumber":"8377388292","additionalInfo":{"beneficiaryBankCode":"014","beneficiaryAccountName":"FHILEA HERMANUS","senderCountryCode":"ID"}}`),
//			},
//			wantErr: false,
//		},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			got, err := generateRequestSignature(tt.args.privateKey, tt.args.httpMethod, tt.args.requestPath, tt.args.accessToken, tt.args.xTimestamp, tt.args.jsonBody)
//			if (err != nil) != tt.wantErr {
//				t.Errorf("generateRequestSignature() error = %v, wantErr %v", err, tt.wantErr)
//				return
//			}
//			t.Logf("Signature: %s", got)
//		})
//	}
//}
