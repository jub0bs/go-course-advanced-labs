package bimap

import "fmt"

type Bimap struct {
	Forward, Inverse map[string]string
}

func New() *Bimap {
	return &Bimap{
		Forward: make(map[string]string),
		Inverse: make(map[string]string),
	}
}

func (bi *Bimap) Store(key string, value string) {
	k, exists := bi.Inverse[value]
	if exists { // value is already associated with k
		delete(bi.Forward, k)
	}
	v, exists := bi.Forward[key]
	if exists { // key is already associated with v
		delete(bi.Inverse, v)
	}
	bi.Forward[key] = value
	bi.Inverse[value] = key
}

func (bi *Bimap) LookupValue(key string) (string, bool) {
	v, ok := bi.Forward[key]
	return v, ok
}

func (bi *Bimap) LookupKey(value string) (string, bool) {
	k, ok := bi.Inverse[value]
	return k, ok
}

func (bi *Bimap) String() string {
	return fmt.Sprintf("bi%v", bi.Forward)
}
