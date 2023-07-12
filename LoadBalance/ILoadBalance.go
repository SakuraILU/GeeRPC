package loadbalance

type ILoadBalance interface {
	Refresh() error
	Get() (string, error)
	GetAll() ([]string, error)
}
