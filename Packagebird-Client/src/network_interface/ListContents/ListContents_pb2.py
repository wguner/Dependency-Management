# -*- coding: utf-8 -*-
# Generated by the protocol buffer compiler.  DO NOT EDIT!
# source: ListContents.proto
"""Generated protocol buffer code."""
from google.protobuf import descriptor as _descriptor
from google.protobuf import message as _message
from google.protobuf import reflection as _reflection
from google.protobuf import symbol_database as _symbol_database
# @@protoc_insertion_point(imports)

_sym_db = _symbol_database.Default()




DESCRIPTOR = _descriptor.FileDescriptor(
  name='ListContents.proto',
  package='listcontent',
  syntax='proto3',
  serialized_options=b'Z\r./listcontent',
  create_key=_descriptor._internal_create_key,
  serialized_pb=b'\n\x12ListContents.proto\x12\x0blistcontent\"Q\n\x0e\x43ontentRequest\x12\x14\n\x0cListProjects\x18\x01 \x01(\x08\x12\x14\n\x0cListPackages\x18\x02 \x01(\x08\x12\x13\n\x0bListMembers\x18\x03 \x01(\x08\"F\n\x0f\x43ontentResponse\x12\x10\n\x08Projects\x18\x01 \x03(\t\x12\x10\n\x08Packages\x18\x02 \x03(\t\x12\x0f\n\x07Members\x18\x03 \x03(\t2^\n\x13ListContentServices\x12G\n\nGetContent\x12\x1b.listcontent.ContentRequest\x1a\x1c.listcontent.ContentResponseB\x0fZ\r./listcontentb\x06proto3'
)




_CONTENTREQUEST = _descriptor.Descriptor(
  name='ContentRequest',
  full_name='listcontent.ContentRequest',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='ListProjects', full_name='listcontent.ContentRequest.ListProjects', index=0,
      number=1, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='ListPackages', full_name='listcontent.ContentRequest.ListPackages', index=1,
      number=2, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='ListMembers', full_name='listcontent.ContentRequest.ListMembers', index=2,
      number=3, type=8, cpp_type=7, label=1,
      has_default_value=False, default_value=False,
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=35,
  serialized_end=116,
)


_CONTENTRESPONSE = _descriptor.Descriptor(
  name='ContentResponse',
  full_name='listcontent.ContentResponse',
  filename=None,
  file=DESCRIPTOR,
  containing_type=None,
  create_key=_descriptor._internal_create_key,
  fields=[
    _descriptor.FieldDescriptor(
      name='Projects', full_name='listcontent.ContentResponse.Projects', index=0,
      number=1, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='Packages', full_name='listcontent.ContentResponse.Packages', index=1,
      number=2, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
    _descriptor.FieldDescriptor(
      name='Members', full_name='listcontent.ContentResponse.Members', index=2,
      number=3, type=9, cpp_type=9, label=3,
      has_default_value=False, default_value=[],
      message_type=None, enum_type=None, containing_type=None,
      is_extension=False, extension_scope=None,
      serialized_options=None, file=DESCRIPTOR,  create_key=_descriptor._internal_create_key),
  ],
  extensions=[
  ],
  nested_types=[],
  enum_types=[
  ],
  serialized_options=None,
  is_extendable=False,
  syntax='proto3',
  extension_ranges=[],
  oneofs=[
  ],
  serialized_start=118,
  serialized_end=188,
)

DESCRIPTOR.message_types_by_name['ContentRequest'] = _CONTENTREQUEST
DESCRIPTOR.message_types_by_name['ContentResponse'] = _CONTENTRESPONSE
_sym_db.RegisterFileDescriptor(DESCRIPTOR)

ContentRequest = _reflection.GeneratedProtocolMessageType('ContentRequest', (_message.Message,), {
  'DESCRIPTOR' : _CONTENTREQUEST,
  '__module__' : 'ListContents_pb2'
  # @@protoc_insertion_point(class_scope:listcontent.ContentRequest)
  })
_sym_db.RegisterMessage(ContentRequest)

ContentResponse = _reflection.GeneratedProtocolMessageType('ContentResponse', (_message.Message,), {
  'DESCRIPTOR' : _CONTENTRESPONSE,
  '__module__' : 'ListContents_pb2'
  # @@protoc_insertion_point(class_scope:listcontent.ContentResponse)
  })
_sym_db.RegisterMessage(ContentResponse)


DESCRIPTOR._options = None

_LISTCONTENTSERVICES = _descriptor.ServiceDescriptor(
  name='ListContentServices',
  full_name='listcontent.ListContentServices',
  file=DESCRIPTOR,
  index=0,
  serialized_options=None,
  create_key=_descriptor._internal_create_key,
  serialized_start=190,
  serialized_end=284,
  methods=[
  _descriptor.MethodDescriptor(
    name='GetContent',
    full_name='listcontent.ListContentServices.GetContent',
    index=0,
    containing_service=None,
    input_type=_CONTENTREQUEST,
    output_type=_CONTENTRESPONSE,
    serialized_options=None,
    create_key=_descriptor._internal_create_key,
  ),
])
_sym_db.RegisterServiceDescriptor(_LISTCONTENTSERVICES)

DESCRIPTOR.services_by_name['ListContentServices'] = _LISTCONTENTSERVICES

# @@protoc_insertion_point(module_scope)