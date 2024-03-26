// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v5.26.0
// source: mflow/apps/task/pb/job_rpc.proto

package task

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	JobRPC_QueryJobTask_FullMethodName        = "/infraboard.mflow.task.JobRPC/QueryJobTask"
	JobRPC_UpdateJobTaskStatus_FullMethodName = "/infraboard.mflow.task.JobRPC/UpdateJobTaskStatus"
	JobRPC_UpdateJobTaskOutput_FullMethodName = "/infraboard.mflow.task.JobRPC/UpdateJobTaskOutput"
	JobRPC_DescribeJobTask_FullMethodName     = "/infraboard.mflow.task.JobRPC/DescribeJobTask"
	JobRPC_WatchJobTaskLog_FullMethodName     = "/infraboard.mflow.task.JobRPC/WatchJobTaskLog"
)

// JobRPCClient is the client API for JobRPC service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type JobRPCClient interface {
	// 查询任务
	QueryJobTask(ctx context.Context, in *QueryJobTaskRequest, opts ...grpc.CallOption) (*JobTaskSet, error)
	// 更新任务状态
	UpdateJobTaskStatus(ctx context.Context, in *UpdateJobTaskStatusRequest, opts ...grpc.CallOption) (*JobTask, error)
	// 更新任务输出结果
	UpdateJobTaskOutput(ctx context.Context, in *UpdateJobTaskOutputRequest, opts ...grpc.CallOption) (*JobTask, error)
	// 任务执行详情
	DescribeJobTask(ctx context.Context, in *DescribeJobTaskRequest, opts ...grpc.CallOption) (*JobTask, error)
	// 查询Task日志
	WatchJobTaskLog(ctx context.Context, in *WatchJobTaskLogRequest, opts ...grpc.CallOption) (JobRPC_WatchJobTaskLogClient, error)
}

type jobRPCClient struct {
	cc grpc.ClientConnInterface
}

func NewJobRPCClient(cc grpc.ClientConnInterface) JobRPCClient {
	return &jobRPCClient{cc}
}

func (c *jobRPCClient) QueryJobTask(ctx context.Context, in *QueryJobTaskRequest, opts ...grpc.CallOption) (*JobTaskSet, error) {
	out := new(JobTaskSet)
	err := c.cc.Invoke(ctx, JobRPC_QueryJobTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobRPCClient) UpdateJobTaskStatus(ctx context.Context, in *UpdateJobTaskStatusRequest, opts ...grpc.CallOption) (*JobTask, error) {
	out := new(JobTask)
	err := c.cc.Invoke(ctx, JobRPC_UpdateJobTaskStatus_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobRPCClient) UpdateJobTaskOutput(ctx context.Context, in *UpdateJobTaskOutputRequest, opts ...grpc.CallOption) (*JobTask, error) {
	out := new(JobTask)
	err := c.cc.Invoke(ctx, JobRPC_UpdateJobTaskOutput_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobRPCClient) DescribeJobTask(ctx context.Context, in *DescribeJobTaskRequest, opts ...grpc.CallOption) (*JobTask, error) {
	out := new(JobTask)
	err := c.cc.Invoke(ctx, JobRPC_DescribeJobTask_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *jobRPCClient) WatchJobTaskLog(ctx context.Context, in *WatchJobTaskLogRequest, opts ...grpc.CallOption) (JobRPC_WatchJobTaskLogClient, error) {
	stream, err := c.cc.NewStream(ctx, &JobRPC_ServiceDesc.Streams[0], JobRPC_WatchJobTaskLog_FullMethodName, opts...)
	if err != nil {
		return nil, err
	}
	x := &jobRPCWatchJobTaskLogClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type JobRPC_WatchJobTaskLogClient interface {
	Recv() (*JobTaskStreamReponse, error)
	grpc.ClientStream
}

type jobRPCWatchJobTaskLogClient struct {
	grpc.ClientStream
}

func (x *jobRPCWatchJobTaskLogClient) Recv() (*JobTaskStreamReponse, error) {
	m := new(JobTaskStreamReponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// JobRPCServer is the server API for JobRPC service.
// All implementations must embed UnimplementedJobRPCServer
// for forward compatibility
type JobRPCServer interface {
	// 查询任务
	QueryJobTask(context.Context, *QueryJobTaskRequest) (*JobTaskSet, error)
	// 更新任务状态
	UpdateJobTaskStatus(context.Context, *UpdateJobTaskStatusRequest) (*JobTask, error)
	// 更新任务输出结果
	UpdateJobTaskOutput(context.Context, *UpdateJobTaskOutputRequest) (*JobTask, error)
	// 任务执行详情
	DescribeJobTask(context.Context, *DescribeJobTaskRequest) (*JobTask, error)
	// 查询Task日志
	WatchJobTaskLog(*WatchJobTaskLogRequest, JobRPC_WatchJobTaskLogServer) error
	mustEmbedUnimplementedJobRPCServer()
}

// UnimplementedJobRPCServer must be embedded to have forward compatible implementations.
type UnimplementedJobRPCServer struct {
}

func (UnimplementedJobRPCServer) QueryJobTask(context.Context, *QueryJobTaskRequest) (*JobTaskSet, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryJobTask not implemented")
}
func (UnimplementedJobRPCServer) UpdateJobTaskStatus(context.Context, *UpdateJobTaskStatusRequest) (*JobTask, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateJobTaskStatus not implemented")
}
func (UnimplementedJobRPCServer) UpdateJobTaskOutput(context.Context, *UpdateJobTaskOutputRequest) (*JobTask, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateJobTaskOutput not implemented")
}
func (UnimplementedJobRPCServer) DescribeJobTask(context.Context, *DescribeJobTaskRequest) (*JobTask, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeJobTask not implemented")
}
func (UnimplementedJobRPCServer) WatchJobTaskLog(*WatchJobTaskLogRequest, JobRPC_WatchJobTaskLogServer) error {
	return status.Errorf(codes.Unimplemented, "method WatchJobTaskLog not implemented")
}
func (UnimplementedJobRPCServer) mustEmbedUnimplementedJobRPCServer() {}

// UnsafeJobRPCServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to JobRPCServer will
// result in compilation errors.
type UnsafeJobRPCServer interface {
	mustEmbedUnimplementedJobRPCServer()
}

func RegisterJobRPCServer(s grpc.ServiceRegistrar, srv JobRPCServer) {
	s.RegisterService(&JobRPC_ServiceDesc, srv)
}

func _JobRPC_QueryJobTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryJobTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobRPCServer).QueryJobTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobRPC_QueryJobTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobRPCServer).QueryJobTask(ctx, req.(*QueryJobTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobRPC_UpdateJobTaskStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateJobTaskStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobRPCServer).UpdateJobTaskStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobRPC_UpdateJobTaskStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobRPCServer).UpdateJobTaskStatus(ctx, req.(*UpdateJobTaskStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobRPC_UpdateJobTaskOutput_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateJobTaskOutputRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobRPCServer).UpdateJobTaskOutput(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobRPC_UpdateJobTaskOutput_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobRPCServer).UpdateJobTaskOutput(ctx, req.(*UpdateJobTaskOutputRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobRPC_DescribeJobTask_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeJobTaskRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JobRPCServer).DescribeJobTask(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: JobRPC_DescribeJobTask_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JobRPCServer).DescribeJobTask(ctx, req.(*DescribeJobTaskRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _JobRPC_WatchJobTaskLog_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(WatchJobTaskLogRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(JobRPCServer).WatchJobTaskLog(m, &jobRPCWatchJobTaskLogServer{stream})
}

type JobRPC_WatchJobTaskLogServer interface {
	Send(*JobTaskStreamReponse) error
	grpc.ServerStream
}

type jobRPCWatchJobTaskLogServer struct {
	grpc.ServerStream
}

func (x *jobRPCWatchJobTaskLogServer) Send(m *JobTaskStreamReponse) error {
	return x.ServerStream.SendMsg(m)
}

// JobRPC_ServiceDesc is the grpc.ServiceDesc for JobRPC service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var JobRPC_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "infraboard.mflow.task.JobRPC",
	HandlerType: (*JobRPCServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "QueryJobTask",
			Handler:    _JobRPC_QueryJobTask_Handler,
		},
		{
			MethodName: "UpdateJobTaskStatus",
			Handler:    _JobRPC_UpdateJobTaskStatus_Handler,
		},
		{
			MethodName: "UpdateJobTaskOutput",
			Handler:    _JobRPC_UpdateJobTaskOutput_Handler,
		},
		{
			MethodName: "DescribeJobTask",
			Handler:    _JobRPC_DescribeJobTask_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "WatchJobTaskLog",
			Handler:       _JobRPC_WatchJobTaskLog_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "mflow/apps/task/pb/job_rpc.proto",
}
