package grpcinfra

type GrpcClients struct {
	UserManagementClient *UserManagementClient
	connManager          *ConnectionManager
}

func NewGrpcClients(userManagementAddr string) (*GrpcClients, error) {
	cm := NewConnectionManager()

	userManagementConn, err := cm.GetConnection(userManagementAddr)
	if err != nil {
		return nil, err
	}

	return &GrpcClients{
		UserManagementClient:   NewUserManagementClient(userManagementConn),
		connManager:   cm,
	}, nil
}

func (gc *GrpcClients) Close() {
	gc.connManager.CloseAll()
}
