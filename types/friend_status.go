package types

type FriendStatus string

const (
	Friend           FriendStatus = "friend"
	BlockFirstSecond FriendStatus = "block_first_second"
	BlockSecondFirst FriendStatus = "block_second_first"
	UnFriend         FriendStatus = "unfriend"
)
