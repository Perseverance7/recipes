package service

import (
	"reflect"
	"testing"
)

func Test_extractWords(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "One ingredient",
			args: args{
				input: "сыр",
			},
			want: []string{"сыр"},
		},
		{
			name: "Two ingredient",
			args: args{
				input: "сыр, укроп",
			},
			want: []string{"сыр", "укроп"},
		},
		{
			name: "No ingredient",
			args: args{
				input: "",
			},
			want: nil,
		},
		{
			name: "Long ingredient",
			args: args{
				input: "сыр пармезан",
			},
			want: []string{"сыр пармезан"},
		},
		{
			name: "Upper case ingredients",
			args: args{
				input: "Паста, картофЕль , морКовЬ",
			},
			want: []string{"паста", "картофель", "морковь"},
		},


	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractWords(tt.args.input); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
