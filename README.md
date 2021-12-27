# object_server
分布式对象存储

api_server需要用到的环境变量
	# object server
	LISTEN_ADDRESS=xxx.xxx.xxx.xxx:9000
	export LISTEN_ADDRESS

data_server需要用到的环境变量：
	DATA_ADDRESS=xxx.xxx.xxx.xxx:9100
	export DATA_ADDRESS

	RABBITMQ_SERVER=amqp://admin:123.com@xxx.xxx.xxx.xxx:5672
	export RABBITMQ_SERVER

	STORAGE_ROOT=/Users/liyuntang/data
	export STORAGE_ROOT

rabbitmq设置流程：
	yum install -y rabbitmq-server.noarch
	service rabbitmq-server start
	rabbitmq-plugins enable rabbitmq_management
	rabbitmqctl add_user admin 123.com
	rabbitmqctl set_user_tags admin administrator
	rabbitmqctl set_permissions -p "/" admin "." "." ".*"
	http://xxx.xxx.xxx.xxx:15672创建apiServers、dataServers两个exchange，type为fanout


功能说明：
    1、支持对象多版本操作

操作手册：
    get(不包含删除的数据)：
	1、获取说有对象列表:operation.Get("", "")
	2、获取某个对象列表:operation.Get("5-liyuntang-2.local-1636514830705816000", "")
	3、获取某个对象内容:operation.Get("5-liyuntang-2.local-1636514830705816000", "MTkyLjE2OC43NC45ODo1OTAwOS0xNjM2NjAwMDQ1NTI2MTYxMD")
    put：
	1、http://192.168.74.98:10000/objects/5-liyuntang-2.local-1636514830705816000
	2、header设置：
	    数据长度：	    req.Header.Set("Content-Length", strconv.Itoa(bit))
	    sha256加密值：  req.Header.Add("digest", secretData)
	    加密方法：	    req.Header.Set("secretMethod", "sha256")
    delete：
	1、删除某个对象的所有版本:operation.Get("5-liyuntang-2.local-1636514830705816000", "")
        2、删除某个对象的指定版本:operation.Get("5-liyuntang-2.local-1636514830705816000", "MTkyLjE2OC43NC45ODo1OTAwOS0xNjM2NjAwMDQ1NTI2MTYxMD")



api说明：

访问该对象的指定版本，默认是最新的那个，如果对象的最新版本是一个删除标记，则返回404
    GET /objects/<object_name>?version=<version_id>
    
上传数据，上传成功以后在元服务中增加一个新的版本，版本号从1开始
    PUT /objects/<object_name>  (Digest:SHA-256=<value>, Content-Length:size)
    
删除对象，软删除，在该对象的元数据中标记删除(size=0,hash="")
    DELETE /objects/<object_name>
    
所有对象的所有版本，相当于列表
    GET /versions/

指定对象的所有版本,获取某个对象的版本列表
    GET /versions/<object_name>

元数据结构：
    {
        name:string
        version:int
        size:int
        hash:string
    }

















