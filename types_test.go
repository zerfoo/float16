package float16

import "testing"

func TestFromInt(t *testing.T) {
	tests := []struct {
		name string
		i    int
		want Float16
	}{
		{"FromInt(0)", 0, PositiveZero},
		{"FromInt(1)", 1, 0x3C00},
		{"FromInt(-1)", -1, 0xBC00},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromInt(tt.i); got != tt.want {
				t.Errorf("FromInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromInt32(t *testing.T) {
	tests := []struct {
		name string
		i    int32
		want Float16
	}{
		{"FromInt32(0)", 0, PositiveZero},
		{"FromInt32(1)", 1, 0x3C00},
		{"FromInt32(-1)", -1, 0xBC00},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromInt32(tt.i); got != tt.want {
				t.Errorf("FromInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFromInt64(t *testing.T) {
	tests := []struct {
		name string
		i    int64
		want Float16
	}{
		{"FromInt64(0)", 0, PositiveZero},
		{"FromInt64(1)", 1, 0x3C00},
		{"FromInt64(-1)", -1, 0xBC00},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := FromInt64(tt.i); got != tt.want {
				t.Errorf("FromInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt(t *testing.T) {
	tests := []struct {
		name string
		f    Float16
		want int
	}{
		{"ToInt(0)", PositiveZero, 0},
		{"ToInt(1.0)", 0x3C00, 1},
		{"ToInt(-1.0)", 0xBC00, -1},
		{"ToInt(1.9)", 0x3F33, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.ToInt(); got != tt.want {
				t.Errorf("ToInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt32(t *testing.T) {
	tests := []struct {
		name string
		f    Float16
		want int32
	}{
		{"ToInt32(0)", PositiveZero, 0},
		{"ToInt32(1.0)", 0x3C00, 1},
		{"ToInt32(-1.0)", 0xBC00, -1},
		{"ToInt32(1.9)", 0x3F33, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.ToInt32(); got != tt.want {
				t.Errorf("ToInt32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToInt64(t *testing.T) {
	tests := []struct {
		name string
		f    Float16
		want int64
	}{
		{"ToInt64(0)", PositiveZero, 0},
		{"ToInt64(1.0)", 0x3C00, 1},
		{"ToInt64(-1.0)", 0xBC00, -1},
		{"ToInt64(1.9)", 0x3F33, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.f.ToInt64(); got != tt.want {
				t.Errorf("ToInt64() = %v, want %v", got, tt.want)
			}
		})
	}
}
