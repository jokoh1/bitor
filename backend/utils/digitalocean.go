package utils

import (
	"context"
	"fmt"

	"github.com/digitalocean/godo"
)

// FetchDigitalOceanSizes fetches available droplet sizes from DigitalOcean for a specific region
func FetchDigitalOceanSizes(apiKey string, region string) ([]map[string]interface{}, error) {
	client := godo.NewFromToken(apiKey)
	ctx := context.Background()

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}

	var allSizes []godo.Size
	for {
		sizes, resp, err := client.Sizes.List(ctx, opt)
		if err != nil {
			return nil, fmt.Errorf("failed to list sizes: %w", err)
		}
		allSizes = append(allSizes, sizes...)

		if resp.Links == nil || resp.Links.IsLastPage() {
			break
		}

		page, err := resp.Links.CurrentPage()
		if err != nil {
			return nil, fmt.Errorf("failed to get current page: %w", err)
		}
		opt.Page = page + 1
	}

	var result []map[string]interface{}
	for _, size := range allSizes {
		for _, availableRegion := range size.Regions {
			if availableRegion == region {
				result = append(result, map[string]interface{}{
					"slug":         size.Slug,
					"price_hourly": size.PriceHourly,
				})
				break
			}
		}
	}

	return result, nil
}
