package app_test

import (
	"context"
	"testing"

	"github.com/docker/docker/api/types"
	"github.com/golang/mock/gomock"
	"github.com/latonaio/container-image-sweeper-kube/internal/app"
	"github.com/latonaio/container-image-sweeper-kube/internal/mock_app"
)

func Test_process(t *testing.T) {
	table := []struct {
		name        string
		cmdArgs     *app.CommandArgs
		prepareMock func(mock *mock_app.MockIDockerClient)
	}{
		{
			name: "check for PruneImages",
			cmdArgs: &app.CommandArgs{
				PruneImages: true,
			},
			prepareMock: func(mock *mock_app.MockIDockerClient) {
				mock.EXPECT().ImageList(gomock.Any(), gomock.Any()).Return(nil, nil)
				mock.EXPECT().ImagesPrune(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "check for PruneBuildCache",
			cmdArgs: &app.CommandArgs{
				PruneBuildCache: true,
			},
			prepareMock: func(mock *mock_app.MockIDockerClient) {
				mock.EXPECT().ImageList(gomock.Any(), gomock.Any()).Return(nil, nil)
				mock.EXPECT().BuildCachePrune(gomock.Any(), gomock.Any())
			},
		},
		{
			name: "check for deletion",
			cmdArgs: &app.CommandArgs{
				RetainCount: 3,
				KeepLatest:  true,
			},
			prepareMock: func(mock *mock_app.MockIDockerClient) {
				mock.EXPECT().ImageList(gomock.Any(), gomock.Any()).Return([]types.ImageSummary{
					// latest は常に保持
					{
						RepoTags: []string{
							"latonaio/test-image:latest",
						},
						Created: 1612100000000, // 一番古いが、latest が保持されるか？
					},

					// 以下の 2 つが削除されるか？
					{
						RepoTags: []string{
							"latonaio/test-image:test1",
						},
						Created: 1612105200000,
					},
					{
						RepoTags: []string{
							"latonaio/test-image:test2",
						},
						Created: 1612105200001,
					},

					// 以下 3 つは保持されるか？
					{
						RepoTags: []string{
							"latonaio/test-image:test3",
						},
						Created: 1612105200002,
					},
					{
						RepoTags: []string{
							"latonaio/test-image:test4",
						},
						Created: 1612105200003,
					},
					{
						RepoTags: []string{
							"latonaio/test-image:test5",
						},
						Created: 1612105200004,
					},
				}, nil)

				mock.EXPECT().ImageRemove(gomock.Any(), "latonaio/test-image:test1", gomock.Any())
				mock.EXPECT().ImageRemove(gomock.Any(), "latonaio/test-image:test2", gomock.Any())
			},
		},
	}

	for _, unit := range table {
		t.Run(unit.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mock := mock_app.NewMockIDockerClient(ctrl)
			unit.prepareMock(mock)
			m := app.NewDockerManagerForTest(mock)

			a := app.NewAppForTest(unit.cmdArgs)
			if err := app.ProcessForTest(a, context.Background(), m); err != nil {
				t.Error(err)
			}
		})
	}
}
