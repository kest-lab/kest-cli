package client

import (
	"context"
	"fmt"
	"time"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCOptions struct {
	Address   string
	Method    string // service/method
	Data      string // JSON
	ProtoPath string
	Timeout   time.Duration
}

type GRPCResponse struct {
	Data     []byte
	Status   string
	Duration time.Duration
}

func ExecuteGRPC(opts GRPCOptions) (*GRPCResponse, error) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	// Parse proto
	p := protoparse.Parser{}
	fds, err := p.ParseFiles(opts.ProtoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to parse proto: %v", err)
	}

	var methodDesc *desc.MethodDescriptor
	// Search in all files
	// This is simplified, usually we'd want to specify the package or service more clearly
	// Format expected: package.Service/Method
	found := false
	for _, fd := range fds {
		for _, sd := range fd.GetServices() {
			for _, md := range sd.GetMethods() {
				fullName := fmt.Sprintf("%s/%s", sd.GetFullyQualifiedName(), md.GetName())
				if fullName == opts.Method || md.GetName() == opts.Method {
					methodDesc = md
					found = true
					break
				}
			}
			if found {
				break
			}
		}
		if found {
			break
		}
	}

	if methodDesc == nil {
		return nil, fmt.Errorf("method not found in proto: %s", opts.Method)
	}

	// Dial
	conn, err := grpc.DialContext(ctx, opts.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to dial: %v", err)
	}
	defer conn.Close()

	stub := grpcdynamic.NewStub(conn)

	// Prepare request
	reqMsg := dynamic.NewMessage(methodDesc.GetInputType())
	err = reqMsg.UnmarshalJSON([]byte(opts.Data))
	if err != nil {
		return nil, fmt.Errorf("failed to parse input JSON as proto: %v", err)
	}

	// Call
	resp, err := stub.InvokeRpc(ctx, methodDesc, reqMsg)
	if err != nil {
		return nil, fmt.Errorf("gRPC call failed: %v", err)
	}

	respMsg := resp.(*dynamic.Message)
	respJSON, err := respMsg.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal response to JSON: %v", err)
	}

	return &GRPCResponse{
		Data:     respJSON,
		Status:   "OK",
		Duration: time.Since(start),
	}, nil
}
