// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'chat_response.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_ChatResponse _$ChatResponseFromJson(Map<String, dynamic> json) =>
    _ChatResponse(
      sessionId: json['session_id'] as String,
      message: json['message'] as String,
      quickReplies: (json['quick_replies'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      products: (json['products'] as List<dynamic>?)
          ?.map((e) => Product.fromJson(e as Map<String, dynamic>))
          .toList(),
      responseType: json['response_type'] as String?,
      searchType: json['search_type'] as String?,
    );

Map<String, dynamic> _$ChatResponseToJson(_ChatResponse instance) =>
    <String, dynamic>{
      'session_id': instance.sessionId,
      'message': instance.message,
      'quick_replies': instance.quickReplies,
      'products': instance.products,
      'response_type': instance.responseType,
      'search_type': instance.searchType,
    };
