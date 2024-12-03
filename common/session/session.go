// Package session provides functions for sessions of incoming requests.
package session

import (
	//"context"
	"math/rand"
	"fmt"

	"github.com/gzjjjfree/hello/common/net"
	//"github.com/gzjjjfree/hello/common/protocol"
)

// ID of a session.
type ID uint32

// NewID generates a new ID. The generated ID is high likely to be unique, but not cryptographically secure.
// NewID 生成一个新的 ID。生成的 ID 很可能是唯一的，但不具备加密安全性
// The generated ID will never be 0.
func NewID() ID {
	fmt.Println("in common-session-session.go func NewID")
	for {
		id := ID(rand.Uint32())
		if id != 0 {
			return id
		}
	}
}

// Inbound is the metadata of an inbound connection.
type Inbound struct {
	// Source address of the inbound connection.
	Source net.Destination
	// Gateway address
	Gateway net.Destination
	// Tag of the inbound proxy that handles the connection.
	Tag string
}

// Outbound is the metadata of an outbound connection.
type Outbound struct {
	// Target address of the outbound connection.
	Target net.Destination
	// Gateway address
	Gateway net.Address
}

// SniffingRequest controls the behavior of content sniffing.
type SniffingRequest struct {
	OverrideDestinationForProtocol []string
	Enabled                        bool
	MetadataOnly                   bool
}

// Content is the metadata of the connection content.
type Content struct {
	// Protocol of current content.
	Protocol string

	SniffingRequest SniffingRequest

	Attributes map[string]string

	SkipDNSResolve bool
}

// Sockopt is the settings for socket connection.
type Sockopt struct {
	// Mark of the socket connection.
	Mark int32
}

// SetAttribute attachs additional string attributes to content.
func (c *Content) SetAttribute(name string, value string) {
	fmt.Println("in common-session-session.go func  (c *Content) SetAttribute")
	if c.Attributes == nil {
		c.Attributes = make(map[string]string)
	}
	c.Attributes[name] = value
}

// Attribute retrieves additional string attributes from content.
func (c *Content) Attribute(name string) string {
	fmt.Println("in common-session-session.go func (c *Content) Attribute: ", c.Attributes[name])
	if c.Attributes == nil {
		return ""
	}
	return c.Attributes[name]
}
