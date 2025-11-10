// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'chat_response.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$ChatResponse {

@JsonKey(name: 'session_id') String get sessionId; String get message;@JsonKey(name: 'quick_replies') List<String>? get quickReplies; List<Product>? get products;@JsonKey(name: 'response_type') String? get responseType;@JsonKey(name: 'search_type') String? get searchType;
/// Create a copy of ChatResponse
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ChatResponseCopyWith<ChatResponse> get copyWith => _$ChatResponseCopyWithImpl<ChatResponse>(this as ChatResponse, _$identity);

  /// Serializes this ChatResponse to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ChatResponse&&(identical(other.sessionId, sessionId) || other.sessionId == sessionId)&&(identical(other.message, message) || other.message == message)&&const DeepCollectionEquality().equals(other.quickReplies, quickReplies)&&const DeepCollectionEquality().equals(other.products, products)&&(identical(other.responseType, responseType) || other.responseType == responseType)&&(identical(other.searchType, searchType) || other.searchType == searchType));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,sessionId,message,const DeepCollectionEquality().hash(quickReplies),const DeepCollectionEquality().hash(products),responseType,searchType);

@override
String toString() {
  return 'ChatResponse(sessionId: $sessionId, message: $message, quickReplies: $quickReplies, products: $products, responseType: $responseType, searchType: $searchType)';
}


}

/// @nodoc
abstract mixin class $ChatResponseCopyWith<$Res>  {
  factory $ChatResponseCopyWith(ChatResponse value, $Res Function(ChatResponse) _then) = _$ChatResponseCopyWithImpl;
@useResult
$Res call({
@JsonKey(name: 'session_id') String sessionId, String message,@JsonKey(name: 'quick_replies') List<String>? quickReplies, List<Product>? products,@JsonKey(name: 'response_type') String? responseType,@JsonKey(name: 'search_type') String? searchType
});




}
/// @nodoc
class _$ChatResponseCopyWithImpl<$Res>
    implements $ChatResponseCopyWith<$Res> {
  _$ChatResponseCopyWithImpl(this._self, this._then);

  final ChatResponse _self;
  final $Res Function(ChatResponse) _then;

/// Create a copy of ChatResponse
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? sessionId = null,Object? message = null,Object? quickReplies = freezed,Object? products = freezed,Object? responseType = freezed,Object? searchType = freezed,}) {
  return _then(_self.copyWith(
sessionId: null == sessionId ? _self.sessionId : sessionId // ignore: cast_nullable_to_non_nullable
as String,message: null == message ? _self.message : message // ignore: cast_nullable_to_non_nullable
as String,quickReplies: freezed == quickReplies ? _self.quickReplies : quickReplies // ignore: cast_nullable_to_non_nullable
as List<String>?,products: freezed == products ? _self.products : products // ignore: cast_nullable_to_non_nullable
as List<Product>?,responseType: freezed == responseType ? _self.responseType : responseType // ignore: cast_nullable_to_non_nullable
as String?,searchType: freezed == searchType ? _self.searchType : searchType // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

}


/// Adds pattern-matching-related methods to [ChatResponse].
extension ChatResponsePatterns on ChatResponse {
/// A variant of `map` that fallback to returning `orElse`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ChatResponse value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ChatResponse() when $default != null:
return $default(_that);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// Callbacks receives the raw object, upcasted.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case final Subclass2 value:
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ChatResponse value)  $default,){
final _that = this;
switch (_that) {
case _ChatResponse():
return $default(_that);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `map` that fallback to returning `null`.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case final Subclass value:
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ChatResponse value)?  $default,){
final _that = this;
switch (_that) {
case _ChatResponse() when $default != null:
return $default(_that);case _:
  return null;

}
}
/// A variant of `when` that fallback to an `orElse` callback.
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return orElse();
/// }
/// ```

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function(@JsonKey(name: 'session_id')  String sessionId,  String message, @JsonKey(name: 'quick_replies')  List<String>? quickReplies,  List<Product>? products, @JsonKey(name: 'response_type')  String? responseType, @JsonKey(name: 'search_type')  String? searchType)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ChatResponse() when $default != null:
return $default(_that.sessionId,_that.message,_that.quickReplies,_that.products,_that.responseType,_that.searchType);case _:
  return orElse();

}
}
/// A `switch`-like method, using callbacks.
///
/// As opposed to `map`, this offers destructuring.
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case Subclass2(:final field2):
///     return ...;
/// }
/// ```

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function(@JsonKey(name: 'session_id')  String sessionId,  String message, @JsonKey(name: 'quick_replies')  List<String>? quickReplies,  List<Product>? products, @JsonKey(name: 'response_type')  String? responseType, @JsonKey(name: 'search_type')  String? searchType)  $default,) {final _that = this;
switch (_that) {
case _ChatResponse():
return $default(_that.sessionId,_that.message,_that.quickReplies,_that.products,_that.responseType,_that.searchType);case _:
  throw StateError('Unexpected subclass');

}
}
/// A variant of `when` that fallback to returning `null`
///
/// It is equivalent to doing:
/// ```dart
/// switch (sealedClass) {
///   case Subclass(:final field):
///     return ...;
///   case _:
///     return null;
/// }
/// ```

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function(@JsonKey(name: 'session_id')  String sessionId,  String message, @JsonKey(name: 'quick_replies')  List<String>? quickReplies,  List<Product>? products, @JsonKey(name: 'response_type')  String? responseType, @JsonKey(name: 'search_type')  String? searchType)?  $default,) {final _that = this;
switch (_that) {
case _ChatResponse() when $default != null:
return $default(_that.sessionId,_that.message,_that.quickReplies,_that.products,_that.responseType,_that.searchType);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ChatResponse implements ChatResponse {
  const _ChatResponse({@JsonKey(name: 'session_id') required this.sessionId, required this.message, @JsonKey(name: 'quick_replies') final  List<String>? quickReplies, final  List<Product>? products, @JsonKey(name: 'response_type') this.responseType, @JsonKey(name: 'search_type') this.searchType}): _quickReplies = quickReplies,_products = products;
  factory _ChatResponse.fromJson(Map<String, dynamic> json) => _$ChatResponseFromJson(json);

@override@JsonKey(name: 'session_id') final  String sessionId;
@override final  String message;
 final  List<String>? _quickReplies;
@override@JsonKey(name: 'quick_replies') List<String>? get quickReplies {
  final value = _quickReplies;
  if (value == null) return null;
  if (_quickReplies is EqualUnmodifiableListView) return _quickReplies;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(value);
}

 final  List<Product>? _products;
@override List<Product>? get products {
  final value = _products;
  if (value == null) return null;
  if (_products is EqualUnmodifiableListView) return _products;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(value);
}

@override@JsonKey(name: 'response_type') final  String? responseType;
@override@JsonKey(name: 'search_type') final  String? searchType;

/// Create a copy of ChatResponse
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ChatResponseCopyWith<_ChatResponse> get copyWith => __$ChatResponseCopyWithImpl<_ChatResponse>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ChatResponseToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ChatResponse&&(identical(other.sessionId, sessionId) || other.sessionId == sessionId)&&(identical(other.message, message) || other.message == message)&&const DeepCollectionEquality().equals(other._quickReplies, _quickReplies)&&const DeepCollectionEquality().equals(other._products, _products)&&(identical(other.responseType, responseType) || other.responseType == responseType)&&(identical(other.searchType, searchType) || other.searchType == searchType));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,sessionId,message,const DeepCollectionEquality().hash(_quickReplies),const DeepCollectionEquality().hash(_products),responseType,searchType);

@override
String toString() {
  return 'ChatResponse(sessionId: $sessionId, message: $message, quickReplies: $quickReplies, products: $products, responseType: $responseType, searchType: $searchType)';
}


}

/// @nodoc
abstract mixin class _$ChatResponseCopyWith<$Res> implements $ChatResponseCopyWith<$Res> {
  factory _$ChatResponseCopyWith(_ChatResponse value, $Res Function(_ChatResponse) _then) = __$ChatResponseCopyWithImpl;
@override @useResult
$Res call({
@JsonKey(name: 'session_id') String sessionId, String message,@JsonKey(name: 'quick_replies') List<String>? quickReplies, List<Product>? products,@JsonKey(name: 'response_type') String? responseType,@JsonKey(name: 'search_type') String? searchType
});




}
/// @nodoc
class __$ChatResponseCopyWithImpl<$Res>
    implements _$ChatResponseCopyWith<$Res> {
  __$ChatResponseCopyWithImpl(this._self, this._then);

  final _ChatResponse _self;
  final $Res Function(_ChatResponse) _then;

/// Create a copy of ChatResponse
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? sessionId = null,Object? message = null,Object? quickReplies = freezed,Object? products = freezed,Object? responseType = freezed,Object? searchType = freezed,}) {
  return _then(_ChatResponse(
sessionId: null == sessionId ? _self.sessionId : sessionId // ignore: cast_nullable_to_non_nullable
as String,message: null == message ? _self.message : message // ignore: cast_nullable_to_non_nullable
as String,quickReplies: freezed == quickReplies ? _self._quickReplies : quickReplies // ignore: cast_nullable_to_non_nullable
as List<String>?,products: freezed == products ? _self._products : products // ignore: cast_nullable_to_non_nullable
as List<Product>?,responseType: freezed == responseType ? _self.responseType : responseType // ignore: cast_nullable_to_non_nullable
as String?,searchType: freezed == searchType ? _self.searchType : searchType // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}


}

// dart format on
