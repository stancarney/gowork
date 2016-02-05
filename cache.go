package gowork

func Cache(ctx Context, name, key string, function func() (interface{}, error)) (interface{}, error) {

	var cache map[string]interface{}
	
	c := ctx.Get(name)
	if c == nil {
		cache = make(map[string]interface{})
		ctx.Put(name, cache)
	} else {
		cache = c.(map[string]interface{})
	}

	value := cache[key]
	if value != nil {
		return value, nil
	}

	value, err := function()
	if err != nil {
		return nil, err
	}

	if value != nil {
		cache[key] = value
	}

	return value, nil
}