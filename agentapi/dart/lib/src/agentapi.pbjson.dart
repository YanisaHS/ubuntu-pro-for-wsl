//
//  Generated code. Do not modify.
//  source: agentapi.proto
//
// @dart = 2.12

// ignore_for_file: annotate_overrides, camel_case_types
// ignore_for_file: constant_identifier_names, library_prefixes
// ignore_for_file: non_constant_identifier_names, prefer_final_fields
// ignore_for_file: unnecessary_import, unnecessary_this, unused_import

import 'dart:convert' as $convert;
import 'dart:core' as $core;
import 'dart:typed_data' as $typed_data;

@$core.Deprecated('Use proAttachInfoDescriptor instead')
const ProAttachInfo$json = {
  '1': 'ProAttachInfo',
  '2': [
    {'1': 'token', '3': 1, '4': 1, '5': 9, '10': 'token'},
  ],
};

/// Descriptor for `ProAttachInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List proAttachInfoDescriptor = $convert.base64Decode(
    'Cg1Qcm9BdHRhY2hJbmZvEhQKBXRva2VuGAEgASgJUgV0b2tlbg==');

@$core.Deprecated('Use subscriptionInfoDescriptor instead')
const SubscriptionInfo$json = {
  '1': 'SubscriptionInfo',
  '2': [
    {'1': 'productId', '3': 1, '4': 1, '5': 9, '10': 'productId'},
    {'1': 'immutable', '3': 2, '4': 1, '5': 8, '10': 'immutable'},
    {'1': 'none', '3': 3, '4': 1, '5': 11, '6': '.google.protobuf.Empty', '9': 0, '10': 'none'},
    {'1': 'user', '3': 4, '4': 1, '5': 11, '6': '.google.protobuf.Empty', '9': 0, '10': 'user'},
    {'1': 'organization', '3': 5, '4': 1, '5': 11, '6': '.google.protobuf.Empty', '9': 0, '10': 'organization'},
    {'1': 'microsoftStore', '3': 6, '4': 1, '5': 11, '6': '.google.protobuf.Empty', '9': 0, '10': 'microsoftStore'},
  ],
  '8': [
    {'1': 'subscriptionType'},
  ],
};

/// Descriptor for `SubscriptionInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List subscriptionInfoDescriptor = $convert.base64Decode(
    'ChBTdWJzY3JpcHRpb25JbmZvEhwKCXByb2R1Y3RJZBgBIAEoCVIJcHJvZHVjdElkEhwKCWltbX'
    'V0YWJsZRgCIAEoCFIJaW1tdXRhYmxlEiwKBG5vbmUYAyABKAsyFi5nb29nbGUucHJvdG9idWYu'
    'RW1wdHlIAFIEbm9uZRIsCgR1c2VyGAQgASgLMhYuZ29vZ2xlLnByb3RvYnVmLkVtcHR5SABSBH'
    'VzZXISPAoMb3JnYW5pemF0aW9uGAUgASgLMhYuZ29vZ2xlLnByb3RvYnVmLkVtcHR5SABSDG9y'
    'Z2FuaXphdGlvbhJACg5taWNyb3NvZnRTdG9yZRgGIAEoCzIWLmdvb2dsZS5wcm90b2J1Zi5FbX'
    'B0eUgAUg5taWNyb3NvZnRTdG9yZUISChBzdWJzY3JpcHRpb25UeXBl');

@$core.Deprecated('Use distroInfoDescriptor instead')
const DistroInfo$json = {
  '1': 'DistroInfo',
  '2': [
    {'1': 'wsl_name', '3': 1, '4': 1, '5': 9, '10': 'wslName'},
    {'1': 'id', '3': 2, '4': 1, '5': 9, '10': 'id'},
    {'1': 'version_id', '3': 3, '4': 1, '5': 9, '10': 'versionId'},
    {'1': 'pretty_name', '3': 4, '4': 1, '5': 9, '10': 'prettyName'},
    {'1': 'pro_attached', '3': 5, '4': 1, '5': 8, '10': 'proAttached'},
    {'1': 'hostname', '3': 6, '4': 1, '5': 9, '10': 'hostname'},
  ],
};

/// Descriptor for `DistroInfo`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List distroInfoDescriptor = $convert.base64Decode(
    'CgpEaXN0cm9JbmZvEhkKCHdzbF9uYW1lGAEgASgJUgd3c2xOYW1lEg4KAmlkGAIgASgJUgJpZB'
    'IdCgp2ZXJzaW9uX2lkGAMgASgJUgl2ZXJzaW9uSWQSHwoLcHJldHR5X25hbWUYBCABKAlSCnBy'
    'ZXR0eU5hbWUSIQoMcHJvX2F0dGFjaGVkGAUgASgIUgtwcm9BdHRhY2hlZBIaCghob3N0bmFtZR'
    'gGIAEoCVIIaG9zdG5hbWU=');

@$core.Deprecated('Use portDescriptor instead')
const Port$json = {
  '1': 'Port',
  '2': [
    {'1': 'port', '3': 1, '4': 1, '5': 13, '10': 'port'},
  ],
};

/// Descriptor for `Port`. Decode as a `google.protobuf.DescriptorProto`.
final $typed_data.Uint8List portDescriptor = $convert.base64Decode(
    'CgRQb3J0EhIKBHBvcnQYASABKA1SBHBvcnQ=');

