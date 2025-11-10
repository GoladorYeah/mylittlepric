// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'saved_search.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_SavedSearch _$SavedSearchFromJson(Map<String, dynamic> json) => _SavedSearch(
  id: json['id'] as String,
  messages: (json['messages'] as List<dynamic>)
      .map((e) => ChatMessage.fromJson(e as Map<String, dynamic>))
      .toList(),
  sessionId: json['session_id'] as String,
  category: json['category'] as String,
  timestamp: (json['timestamp'] as num).toInt(),
);

Map<String, dynamic> _$SavedSearchToJson(_SavedSearch instance) =>
    <String, dynamic>{
      'id': instance.id,
      'messages': instance.messages,
      'session_id': instance.sessionId,
      'category': instance.category,
      'timestamp': instance.timestamp,
    };
