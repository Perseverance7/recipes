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
				input: "cheese",
			},
			want: []string{"cheese"},
		},
		{
			name: "Two ingredient",
			args: args{
				input: "cheese, potato",
			},
			want: []string{"cheese", "potato"},
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
				input: "parmesan cheese",
			},
			want: []string{"parmesan cheese"},
		},
		{
			name: "Upper case ingredients",
			args: args{
				input: "pasTA, tomaTo , Broccoli",
			},
			want: []string{"pasta", "tomato", "broccoli"},
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
