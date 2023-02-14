/*
 * +----------------------------------------------------------------------
 *  | kungfu [ A FAST GAME FRAMEWORK ]
 *  +----------------------------------------------------------------------
 *  | Copyright (c) 2023-2029 All rights reserved.
 *  +----------------------------------------------------------------------
 *  | Licensed ( http:www.apache.org/licenses/LICENSE-2.0 )
 *  +----------------------------------------------------------------------
 *  | Author: jqiris <1920624985@qq.com>
 *  +----------------------------------------------------------------------
 */

package cloud

import (
	"context"
	"net/http"
	"net/url"
	"strings"

	"github.com/jqiris/kungfu/v2/config"
	"github.com/jqiris/kungfu/v2/logger"
	cdn "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cdn/v20180606"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentyun/cos-go-sdk-v5"
)

type ObsClient struct {
	cosClient *cos.Client
	cdnClient *cdn.Client
}

func NewObsClient(cfg config.TecentOBS) *ObsClient {
	u, _ := url.Parse(cfg.BulletUrl)
	su, _ := url.Parse(cfg.ServiceUrl)
	b := &cos.BaseURL{BucketURL: u, ServiceURL: su}
	// 1.永久密钥
	cosClient := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  cfg.SecretId,
			SecretKey: cfg.SecretKey,
		},
	})
	//cdn 操作
	credential := common.NewCredential(
		cfg.SecretId,
		cfg.SecretKey,
	)
	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "cdn.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	cdnClient, err := cdn.NewClient(credential, "", cpf)
	if err != nil {
		logger.Error(err)
	}
	return &ObsClient{
		cosClient: cosClient,
		cdnClient: cdnClient,
	}
}

func (c *ObsClient) PutString(key, val string) error {
	f := strings.NewReader(val)
	_, err := c.cosClient.Object.Put(context.Background(), key, f, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ObsClient) PutFile(key, path string) error {
	// 2.通过本地文件上传对象
	_, err := c.cosClient.Object.PutFromFile(context.Background(), key, path, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ObsClient) CreateDir(key string) error {
	f := strings.NewReader("")
	_, err := c.cosClient.Object.Put(context.Background(), key, f, nil)
	if err != nil {
		return err
	}
	return nil
}

func (c *ObsClient) Delete(key string) error {
	_, err := c.cosClient.Object.Delete(context.Background(), key)
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
		v, _, err := c.cosClient.Bucket.Get(context.Background(), opt)
		if err != nil {
			return err
		}
		for _, content := range v.Contents {
			_, err = c.cosClient.Object.Delete(context.Background(), content.Key)
			if err != nil {
				return err
			}
		}
		isTruncated = v.IsTruncated
		marker = v.NextMarker
	}
	return nil
}

func (c *ObsClient) FlushCdn(cdnUrl, flushType string) error {
	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := cdn.NewPurgePathCacheRequest()

	request.Paths = common.StringPtrs([]string{cdnUrl})
	request.FlushType = common.StringPtr(flushType)

	// 返回的resp是一个PurgePathCacheResponse的实例，与请求对象对应
	response, err := c.cdnClient.PurgePathCache(request)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		return err
	}
	if err != nil {
		return err
	}
	// 输出json格式的字符串回包
	logger.Infof("FlushCdn resp:%s", response.ToJsonString())
	return nil
}
