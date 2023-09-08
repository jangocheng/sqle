package register

import (
	"context"
	"fmt"

	dmsV1 "github.com/actiontech/dms/pkg/dms-common/api/dms/v1"
	pkgHttp "github.com/actiontech/dms/pkg/dms-common/pkg/http"
)

// RegisterDMSProxyTarget 向DMS注册反向代理，将proxyPrefix开头的请求转发到自身服务
// eg: name = sqle; url = http://10.1.2.1:5432; proxyPrefix = /v1/sqle 表示要求DMS将/v1/sqle开头的请求转发到sqle服务所在地址 http://10.1.2.1:5432
func RegisterDMSProxyTarget(ctx context.Context, dmsAddr, targetName, targetAddr, version string, proxyUrlPrefixs []string) error {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}
	reqBody := struct {
		DMSProxyTarget *dmsV1.DMSProxyTarget `json:"dms_proxy_target"`
	}{
		DMSProxyTarget: &dmsV1.DMSProxyTarget{
			Name:            targetName,
			Addr:            targetAddr,
			Version:         version,
			ProxyUrlPrefixs: proxyUrlPrefixs,
		},
	}

	reply := &dmsV1.RegisterDMSProxyTargetReply{}

	dmsUrl := fmt.Sprintf("%s%s", dmsAddr, dmsV1.GetProxyRouter())

	if err := pkgHttp.POST(ctx, dmsUrl, header, reqBody, reply); err != nil {
		return fmt.Errorf("failed to register dms proxy target %v: %v", dmsUrl, err)
	}
	if reply.Code != 0 {
		return fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Msg)
	}

	return nil
}

// RegisterDMSPlugin 向DMS注册校验插件，DMS会在对应操作时调用插件进行校验。注意：注册的插件接口需要服务自己实现
func RegisterDMSPlugin(ctx context.Context, dmsAddr string, plugin *dmsV1.Plugin) error {
	header := map[string]string{
		"Authorization": pkgHttp.DefaultDMSToken,
	}
	reqBody := struct {
		Plugin *dmsV1.Plugin `json:"plugin"`
	}{
		Plugin: plugin,
	}

	reply := &dmsV1.RegisterDMSPluginReply{}

	dmsUrl := fmt.Sprintf("%s%s", dmsAddr, dmsV1.GetPluginRouter())

	if err := pkgHttp.POST(ctx, dmsUrl, header, reqBody, reply); err != nil {
		return fmt.Errorf("failed to register dms plugin %v: %v", dmsUrl, err)
	}
	if reply.Code != 0 {
		return fmt.Errorf("http reply code(%v) error: %v", reply.Code, reply.Msg)
	}

	return nil
}