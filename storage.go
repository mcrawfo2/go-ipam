package ipam

// Storage is a interface to store ipam objects.
type Storage interface {
	CreatePrefix(prefix Prefix, tenantid string) (Prefix, error)
	ReadPrefix(prefix string, tenantid string) (Prefix, error)
	ReadAllPrefixes(tenantid string) ([]Prefix, error)
	UpdatePrefix(prefix Prefix, tenantid string) (Prefix, error)
	DeletePrefix(prefix Prefix, tenantid string) (Prefix, error)
}

// OptimisticLockError indicates that the operation could not be executed because the dataset to update has changed in the meantime.
// clients can decide to read the current dataset and retry the operation.
type OptimisticLockError struct {
	msg string
}

func (o OptimisticLockError) Error() string {
	return o.msg
}

func newOptimisticLockError(msg string) OptimisticLockError {
	return OptimisticLockError{msg: msg}
}
