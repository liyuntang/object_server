package operation

func (o *objectInfo)getToken() error {
	// 获取token
	server, token, err := post(o.name, o.path, o.size, o.isBigObject)
	if err != nil {
		return err
	}
	o.apiServer = server
	o.token = token
	return nil
}

