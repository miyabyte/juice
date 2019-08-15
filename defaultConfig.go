package juice

import "time"

const HeartbeatCheckInterval = time.Second * 15
const HeartbeatIdleTime = time.Second * 30
const ReadBufferSize = 1024
const WriteBufferSize = 1024
const RpcChat = "127.0.0.1:50051"
