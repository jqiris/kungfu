package cloud

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type ObsClient struct {
	client *cos.Client
}

func NewObsClient(cfg config.TecentOBS) *ObsClient {
	u, _ := url.Parse(cfg.BulletUrl)
	su, _ := url.Parse(cfg.ServiceUrl)
	b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
	// 1.永久密钥
	client := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretId,
			SecretKey: cfg.SecretKey,
		},
	})
	return &ObsClient{client}
}

func (c *ObsClient) PutString(key, val string) error {
	f := strings.NewReader(val)
	_, err := c.client.Object.Put(context.Background(), key, f, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ObsClient) PutFile(key, path string) error {
	// 2.通过本地文件上传对象
	_, err := c.client.Object.PutFromFile(context.Background(), key, path, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ObsClient) CreateDir(key string) error {
	f := strings.NewReader("")
	_, err := c.client.Object.Put(context.Background(), key, f, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ObsClient) Delete(key string) error {
	_, err := c.client.Object.Delete(context.Background(), key)
	if err != nil {
		return err
	}
	return nil
}

func (c *ObsClient) DeleteDir(key string) error {
	var marker string
	opt := &cos.BucketGetOptions{
		Prefix:  key,
		MaxKeys: 1000,
	}
	isTruncated := true
	for isTruncated {
		opt.Marker = marker
		v, _, err := c.client.Bucket.Get(context.Background(), opt)
		if err != nil {
			return err
		}
		for _, content := range v.Contents {
			_, err = c.client.Object.Delete(context.Background(), content.Key)
			if err != nil {
				return err
			}
		}
		isTruncated = v.IsTruncated
		marker = v.NextMarker
	}
	return nil
}
