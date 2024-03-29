package ShowroomCGIlib

import "testing"

func TestIsAllowIp(t *testing.T) {
	type args struct {
		sip string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test191",
			args: args{"85.208.96.191" },
			want: true,	
		},
		{
			name: "test192",
			args: args{"85.208.96.192" },
			want: false,	
		},
		{
			name: "test223",
			args: args{"85.208.96.223" },
			want: false,	
		},
		{
			name: "test224",
			args: args{"85.208.96.224" },
			want: true,	
		},
	}

	err := LoadDenyIp("../DenyIp.txt")
	if err != nil {
		t.Errorf("LoadDenyIp() = %v", err)
		return
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsAllowIp(tt.args.sip); got != tt.want {
				t.Errorf("IsAllowIp() = %v, want %v", got, tt.want)
			}
		})
	}
}
