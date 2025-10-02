package maps

type GoMap struct {
	dict map[string]string
}

func NewGoMap() *GoMap {
	return &GoMap{dict: make(map[string]string)}
}

func (g *GoMap) Put(key string, value string) {
	g.dict[key] = value
}

func (g *GoMap) Get(key string) (string, bool) {
	value, found := g.dict[key]
	return value, found
}

func (g *GoMap) Contains(key string) bool {
	_, found := g.dict[key]
	return found
}
