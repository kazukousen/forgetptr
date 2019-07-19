package a

type foo struct {
	name    string
	counter int
}

func (f *foo) _(name string) {
	f.name = name
	f.counter++
}

func (f foo) _(name string) {
	f.name = name // want "this statement can not modify the value"
	f.counter++   // want "this statement can not modify the value"
}

func (f foo) _(name string) foo {
	f.name = name
	f.counter++
	return f
}

func (f foo) _(name string) (int, error) {
	f.name = name // want "this statement can not modify the value"
	return f.counter, nil
}
