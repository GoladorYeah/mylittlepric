// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'session.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_SearchState _$SearchStateFromJson(Map<String, dynamic> json) => _SearchState(
  category: json['category'] as String?,
  status: json['status'] as String?,
  lastProduct: json['last_product'] == null
      ? null
      : LastProduct.fromJson(json['last_product'] as Map<String, dynamic>),
);

Map<String, dynamic> _$SearchStateToJson(_SearchState instance) =>
    <String, dynamic>{
      'category': instance.category,
      'status': instance.status,
      'last_product': instance.lastProduct,
    };

_LastProduct _$LastProductFromJson(Map<String, dynamic> json) =>
    _LastProduct(name: json['name'] as String, price: json['price'] as String);

Map<String, dynamic> _$LastProductToJson(_LastProduct instance) =>
    <String, dynamic>{'name': instance.name, 'price': instance.price};

_SessionResponse _$SessionResponseFromJson(Map<String, dynamic> json) =>
    _SessionResponse(
      sessionId: json['session_id'] as String,
      messages: (json['messages'] as List<dynamic>)
          .map((e) => SessionMessage.fromJson(e as Map<String, dynamic>))
          .toList(),
      searchState: json['search_state'] == null
          ? null
          : SearchState.fromJson(json['search_state'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$SessionResponseToJson(_SessionResponse instance) =>
    <String, dynamic>{
      'session_id': instance.sessionId,
      'messages': instance.messages,
      'search_state': instance.searchState,
    };
