package juice

import "time"

const HeartbeatCheckInterval = time.Second * 60 * 3 * 9999
const HeartbeatIdleTime = time.Second * 60 * 3 * 9999
const ReadBufferSize = 1024
const WriteBufferSize = 1024
const RpcChat = "127.0.0.1:50051"
