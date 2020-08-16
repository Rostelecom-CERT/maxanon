package storage

type Data struct {
	IP                string `bson:"ip"`
	Anonymous         bool   `bson:"anonymous"`
	AnonymousVPN      bool   `bson:"anonymous_vpn"`
	IsHostingProvider bool   `bson:"is_hosting_provider"`
	IsPublicProxy     bool   `bson:"is_public_proxy"`
	IsTorExitNode     bool   `bson:"is_tor_exit_node"`
}

type Database interface {
	Open(url string) error
	Insert(data Data) error
	InsertBulk(bulk []interface{}) error
	Get(ip string) (*Data, error)
	Exist(collName string) (bool, error)
}
