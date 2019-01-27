package utils

import "testing"

func TestCmp(t *testing.T) {
	type args struct {
		a      int32
		b      int32
		offset int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "Cmp with positive",
			args: args{100, 109, 0},
			want: false,
		},
		{
			name: "Cmp with positive and offset",
			args: args{100, 109, 10},
			want: true,
		},
		{
			name: "Cmp with positive and offset fails",
			args: args{100, 111, 10},
			want: false,
		},
		{
			name: "Cmp with negative and offset fails",
			args: args{100, -100, 10},
			want: false,
		},
		{
			name: "Cmp with positive 2",
			args: args{1583, 1590, 10},
			want: true,
		},
		{
			name: "Cmp with positive fails",
			args: args{1573, 1590, 10},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Cmp(tt.args.a, tt.args.b, tt.args.offset); got != tt.want {
				t.Errorf("Cmp() = %v, want %v", got, tt.want)
			}
		})
	}
}
