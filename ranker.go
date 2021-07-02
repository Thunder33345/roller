package ranker

type Ranker struct {
	provider  GroupProvider
	processor Processor
	judge     Comparator
}

func (r *Ranker) Process(p RawPermissionList) (PermissionList, error) {
	return r.processor.Process(p, r.provider)
}

func (r *Ranker) HasPermission(p PermissionList, node string) bool {
	return r.judge.HasPermission(p, node)
}

func (r *Ranker) HasPermissionWithLevel(p PermissionList, node string, level int) bool {
	return r.judge.HasPermissionWithLevel(p, node, level)
}

func (r *Ranker) IsHigherLevel(target PermissionList, subject PermissionList) bool {
	return r.judge.IsHigherLevel(target, subject)
}

type WrappedRanker struct {
	provider  DataProvider
	processor WrappedProcessor
	judge     Comparator
}

func (r *WrappedRanker) GetPermission(uid string) (WrappedPermissionList, error) {
	p, e := r.processor.Process(uid)
	if e != nil {
		return WrappedPermissionList{}, e
	}
	return WrappedPermissionList{permissible: p, judge: r.judge}, e
}
