package ranker

type Ranker struct {
	provider  GroupProvider
	processor Processor
	judge     Judge
}

func (r *Ranker) Process(p RawPermissible) (Permissible, error) {
	return r.processor.Process(p, r.provider)
}

func (r *Ranker) HasPermission(p Permissible, node string) bool {
	return r.judge.HasPermission(p, node)
}

func (r *Ranker) HasPermissionWithLevel(p Permissible, node string, level int) bool {
	return r.judge.HasPermissionWithLevel(p, node, level)
}

func (r *Ranker) IsHigherLevel(target Permissible, subject Permissible) bool {
	return r.judge.IsHigherLevel(target, subject)
}
