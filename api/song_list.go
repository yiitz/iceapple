package api

func PersonalFM() []interface{} {
	r := post(httpRoot+"/weapi/v1/radio/get", nil)
	if r == nil || int(r["code"].(float64)) != 200 {
		return nil
	}

	songs := r["data"].([]interface{})
	if len(songs) == 0 {
		return nil
	}
	return songs
}