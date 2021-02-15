package nolan

import (
	"context"
	"io"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/bdkiran/nolan/nolan/config"
	"github.com/bdkiran/nolan/protocol"
	"github.com/davecgh/go-spew/spew"
	"go.uber.org/zap"
)

type contextKey string

var (
	serverVerboseLogs    bool
	requestQueueSpanKey  = contextKey("request queue span key")
	responseQueueSpanKey = contextKey("response queue span key")
)

func init() {
	spew.Config.Indent = ""

	e := os.Getenv("JOCKODEBUG")
	if strings.Contains(e, "server=1") {
		serverVerboseLogs = true
	}
}

// Handler is the interface that wraps the Broker's methods.
type Handler interface {
	Run(context.Context, <-chan *Context, chan<- *Context)
	Leave() error
	Shutdown() error
}

// Server is used to handle the TCP connections, decode requests,
// defer to the broker, and encode the responses.
type Server struct {
	config       *config.Config
	protocolLn   *net.TCPListener
	handler      Handler
	shutdown     bool
	shutdownCh   chan struct{}
	shutdownLock sync.Mutex
	requestCh    chan *Context
	responseCh   chan *Context
	close        func() error
	//metrics      *Metrics
	//tracer     opentracing.Tracer
}

//NewServer creates and returns a new sever struct
func NewServer(config *config.Config, handler Handler) *Server {
	s := &Server{
		config:  config,
		handler: handler,
		//metrics:    metrics,
		shutdownCh: make(chan struct{}),
		requestCh:  make(chan *Context, 1024),
		responseCh: make(chan *Context, 1024),
		//tracer:     tracer,
		//close: close,
	}
	return s
}

// Start starts the service.
func (s *Server) Start(ctx context.Context) error {
	zap.S().Debugw("Starting nolan server...")
	//Resolve the tcp address from the server configuration
	protocolAddr, err := net.ResolveTCPAddr("tcp", s.config.Addr)
	if err != nil {
		return err
	}
	//Start listening on that tcp address
	if s.protocolLn, err = net.ListenTCP("tcp", protocolAddr); err != nil {
		return err
	}

	//This go routine accepts connections and handles requests to the server
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case <-s.shutdownCh:
				break
			default:
				conn, err := s.protocolLn.Accept()
				if err != nil {
					zap.S().Errorf("server/%d: listener accept error: %s", s.config.ID, err)
					continue
				}

				go s.handleRequest(conn)
			}
		}
	}()

	//This go routine handles responses that the server sends out
	go func() {
		for {
			select {
			case <-ctx.Done():
				break
			case <-s.shutdownCh:
				break
			case respCtx := <-s.responseCh:
				// if queueSpan, ok := respCtx.Value(responseQueueSpanKey).(opentracing.Span); ok {
				// 	queueSpan.Finish()
				// }
				if err := s.handleResponse(respCtx); err != nil {
					zap.S().Errorf("server/%d: listener accept error: %s", s.config.ID, err)
				}
			}
		}
	}()
	zap.S().Debugf("server/%d: Running server", s.config.ID)
	//Runs the broker...
	go s.handler.Run(ctx, s.requestCh, s.responseCh)

	return nil
}

//Leave server leaves the cluser?
func (s *Server) Leave() error {
	return s.handler.Leave()
}

// Shutdown closes the service.
func (s *Server) Shutdown() error {
	s.shutdownLock.Lock()
	defer s.shutdownLock.Unlock()

	if s.shutdown {
		return nil
	}

	s.shutdown = true
	close(s.shutdownCh)

	if err := s.handler.Shutdown(); err != nil {
		return err
	}
	if err := s.protocolLn.Close(); err != nil {
		return err
	}

	//s.close()

	return nil
}

func (s *Server) handleRequest(conn net.Conn) {
	defer conn.Close()

	for {
		p := make([]byte, 4)
		_, err := io.ReadFull(conn, p[:])
		if err == io.EOF {
			break
		}
		if err != nil {
			zap.S().Errorf("conn read error: %s", err)
			break
		}

		//span := s.tracer.StartSpan("request")
		//decodeSpan := s.tracer.StartSpan("server: decode request", opentracing.ChildOf(span.Context()))

		size := protocol.Encoding.Uint32(p)
		if size == 0 {
			break // TODO: should this even happen?
		}

		b := make([]byte, size+4) //+4 since we're going to copy the size into b
		copy(b, p)

		if _, err = io.ReadFull(conn, b[4:]); err != nil {
			// TODO: handle request
			//span.LogKV("msg", "failed to read from connection", "err", err)
			//span.Finish()
			zap.S().Panic(err)
			//panic(err)
		}

		d := protocol.NewDecoder(b)
		header := new(protocol.RequestHeader)
		if err := header.Decode(d); err != nil {
			// TODO: handle err
			//span.LogKV("msg", "failed to decode header", "err", err)
			//span.Finish()
			zap.S().Panic(err)
		}

		//span.SetTag("api_key", header.APIKey)
		//span.SetTag("correlation_id", header.CorrelationID)
		//span.SetTag("client_id", header.ClientID)
		//span.SetTag("size", size)
		//span.SetTag("node_id", s.config.ID) // can I set this globally for the tracer?
		//span.SetTag("addr", s.config.Addr)

		var req protocol.VersionedDecoder

		switch header.APIKey {
		case protocol.ProduceKey:
			req = &protocol.ProduceRequest{}
		case protocol.FetchKey:
			req = &protocol.FetchRequest{}
		case protocol.OffsetsKey:
			req = &protocol.OffsetsRequest{}
		case protocol.MetadataKey:
			req = &protocol.MetadataRequest{}
		case protocol.LeaderAndISRKey:
			req = &protocol.LeaderAndISRRequest{}
		case protocol.StopReplicaKey:
			req = &protocol.StopReplicaRequest{}
		case protocol.UpdateMetadataKey:
			req = &protocol.UpdateMetadataRequest{}
		case protocol.ControlledShutdownKey:
			req = &protocol.ControlledShutdownRequest{}
		case protocol.OffsetCommitKey:
			req = &protocol.OffsetCommitRequest{}
		case protocol.OffsetFetchKey:
			req = &protocol.OffsetFetchRequest{}
		case protocol.FindCoordinatorKey:
			req = &protocol.FindCoordinatorRequest{}
		case protocol.JoinGroupKey:
			req = &protocol.JoinGroupRequest{}
		case protocol.HeartbeatKey:
			req = &protocol.HeartbeatRequest{}
		case protocol.LeaveGroupKey:
			req = &protocol.LeaveGroupRequest{}
		case protocol.SyncGroupKey:
			req = &protocol.SyncGroupRequest{}
		case protocol.DescribeGroupsKey:
			req = &protocol.DescribeGroupsRequest{}
		case protocol.ListGroupsKey:
			req = &protocol.ListGroupsRequest{}
		case protocol.SaslHandshakeKey:
			req = &protocol.SaslHandshakeRequest{}
		case protocol.APIVersionsKey:
			req = &protocol.APIVersionsRequest{}
		case protocol.CreateTopicsKey:
			req = &protocol.CreateTopicRequests{}
		case protocol.DeleteTopicsKey:
			req = &protocol.DeleteTopicsRequest{}
		}

		if err := req.Decode(d, header.APIVersion); err != nil {
			zap.S().Errorf("server/%d: %s: decode request failed: %s", s.config.ID, header, err)
			//span.LogKV("msg", "failed to decode request", "err", err)
			//span.Finish()
			zap.S().Panic(err)
		}

		//decodeSpan.Finish()

		//ctx := opentracing.ContextWithSpan(context.Background(), span)
		//queueSpan := s.tracer.StartSpan("server: queue request", opentracing.ChildOf(span.Context()))
		//ctx = context.WithValue(ctx, requestQueueSpanKey, queueSpan)

		reqCtx := &Context{
			//parent: ctx,
			header: header,
			req:    req,
			conn:   conn,
		}
		zap.S().Debugf("server/%d: handle request: %s", s.config.ID, reqCtx)

		s.requestCh <- reqCtx
	}
}

func (s *Server) handleResponse(respCtx *Context) error {
	//psp := opentracing.SpanFromContext(respCtx)
	//sp := s.tracer.StartSpan("server: handle response", opentracing.ChildOf(psp.Context()))
	zap.S().Debugf("server/%d: handle response: %s", s.config.ID, respCtx)

	//defer psp.Finish()
	//defer sp.Finish()

	b, err := protocol.Encode(respCtx.res.(protocol.Encoder))
	if err != nil {
		return err
	}
	_, err = respCtx.conn.Write(b)
	return err
}

// Addr returns the address on which the Server is listening
func (s *Server) Addr() net.Addr {
	return s.protocolLn.Addr()
}

//ID returns the server config ID
func (s *Server) ID() int32 {
	return s.config.ID
}
