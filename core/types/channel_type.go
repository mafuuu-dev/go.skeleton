package types

import "strconv"

type ChannelName string
type ChannelDataType string

type ChannelPersonal struct {
	id   *int64
	uuid *string
}

func NewChannelPersonalId(ID int64) ChannelPersonal {
	return ChannelPersonal{id: &ID, uuid: nil}
}

func NewChannelPersonalUUID(UUID string) ChannelPersonal {
	return ChannelPersonal{id: nil, uuid: &UUID}
}

func (personal *ChannelPersonal) Get() *string {
	if personal.id != nil {
		id := strconv.FormatInt(*personal.id, 10)
		return &id
	}

	return personal.uuid
}

type Channel struct {
	Name     ChannelName
	DataType ChannelDataType
	Personal *ChannelPersonal
}

func (channel *Channel) WithPersonal(personal ChannelPersonal) Channel {
	channel.Personal = &personal
	return *channel
}
