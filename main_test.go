package getchromecookies

import (
	"fmt"
	"reflect"
	"testing"
)

func TestChrome_GetCookie(t *testing.T) {
	tests := []struct {
		name    string
		c       *Chrome
		wantS   []byte
		wantErr bool
	}{
		{
			name: "1",
			c: &Chrome{
				Path:    locateChrome(),
				DataDir: `C:\Users\xmdhs\AppData\Local\Google\Chrome\User Data`,
				Web:     "https://www.mcbbs.net",
			},
			wantS:   nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotS, err := tt.c.GetCookie()
			fmt.Println(gotS)
			if (err != nil) != tt.wantErr {
				t.Errorf("Chrome.GetCookie() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotS, tt.wantS) {
				t.Errorf("Chrome.GetCookie() = %v, want %v", gotS, tt.wantS)
			}
		})
	}
}
