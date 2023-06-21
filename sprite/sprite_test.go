package sprite

import (
	"reflect"
	"testing"
)

func MakeSpriteList() SpriteList {
	return SpriteList{
		new(Sprite),
		new(Sprite),
		new(Sprite),
		new(Sprite),
		new(Sprite),
	}
}

func TestSpriteList_BringForwards(t *testing.T) {
	type args struct {
		s *Sprite
	}
	l := MakeSpriteList()
	tests := []struct {
		name string
		list SpriteList
		args args
		want SpriteList
	}{
		{
			name: "normal",
			list: l,
			args: args{l[2]},
			want: []*Sprite{l[0], l[2], l[1], l[3], l[4]},
		},
		{
			name: "beginning",
			list: l,
			args: args{l[0]},
			want: l,
		},
		{
			name: "end",
			list: l,
			args: args{l[4]},
			want: []*Sprite{l[0], l[1], l[2], l[4], l[3]},
		},
		{
			name: "missing",
			list: l,
			args: args{new(Sprite)},
			want: l,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.list.BringForwards(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SpriteList.BringForwards() = %v, want %v", got, tt.want)
			}
		})
	}
}
