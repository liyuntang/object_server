package es

func AddVersion(name, hash string, size int64) error {

	version, err:= SearchLatestVersion(name)
	if err != nil {
		return err
	}
	//fmt.Println(">>>>>>>", name, hash, version)
	return PutMetadata(name, version+1, size, hash)
}
