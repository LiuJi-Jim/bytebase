// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             (unknown)
// source: v1/database_service.proto

package v1

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

// DatabaseServiceClient is the client API for DatabaseService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DatabaseServiceClient interface {
	GetDatabase(ctx context.Context, in *GetDatabaseRequest, opts ...grpc.CallOption) (*Database, error)
	ListDatabases(ctx context.Context, in *ListDatabasesRequest, opts ...grpc.CallOption) (*ListDatabasesResponse, error)
	UpdateDatabase(ctx context.Context, in *UpdateDatabaseRequest, opts ...grpc.CallOption) (*Database, error)
	BatchUpdateDatabases(ctx context.Context, in *BatchUpdateDatabasesRequest, opts ...grpc.CallOption) (*BatchUpdateDatabasesResponse, error)
	GetDatabaseMetadata(ctx context.Context, in *GetDatabaseMetadataRequest, opts ...grpc.CallOption) (*DatabaseMetadata, error)
	GetDatabaseSchema(ctx context.Context, in *GetDatabaseSchemaRequest, opts ...grpc.CallOption) (*DatabaseSchema, error)
	GetBackupSetting(ctx context.Context, in *GetBackupSettingRequest, opts ...grpc.CallOption) (*BackupSetting, error)
	UpdateBackupSetting(ctx context.Context, in *UpdateBackupSettingRequest, opts ...grpc.CallOption) (*BackupSetting, error)
}

type databaseServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDatabaseServiceClient(cc grpc.ClientConnInterface) DatabaseServiceClient {
	return &databaseServiceClient{cc}
}

func (c *databaseServiceClient) GetDatabase(ctx context.Context, in *GetDatabaseRequest, opts ...grpc.CallOption) (*Database, error) {
	out := new(Database)
	err := c.cc.Invoke(ctx, "/bytebase.v1.DatabaseService/GetDatabase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) ListDatabases(ctx context.Context, in *ListDatabasesRequest, opts ...grpc.CallOption) (*ListDatabasesResponse, error) {
	out := new(ListDatabasesResponse)
	err := c.cc.Invoke(ctx, "/bytebase.v1.DatabaseService/ListDatabases", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) UpdateDatabase(ctx context.Context, in *UpdateDatabaseRequest, opts ...grpc.CallOption) (*Database, error) {
	out := new(Database)
	err := c.cc.Invoke(ctx, "/bytebase.v1.DatabaseService/UpdateDatabase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) BatchUpdateDatabases(ctx context.Context, in *BatchUpdateDatabasesRequest, opts ...grpc.CallOption) (*BatchUpdateDatabasesResponse, error) {
	out := new(BatchUpdateDatabasesResponse)
	err := c.cc.Invoke(ctx, "/bytebase.v1.DatabaseService/BatchUpdateDatabases", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) GetDatabaseMetadata(ctx context.Context, in *GetDatabaseMetadataRequest, opts ...grpc.CallOption) (*DatabaseMetadata, error) {
	out := new(DatabaseMetadata)
	err := c.cc.Invoke(ctx, "/bytebase.v1.DatabaseService/GetDatabaseMetadata", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) GetDatabaseSchema(ctx context.Context, in *GetDatabaseSchemaRequest, opts ...grpc.CallOption) (*DatabaseSchema, error) {
	out := new(DatabaseSchema)
	err := c.cc.Invoke(ctx, "/bytebase.v1.DatabaseService/GetDatabaseSchema", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) GetBackupSetting(ctx context.Context, in *GetBackupSettingRequest, opts ...grpc.CallOption) (*BackupSetting, error) {
	out := new(BackupSetting)
	err := c.cc.Invoke(ctx, "/bytebase.v1.DatabaseService/GetBackupSetting", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) UpdateBackupSetting(ctx context.Context, in *UpdateBackupSettingRequest, opts ...grpc.CallOption) (*BackupSetting, error) {
	out := new(BackupSetting)
	err := c.cc.Invoke(ctx, "/bytebase.v1.DatabaseService/UpdateBackupSetting", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatabaseServiceServer is the server API for DatabaseService service.
// All implementations must embed UnimplementedDatabaseServiceServer
// for forward compatibility
type DatabaseServiceServer interface {
	GetDatabase(context.Context, *GetDatabaseRequest) (*Database, error)
	ListDatabases(context.Context, *ListDatabasesRequest) (*ListDatabasesResponse, error)
	UpdateDatabase(context.Context, *UpdateDatabaseRequest) (*Database, error)
	BatchUpdateDatabases(context.Context, *BatchUpdateDatabasesRequest) (*BatchUpdateDatabasesResponse, error)
	GetDatabaseMetadata(context.Context, *GetDatabaseMetadataRequest) (*DatabaseMetadata, error)
	GetDatabaseSchema(context.Context, *GetDatabaseSchemaRequest) (*DatabaseSchema, error)
	GetBackupSetting(context.Context, *GetBackupSettingRequest) (*BackupSetting, error)
	UpdateBackupSetting(context.Context, *UpdateBackupSettingRequest) (*BackupSetting, error)
	mustEmbedUnimplementedDatabaseServiceServer()
}

// UnimplementedDatabaseServiceServer must be embedded to have forward compatible implementations.
type UnimplementedDatabaseServiceServer struct {
}

func (UnimplementedDatabaseServiceServer) GetDatabase(context.Context, *GetDatabaseRequest) (*Database, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatabase not implemented")
}
func (UnimplementedDatabaseServiceServer) ListDatabases(context.Context, *ListDatabasesRequest) (*ListDatabasesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDatabases not implemented")
}
func (UnimplementedDatabaseServiceServer) UpdateDatabase(context.Context, *UpdateDatabaseRequest) (*Database, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDatabase not implemented")
}
func (UnimplementedDatabaseServiceServer) BatchUpdateDatabases(context.Context, *BatchUpdateDatabasesRequest) (*BatchUpdateDatabasesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchUpdateDatabases not implemented")
}
func (UnimplementedDatabaseServiceServer) GetDatabaseMetadata(context.Context, *GetDatabaseMetadataRequest) (*DatabaseMetadata, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatabaseMetadata not implemented")
}
func (UnimplementedDatabaseServiceServer) GetDatabaseSchema(context.Context, *GetDatabaseSchemaRequest) (*DatabaseSchema, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatabaseSchema not implemented")
}
func (UnimplementedDatabaseServiceServer) GetBackupSetting(context.Context, *GetBackupSettingRequest) (*BackupSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetBackupSetting not implemented")
}
func (UnimplementedDatabaseServiceServer) UpdateBackupSetting(context.Context, *UpdateBackupSettingRequest) (*BackupSetting, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateBackupSetting not implemented")
}
func (UnimplementedDatabaseServiceServer) mustEmbedUnimplementedDatabaseServiceServer() {}

// UnsafeDatabaseServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DatabaseServiceServer will
// result in compilation errors.
type UnsafeDatabaseServiceServer interface {
	mustEmbedUnimplementedDatabaseServiceServer()
}

func RegisterDatabaseServiceServer(s grpc.ServiceRegistrar, srv DatabaseServiceServer) {
	s.RegisterService(&DatabaseService_ServiceDesc, srv)
}

func _DatabaseService_GetDatabase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatabaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).GetDatabase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bytebase.v1.DatabaseService/GetDatabase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).GetDatabase(ctx, req.(*GetDatabaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_ListDatabases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDatabasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).ListDatabases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bytebase.v1.DatabaseService/ListDatabases",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).ListDatabases(ctx, req.(*ListDatabasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_UpdateDatabase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDatabaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).UpdateDatabase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bytebase.v1.DatabaseService/UpdateDatabase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).UpdateDatabase(ctx, req.(*UpdateDatabaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_BatchUpdateDatabases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchUpdateDatabasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).BatchUpdateDatabases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bytebase.v1.DatabaseService/BatchUpdateDatabases",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).BatchUpdateDatabases(ctx, req.(*BatchUpdateDatabasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_GetDatabaseMetadata_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatabaseMetadataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).GetDatabaseMetadata(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bytebase.v1.DatabaseService/GetDatabaseMetadata",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).GetDatabaseMetadata(ctx, req.(*GetDatabaseMetadataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_GetDatabaseSchema_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatabaseSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).GetDatabaseSchema(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bytebase.v1.DatabaseService/GetDatabaseSchema",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).GetDatabaseSchema(ctx, req.(*GetDatabaseSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_GetBackupSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetBackupSettingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).GetBackupSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bytebase.v1.DatabaseService/GetBackupSetting",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).GetBackupSetting(ctx, req.(*GetBackupSettingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_UpdateBackupSetting_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateBackupSettingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).UpdateBackupSetting(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/bytebase.v1.DatabaseService/UpdateBackupSetting",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).UpdateBackupSetting(ctx, req.(*UpdateBackupSettingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DatabaseService_ServiceDesc is the grpc.ServiceDesc for DatabaseService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DatabaseService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "bytebase.v1.DatabaseService",
	HandlerType: (*DatabaseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetDatabase",
			Handler:    _DatabaseService_GetDatabase_Handler,
		},
		{
			MethodName: "ListDatabases",
			Handler:    _DatabaseService_ListDatabases_Handler,
		},
		{
			MethodName: "UpdateDatabase",
			Handler:    _DatabaseService_UpdateDatabase_Handler,
		},
		{
			MethodName: "BatchUpdateDatabases",
			Handler:    _DatabaseService_BatchUpdateDatabases_Handler,
		},
		{
			MethodName: "GetDatabaseMetadata",
			Handler:    _DatabaseService_GetDatabaseMetadata_Handler,
		},
		{
			MethodName: "GetDatabaseSchema",
			Handler:    _DatabaseService_GetDatabaseSchema_Handler,
		},
		{
			MethodName: "GetBackupSetting",
			Handler:    _DatabaseService_GetBackupSetting_Handler,
		},
		{
			MethodName: "UpdateBackupSetting",
			Handler:    _DatabaseService_UpdateBackupSetting_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "v1/database_service.proto",
}
