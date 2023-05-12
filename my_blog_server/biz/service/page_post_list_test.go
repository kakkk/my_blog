package service

import (
	"testing"

	"my_blog/biz/common/config"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
)

func Test_getIDsByPage(t *testing.T) {
	patches := gomonkey.ApplyFunc(config.GetPageListSize, func() int {
		return 5
	})
	defer patches.Reset()
	type args struct {
		ids  []int64
		page int64
	}
	tests := []struct {
		name        string
		args        args
		want        []int64
		wantHasMore bool
	}{
		{
			name: "normal page 1",
			args: args{
				ids:  []int64{1, 2, 3, 4, 5, 6, 7, 8, 9},
				page: 1,
			},
			want:        []int64{1, 2, 3, 4, 5},
			wantHasMore: true,
		},
		{
			name: "len less than 5, page 1",
			args: args{
				ids:  []int64{1, 2, 3, 4},
				page: 1,
			},
			want:        []int64{1, 2, 3, 4},
			wantHasMore: false,
		},
		{
			name: "len is 5",
			args: args{
				ids:  []int64{1, 2, 3, 4, 5},
				page: 1,
			},
			want:        []int64{1, 2, 3, 4, 5},
			wantHasMore: false,
		},
		{
			name: "normal page 2",
			args: args{
				ids:  []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				page: 2,
			},
			want:        []int64{6, 7, 8, 9, 10},
			wantHasMore: true,
		},
		{
			name: "normal page 3",
			args: args{
				ids:  []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13},
				page: 3,
			},
			want:        []int64{11, 12, 13},
			wantHasMore: false,
		},
		{
			name: "just page 3",
			args: args{
				ids:  []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15},
				page: 3,
			},
			want:        []int64{11, 12, 13, 14, 15},
			wantHasMore: false,
		},
		{
			name: "page out of index",
			args: args{
				ids:  []int64{1, 2, 3, 4, 5, 6},
				page: 100,
			},
			want:        []int64{},
			wantHasMore: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotHasMore := getIDsByPage(tt.args.ids, tt.args.page)
			assert.EqualValues(t, tt.want, got)
			assert.Equal(t, tt.wantHasMore, gotHasMore)
		})
	}
}

func Fuzz_getIDsByPage(f *testing.F) {
	f.Add([]byte{1, 2, 3}, int64(1))
	f.Fuzz(func(t *testing.T, bytes []byte, page int64) {
		var ids []int64
		for _, b := range bytes {
			ids = append(ids, int64(b))
		}
		assert.NotPanics(t, func() {
			_, _ = getIDsByPage(ids, page)
		})
	})
}
