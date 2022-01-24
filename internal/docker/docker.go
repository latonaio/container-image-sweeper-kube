package docker

import (
	"context"
	"sort"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
)

type DockerClient struct {
	client *client.Client
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

func NewDockerClient(ops ...client.Opt) (*DockerClient, error) {
	c, err := client.NewClientWithOpts(ops...)
	if err != nil {
		return nil, err
	}

	return &DockerClient{
		client: c,
	}, nil
}

func (c *DockerClient) getImageSummaries(ctx context.Context) ([]types.ImageSummary, error) {
	imgSummaries, err := c.client.ImageList(ctx, types.ImageListOptions{All: true})
	if err != nil {
		return nil, err
	}

	return imgSummaries, nil
}

func (c *DockerClient) GetImages(ctx context.Context) (map[string][]*TagInfo, error) {
	imgSummaries, err := c.getImageSummaries(ctx)
	if err != nil {
		return nil, err
	}

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
				panic("invalid repoTag passed")
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

func (c *DockerClient) DeleteImage(ctx context.Context, imageID string) error {
	_, err := c.client.ImageRemove(ctx, imageID, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	return err
}

func (c *DockerClient) PruneImages(ctx context.Context) error {
	_, err := c.client.ImagesPrune(ctx, filters.Args{})
	return err
}

func (c *DockerClient) PruneBuildCache(ctx context.Context) error {
	_, err := c.client.BuildCachePrune(ctx, types.BuildCachePruneOptions{})
	return err
}
