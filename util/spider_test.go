package util

import "testing"

func TestGetWsFromChrome(t *testing.T) {
	type args struct {
		url string
	}
	tests := []struct {
		name    string
		args    args
		wantWs  string
		wantErr bool
	}{
		{
			name:    "get from chrome",
			args:    args{"http://10.100.103.211:9222/json"},
			wantWs:  "ws://10.100.103.211:9222/devtools/page/25CB6E28CC7ADD88F97194A4E70120DC",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotWs, err := GetWsFromChrome(tt.args.url)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetWsFromChrome() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotWs != tt.wantWs {
				t.Errorf("GetWsFromChrome() gotWs = %v, want %v", gotWs, tt.wantWs)
			}
		})
	}
}
