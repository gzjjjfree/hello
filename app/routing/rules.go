package routing

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	core "github.com/gzjjjfree/hello"
	"github.com/gzjjjfree/hello/common"
	"github.com/gzjjjfree/hello/common/session"

	//"github.com/gzjjjfree/hello/proxy/vmess/encoding"
	"github.com/gzjjjfree/hello/app/dispatcher"
	"github.com/gzjjjfree/hello/common/buf"
	"github.com/gzjjjfree/hello/common/net"
	"github.com/gzjjjfree/hello/features"
	//"github.com/gzjjjfree/hello/proxy/vmess/outbound"
	"github.com/gzjjjfree/hello/transport"
	"github.com/gzjjjfree/hello/transport/pipe"
)

type Dispatcher interface {
	features.Feature

	// Dispatch returns a Ray for transporting data for the given request.
	// Dispatch 返回一个用于传输给定请求的数据的 Ray
	Dispatch(ctx context.Context, dest net.Destination) (*transport.Link, error)
}

// DispatcherType returns the type of Dispatcher interface. Can be used to implement common.HasType.
// DispatcherType 返回 Dispatcher 接口的类型。可用于实现 common.HasType。
// v2ray:api:stable
func DispatcherType() interface{} {
	return (*Dispatcher)(nil)
}

type Handler struct {
	//access          sync.RWMutex
	//clients  *clientsConfig
	//tag      string
	ctx context.Context
}

//type CloserHandler struct {
//    Handler
//}

//func (c *CloserHandler) Close() error {
//    return c.Handler.Close()
//}

func (d *Handler) Dispatch(ctx context.Context, destination net.Destination) (*transport.Link, error) {
	fmt.Println("in app-routing-rules.go func (d *DefaultDispatcher) Dispatch destination is: ", destination)
	if !destination.IsValid() {
		panic("Dispatcher: Invalid destination.")
	}
	ob := &session.Outbound{
		Target: destination,
	}
	ctx = session.ContextWithOutbound(ctx, ob)

	inbound, outbound := d.getLink(ctx)
	content := session.ContentFromContext(ctx)
	if content == nil {
		content = new(session.Content)
		ctx = session.ContextWithContent(ctx, content)
	}
	sniffingRequest := content.SniffingRequest
	switch {
	case !sniffingRequest.Enabled:
		fmt.Println("in app-routing-rules.go func (d *DefaultDispatcher) Dispatch !sniffingRequest.Enabled:")
		go d.routedDispatch(ctx, outbound, destination)
	case destination.Network != net.Network_TCP:
		fmt.Println("in app-routing-rules.go func (d *DefaultDispatcher) Dispatch destination.Network != net.Network_TCP")
		// Only metadata sniff will be used for non tcp connection
		// 非 TCP 连接仅使用元数据嗅探
		result, err := sniffer(ctx, nil, true)
		if err == nil {
			content.Protocol = result.Protocol()
			if shouldOverride(result, sniffingRequest.OverrideDestinationForProtocol) {
				domain := result.Domain()
				fmt.Println("sniffed domain: ", domain)
				destination.Address = net.ParseAddress(domain)
				ob.Target = destination
			}
		}
		go d.routedDispatch(ctx, outbound, destination)
	default:
		fmt.Println("in app-routing-rules.go func (d *DefaultDispatcher) Dispatch default:")
		go func() {
			cReader := &cachedReader{
				reader: outbound.Reader.(*pipe.Reader),
			}
			outbound.Reader = cReader
			result, err := sniffer(ctx, cReader, sniffingRequest.MetadataOnly)
			if err == nil {
				content.Protocol = result.Protocol()
			}
			if err == nil && shouldOverride(result, sniffingRequest.OverrideDestinationForProtocol) {
				domain := result.Domain()
				fmt.Println("sniffed domain: ", domain)
				destination.Address = net.ParseAddress(domain)
				ob.Target = destination
				fmt.Println("in app-routing-rules.go func (d *DefaultDispatcher) Dispatch ob.Target: ", ob.Target)
			}
			d.routedDispatch(ctx, outbound, destination)
		}()
	}
	return inbound, nil
}

func (d *Handler) getLink(ctx context.Context) (*transport.Link, *transport.Link) {
	fmt.Println("in app-routing-rules.go func (d *DefaultDispatcher) getLink")
	opt := pipe.OptionsFromContext(ctx)
	// 设置缓冲区限制
	uplinkReader, uplinkWriter := pipe.New(opt...)
	downlinkReader, downlinkWriter := pipe.New(opt...)

	inboundLink := &transport.Link{
		Reader: downlinkReader,
		Writer: uplinkWriter,
	}

	outboundLink := &transport.Link{
		Reader: uplinkReader,
		Writer: downlinkWriter,
	}

	//sessionInbound := session.InboundFromContext(ctx)

	return inboundLink, outboundLink
}

func (d *Handler) routedDispatch(ctx context.Context, link *transport.Link, destination net.Destination) {
	 fmt.Println("in app-routing-rules.go func (d *DefaultDispatcher) routedDispatch")
	// var handler outbound.Handler
// 
	// if forcedOutboundTag := session.GetForcedOutboundTagFromContext(ctx); forcedOutboundTag != "" {
		// fmt.Println("in app-dispatcher-default.go func (d *DefaultDispatcher) routedDispatch forcedOutboundTag is: ", forcedOutboundTag)
		// ctx = session.SetForcedOutboundTagToContext(ctx, "")
		// if h := d.ohm.GetHandler(forcedOutboundTag); h != nil {
			// fmt.Println("taking platform initialized detour [", forcedOutboundTag, "] for [", destination, "]")
			// handler = h
		// } else {
			// fmt.Println("non existing tag for platform initialized detour: ", forcedOutboundTag)
			// common.Close(link.Writer)
			// common.Interrupt(link.Reader)
			// return
		// }
	// } else if d.router != nil {
		// if route, err := d.router.PickRoute(routing_session.AsRoutingContext(ctx)); err == nil {
			// tag := route.GetOutboundTag()
			// fmt.Println("in app-dispatcher-default.go func (d *DefaultDispatcher) routedDispatch Tag is: ", tag)
			// if h := d.ohm.GetHandler(tag); h != nil {
				// fmt.Println("taking detour [", tag, "] for [", destination, "]")
				// handler = h
			// } else {
				// fmt.Println("non existing tag: ", tag)
			// }
		// } else {
			// fmt.Println("default route for ", destination)
		// }
	// }
// 
	// if handler == nil {
		// handler = d.ohm.GetDefaultHandler()
	// }
// 
	// if handler == nil {
		// fmt.Println("default outbound handler not exist")
		// common.Close(link.Writer)
		// common.Interrupt(link.Reader)
		// return
	// }

//	handler.Dispatch(ctx, link)

}

func sniffer(ctx context.Context, cReader *cachedReader, metadataOnly bool) (SniffResult, error) {
	fmt.Println("in app-routing-rules.go func sniffer")
	payload := buf.New()
	defer payload.Release()

	sniffer := dispatcher.NewSniffer(ctx)

	//metaresult, metadataErr := sniffer.SniffMetadata(ctx)
	metaresult := SniffResult(nil)
	metadataErr := errors.New("NOT have sniff")
	//fmt.Println("in app-routing-rules.go func sniffer metaresult is: ", metaresult)
	if metadataOnly {
		return metaresult, metadataErr
	}
	//fmt.Println("in app-routing-rules.go func sniffer 不会执行")
	contentResult, contentErr := func() (SniffResult, error) {
		totalAttempt := 0
		for {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			default:
				totalAttempt++
				if totalAttempt > 2 {
					return nil, errSniffingTimeout
				}

				cReader.Cache(payload)
				if !payload.IsEmpty() {
					result, err := sniffer.Sniff(ctx, payload.Bytes())
					fmt.Println("in app-routing-rules.go func sniffer result is: ", result.Domain(), result.Protocol())
					if err != common.ErrNoClue {
						return result, err
					}
				}
				if payload.IsFull() {
					return nil, errUnknownContent
				}
			}
		}
	}()
	if contentErr != nil && metadataErr == nil {
		return metaresult, nil
	}
	if contentErr == nil && metadataErr == nil {
		return CompositeResult(metaresult, contentResult), nil
	}
	return contentResult, contentErr
}

func CompositeResult(domainResult SniffResult, protocolResult SniffResult) SniffResult {
	return &compositeResult{domainResult: domainResult, protocolResult: protocolResult}
}

type compositeResult struct {
	domainResult   SniffResult
	protocolResult SniffResult
}

func (c compositeResult) Protocol() string {
	return c.protocolResult.Protocol()
}

func (c compositeResult) Domain() string {
	return c.domainResult.Domain()
}

func (c compositeResult) ProtocolForDomainResult() string {
	return c.domainResult.Protocol()
}

var errUnknownContent = errors.New("unknown content")

var (
	errSniffingTimeout = errors.New("timeout on sniffing")
)

type cachedReader struct {
	sync.Mutex
	reader *pipe.Reader
	cache  buf.MultiBuffer
}

func (r *cachedReader) Cache(b *buf.Buffer) {
	//fmt.Println("in app-routing-rules.go func (r *cachedReader) Cache")
	mb, _ := r.reader.ReadMultiBufferTimeout(time.Millisecond * 100)
	r.Lock()
	if !mb.IsEmpty() {
		r.cache, _ = buf.MergeMulti(r.cache, mb)
	}
	b.Clear()
	rawBytes := b.Extend(buf.Size)
	n := r.cache.Copy(rawBytes)
	b.Resize(0, int32(n))
	r.Unlock()
}

func (r *cachedReader) readInternal() buf.MultiBuffer {
	//fmt.Println("in app-routing-rules.go func (r *cachedReader) readInternal()")
	r.Lock()
	defer r.Unlock()

	if r.cache != nil && !r.cache.IsEmpty() {
		mb := r.cache
		r.cache = nil
		return mb
	}

	return nil
}

func (r *cachedReader) ReadMultiBuffer() (buf.MultiBuffer, error) {
	//fmt.Println("in app-routing-rules.go func (r *cachedReader) ReadMultiBuffer()")
	mb := r.readInternal()
	if mb != nil {
		return mb, nil
	}

	return r.reader.ReadMultiBuffer()
}

func (r *cachedReader) ReadMultiBufferTimeout(timeout time.Duration) (buf.MultiBuffer, error) {
	fmt.Println("in app-routing-rules.go func (r *cachedReader) ReadMultiBufferTimeout")
	mb := r.readInternal()
	if mb != nil {
		return mb, nil
	}

	return r.reader.ReadMultiBufferTimeout(timeout)
}

func (r *cachedReader) Interrupt() {
	fmt.Println("in app-routing-rules.go func (r *cachedReader) Interrupt()")
	r.Lock()
	if r.cache != nil {
		r.cache = buf.ReleaseMulti(r.cache)
	}
	r.Unlock()
	r.reader.Interrupt()
}

func shouldOverride(result SniffResult, domainOverride []string) bool {
	fmt.Println("in app-routing-rules.go func shouldOverride")
	protocolString := result.Protocol()
	if resComp, ok := result.(SnifferResultComposite); ok {
		protocolString = resComp.ProtocolForDomainResult()
	}
	for _, p := range domainOverride {
		if strings.HasPrefix(protocolString, p) {
			return true
		}
	}
	return false
}

type SnifferResultComposite interface {
	ProtocolForDomainResult() string
}

type SniffResult interface {
	Protocol() string
	Domain() string
}

func New(ctx context.Context, config *core.RoutingHandlerConfig) (*Handler, error) {
	//v := core.MustFromContext(ctx)
	//var rulesTag core.Tag = "rulesTag"
	handler := &Handler{

		ctx: session.ContextWithRouting(ctx, config),
	}

	return handler, nil
}

func (handler *Handler) Start() error {
	//fmt.Println("in proxy-vmess-inbound-inbound.go func (handler *Handler) Start()")
	return nil
}

func (handler *Handler) Close() error {
	fmt.Println("in app-routing-rules.go func (handler *Handler) Close()")
	return nil
}

func (handler *Handler) Type() interface{} {
	fmt.Println("in proxy-vmess-inbound-inbound.go func (handler *Handler) Type()")
	return (*Handler)(nil)
}

func RouterType() interface{} {
	return (*Handler)(nil)
}

func (handler *Handler) Getctx() context.Context {
	return handler.ctx
}

func init() {
	fmt.Println("in app-routing-rules.go func init()")
	common.RegisterConfig((*core.RoutingHandlerConfig)(nil), func(ctx context.Context, config interface{}) (interface{}, error) {
		//d := new(Handler)
		//return d, nil
		return New(ctx, config.(*core.RoutingHandlerConfig))
	})
}
