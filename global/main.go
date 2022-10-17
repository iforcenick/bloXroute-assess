package global

const MQ_SERVER_URL = "amqp://guest:guest@localhost:5672/"
const MSG_QUEUE_NAME = "message-queue"

const ADD_ITEM = "add"
const REMOVE_ITEM = "remove"
const GET_ITEM = "get"
const GET_ALL_ITEMS = "getall"

type MessageBody struct {
	Command string
	Key     string
	Value   string
}

type KeyValuePair struct {
	Key   string
	Value string
}

type AddItemResponse struct {
	Success bool
}

type RemoveItemResponse struct {
	Success bool
}

type GetItemResponse struct {
	Success bool
	Item    string
}

type GetAllItemsResponse struct {
	Items []KeyValuePair
}
