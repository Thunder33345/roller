package wrapper

import "ranker"

type WrappedPermissible struct {
	permissible ranker.Permissible
	judge       ranker.Judge
}

func (w WrappedPermissible) HasPermission(node string) bool {
	return w.judge.HasPermission(w.permissible, node)
}

func (w WrappedPermissible) HasPermissionWithLevel(node string, level int) bool {
	return w.judge.HasPermissionWithLevel(w.permissible, node, level)
}

func (w WrappedPermissible) IsHigherLevel(subject ranker.Permissible) bool {
	return w.judge.IsHigherLevel(w.permissible, subject)
}
