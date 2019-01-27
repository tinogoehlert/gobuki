package commands

import (
	"reflect"
	"testing"
)

func TestSoundSequenceCmd(t *testing.T) {
	type args struct {
		SoundSequence SoundSequence
	}
	tests := []struct {
		name string
		args args
		want Command
	}{
		{
			name: "test ON sequence",
			args: args{0x00},
			want: Command{
				ID:   0x04,
				Data: []byte{0x00},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SoundSequenceCmd(tt.args.SoundSequence); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SoundSequenceCmd() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCommand_Serialize(t *testing.T) {
	type fields struct {
		ID   byte
		Data []byte
	}
	tests := []struct {
		name   string
		fields fields
		want   []byte
	}{
		{
			name: "serialize SoundSequence",
			fields: fields{
				ID:   0x04,
				Data: []byte{0x00},
			},
			want: []byte{0xAA, 0x55, 0x03, 0x04, 0x01, 0x00, 0x06},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Command{
				ID:   tt.fields.ID,
				Data: tt.fields.Data,
			}
			if got := c.Serialize(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Command.Serialize() = %v, want %v", got, tt.want)
			}
		})
	}
}
