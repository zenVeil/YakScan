yakit.AutoInitYakit()

# 参数设置
target = cli.Text("target", cli.setHelp("输入想要测试的目标：支持 IP / IP:Port / URL/ 域名 等，可用逗号分隔，也可用换行分隔"), cli.setVerboseName("检测目标"),cli.setRequired(true))
timeoutSeconds = cli.Int("timeout_seconds", cli.setDefault(5), cli.setHelp("每个请求的超时时间"), cli.setVerboseName("请求超时"))
delaySeconds = cli.Float("delay_seconds", cli.setDefault(0.200000),cli.setHelp("参数默认值为0，如果大于0，自动设置并发为1，每次请求将会 sleep 配置的时间"),cli.setVerboseName("延迟请求"))
concurrent = cli.Int("concurrent", cli.setDefault(5),cli.setVerboseName("并发数"))
scanType = cli.StringSlice("扫描类型", cli.setMultipleSelect(true), cli.setHelp("请选择要扫描的类型, 可多选"), cli.setSelectOption("Swagger", "Swagger"), cli.setRequired(true), cli.setDefault(["Swagger"]))

cli.check()

targets = []
scanPath = [] # 扫描路径字典，维护唯一性

siteMap = sync.NewMap() # 记录已经经过目标检查的域名
abandonedMap = sync.NewMap() # 如果对应站点显示无法连接时, 应当不在扫描
# 标准 404 判断
standard404Map = sync.NewMap() # 网站标准 404 模板页面

prehttpRequestBody = `GET /{{rs}} HTTP/1.1
Host: {{params(target)}}
User-Agent:  Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36
`
httpRequestBody = `GET {{params(path)}} HTTP/1.1
Host: {{params(target)}}
User-Agent:  Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/135.0.0.0 Safari/537.36
`

dict = make(map[string][]string)
dict["Swagger"] = ['/actuator/openapi','/api/openapi','/api/openapi.json','/api/openapi.yaml','/api/openapi.yml','/api/v1/openapi.json','/openapi','/openapi.json','/openapi.yaml','/openapi.yml','/openapi/index.html','/openapi/v2','/openapi/v2/api-docs','/swagger/openapi.json','/swagger/openapi.yaml','/actuator/swagger','/health/swagger','/actuator/swagger-ui','/admin/swagger-ui.html','/api/swagger-ui','/api/swagger-ui.html','/api/swagger-ui.html/','/api/swagger-ui.json','/api/swagger-ui/oauth2-redirect.html','/api/swaggerui/','/api/v1/swagger-ui','/api/v1/swagger-ui.html','/api/v1/swagger-ui.json','/api/v2/swagger-ui','/api/v2/swagger-ui.html','/api/v2/swagger-ui.json','/api/v3/swagger-ui','/api/v3/swagger-ui.html','/api/v3/swagger-ui.json','/console/swagger-ui.html','/docs/swagger-ui','/documentation/swagger-ui','/documentation/swagger-ui.html','/internal/swagger-ui.html','/libs/swaggerui','/libs/swaggerui/','/platform-api/swagger-ui.html','/platform/swagger-ui.html','/portal/swagger-ui.html','/swagger-ui','/swagger-ui-custom.html','/swagger-ui.html','/swagger-ui.html/','/swagger-ui.html/v2/api-docs','/swagger-ui.html/v3/api-docs','/swagger-ui.json','/swagger-ui.yaml','/swagger-ui.yml','/swagger-ui/','/swagger-ui/index','/swagger-ui/index.htm','/swagger-ui/index.html','/swagger-ui/oauth2-redirect.html','/swagger-ui/static','/swagger-ui/static/index.html','/swagger-ui/swagger-ui.html','/swagger-ui/swagger.json','/swagger/swagger-ui.html','/swaggerui','/swaggerui.htm','/swaggerui.html','/user/swagger-ui.html/','/admin/swagger','/api','/api-doc','/api.html','/api.json','/api.yaml','/api.yml','/api/doc','/api/docs','/api/docs-json','/api/help','/api/schema.json','/api/swagger','/api/swagger-config.json','/api/swagger-docs','/api/swagger/','/api/swagger/docs','/api/swagger/static/index.html','/api/swagger/ui/','/api/swagger_docs','/api/v1/','/api/v1/swagger','/api/v1/swagger/','/api/v2','/api/v2/swagger','/api/v2/swagger/','/api/v3','/api/v3/swagger','/api/v3/swagger/','/apis','/backend/swagger','/developer/docs','/docs','/docs/swagger','/swagge-ui','/swagger','/swagger-config.json','/swagger-docs','/swagger-editor/','/swagger/','/swagger/index','/swagger/index.html','/swagger/static/index.html','/swagger/ui/','/swagger/ui/index','/swagger_docs','/swaggerdoc','/swaggers','/swaggeru','/swgger','/tools/swagger','/web/swagger','/api-docs','/api-docs.json','/api-docs/index.html','/api-docs/v2','/api-docs/v3','/api/api-docs','/api/apidocs','/api/spec','/api/specification','/api/swagger/api-docs','/api/swagger/v2/api-docs','/api/swagger/v3/api-docs','/api/v1/api-docs','/api/v1/apidocs','/api/v1/documentation','/api/v1/spec.yaml','/api/v2/api-docs','/api/v2/apidocs','/api/v2/documentation','/api/v3/api-docs','/api/v3/apidocs','/api/v3/documentation','/apidocs','/dev/api-docs','/developer/api-docs','/mobile/api-docs','/spec','/spec.json','/swagger-v1/api-docs','/swagger-v2/api-docs','/swagger-v3/api-docs','/swagger/api-docs','/swagger/apidocs','/swagger/v2/api-docs','/swagger/v3/api-docs','/v2/api-docs','/v3/api-docs','/v4/api-docs','/api-docs/swagger.json','/api/gateway/swagger.json','/api/rest/swagger.json','/api/swagger.json','/api/swagger.yaml','/api/swagger.yml','/api/v1/swagger.json','/api/v2/swagger.json','/api/v3/swagger.json','/apidocs/swagger.json','/devportal/swagger.json','/doc/swagger.json','/docs/swagger.json','/documentation/swagger.json','/documentation/swagger.yaml','/documentation/swagger.yml','/external/swagger.json','/integration/swagger.json','/openapi/swagger.json','/openapi/swagger.yaml','/rest/swagger.json','/schema/swagger.json','/services/swagger.yaml','/spec/swagger.yaml','/swagger.config.json','/swagger.json','/swagger.yaml','/swagger.yml','/swagger/static/swagger.json','/swagger/swagger.json','/swagger/swagger.yaml','/swagger/v1/swagger.json','/v1/swagger.json','/api/explorer','/explorer','/redoc','/api/swagger-resources','/api/swagger/configuration/ui','/swagger-resources','/swagger-resources/configuration/oauth2','/swagger-resources/configuration/security','/swagger-resources/configuration/security/','/swagger-resources/configuration/ui','/swagger-resources/configuration/ui/','/swagger/configuration/security','/swagger/configuration/ui']



counter = 0

# 用于参数的校验以及相关设置
validateParams = func(){
    
    # 处理目标站点信息
    target = str.TrimSpace(target)

    if target == "" {
        yakit.Error("目标为空，请输入合理的目标")
        return false
    } else if str.Contains(target, "\n") || str.Contains(target, ",") {
        targets = x.Map(str.ParseStringToHosts(str.Join(str.Split(target, "\n"), ",")), func(i){return str.TrimSpace(i)})
    } else {
        targets = [target]
    }

    if delaySeconds > 0 {
        concurrent = 1
    }

    if concurrent <= 0 {
        concurrent = 5
    }

    if timeoutSeconds < 0 {
        timeoutSeconds = 5
    }

    return true
}

valiStatus = validateParams()
swg = sync.NewSizedWaitGroup(concurrent)
lock = sync.NewLock()
enableTableOnce = sync.NewOnce()

defer swg.Wait()
defer func(){yakit.SetProgress(1.0)}()

# 将扫描的结果输出到表格中进行展示
addTable = func(url, statusCode, contentLength, title) {
    
    enableTableOnce.Do(func(){
        yakit.EnableTable("扫描结果", ["URL", "StatusCode", "ContentLength", "Title"])
    })

    data = make(map[string]var)
    data["URL"] = url
    data["StatusCode"] = statusCode
    data["ContentLength"] = contentLength
    data["Title"] = title
    yakit.TableData("扫描结果", data)
}

counter = 0
# 对目标站点发起扫描
handleResult = func(host, path, ishttps){
    counter = counter + 1
    yakit.StatusCard("Requests", sprint(counter))
    swg.Add()
    go func{
        defer swg.Done()
        defer func{err = recover(); if err != nil {yakit.Error("Recover from panic: %v", err)}}
        defer func{ if delaySeconds > 0 { sleep(delaySeconds) } }
        resp, req, err = poc.HTTP(httpRequestBody, poc.timeout(timeoutSeconds), poc.https(ishttps), poc.params({"target": host, "path": path}))
        

        if err != nil && sprint(err) != "<nil>" {
            err = sprint(err)
            println(err)
            if str.MatchAnyOfRegexp(err, "no such host", "reset by remote peer") {
                abandonedMap.Store(host, err)
            }
            return
        }

        rspIns, _ := poc.ParseBytesToHTTPResponse(resp)
        if rspIns != nil {
            yakit.Info("[%v] %v: %v", rspIns.StatusCode, host, path)

            url, _ = str.ExtractURLFromHTTPRequestRaw(req, ishttps)
            _, body = str.SplitHTTPHeadersAndBodyFromPacket(resp)
            title = str.ExtractTitle(body)
            _, check404 = standard404Map.Load(host)

            if check404 &&  rspIns.StatusCode == 404{
                return
            } else {
                addTable(url.String(), rspIns.StatusCode, len(body), title)
                return
            }
            
            // 啥都没有的话，不是404就报
            if rspIns.StatusCode != 404 {
                addTable(url.String(), rspIns.StatusCode, len(body), title)
            }
        }
    }
}

checkForAddr = func(params, prePath, ishttps){

    # 向指定网站发送随机请求, 校验是否返回一些 404 页面, 将其进行记录, 提高检测率, 每个站点仅会检测一次
    siteKey = params["host"] + prePath
    _, ok = siteMap.Load(siteKey)

    if ishttps && params["host"].HasSuffix(":443") {
        params["host"], _ = str.CutSuffix(params["host"], ":443")
    } else if !ishttps && params["host"].HasSuffix(":80") {
        params["host"], _ = str.CutSuffix(params["host"], ":80")
    }

    if !ok {
        resp, req, err = poc.HTTP(prehttpRequestBody, poc.timeout(timeoutSeconds), poc.https(ishttps), poc.params({"target": params["host"]}))

        if err != nil && sprint(err) != "<nil>" {
            yakit.Info("Target: %v 无法连接：%v", params["host"], err)
            return
        }

        rspIns, _ = poc.ParseBytesToHTTPResponse(resp)

        if rspIns != nil && rspIns.StatusCode == 404{
            standard404Map.Store(params["host"], true)
        }
        siteMap.Store(siteKey, true)
    }

     # 根据需要扫描的类型的选择对应的扫描字典进行扫描
    for Type in scanType {
        for scanApiPath in dict[Type] {
            err, ok = abandonedMap.Load(params["host"])
            if ok {
                yakit.Info("Target: %v 无法连接：%v", params["host"], err)
                break
            }
            path = params["path"] + scanApiPath

            handleResult(params["host"], path, ishttps)
        }
    }
}

extractFlows = func(target) {
    scanPath = [] # 重置扫描路径字典，避免因为多个目标重复路径造成漏扫

    flows := yakit.QueryHTTPFlowsByKeyword(target)

    pattern = `(?:https?://)?[^/]+(/[^?#]*)`
    prePath = re.MustCompile(pattern).FindStringSubmatch(target)
    if len(prePath) > 1 {prePath = str.TrimSuffix(prePath[1], "/")} else {prePath = ""}

    for flow in flows{
        # 获取指定目标站点下所有历史记录中包含的路径
        disPath = str.Split(str.TrimPrefix(flow.Path, prePath), "/")[1]

        # 排除 .js 、.html 这类内容
        if !str.Contains(disPath, ".") {
            # 对提取的路径进行处理
            if str.Contains(disPath, "?") {disPath =  str.Split(disPath, "?")[0]}

            # 确保扫描路径的开头为 "/"
            if len(disPath) > 1 && !str.StartsWith(disPath, "/"){disPath = "/" + disPath}else {disPath = "/"}

            if disPath not in scanPath{ 
                scanPath.Push(disPath)
                params = make(map[string]string)

                params["path"] = prePath + disPath
                params["host"] = str.ExtractHostPort(flow.Url)

                checkForAddr(params, prePath, flow.IsHTTPS)
            }
        }
    }
}

if valiStatus{
    for target in targets{
        # 遍历每一个目标站点
        extractFlows(target)
    }
}

