package grpcinfra

import (
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ConnectionManager struct {
	mu          sync.Mutex
	connections map[string]*grpc.ClientConn
}

func NewConnectionManager() *ConnectionManager {
	return &ConnectionManager{
		connections: make(map[string]*grpc.ClientConn),
	}
}

func (cm *ConnectionManager) GetConnection(targetUrl string) (*grpc.ClientConn, error) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	var opts []grpc.DialOption

	if conn, exists := cm.connections[targetUrl]; exists {
		return conn, nil
	}

	opts = append(opts,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

	// Create gRPC client connection
	conn, err := grpc.NewClient(targetUrl, opts...)
	if err != nil {
		return nil, err
	}

	cm.connections[targetUrl] = conn
	return conn, nil
}

func (cm *ConnectionManager) CloseAll() {
	cm.mu.Lock()
	defer cm.mu.Unlock()

	for target, conn := range cm.connections {
		conn.Close()
		delete(cm.connections, target)
	}
}
