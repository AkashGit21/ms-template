// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: internal/proto-files/movie.proto

package moviepb

import (
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type Tag int32

const (
	Tag_UNDEFINED_TAG Tag = 0
	Tag_Action        Tag = 1
	Tag_Adventure     Tag = 2
	Tag_Fantasy       Tag = 3
)

// Enum value maps for Tag.
var (
	Tag_name = map[int32]string{
		0: "UNDEFINED_TAG",
		1: "Action",
		2: "Adventure",
		3: "Fantasy",
	}
	Tag_value = map[string]int32{
		"UNDEFINED_TAG": 0,
		"Action":        1,
		"Adventure":     2,
		"Fantasy":       3,
	}
)

func (x Tag) Enum() *Tag {
	p := new(Tag)
	*p = x
	return p
}

func (x Tag) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Tag) Descriptor() protoreflect.EnumDescriptor {
	return file_internal_proto_files_movie_proto_enumTypes[0].Descriptor()
}

func (Tag) Type() protoreflect.EnumType {
	return &file_internal_proto_files_movie_proto_enumTypes[0]
}

func (x Tag) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Tag.Descriptor instead.
func (Tag) EnumDescriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{0}
}

type ListMoviesRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The page_token for that page to be uniquely identified
	PageToken string `protobuf:"bytes,1,opt,name=page_token,json=pageToken,proto3" json:"page_token,omitempty"`
	// The maximum number of objects to return. Server may return fewer objects
	// than requested. If unspecified, server will pick an appropriate default.
	PageSize int32 `protobuf:"varint,2,opt,name=page_size,json=pageSize,proto3" json:"page_size,omitempty"`
}

func (x *ListMoviesRequest) Reset() {
	*x = ListMoviesRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMoviesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMoviesRequest) ProtoMessage() {}

func (x *ListMoviesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMoviesRequest.ProtoReflect.Descriptor instead.
func (*ListMoviesRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{0}
}

func (x *ListMoviesRequest) GetPageToken() string {
	if x != nil {
		return x.PageToken
	}
	return ""
}

func (x *ListMoviesRequest) GetPageSize() int32 {
	if x != nil {
		return x.PageSize
	}
	return 0
}

type ListMoviesResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Movies        []*Movie `protobuf:"bytes,1,rep,name=movies,proto3" json:"movies,omitempty"`
	NextPageToken string   `protobuf:"bytes,2,opt,name=next_page_token,json=nextPageToken,proto3" json:"next_page_token,omitempty"`
}

func (x *ListMoviesResponse) Reset() {
	*x = ListMoviesResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListMoviesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListMoviesResponse) ProtoMessage() {}

func (x *ListMoviesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListMoviesResponse.ProtoReflect.Descriptor instead.
func (*ListMoviesResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{1}
}

func (x *ListMoviesResponse) GetMovies() []*Movie {
	if x != nil {
		return x.Movies
	}
	return nil
}

func (x *ListMoviesResponse) GetNextPageToken() string {
	if x != nil {
		return x.NextPageToken
	}
	return ""
}

type GetMovieRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *GetMovieRequest) Reset() {
	*x = GetMovieRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetMovieRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetMovieRequest) ProtoMessage() {}

func (x *GetMovieRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetMovieRequest.ProtoReflect.Descriptor instead.
func (*GetMovieRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{2}
}

func (x *GetMovieRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type CreateMovieRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Movie *Movie `protobuf:"bytes,1,opt,name=movie,proto3" json:"movie,omitempty"`
}

func (x *CreateMovieRequest) Reset() {
	*x = CreateMovieRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMovieRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMovieRequest) ProtoMessage() {}

func (x *CreateMovieRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMovieRequest.ProtoReflect.Descriptor instead.
func (*CreateMovieRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{3}
}

func (x *CreateMovieRequest) GetMovie() *Movie {
	if x != nil {
		return x.Movie
	}
	return nil
}

type CreateMovieResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *CreateMovieResponse) Reset() {
	*x = CreateMovieResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateMovieResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateMovieResponse) ProtoMessage() {}

func (x *CreateMovieResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateMovieResponse.ProtoReflect.Descriptor instead.
func (*CreateMovieResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{4}
}

func (x *CreateMovieResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type UpdateMovieRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id    string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Movie *Movie `protobuf:"bytes,2,opt,name=movie,proto3" json:"movie,omitempty"`
}

func (x *UpdateMovieRequest) Reset() {
	*x = UpdateMovieRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMovieRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMovieRequest) ProtoMessage() {}

func (x *UpdateMovieRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMovieRequest.ProtoReflect.Descriptor instead.
func (*UpdateMovieRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{5}
}

func (x *UpdateMovieRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *UpdateMovieRequest) GetMovie() *Movie {
	if x != nil {
		return x.Movie
	}
	return nil
}

type UpdateMovieResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *UpdateMovieResponse) Reset() {
	*x = UpdateMovieResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateMovieResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateMovieResponse) ProtoMessage() {}

func (x *UpdateMovieResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateMovieResponse.ProtoReflect.Descriptor instead.
func (*UpdateMovieResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{6}
}

func (x *UpdateMovieResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type PartialUpdateMovieRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique ID for the Movie
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Brief summary of the Movie
	Summary string `protobuf:"bytes,3,opt,name=summary,proto3" json:"summary,omitempty"`
	// Group of actors who make up a Film or stage play
	Cast []string `protobuf:"bytes,5,rep,name=cast,proto3" json:"cast,omitempty"`
	// Tags related to the Movie
	Tags []Tag `protobuf:"varint,6,rep,packed,name=tags,proto3,enum=movie.Tag" json:"tags,omitempty"`
	// Director of film
	Director string `protobuf:"bytes,7,opt,name=director,proto3" json:"director,omitempty"`
	// The author(s) of film
	Writers []string `protobuf:"bytes,8,rep,name=writers,proto3" json:"writers,omitempty"`
}

func (x *PartialUpdateMovieRequest) Reset() {
	*x = PartialUpdateMovieRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartialUpdateMovieRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartialUpdateMovieRequest) ProtoMessage() {}

func (x *PartialUpdateMovieRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartialUpdateMovieRequest.ProtoReflect.Descriptor instead.
func (*PartialUpdateMovieRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{7}
}

func (x *PartialUpdateMovieRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *PartialUpdateMovieRequest) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *PartialUpdateMovieRequest) GetCast() []string {
	if x != nil {
		return x.Cast
	}
	return nil
}

func (x *PartialUpdateMovieRequest) GetTags() []Tag {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *PartialUpdateMovieRequest) GetDirector() string {
	if x != nil {
		return x.Director
	}
	return ""
}

func (x *PartialUpdateMovieRequest) GetWriters() []string {
	if x != nil {
		return x.Writers
	}
	return nil
}

type PartialUpdateMovieResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *PartialUpdateMovieResponse) Reset() {
	*x = PartialUpdateMovieResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartialUpdateMovieResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartialUpdateMovieResponse) ProtoMessage() {}

func (x *PartialUpdateMovieResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartialUpdateMovieResponse.ProtoReflect.Descriptor instead.
func (*PartialUpdateMovieResponse) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{8}
}

func (x *PartialUpdateMovieResponse) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type DeleteMovieRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
}

func (x *DeleteMovieRequest) Reset() {
	*x = DeleteMovieRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeleteMovieRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeleteMovieRequest) ProtoMessage() {}

func (x *DeleteMovieRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeleteMovieRequest.ProtoReflect.Descriptor instead.
func (*DeleteMovieRequest) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{9}
}

func (x *DeleteMovieRequest) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type Movie struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unique ID for the Movie
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Name of the Movie
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Brief summary of the Movie
	Summary string `protobuf:"bytes,3,opt,name=summary,proto3" json:"summary,omitempty"`
	// Group of actors who make up a Film or stage play
	Cast []string `protobuf:"bytes,5,rep,name=cast,proto3" json:"cast,omitempty"`
	// Tags related to the Movie
	Tags []Tag `protobuf:"varint,6,rep,packed,name=tags,proto3,enum=movie.Tag" json:"tags,omitempty"`
	// Director of film
	Director string `protobuf:"bytes,7,opt,name=director,proto3" json:"director,omitempty"`
	// The author(s) of film
	Writers []string `protobuf:"bytes,8,rep,name=writers,proto3" json:"writers,omitempty"`
}

func (x *Movie) Reset() {
	*x = Movie{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_proto_files_movie_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Movie) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Movie) ProtoMessage() {}

func (x *Movie) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_files_movie_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Movie.ProtoReflect.Descriptor instead.
func (*Movie) Descriptor() ([]byte, []int) {
	return file_internal_proto_files_movie_proto_rawDescGZIP(), []int{10}
}

func (x *Movie) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Movie) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Movie) GetSummary() string {
	if x != nil {
		return x.Summary
	}
	return ""
}

func (x *Movie) GetCast() []string {
	if x != nil {
		return x.Cast
	}
	return nil
}

func (x *Movie) GetTags() []Tag {
	if x != nil {
		return x.Tags
	}
	return nil
}

func (x *Movie) GetDirector() string {
	if x != nil {
		return x.Director
	}
	return ""
}

func (x *Movie) GetWriters() []string {
	if x != nil {
		return x.Writers
	}
	return nil
}

var File_internal_proto_files_movie_proto protoreflect.FileDescriptor

var file_internal_proto_files_movie_proto_rawDesc = []byte{
	0x0a, 0x20, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2d, 0x66, 0x69, 0x6c, 0x65, 0x73, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x1a, 0x23, 0x70, 0x6b, 0x67, 0x2f, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x5f,
	0x62, 0x65, 0x68, 0x61, 0x76, 0x69, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x20,
	0x70, 0x6b, 0x67, 0x2f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x61,
	0x6e, 0x6e, 0x6f, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x4f, 0x0a,
	0x11, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x70, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b, 0x65,
	0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x73, 0x69, 0x7a, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x61, 0x67, 0x65, 0x53, 0x69, 0x7a, 0x65, 0x22, 0x62,
	0x0a, 0x12, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x24, 0x0a, 0x06, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x4d, 0x6f, 0x76,
	0x69, 0x65, 0x52, 0x06, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x12, 0x26, 0x0a, 0x0f, 0x6e, 0x65,
	0x78, 0x74, 0x5f, 0x70, 0x61, 0x67, 0x65, 0x5f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0d, 0x6e, 0x65, 0x78, 0x74, 0x50, 0x61, 0x67, 0x65, 0x54, 0x6f, 0x6b,
	0x65, 0x6e, 0x22, 0x21, 0x0a, 0x0f, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x38, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d,
	0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x22, 0x0a, 0x05, 0x6d,
	0x6f, 0x76, 0x69, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x2e, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x05, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x22,
	0x25, 0x0a, 0x13, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0x48, 0x0a, 0x12, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x22, 0x0a, 0x05,
	0x6d, 0x6f, 0x76, 0x69, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x6d, 0x6f,
	0x76, 0x69, 0x65, 0x2e, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x05, 0x6d, 0x6f, 0x76, 0x69, 0x65,
	0x22, 0x25, 0x0a, 0x13, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xb9, 0x01, 0x0a, 0x19, 0x50, 0x61, 0x72, 0x74,
	0x69, 0x61, 0x6c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1d, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x07, 0x73, 0x75, 0x6d,
	0x6d, 0x61, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x61, 0x73, 0x74, 0x18, 0x05, 0x20, 0x03,
	0x28, 0x09, 0x52, 0x04, 0x63, 0x61, 0x73, 0x74, 0x12, 0x23, 0x0a, 0x04, 0x74, 0x61, 0x67, 0x73,
	0x18, 0x06, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x54,
	0x61, 0x67, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x74, 0x61, 0x67, 0x73, 0x12, 0x1a, 0x0a,
	0x08, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x18, 0x0a, 0x07, 0x77, 0x72, 0x69,
	0x74, 0x65, 0x72, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x77, 0x72, 0x69, 0x74,
	0x65, 0x72, 0x73, 0x22, 0x2c, 0x0a, 0x1a, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x55, 0x70,
	0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69,
	0x64, 0x22, 0x24, 0x0a, 0x12, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x22, 0xbe, 0x01, 0x0a, 0x05, 0x4d, 0x6f, 0x76, 0x69,
	0x65, 0x12, 0x13, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0,
	0x41, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x17, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x07, 0x73, 0x75, 0x6d, 0x6d, 0x61, 0x72, 0x79, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x61, 0x73,
	0x74, 0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x63, 0x61, 0x73, 0x74, 0x12, 0x23, 0x0a,
	0x04, 0x74, 0x61, 0x67, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0e, 0x32, 0x0a, 0x2e, 0x6d, 0x6f,
	0x76, 0x69, 0x65, 0x2e, 0x54, 0x61, 0x67, 0x42, 0x03, 0xe0, 0x41, 0x02, 0x52, 0x04, 0x74, 0x61,
	0x67, 0x73, 0x12, 0x1a, 0x0a, 0x08, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x18,
	0x0a, 0x07, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x73, 0x18, 0x08, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x07, 0x77, 0x72, 0x69, 0x74, 0x65, 0x72, 0x73, 0x2a, 0x40, 0x0a, 0x03, 0x54, 0x61, 0x67, 0x12,
	0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x5f, 0x54, 0x41, 0x47,
	0x10, 0x00, 0x12, 0x0a, 0x0a, 0x06, 0x41, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x10, 0x01, 0x12, 0x0d,
	0x0a, 0x09, 0x41, 0x64, 0x76, 0x65, 0x6e, 0x74, 0x75, 0x72, 0x65, 0x10, 0x02, 0x12, 0x0b, 0x0a,
	0x07, 0x46, 0x61, 0x6e, 0x74, 0x61, 0x73, 0x79, 0x10, 0x03, 0x32, 0xc9, 0x04, 0x0a, 0x0c, 0x4d,
	0x6f, 0x76, 0x69, 0x65, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x55, 0x0a, 0x0a, 0x4c,
	0x69, 0x73, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x12, 0x18, 0x2e, 0x6d, 0x6f, 0x76, 0x69,
	0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x19, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74,
	0x4d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x12,
	0x82, 0xd3, 0xe4, 0x93, 0x02, 0x0c, 0x12, 0x0a, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x69,
	0x65, 0x73, 0x12, 0x49, 0x0a, 0x08, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x16,
	0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x47, 0x65, 0x74, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0c, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x4d,
	0x6f, 0x76, 0x69, 0x65, 0x22, 0x17, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x11, 0x12, 0x0f, 0x2f, 0x76,
	0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x12, 0x5f, 0x0a,
	0x0b, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x19, 0x2e, 0x6d,
	0x6f, 0x76, 0x69, 0x65, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x22, 0x19, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x13, 0x22, 0x0a, 0x2f, 0x76, 0x31,
	0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x3a, 0x05, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x64,
	0x0a, 0x0b, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x19, 0x2e,
	0x6d, 0x6f, 0x76, 0x69, 0x65, 0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1a, 0x2e, 0x6d, 0x6f, 0x76, 0x69, 0x65,
	0x2e, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x1e, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x18, 0x1a, 0x0f, 0x2f, 0x76,
	0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x3a, 0x05, 0x6d,
	0x6f, 0x76, 0x69, 0x65, 0x12, 0x75, 0x0a, 0x12, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x55,
	0x70, 0x64, 0x61, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x20, 0x2e, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x21, 0x2e, 0x6d,
	0x6f, 0x76, 0x69, 0x65, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x55, 0x70, 0x64, 0x61,
	0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x1a, 0x82, 0xd3, 0xe4, 0x93, 0x02, 0x14, 0x32, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x3a, 0x01, 0x2a, 0x12, 0x59, 0x0a, 0x0b, 0x44,
	0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x12, 0x19, 0x2e, 0x6d, 0x6f, 0x76,
	0x69, 0x65, 0x2e, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x4d, 0x6f, 0x76, 0x69, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x17, 0x82,
	0xd3, 0xe4, 0x93, 0x02, 0x11, 0x2a, 0x0f, 0x2f, 0x76, 0x31, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65,
	0x73, 0x2f, 0x7b, 0x69, 0x64, 0x7d, 0x42, 0x1e, 0x5a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e,
	0x61, 0x6c, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x2f, 0x6d, 0x6f, 0x76, 0x69, 0x65, 0x3b, 0x20, 0x6d,
	0x6f, 0x76, 0x69, 0x65, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_proto_files_movie_proto_rawDescOnce sync.Once
	file_internal_proto_files_movie_proto_rawDescData = file_internal_proto_files_movie_proto_rawDesc
)

func file_internal_proto_files_movie_proto_rawDescGZIP() []byte {
	file_internal_proto_files_movie_proto_rawDescOnce.Do(func() {
		file_internal_proto_files_movie_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_proto_files_movie_proto_rawDescData)
	})
	return file_internal_proto_files_movie_proto_rawDescData
}

var file_internal_proto_files_movie_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_internal_proto_files_movie_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_internal_proto_files_movie_proto_goTypes = []interface{}{
	(Tag)(0),                           // 0: movie.Tag
	(*ListMoviesRequest)(nil),          // 1: movie.ListMoviesRequest
	(*ListMoviesResponse)(nil),         // 2: movie.ListMoviesResponse
	(*GetMovieRequest)(nil),            // 3: movie.GetMovieRequest
	(*CreateMovieRequest)(nil),         // 4: movie.CreateMovieRequest
	(*CreateMovieResponse)(nil),        // 5: movie.CreateMovieResponse
	(*UpdateMovieRequest)(nil),         // 6: movie.UpdateMovieRequest
	(*UpdateMovieResponse)(nil),        // 7: movie.UpdateMovieResponse
	(*PartialUpdateMovieRequest)(nil),  // 8: movie.PartialUpdateMovieRequest
	(*PartialUpdateMovieResponse)(nil), // 9: movie.PartialUpdateMovieResponse
	(*DeleteMovieRequest)(nil),         // 10: movie.DeleteMovieRequest
	(*Movie)(nil),                      // 11: movie.Movie
	(*empty.Empty)(nil),                // 12: google.protobuf.Empty
}
var file_internal_proto_files_movie_proto_depIdxs = []int32{
	11, // 0: movie.ListMoviesResponse.movies:type_name -> movie.Movie
	11, // 1: movie.CreateMovieRequest.movie:type_name -> movie.Movie
	11, // 2: movie.UpdateMovieRequest.movie:type_name -> movie.Movie
	0,  // 3: movie.PartialUpdateMovieRequest.tags:type_name -> movie.Tag
	0,  // 4: movie.Movie.tags:type_name -> movie.Tag
	1,  // 5: movie.MovieService.ListMovies:input_type -> movie.ListMoviesRequest
	3,  // 6: movie.MovieService.GetMovie:input_type -> movie.GetMovieRequest
	4,  // 7: movie.MovieService.CreateMovie:input_type -> movie.CreateMovieRequest
	6,  // 8: movie.MovieService.UpdateMovie:input_type -> movie.UpdateMovieRequest
	8,  // 9: movie.MovieService.PartialUpdateMovie:input_type -> movie.PartialUpdateMovieRequest
	10, // 10: movie.MovieService.DeleteMovie:input_type -> movie.DeleteMovieRequest
	2,  // 11: movie.MovieService.ListMovies:output_type -> movie.ListMoviesResponse
	11, // 12: movie.MovieService.GetMovie:output_type -> movie.Movie
	5,  // 13: movie.MovieService.CreateMovie:output_type -> movie.CreateMovieResponse
	7,  // 14: movie.MovieService.UpdateMovie:output_type -> movie.UpdateMovieResponse
	9,  // 15: movie.MovieService.PartialUpdateMovie:output_type -> movie.PartialUpdateMovieResponse
	12, // 16: movie.MovieService.DeleteMovie:output_type -> google.protobuf.Empty
	11, // [11:17] is the sub-list for method output_type
	5,  // [5:11] is the sub-list for method input_type
	5,  // [5:5] is the sub-list for extension type_name
	5,  // [5:5] is the sub-list for extension extendee
	0,  // [0:5] is the sub-list for field type_name
}

func init() { file_internal_proto_files_movie_proto_init() }
func file_internal_proto_files_movie_proto_init() {
	if File_internal_proto_files_movie_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_proto_files_movie_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListMoviesRequest); i {
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
		file_internal_proto_files_movie_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListMoviesResponse); i {
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
		file_internal_proto_files_movie_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetMovieRequest); i {
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
		file_internal_proto_files_movie_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMovieRequest); i {
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
		file_internal_proto_files_movie_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateMovieResponse); i {
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
		file_internal_proto_files_movie_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMovieRequest); i {
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
		file_internal_proto_files_movie_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateMovieResponse); i {
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
		file_internal_proto_files_movie_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartialUpdateMovieRequest); i {
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
		file_internal_proto_files_movie_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartialUpdateMovieResponse); i {
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
		file_internal_proto_files_movie_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeleteMovieRequest); i {
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
		file_internal_proto_files_movie_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Movie); i {
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
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_internal_proto_files_movie_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_files_movie_proto_goTypes,
		DependencyIndexes: file_internal_proto_files_movie_proto_depIdxs,
		EnumInfos:         file_internal_proto_files_movie_proto_enumTypes,
		MessageInfos:      file_internal_proto_files_movie_proto_msgTypes,
	}.Build()
	File_internal_proto_files_movie_proto = out.File
	file_internal_proto_files_movie_proto_rawDesc = nil
	file_internal_proto_files_movie_proto_goTypes = nil
	file_internal_proto_files_movie_proto_depIdxs = nil
}
