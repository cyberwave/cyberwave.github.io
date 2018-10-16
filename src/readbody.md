1. 读取 `GET` 和 `POST` 方法
```golang

func ReadReq(req *http.Request, v interface{}) error {
	var buf []byte
	var err error
	defer req.Body.Close()
	if req.Method == "POST" {
		buf, err = ioutil.ReadAll(req.Body)
		if err != nil {
			return err
		}
	} else if req.Method == "GET" {
		v = req.URL.RawQuery
		return nil
	} else {
		return errors.New("unsupported method:" + req.Method)
	}

	temp := api.RealNameAuthInfo{}
	err = json.Unmarshal(buf, &temp)
	if err != nil {
		logger.Errorf("json unmarshal err:%s", err.Error())
	} else {
		logger.Debugf("temp:%+v", temp)
	}
	logger.Debugf("parse http body:%s", string(buf))
	fmt.Printf("req.Body:%+v\n", string(buf))
	return nil
}
```

2. 读取 `json` body
```golang
func ReadRequest(r *http.Request, v interface{}) error {
  decoder := json.NewDecoder(r.Body)
  decode.UserNumber()
  
  defer r.Body.Close()
  
  return decoder.Decode(v)
}
```
