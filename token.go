package bricked

type TokenList struct {
	TokenInfos []struct {
		TokenID      string `json:"token_id"`
		CreationTime int64  `json:"creation_time"`
		ExpiryTime   int64  `json:"expiry_time"`
		Comment      string `json:"comment"`
	} `json:"token_infos"`
}

func (c *databricks) TokenList() TokenList {
	c.assertVersion("2.0")

	var resp TokenList
	c.GET("token/list", &resp)
	return resp
}
