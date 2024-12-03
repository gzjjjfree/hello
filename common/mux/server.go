package mux

import (
	"context"
	"io"
	"errors"
	"fmt"

	//core "github.com/gzjjjfree/hello"
	"github.com/gzjjjfree/hello/common"
	"github.com/gzjjjfree/hello/common/buf"
	//"github.com/gzjjjfree/hello/common/log"
	"github.com/gzjjjfree/hello/common/net"
	"github.com/gzjjjfree/hello/common/protocol"
	//"github.com/gzjjjfree/hello/common/session"
	"github.com/gzjjjfree/hello/app/routing"
	"github.com/gzjjjfree/hello/transport"
	"github.com/gzjjjfree/hello/transport/pipe"
)

type Server struct {
	dispatcher routing.Dispatcher
	
}

func (s *Server) Getctx() context.Context {
	return nil
}

// NewServer creates a new mux.Server.
// NewServer 创建一个新的 mux.Server
func NewServer(ctx context.Context) *Server {
	fmt.Println("in common-mux-server.go func NewServer")
	//v := core.MustFromContext(ctx) 
	var r routing.Handler
	var d routing.Dispatcher = &r
	s := &Server{
		dispatcher: d,
	}
	//d := &routing.Handler{}
	//core.RequireFeatures(ctx, func(d routing.Handler) {
	//	s.dispatcher = d
	//})
	return s
}

// Type implements common.HasType.
// 类型实现 common.HasType。
func (s *Server) Type() interface{} {
	return s.dispatcher.Type()
}

// Dispatch implements routing.Dispatcher
func (s *Server) Dispatch(ctx context.Context, dest net.Destination) (*transport.Link, error) {
	fmt.Println("in common-mux-server.go func (s *Server) Dispatch")
	if dest.Address != muxCoolAddress {
		return s.dispatcher.Dispatch(ctx, dest)
	}

	opts := pipe.OptionsFromContext(ctx)
	uplinkReader, uplinkWriter := pipe.New(opts...)
	downlinkReader, downlinkWriter := pipe.New(opts...)

	_, err := NewServerWorker(ctx, s.dispatcher, &transport.Link{
		Reader: uplinkReader,
		Writer: downlinkWriter,
	})
	if err != nil {
		return nil, err
	}

	return &transport.Link{Reader: downlinkReader, Writer: uplinkWriter}, nil
}

// Start implements common.Runnable.
func (s *Server) Start() error {
	fmt.Println("in common-mux-server.go func Start()")
	return nil
}

// Close implements common.Closable.
func (s *Server) Close() error {
	fmt.Println("in common-mux-server.go func (s *Server) Close()")
	return nil
}

type ServerWorker struct {
	dispatcher     routing.Dispatcher
	link           *transport.Link
	sessionManager *SessionManager
}

func NewServerWorker(ctx context.Context, d routing.Dispatcher, link *transport.Link) (*ServerWorker, error) {
	fmt.Println("in common-mux-server.go func NewServerWorker")
	worker := &ServerWorker{
		dispatcher:     d,
		link:           link,
		sessionManager: NewSessionManager(),
	}
	go worker.run(ctx)
	return worker, nil
}

func handle(ctx context.Context, s *Session, output buf.Writer) {
	fmt.Println("in common-mux-server.go func handle")
	writer := NewResponseWriter(s.ID, output, s.transferType)
	if err := buf.Copy(s.input, writer); err != nil {
		fmt.Println("session  ends.")
		writer.hasError = true
	}

	writer.Close()
	s.Close()
}

func (w *ServerWorker) ActiveConnections() uint32 {
	return uint32(w.sessionManager.Size())
}

func (w *ServerWorker) Closed() bool {
	return w.sessionManager.Closed()
}

func (w *ServerWorker) handleStatusKeepAlive(meta *FrameMetadata, reader *buf.BufferedReader) error {
	fmt.Println("in common-mux-server.go func (w *ServerWorker) handleStatusKeepAlive")
	if meta.Option.Has(OptionData) {
		return buf.Copy(NewStreamReader(reader), buf.Discard)
	}
	return nil
}

func (w *ServerWorker) handleStatusNew(ctx context.Context, meta *FrameMetadata, reader *buf.BufferedReader) error {
	fmt.Println("in common-mux-server.go func (w *ServerWorker) handleStatusNew")
	fmt.Println("received request for ")
	
	link, err := w.dispatcher.Dispatch(ctx, meta.Target)
	if err != nil {
		if meta.Option.Has(OptionData) {
			buf.Copy(NewStreamReader(reader), buf.Discard)
		}
		return errors.New("failed to dispatch request")
	}
	s := &Session{
		input:        link.Reader,
		output:       link.Writer,
		parent:       w.sessionManager,
		ID:           meta.SessionID,
		transferType: protocol.TransferTypeStream,
	}
	if meta.Target.Network == net.Network_UDP {
		s.transferType = protocol.TransferTypePacket
	}
	w.sessionManager.Add(s)
	go handle(ctx, s, w.link.Writer)
	if !meta.Option.Has(OptionData) {
		return nil
	}

	rr := s.NewReader(reader)
	if err := buf.Copy(rr, s.output); err != nil {
		buf.Copy(rr, buf.Discard)
		common.Interrupt(s.input)
		return s.Close()
	}
	return nil
}

func (w *ServerWorker) handleStatusKeep(meta *FrameMetadata, reader *buf.BufferedReader) error {
	fmt.Println("in common-mux-server.go func (w *ServerWorker) handleStatusKeep")
	if !meta.Option.Has(OptionData) {
		return nil
	}

	s, found := w.sessionManager.Get(meta.SessionID)
	if !found {
		// Notify remote peer to close this session.
		closingWriter := NewResponseWriter(meta.SessionID, w.link.Writer, protocol.TransferTypeStream)
		closingWriter.Close()

		return buf.Copy(NewStreamReader(reader), buf.Discard)
	}

	rr := s.NewReader(reader)
	err := buf.Copy(rr, s.output)

	if err != nil && buf.IsWriteError(err) {
		fmt.Println("failed to write to downstream writer. closing session ")

		// Notify remote peer to close this session.
		closingWriter := NewResponseWriter(meta.SessionID, w.link.Writer, protocol.TransferTypeStream)
		closingWriter.Close()

		drainErr := buf.Copy(rr, buf.Discard)
		common.Interrupt(s.input)
		s.Close()
		return drainErr
	}

	return err
}

func (w *ServerWorker) handleStatusEnd(meta *FrameMetadata, reader *buf.BufferedReader) error {
	fmt.Println("in common-mux-server.go func (w *ServerWorker) handleStatusEnd")
	if s, found := w.sessionManager.Get(meta.SessionID); found {
		if meta.Option.Has(OptionError) {
			common.Interrupt(s.input)
			common.Interrupt(s.output)
		}
		s.Close()
	}
	if meta.Option.Has(OptionData) {
		return buf.Copy(NewStreamReader(reader), buf.Discard)
	}
	return nil
}

func (w *ServerWorker) handleFrame(ctx context.Context, reader *buf.BufferedReader) error {
	fmt.Println("in common-mux-server.go func (w *ServerWorker) handleFrame")
	var meta FrameMetadata
	err := meta.Unmarshal(reader)
	if err != nil {
		return errors.New("failed to read metadata")
	}

	switch meta.SessionStatus {
	case SessionStatusKeepAlive:
		err = w.handleStatusKeepAlive(&meta, reader)
	case SessionStatusEnd:
		err = w.handleStatusEnd(&meta, reader)
	case SessionStatusNew:
		err = w.handleStatusNew(ctx, &meta, reader)
	case SessionStatusKeep:
		err = w.handleStatusKeep(&meta, reader)
	default:
		//status := meta.SessionStatus
		return errors.New("unknown status: ")
	}

	if err != nil {
		return errors.New("failed to process data")
	}
	return nil
}

func (w *ServerWorker) run(ctx context.Context) {
	fmt.Println("in common-mux-server.go func (w *ServerWorker) run")
	input := w.link.Reader
	reader := &buf.BufferedReader{Reader: input}

	defer w.sessionManager.Close()

	for {
		select {
		case <-ctx.Done():
			return
		default:
			err := w.handleFrame(ctx, reader)
			if err != nil {
				if err != io.EOF {
					fmt.Println("unexpected EOF")
					common.Interrupt(input)
				}
				return
			}
		}
	}
}
