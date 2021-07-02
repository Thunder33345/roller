package wrapper

import (
	"ranker"
)

type Ranker struct {
	provider  Provider
	processor CachedProcessor
	judge     ranker.Judge
}

func (r *Ranker) GetPermissible(uid string) (WrappedPermissible, error) {
	p, e := r.processor.Process(uid)
	if e != nil {
		return WrappedPermissible{}, e
	}
	return WrappedPermissible{permissible: p, judge: r.judge}, e
}
