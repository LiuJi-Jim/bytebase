// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0
// 	protoc        (unknown)
// source: store/activity.proto

package store

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type int32

const (
	ActivityIssueCommentCreatePayload_ExternalApprovalEvent_TYPE_UNSPECIFIED ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type = 0
	ActivityIssueCommentCreatePayload_ExternalApprovalEvent_TYPE_FEISHU      ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type = 1
)

// Enum value maps for ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type.
var (
	ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type_name = map[int32]string{
		0: "TYPE_UNSPECIFIED",
		1: "TYPE_FEISHU",
	}
	ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type_value = map[string]int32{
		"TYPE_UNSPECIFIED": 0,
		"TYPE_FEISHU":      1,
	}
)

func (x ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type) Enum() *ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type {
	p := new(ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type)
	*p = x
	return p
}

func (x ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_store_activity_proto_enumTypes[0].Descriptor()
}

func (ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type) Type() protoreflect.EnumType {
	return &file_store_activity_proto_enumTypes[0]
}

func (x ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type.Descriptor instead.
func (ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type) EnumDescriptor() ([]byte, []int) {
	return file_store_activity_proto_rawDescGZIP(), []int{1, 1, 0}
}

type ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action int32

const (
	ActivityIssueCommentCreatePayload_ExternalApprovalEvent_ACTION_UNSPECIFIED ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action = 0
	ActivityIssueCommentCreatePayload_ExternalApprovalEvent_ACTION_APPROVE     ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action = 1
	ActivityIssueCommentCreatePayload_ExternalApprovalEvent_ACTION_REJECT      ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action = 2
)

// Enum value maps for ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action.
var (
	ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action_name = map[int32]string{
		0: "ACTION_UNSPECIFIED",
		1: "ACTION_APPROVE",
		2: "ACTION_REJECT",
	}
	ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action_value = map[string]int32{
		"ACTION_UNSPECIFIED": 0,
		"ACTION_APPROVE":     1,
		"ACTION_REJECT":      2,
	}
)

func (x ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action) Enum() *ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action {
	p := new(ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action)
	*p = x
	return p
}

func (x ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action) Descriptor() protoreflect.EnumDescriptor {
	return file_store_activity_proto_enumTypes[1].Descriptor()
}

func (ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action) Type() protoreflect.EnumType {
	return &file_store_activity_proto_enumTypes[1]
}

func (x ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action.Descriptor instead.
func (ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action) EnumDescriptor() ([]byte, []int) {
	return file_store_activity_proto_rawDescGZIP(), []int{1, 1, 1}
}

type ActivityIssueCommentCreatePayload_ApprovalEvent_Status int32

const (
	ActivityIssueCommentCreatePayload_ApprovalEvent_STATUS_UNSPECIFIED ActivityIssueCommentCreatePayload_ApprovalEvent_Status = 0
	ActivityIssueCommentCreatePayload_ApprovalEvent_PENDING            ActivityIssueCommentCreatePayload_ApprovalEvent_Status = 1
	ActivityIssueCommentCreatePayload_ApprovalEvent_APPROVED           ActivityIssueCommentCreatePayload_ApprovalEvent_Status = 2
)

// Enum value maps for ActivityIssueCommentCreatePayload_ApprovalEvent_Status.
var (
	ActivityIssueCommentCreatePayload_ApprovalEvent_Status_name = map[int32]string{
		0: "STATUS_UNSPECIFIED",
		1: "PENDING",
		2: "APPROVED",
	}
	ActivityIssueCommentCreatePayload_ApprovalEvent_Status_value = map[string]int32{
		"STATUS_UNSPECIFIED": 0,
		"PENDING":            1,
		"APPROVED":           2,
	}
)

func (x ActivityIssueCommentCreatePayload_ApprovalEvent_Status) Enum() *ActivityIssueCommentCreatePayload_ApprovalEvent_Status {
	p := new(ActivityIssueCommentCreatePayload_ApprovalEvent_Status)
	*p = x
	return p
}

func (x ActivityIssueCommentCreatePayload_ApprovalEvent_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ActivityIssueCommentCreatePayload_ApprovalEvent_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_store_activity_proto_enumTypes[2].Descriptor()
}

func (ActivityIssueCommentCreatePayload_ApprovalEvent_Status) Type() protoreflect.EnumType {
	return &file_store_activity_proto_enumTypes[2]
}

func (x ActivityIssueCommentCreatePayload_ApprovalEvent_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ActivityIssueCommentCreatePayload_ApprovalEvent_Status.Descriptor instead.
func (ActivityIssueCommentCreatePayload_ApprovalEvent_Status) EnumDescriptor() ([]byte, []int) {
	return file_store_activity_proto_rawDescGZIP(), []int{1, 2, 0}
}

// ActivityIssueCreatePayload is the payloads for creating issues.
// These payload types are only used when marshalling to the json format for saving into the database.
// So we annotate with json tag using camelCase naming which is consistent with normal
// json naming convention. More importantly, frontend code can simply use JSON.parse to
// convert to the expected struct there.
type ActivityIssueCreatePayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Used by inbox to display info without paying the join cost
	IssueName string `protobuf:"bytes,1,opt,name=issue_name,json=issueName,proto3" json:"issue_name,omitempty"`
}

func (x *ActivityIssueCreatePayload) Reset() {
	*x = ActivityIssueCreatePayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_activity_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivityIssueCreatePayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivityIssueCreatePayload) ProtoMessage() {}

func (x *ActivityIssueCreatePayload) ProtoReflect() protoreflect.Message {
	mi := &file_store_activity_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivityIssueCreatePayload.ProtoReflect.Descriptor instead.
func (*ActivityIssueCreatePayload) Descriptor() ([]byte, []int) {
	return file_store_activity_proto_rawDescGZIP(), []int{0}
}

func (x *ActivityIssueCreatePayload) GetIssueName() string {
	if x != nil {
		return x.IssueName
	}
	return ""
}

// ActivityIssueCommentCreatePayload is the payloads for creating issue comments.
type ActivityIssueCommentCreatePayload struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Types that are assignable to Event:
	//
	//	*ActivityIssueCommentCreatePayload_ExternalApprovalEvent_
	//	*ActivityIssueCommentCreatePayload_TaskRollbackBy_
	//	*ActivityIssueCommentCreatePayload_ApprovalEvent_
	Event isActivityIssueCommentCreatePayload_Event `protobuf_oneof:"event"`
	// Used by inbox to display info without paying the join cost
	IssueName string `protobuf:"bytes,4,opt,name=issue_name,json=issueName,proto3" json:"issue_name,omitempty"`
}

func (x *ActivityIssueCommentCreatePayload) Reset() {
	*x = ActivityIssueCommentCreatePayload{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_activity_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivityIssueCommentCreatePayload) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivityIssueCommentCreatePayload) ProtoMessage() {}

func (x *ActivityIssueCommentCreatePayload) ProtoReflect() protoreflect.Message {
	mi := &file_store_activity_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivityIssueCommentCreatePayload.ProtoReflect.Descriptor instead.
func (*ActivityIssueCommentCreatePayload) Descriptor() ([]byte, []int) {
	return file_store_activity_proto_rawDescGZIP(), []int{1}
}

func (m *ActivityIssueCommentCreatePayload) GetEvent() isActivityIssueCommentCreatePayload_Event {
	if m != nil {
		return m.Event
	}
	return nil
}

func (x *ActivityIssueCommentCreatePayload) GetExternalApprovalEvent() *ActivityIssueCommentCreatePayload_ExternalApprovalEvent {
	if x, ok := x.GetEvent().(*ActivityIssueCommentCreatePayload_ExternalApprovalEvent_); ok {
		return x.ExternalApprovalEvent
	}
	return nil
}

func (x *ActivityIssueCommentCreatePayload) GetTaskRollbackBy() *ActivityIssueCommentCreatePayload_TaskRollbackBy {
	if x, ok := x.GetEvent().(*ActivityIssueCommentCreatePayload_TaskRollbackBy_); ok {
		return x.TaskRollbackBy
	}
	return nil
}

func (x *ActivityIssueCommentCreatePayload) GetApprovalEvent() *ActivityIssueCommentCreatePayload_ApprovalEvent {
	if x, ok := x.GetEvent().(*ActivityIssueCommentCreatePayload_ApprovalEvent_); ok {
		return x.ApprovalEvent
	}
	return nil
}

func (x *ActivityIssueCommentCreatePayload) GetIssueName() string {
	if x != nil {
		return x.IssueName
	}
	return ""
}

type isActivityIssueCommentCreatePayload_Event interface {
	isActivityIssueCommentCreatePayload_Event()
}

type ActivityIssueCommentCreatePayload_ExternalApprovalEvent_ struct {
	ExternalApprovalEvent *ActivityIssueCommentCreatePayload_ExternalApprovalEvent `protobuf:"bytes,1,opt,name=external_approval_event,json=externalApprovalEvent,proto3,oneof"`
}

type ActivityIssueCommentCreatePayload_TaskRollbackBy_ struct {
	TaskRollbackBy *ActivityIssueCommentCreatePayload_TaskRollbackBy `protobuf:"bytes,2,opt,name=task_rollback_by,json=taskRollbackBy,proto3,oneof"`
}

type ActivityIssueCommentCreatePayload_ApprovalEvent_ struct {
	ApprovalEvent *ActivityIssueCommentCreatePayload_ApprovalEvent `protobuf:"bytes,3,opt,name=approval_event,json=approvalEvent,proto3,oneof"`
}

func (*ActivityIssueCommentCreatePayload_ExternalApprovalEvent_) isActivityIssueCommentCreatePayload_Event() {
}

func (*ActivityIssueCommentCreatePayload_TaskRollbackBy_) isActivityIssueCommentCreatePayload_Event() {
}

func (*ActivityIssueCommentCreatePayload_ApprovalEvent_) isActivityIssueCommentCreatePayload_Event() {
}

// TaskRollbackBy records an issue rollback activity.
// The task with taskID in IssueID is rollbacked by the task with RollbackByTaskID in RollbackByIssueID.
type ActivityIssueCommentCreatePayload_TaskRollbackBy struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	IssueId           int64 `protobuf:"varint,1,opt,name=issue_id,json=issueId,proto3" json:"issue_id,omitempty"`
	TaskId            int64 `protobuf:"varint,2,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
	RollbackByIssueId int64 `protobuf:"varint,3,opt,name=rollback_by_issue_id,json=rollbackByIssueId,proto3" json:"rollback_by_issue_id,omitempty"`
	RollbackByTaskId  int64 `protobuf:"varint,4,opt,name=rollback_by_task_id,json=rollbackByTaskId,proto3" json:"rollback_by_task_id,omitempty"`
}

func (x *ActivityIssueCommentCreatePayload_TaskRollbackBy) Reset() {
	*x = ActivityIssueCommentCreatePayload_TaskRollbackBy{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_activity_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivityIssueCommentCreatePayload_TaskRollbackBy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivityIssueCommentCreatePayload_TaskRollbackBy) ProtoMessage() {}

func (x *ActivityIssueCommentCreatePayload_TaskRollbackBy) ProtoReflect() protoreflect.Message {
	mi := &file_store_activity_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivityIssueCommentCreatePayload_TaskRollbackBy.ProtoReflect.Descriptor instead.
func (*ActivityIssueCommentCreatePayload_TaskRollbackBy) Descriptor() ([]byte, []int) {
	return file_store_activity_proto_rawDescGZIP(), []int{1, 0}
}

func (x *ActivityIssueCommentCreatePayload_TaskRollbackBy) GetIssueId() int64 {
	if x != nil {
		return x.IssueId
	}
	return 0
}

func (x *ActivityIssueCommentCreatePayload_TaskRollbackBy) GetTaskId() int64 {
	if x != nil {
		return x.TaskId
	}
	return 0
}

func (x *ActivityIssueCommentCreatePayload_TaskRollbackBy) GetRollbackByIssueId() int64 {
	if x != nil {
		return x.RollbackByIssueId
	}
	return 0
}

func (x *ActivityIssueCommentCreatePayload_TaskRollbackBy) GetRollbackByTaskId() int64 {
	if x != nil {
		return x.RollbackByTaskId
	}
	return 0
}

type ActivityIssueCommentCreatePayload_ExternalApprovalEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type      ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type   `protobuf:"varint,1,opt,name=type,proto3,enum=bytebase.store.ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type" json:"type,omitempty"`
	Action    ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action `protobuf:"varint,2,opt,name=action,proto3,enum=bytebase.store.ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action" json:"action,omitempty"`
	StageName string                                                         `protobuf:"bytes,3,opt,name=stage_name,json=stageName,proto3" json:"stage_name,omitempty"`
}

func (x *ActivityIssueCommentCreatePayload_ExternalApprovalEvent) Reset() {
	*x = ActivityIssueCommentCreatePayload_ExternalApprovalEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_activity_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivityIssueCommentCreatePayload_ExternalApprovalEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivityIssueCommentCreatePayload_ExternalApprovalEvent) ProtoMessage() {}

func (x *ActivityIssueCommentCreatePayload_ExternalApprovalEvent) ProtoReflect() protoreflect.Message {
	mi := &file_store_activity_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivityIssueCommentCreatePayload_ExternalApprovalEvent.ProtoReflect.Descriptor instead.
func (*ActivityIssueCommentCreatePayload_ExternalApprovalEvent) Descriptor() ([]byte, []int) {
	return file_store_activity_proto_rawDescGZIP(), []int{1, 1}
}

func (x *ActivityIssueCommentCreatePayload_ExternalApprovalEvent) GetType() ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type {
	if x != nil {
		return x.Type
	}
	return ActivityIssueCommentCreatePayload_ExternalApprovalEvent_TYPE_UNSPECIFIED
}

func (x *ActivityIssueCommentCreatePayload_ExternalApprovalEvent) GetAction() ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action {
	if x != nil {
		return x.Action
	}
	return ActivityIssueCommentCreatePayload_ExternalApprovalEvent_ACTION_UNSPECIFIED
}

func (x *ActivityIssueCommentCreatePayload_ExternalApprovalEvent) GetStageName() string {
	if x != nil {
		return x.StageName
	}
	return ""
}

type ActivityIssueCommentCreatePayload_ApprovalEvent struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The new status.
	Status ActivityIssueCommentCreatePayload_ApprovalEvent_Status `protobuf:"varint,1,opt,name=status,proto3,enum=bytebase.store.ActivityIssueCommentCreatePayload_ApprovalEvent_Status" json:"status,omitempty"`
}

func (x *ActivityIssueCommentCreatePayload_ApprovalEvent) Reset() {
	*x = ActivityIssueCommentCreatePayload_ApprovalEvent{}
	if protoimpl.UnsafeEnabled {
		mi := &file_store_activity_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ActivityIssueCommentCreatePayload_ApprovalEvent) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ActivityIssueCommentCreatePayload_ApprovalEvent) ProtoMessage() {}

func (x *ActivityIssueCommentCreatePayload_ApprovalEvent) ProtoReflect() protoreflect.Message {
	mi := &file_store_activity_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ActivityIssueCommentCreatePayload_ApprovalEvent.ProtoReflect.Descriptor instead.
func (*ActivityIssueCommentCreatePayload_ApprovalEvent) Descriptor() ([]byte, []int) {
	return file_store_activity_proto_rawDescGZIP(), []int{1, 2}
}

func (x *ActivityIssueCommentCreatePayload_ApprovalEvent) GetStatus() ActivityIssueCommentCreatePayload_ApprovalEvent_Status {
	if x != nil {
		return x.Status
	}
	return ActivityIssueCommentCreatePayload_ApprovalEvent_STATUS_UNSPECIFIED
}

var File_store_activity_proto protoreflect.FileDescriptor

var file_store_activity_proto_rawDesc = []byte{
	0x0a, 0x14, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2f, 0x61, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65,
	0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x22, 0x3b, 0x0a, 0x1a, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69,
	0x74, 0x79, 0x49, 0x73, 0x73, 0x75, 0x65, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x61, 0x79,
	0x6c, 0x6f, 0x61, 0x64, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x73, 0x75, 0x65, 0x5f, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65, 0x4e,
	0x61, 0x6d, 0x65, 0x22, 0xf8, 0x08, 0x0a, 0x21, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79,
	0x49, 0x73, 0x73, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x12, 0x81, 0x01, 0x0a, 0x17, 0x65, 0x78,
	0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x5f, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x5f,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x47, 0x2e, 0x62, 0x79,
	0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x63, 0x74,
	0x69, 0x76, 0x69, 0x74, 0x79, 0x49, 0x73, 0x73, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e,
	0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x45,
	0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x15, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c,
	0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x6c, 0x0a,
	0x10, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x72, 0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x62,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x40, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61,
	0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74,
	0x79, 0x49, 0x73, 0x73, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x54, 0x61, 0x73, 0x6b, 0x52,
	0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x42, 0x79, 0x48, 0x00, 0x52, 0x0e, 0x74, 0x61, 0x73,
	0x6b, 0x52, 0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x42, 0x79, 0x12, 0x68, 0x0a, 0x0e, 0x61,
	0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x5f, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x3f, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x49, 0x73, 0x73,
	0x75, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x48, 0x00, 0x52, 0x0d, 0x61, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x69, 0x73, 0x73, 0x75, 0x65, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x4e, 0x61, 0x6d, 0x65, 0x1a, 0xa4, 0x01, 0x0a, 0x0e, 0x54, 0x61, 0x73, 0x6b, 0x52, 0x6f, 0x6c,
	0x6c, 0x62, 0x61, 0x63, 0x6b, 0x42, 0x79, 0x12, 0x19, 0x0a, 0x08, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x49, 0x64, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x12, 0x2f, 0x0a, 0x14, 0x72,
	0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x62, 0x79, 0x5f, 0x69, 0x73, 0x73, 0x75, 0x65,
	0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x11, 0x72, 0x6f, 0x6c, 0x6c, 0x62,
	0x61, 0x63, 0x6b, 0x42, 0x79, 0x49, 0x73, 0x73, 0x75, 0x65, 0x49, 0x64, 0x12, 0x2d, 0x0a, 0x13,
	0x72, 0x6f, 0x6c, 0x6c, 0x62, 0x61, 0x63, 0x6b, 0x5f, 0x62, 0x79, 0x5f, 0x74, 0x61, 0x73, 0x6b,
	0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x72, 0x6f, 0x6c, 0x6c, 0x62,
	0x61, 0x63, 0x6b, 0x42, 0x79, 0x54, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x1a, 0xf8, 0x02, 0x0a, 0x15,
	0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c,
	0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x60, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0e, 0x32, 0x4c, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61, 0x73, 0x65, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74, 0x79, 0x49, 0x73, 0x73,
	0x75, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50,
	0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x41,
	0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x66, 0x0a, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x4e, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62, 0x61,
	0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69, 0x74,
	0x79, 0x49, 0x73, 0x73, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x72,
	0x6e, 0x61, 0x6c, 0x41, 0x70, 0x70, 0x72, 0x6f, 0x76, 0x61, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74,
	0x2e, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x06, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x1d, 0x0a, 0x0a, 0x73, 0x74, 0x61, 0x67, 0x65, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x73, 0x74, 0x61, 0x67, 0x65, 0x4e, 0x61, 0x6d, 0x65, 0x22, 0x2d,
	0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x14, 0x0a, 0x10, 0x54, 0x59, 0x50, 0x45, 0x5f, 0x55,
	0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x46, 0x45, 0x49, 0x53, 0x48, 0x55, 0x10, 0x01, 0x22, 0x47, 0x0a,
	0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x16, 0x0a, 0x12, 0x41, 0x43, 0x54, 0x49, 0x4f,
	0x4e, 0x5f, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12,
	0x12, 0x0a, 0x0e, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x41, 0x50, 0x50, 0x52, 0x4f, 0x56,
	0x45, 0x10, 0x01, 0x12, 0x11, 0x0a, 0x0d, 0x41, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x52, 0x45,
	0x4a, 0x45, 0x43, 0x54, 0x10, 0x02, 0x1a, 0xac, 0x01, 0x0a, 0x0d, 0x41, 0x70, 0x70, 0x72, 0x6f,
	0x76, 0x61, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x5e, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x46, 0x2e, 0x62, 0x79, 0x74, 0x65, 0x62,
	0x61, 0x73, 0x65, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x65, 0x2e, 0x41, 0x63, 0x74, 0x69, 0x76, 0x69,
	0x74, 0x79, 0x49, 0x73, 0x73, 0x75, 0x65, 0x43, 0x6f, 0x6d, 0x6d, 0x65, 0x6e, 0x74, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x50, 0x61, 0x79, 0x6c, 0x6f, 0x61, 0x64, 0x2e, 0x41, 0x70, 0x70, 0x72,
	0x6f, 0x76, 0x61, 0x6c, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73,
	0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75, 0x73, 0x22, 0x3b, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x12, 0x16, 0x0a, 0x12, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55, 0x4e, 0x53,
	0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0b, 0x0a, 0x07, 0x50, 0x45,
	0x4e, 0x44, 0x49, 0x4e, 0x47, 0x10, 0x01, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x50, 0x50, 0x52, 0x4f,
	0x56, 0x45, 0x44, 0x10, 0x02, 0x42, 0x07, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x42, 0x14,
	0x5a, 0x12, 0x67, 0x65, 0x6e, 0x65, 0x72, 0x61, 0x74, 0x65, 0x64, 0x2d, 0x67, 0x6f, 0x2f, 0x73,
	0x74, 0x6f, 0x72, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_store_activity_proto_rawDescOnce sync.Once
	file_store_activity_proto_rawDescData = file_store_activity_proto_rawDesc
)

func file_store_activity_proto_rawDescGZIP() []byte {
	file_store_activity_proto_rawDescOnce.Do(func() {
		file_store_activity_proto_rawDescData = protoimpl.X.CompressGZIP(file_store_activity_proto_rawDescData)
	})
	return file_store_activity_proto_rawDescData
}

var file_store_activity_proto_enumTypes = make([]protoimpl.EnumInfo, 3)
var file_store_activity_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_store_activity_proto_goTypes = []interface{}{
	(ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Type)(0),   // 0: bytebase.store.ActivityIssueCommentCreatePayload.ExternalApprovalEvent.Type
	(ActivityIssueCommentCreatePayload_ExternalApprovalEvent_Action)(0), // 1: bytebase.store.ActivityIssueCommentCreatePayload.ExternalApprovalEvent.Action
	(ActivityIssueCommentCreatePayload_ApprovalEvent_Status)(0),         // 2: bytebase.store.ActivityIssueCommentCreatePayload.ApprovalEvent.Status
	(*ActivityIssueCreatePayload)(nil),                                  // 3: bytebase.store.ActivityIssueCreatePayload
	(*ActivityIssueCommentCreatePayload)(nil),                           // 4: bytebase.store.ActivityIssueCommentCreatePayload
	(*ActivityIssueCommentCreatePayload_TaskRollbackBy)(nil),            // 5: bytebase.store.ActivityIssueCommentCreatePayload.TaskRollbackBy
	(*ActivityIssueCommentCreatePayload_ExternalApprovalEvent)(nil),     // 6: bytebase.store.ActivityIssueCommentCreatePayload.ExternalApprovalEvent
	(*ActivityIssueCommentCreatePayload_ApprovalEvent)(nil),             // 7: bytebase.store.ActivityIssueCommentCreatePayload.ApprovalEvent
}
var file_store_activity_proto_depIdxs = []int32{
	6, // 0: bytebase.store.ActivityIssueCommentCreatePayload.external_approval_event:type_name -> bytebase.store.ActivityIssueCommentCreatePayload.ExternalApprovalEvent
	5, // 1: bytebase.store.ActivityIssueCommentCreatePayload.task_rollback_by:type_name -> bytebase.store.ActivityIssueCommentCreatePayload.TaskRollbackBy
	7, // 2: bytebase.store.ActivityIssueCommentCreatePayload.approval_event:type_name -> bytebase.store.ActivityIssueCommentCreatePayload.ApprovalEvent
	0, // 3: bytebase.store.ActivityIssueCommentCreatePayload.ExternalApprovalEvent.type:type_name -> bytebase.store.ActivityIssueCommentCreatePayload.ExternalApprovalEvent.Type
	1, // 4: bytebase.store.ActivityIssueCommentCreatePayload.ExternalApprovalEvent.action:type_name -> bytebase.store.ActivityIssueCommentCreatePayload.ExternalApprovalEvent.Action
	2, // 5: bytebase.store.ActivityIssueCommentCreatePayload.ApprovalEvent.status:type_name -> bytebase.store.ActivityIssueCommentCreatePayload.ApprovalEvent.Status
	6, // [6:6] is the sub-list for method output_type
	6, // [6:6] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_store_activity_proto_init() }
func file_store_activity_proto_init() {
	if File_store_activity_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_store_activity_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivityIssueCreatePayload); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_store_activity_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivityIssueCommentCreatePayload); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_store_activity_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivityIssueCommentCreatePayload_TaskRollbackBy); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_store_activity_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivityIssueCommentCreatePayload_ExternalApprovalEvent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_store_activity_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ActivityIssueCommentCreatePayload_ApprovalEvent); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_store_activity_proto_msgTypes[1].OneofWrappers = []interface{}{
		(*ActivityIssueCommentCreatePayload_ExternalApprovalEvent_)(nil),
		(*ActivityIssueCommentCreatePayload_TaskRollbackBy_)(nil),
		(*ActivityIssueCommentCreatePayload_ApprovalEvent_)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_store_activity_proto_rawDesc,
			NumEnums:      3,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_store_activity_proto_goTypes,
		DependencyIndexes: file_store_activity_proto_depIdxs,
		EnumInfos:         file_store_activity_proto_enumTypes,
		MessageInfos:      file_store_activity_proto_msgTypes,
	}.Build()
	File_store_activity_proto = out.File
	file_store_activity_proto_rawDesc = nil
	file_store_activity_proto_goTypes = nil
	file_store_activity_proto_depIdxs = nil
}
