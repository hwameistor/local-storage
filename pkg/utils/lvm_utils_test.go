package utils

import (
	"testing"
)

func TestConvertLVMBytesToNumeric(t *testing.T) {
	type args struct {
		lvmbyte string
	}
	var lvmb = "10240B"
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			args:    args{lvmbyte: lvmb},
			want:    10240,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertLVMBytesToNumeric(tt.args.lvmbyte)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertLVMBytesToNumeric() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ConvertLVMBytesToNumeric() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertNumericToLVMBytes(t *testing.T) {
	type args struct {
		num int64
	}
	var num = int64(4194305)
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			args: args{num: num},
			want: "5242880B",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ConvertNumericToLVMBytes(tt.args.num); got != tt.want {
				t.Errorf("ConvertNumericToLVMBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNumericToLVMBytes(t *testing.T) {
	type args struct {
		bytes int64
	}
	var num = int64(4194305)
	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			args: args{bytes: num},
			want: 5242880,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NumericToLVMBytes(tt.args.bytes); got != tt.want {
				t.Errorf("NumericToLVMBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}
