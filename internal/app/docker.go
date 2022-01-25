package app

import (
	"context"
	"errors"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

// モック用インタフェース
type IDockerClient interface {
	ImageList(ctx context.Context, options types.ImageListOptions) ([]types.ImageSummary, error)
	ImageRemove(ctx context.Context, imageID string, options types.ImageRemoveOptions) ([]types.ImageDeleteResponseItem, error)
	ImagesPrune(ctx context.Context, pruneFilters filters.Args) (types.ImagesPruneReport, error)
	BuildCachePrune(ctx context.Context, opts types.BuildCachePruneOptions) (*types.BuildCachePruneReport, error)
}

type DockerManager struct {
	client IDockerClient
}

type TagInfo struct {
	Tag     string
	Summary *types.ImageSummary
}

// ソート用ダミー型
type tagInfosForSort []*TagInfo

func (a tagInfosForSort) Len() int      { return len(a) }
func (a tagInfosForSort) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a tagInfosForSort) Less(i, j int) bool {
	return a[j].Summary.Created < a[i].Summary.Created
}

func NewDockerManager(ops ...client.Opt) (*DockerManager, error) {
	c, err := client.NewClientWithOpts(ops...)
	if err != nil {
		return nil, err
	}

	return &DockerManager{
		client: c,
	}, nil
}

func (m *DockerManager) getImageSummaries(ctx context.Context) ([]types.ImageSummary, error) {
	return m.client.ImageList(ctx, types.ImageListOptions{All: true})
}

func (m *DockerManager) getImages(imgSummaries []types.ImageSummary) (map[string][]*TagInfo, error) {
	imgMap := map[string][]*TagInfo{}
	for _, imgSummary := range imgSummaries {
		imgSummary := imgSummary

		// タグのないイメージは無視
		if len(imgSummary.RepoTags) == 1 && imgSummary.RepoTags[0] == "<none>:<none>" {
			continue
		}

		for _, repoTag := range imgSummary.RepoTags {
			// 例: "latonaio/example-image:202201010000"
			repoTagSplit := strings.SplitN(repoTag, ":", 2)
			if len(repoTagSplit) != 2 {
				return nil, errors.New("invalid repoTag passed")
			}

			// 例: "latonaio/example-image"
			repoName := repoTagSplit[0]
			// 例: "202201010000"
			tagName := repoTagSplit[1]

			if _, ok := imgMap[repoName]; !ok {
				imgMap[repoName] = make([]*TagInfo, 0)
			}
			imgMap[repoName] = append(imgMap[repoName], &TagInfo{
				Tag:     tagName,
				Summary: &imgSummary,
			})

			// ビルド日時順 (降順) にソート
			sort.Sort(tagInfosForSort(imgMap[repoName]))
		}
	}

	return imgMap, nil
}

func (m *DockerManager) GetImages(ctx context.Context) (map[string][]*TagInfo, error) {
	imgSummaries, err := m.getImageSummaries(ctx)
	if err != nil {
		return nil, err
	}

	imgMap, err := m.getImages(imgSummaries)
	if err != nil {
		return nil, err
	}

	return imgMap, nil
}

func (m *DockerManager) DeleteImage(ctx context.Context, imageID string) error {
	_, err := m.client.ImageRemove(ctx, imageID, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	return err
}

func (m *DockerManager) PruneImages(ctx context.Context) error {
	_, err := m.client.ImagesPrune(ctx, filters.Args{})
	return err
}

func (m *DockerManager) PruneBuildCache(ctx context.Context) error {
	_, err := m.client.BuildCachePrune(ctx, types.BuildCachePruneOptions{})
	return err
}
