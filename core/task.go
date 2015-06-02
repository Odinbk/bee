package core

type Runner interface {
	Run() (response interface{}, err error)
}

type Task struct {
	ID int
}

func (t Task) Run() (response interface{}, err error) {
	return response, err
}

type HttpTask struct {
	ID     int
	Engine Engine
	URL    string
	Header map[string]string
	Params map[string]string
	SSL    bool
}

func (t HttpTask) Run() (response interface{}, err error) {
	response, err = t.Engine.Get(t.URL, t.Params, t.Header, t.SSL)
	return response, err
}
