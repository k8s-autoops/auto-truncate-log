package main

import "testing"

func Test_isActiveLogFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test-1",
			args: args{
				path: "hello.log",
			},
			want: true,
		},
		{
			name: "test-2",
			args: args{
				path: "hello.log.gz",
			},
			want: false,
		},
		{
			name: "test-3",
			args: args{
				path: "ROT110.hello.log",
			},
			want: false,
		},
		{
			name: "test-4",
			args: args{
				path: ".hello.10.log",
			},
			want: false,
		},
		{
			name: "test-5",
			args: args{
				path: ".hello-2020-10-11.log",
			},
			want: false,
		},
		{
			name: "test-6",
			args: args{
				path: ".hello-2020-10-11-123.log",
			},
			want: false,
		},
		{
			name: "test-7",
			args: args{
				path: ".hello.log-2020-10-11-123",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isActiveLogFile(tt.args.path); got != tt.want {
				t.Errorf("isActiveLogFile() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_isHistoryLogFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test-1",
			args: args{
				path: "hello.log",
			},
			want: false,
		},
		{
			name: "test-2",
			args: args{
				path: "hello.log.gz",
			},
			want: true,
		},
		{
			name: "test-3",
			args: args{
				path: "ROT110.hello.log",
			},
			want: true,
		},
		{
			name: "test-4",
			args: args{
				path: ".hello.10.log",
			},
			want: true,
		},
		{
			name: "test-5",
			args: args{
				path: ".hello-2020-10-11.log",
			},
			want: true,
		},
		{
			name: "test-6",
			args: args{
				path: ".hello-2020-10-11-123.log",
			},
			want: true,
		},
		{
			name: "test-7",
			args: args{
				path: ".hello.log-2020-10-11-123",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isHistoryLogFile(tt.args.path); got != tt.want {
				t.Errorf("isHistoryLogFile() = %v, want %v", got, tt.want)
			}
		})
	}
}
