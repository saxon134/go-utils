package saRedis

func (r Redis) Do(command string, args ...interface{}) (res interface{}, err error) {
	c := r.Pool.Get()
	defer c.Close()

	res, err = c.Do(command, args...)
	return
}
