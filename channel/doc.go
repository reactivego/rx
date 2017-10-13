// Package channel provides a multi-sender merging&multicasting channel.
//
// The channel can be used by multiple senders to simultaneously send messages
// to the channel that get merged and buffered and are then multicasted to
// multiple concurrent receivers. There is support for a replay buffer to replay
// messages received in the past to newly connecting receivers. Messages are
// timestamped, so during receiving the maximum age of the messages will cause
// old messages to be skipped.
//
// This package is using generics from jig library:
//		github.com/reactivego/channel
// The generics from channel are specialized on interface{}.
package channel

import _ "github.com/reactivego/channel"
