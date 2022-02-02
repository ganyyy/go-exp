package fuzzing

import (
	"testing"
	"unicode/utf8"
)

func FuzzReverse(f *testing.F) {
	var testCases = []string{
		"hello world!",
		"123",
		//"甘",
	}

	for _, tc := range testCases {
		f.Add(tc)
	}

	f.Fuzz(func(t *testing.T, orig string) {
		t.Logf("origin %q", orig)
		rev, err1 := Reverse(orig)
		if err1 != nil {
			if err1 != ErrNotValidString {
				t.Logf("Input %q got unexcept error %v", orig, err1)
			}
			return
		}
		doubleRev, err2 := Reverse(rev)
		if err2 != nil {
			if err2 != ErrNotValidString {
				t.Logf("DoubleReverse %q got unexcept error %v", rev, err2)
			}
			return
		}

		if orig != doubleRev {
			t.Errorf("Before %q, After %q", orig, doubleRev)
		}
		if utf8.ValidString(orig) && !utf8.ValidString(doubleRev) {
			t.Errorf("Reverse %q produce invalid UTF-8 string %q", orig, rev)
		}
	})
}

func TestReverse(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Normal",
			args: args{
				s: "123",
			},
			want:    "321",
			wantErr: false,
		},

		{
			name: "RuneNormal",
			args: args{
				s: "甘2",
			},
			want:    "2甘",
			wantErr: false,
		},

		{
			name: "RuneError",
			args: args{
				s: "\x94",
			},
			want:    "\x94",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Reverse(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("Reverse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Reverse() got = %v, want %v", got, tt.want)
			}
		})
	}
}
