package types

type FriendStatus string

const (
	Friend           FriendStatus = "friend"
	BlockBoth        FriendStatus = "block_both"
	BlockFirstSecond FriendStatus = "block_first_second"
	BlockSecondFirst FriendStatus = "block_second_first"
)
