// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'chat_message.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_ChatMessage _$ChatMessageFromJson(Map<String, dynamic> json) => _ChatMessage(
  id: json['id'] as String,
  role: $enumDecode(_$MessageRoleEnumMap, json['role']),
  content: json['content'] as String,
  timestamp: (json['timestamp'] as num).toInt(),
  quickReplies: (json['quick_replies'] as List<dynamic>?)
      ?.map((e) => e as String)
      .toList(),
  products: (json['products'] as List<dynamic>?)
      ?.map((e) => Product.fromJson(e as Map<String, dynamic>))
      .toList(),
  searchType: json['search_type'] as String?,
  isLocal: json['isLocal'] as bool? ?? false,
);

Map<String, dynamic> _$ChatMessageToJson(_ChatMessage instance) =>
    <String, dynamic>{
      'id': instance.id,
      'role': _$MessageRoleEnumMap[instance.role]!,
      'content': instance.content,
      'timestamp': instance.timestamp,
      'quick_replies': instance.quickReplies,
      'products': instance.products,
      'search_type': instance.searchType,
      'isLocal': instance.isLocal,
    };

const _$MessageRoleEnumMap = {
  MessageRole.user: 'user',
  MessageRole.assistant: 'assistant',
};

_SessionMessage _$SessionMessageFromJson(Map<String, dynamic> json) =>
    _SessionMessage(
      role: json['role'] as String,
      content: json['content'] as String,
      timestamp: json['timestamp'] as String?,
      quickReplies: (json['quick_replies'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      products: (json['products'] as List<dynamic>?)
          ?.map((e) => Product.fromJson(e as Map<String, dynamic>))
          .toList(),
      searchType: json['search_type'] as String?,
    );

Map<String, dynamic> _$SessionMessageToJson(_SessionMessage instance) =>
    <String, dynamic>{
      'role': instance.role,
      'content': instance.content,
      'timestamp': instance.timestamp,
      'quick_replies': instance.quickReplies,
      'products': instance.products,
      'search_type': instance.searchType,
    };
