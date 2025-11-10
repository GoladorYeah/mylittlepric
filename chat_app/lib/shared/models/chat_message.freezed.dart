// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'chat_message.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$ChatMessage {

 String get id; MessageRole get role; String get content; int get timestamp;@JsonKey(name: 'quick_replies') List<String>? get quickReplies; List<Product>? get products;@JsonKey(name: 'search_type') String? get searchType; bool get isLocal;
/// Create a copy of ChatMessage
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ChatMessageCopyWith<ChatMessage> get copyWith => _$ChatMessageCopyWithImpl<ChatMessage>(this as ChatMessage, _$identity);

  /// Serializes this ChatMessage to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ChatMessage&&(identical(other.id, id) || other.id == id)&&(identical(other.role, role) || other.role == role)&&(identical(other.content, content) || other.content == content)&&(identical(other.timestamp, timestamp) || other.timestamp == timestamp)&&const DeepCollectionEquality().equals(other.quickReplies, quickReplies)&&const DeepCollectionEquality().equals(other.products, products)&&(identical(other.searchType, searchType) || other.searchType == searchType)&&(identical(other.isLocal, isLocal) || other.isLocal == isLocal));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,role,content,timestamp,const DeepCollectionEquality().hash(quickReplies),const DeepCollectionEquality().hash(products),searchType,isLocal);

@override
String toString() {
  return 'ChatMessage(id: $id, role: $role, content: $content, timestamp: $timestamp, quickReplies: $quickReplies, products: $products, searchType: $searchType, isLocal: $isLocal)';
}


}

/// @nodoc
abstract mixin class $ChatMessageCopyWith<$Res>  {
  factory $ChatMessageCopyWith(ChatMessage value, $Res Function(ChatMessage) _then) = _$ChatMessageCopyWithImpl;
@useResult
$Res call({
 String id, MessageRole role, String content, int timestamp,@JsonKey(name: 'quick_replies') List<String>? quickReplies, List<Product>? products,@JsonKey(name: 'search_type') String? searchType, bool isLocal
});




}
/// @nodoc
class _$ChatMessageCopyWithImpl<$Res>
    implements $ChatMessageCopyWith<$Res> {
  _$ChatMessageCopyWithImpl(this._self, this._then);

  final ChatMessage _self;
  final $Res Function(ChatMessage) _then;

/// Create a copy of ChatMessage
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? id = null,Object? role = null,Object? content = null,Object? timestamp = null,Object? quickReplies = freezed,Object? products = freezed,Object? searchType = freezed,Object? isLocal = null,}) {
  return _then(_self.copyWith(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,role: null == role ? _self.role : role // ignore: cast_nullable_to_non_nullable
as MessageRole,content: null == content ? _self.content : content // ignore: cast_nullable_to_non_nullable
as String,timestamp: null == timestamp ? _self.timestamp : timestamp // ignore: cast_nullable_to_non_nullable
as int,quickReplies: freezed == quickReplies ? _self.quickReplies : quickReplies // ignore: cast_nullable_to_non_nullable
as List<String>?,products: freezed == products ? _self.products : products // ignore: cast_nullable_to_non_nullable
as List<Product>?,searchType: freezed == searchType ? _self.searchType : searchType // ignore: cast_nullable_to_non_nullable
as String?,isLocal: null == isLocal ? _self.isLocal : isLocal // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}

}


/// Adds pattern-matching-related methods to [ChatMessage].
extension ChatMessagePatterns on ChatMessage {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ChatMessage value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ChatMessage() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ChatMessage value)  $default,){
final _that = this;
switch (_that) {
case _ChatMessage():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ChatMessage value)?  $default,){
final _that = this;
switch (_that) {
case _ChatMessage() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String id,  MessageRole role,  String content,  int timestamp, @JsonKey(name: 'quick_replies')  List<String>? quickReplies,  List<Product>? products, @JsonKey(name: 'search_type')  String? searchType,  bool isLocal)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ChatMessage() when $default != null:
return $default(_that.id,_that.role,_that.content,_that.timestamp,_that.quickReplies,_that.products,_that.searchType,_that.isLocal);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String id,  MessageRole role,  String content,  int timestamp, @JsonKey(name: 'quick_replies')  List<String>? quickReplies,  List<Product>? products, @JsonKey(name: 'search_type')  String? searchType,  bool isLocal)  $default,) {final _that = this;
switch (_that) {
case _ChatMessage():
return $default(_that.id,_that.role,_that.content,_that.timestamp,_that.quickReplies,_that.products,_that.searchType,_that.isLocal);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String id,  MessageRole role,  String content,  int timestamp, @JsonKey(name: 'quick_replies')  List<String>? quickReplies,  List<Product>? products, @JsonKey(name: 'search_type')  String? searchType,  bool isLocal)?  $default,) {final _that = this;
switch (_that) {
case _ChatMessage() when $default != null:
return $default(_that.id,_that.role,_that.content,_that.timestamp,_that.quickReplies,_that.products,_that.searchType,_that.isLocal);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ChatMessage implements ChatMessage {
  const _ChatMessage({required this.id, required this.role, required this.content, required this.timestamp, @JsonKey(name: 'quick_replies') final  List<String>? quickReplies, final  List<Product>? products, @JsonKey(name: 'search_type') this.searchType, this.isLocal = false}): _quickReplies = quickReplies,_products = products;
  factory _ChatMessage.fromJson(Map<String, dynamic> json) => _$ChatMessageFromJson(json);

@override final  String id;
@override final  MessageRole role;
@override final  String content;
@override final  int timestamp;
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

@override@JsonKey(name: 'search_type') final  String? searchType;
@override@JsonKey() final  bool isLocal;

/// Create a copy of ChatMessage
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ChatMessageCopyWith<_ChatMessage> get copyWith => __$ChatMessageCopyWithImpl<_ChatMessage>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ChatMessageToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ChatMessage&&(identical(other.id, id) || other.id == id)&&(identical(other.role, role) || other.role == role)&&(identical(other.content, content) || other.content == content)&&(identical(other.timestamp, timestamp) || other.timestamp == timestamp)&&const DeepCollectionEquality().equals(other._quickReplies, _quickReplies)&&const DeepCollectionEquality().equals(other._products, _products)&&(identical(other.searchType, searchType) || other.searchType == searchType)&&(identical(other.isLocal, isLocal) || other.isLocal == isLocal));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,id,role,content,timestamp,const DeepCollectionEquality().hash(_quickReplies),const DeepCollectionEquality().hash(_products),searchType,isLocal);

@override
String toString() {
  return 'ChatMessage(id: $id, role: $role, content: $content, timestamp: $timestamp, quickReplies: $quickReplies, products: $products, searchType: $searchType, isLocal: $isLocal)';
}


}

/// @nodoc
abstract mixin class _$ChatMessageCopyWith<$Res> implements $ChatMessageCopyWith<$Res> {
  factory _$ChatMessageCopyWith(_ChatMessage value, $Res Function(_ChatMessage) _then) = __$ChatMessageCopyWithImpl;
@override @useResult
$Res call({
 String id, MessageRole role, String content, int timestamp,@JsonKey(name: 'quick_replies') List<String>? quickReplies, List<Product>? products,@JsonKey(name: 'search_type') String? searchType, bool isLocal
});




}
/// @nodoc
class __$ChatMessageCopyWithImpl<$Res>
    implements _$ChatMessageCopyWith<$Res> {
  __$ChatMessageCopyWithImpl(this._self, this._then);

  final _ChatMessage _self;
  final $Res Function(_ChatMessage) _then;

/// Create a copy of ChatMessage
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? id = null,Object? role = null,Object? content = null,Object? timestamp = null,Object? quickReplies = freezed,Object? products = freezed,Object? searchType = freezed,Object? isLocal = null,}) {
  return _then(_ChatMessage(
id: null == id ? _self.id : id // ignore: cast_nullable_to_non_nullable
as String,role: null == role ? _self.role : role // ignore: cast_nullable_to_non_nullable
as MessageRole,content: null == content ? _self.content : content // ignore: cast_nullable_to_non_nullable
as String,timestamp: null == timestamp ? _self.timestamp : timestamp // ignore: cast_nullable_to_non_nullable
as int,quickReplies: freezed == quickReplies ? _self._quickReplies : quickReplies // ignore: cast_nullable_to_non_nullable
as List<String>?,products: freezed == products ? _self._products : products // ignore: cast_nullable_to_non_nullable
as List<Product>?,searchType: freezed == searchType ? _self.searchType : searchType // ignore: cast_nullable_to_non_nullable
as String?,isLocal: null == isLocal ? _self.isLocal : isLocal // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}


}


/// @nodoc
mixin _$SessionMessage {

 String get role; String get content; String? get timestamp;@JsonKey(name: 'quick_replies') List<String>? get quickReplies; List<Product>? get products;@JsonKey(name: 'search_type') String? get searchType;
/// Create a copy of SessionMessage
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$SessionMessageCopyWith<SessionMessage> get copyWith => _$SessionMessageCopyWithImpl<SessionMessage>(this as SessionMessage, _$identity);

  /// Serializes this SessionMessage to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is SessionMessage&&(identical(other.role, role) || other.role == role)&&(identical(other.content, content) || other.content == content)&&(identical(other.timestamp, timestamp) || other.timestamp == timestamp)&&const DeepCollectionEquality().equals(other.quickReplies, quickReplies)&&const DeepCollectionEquality().equals(other.products, products)&&(identical(other.searchType, searchType) || other.searchType == searchType));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,role,content,timestamp,const DeepCollectionEquality().hash(quickReplies),const DeepCollectionEquality().hash(products),searchType);

@override
String toString() {
  return 'SessionMessage(role: $role, content: $content, timestamp: $timestamp, quickReplies: $quickReplies, products: $products, searchType: $searchType)';
}


}

/// @nodoc
abstract mixin class $SessionMessageCopyWith<$Res>  {
  factory $SessionMessageCopyWith(SessionMessage value, $Res Function(SessionMessage) _then) = _$SessionMessageCopyWithImpl;
@useResult
$Res call({
 String role, String content, String? timestamp,@JsonKey(name: 'quick_replies') List<String>? quickReplies, List<Product>? products,@JsonKey(name: 'search_type') String? searchType
});




}
/// @nodoc
class _$SessionMessageCopyWithImpl<$Res>
    implements $SessionMessageCopyWith<$Res> {
  _$SessionMessageCopyWithImpl(this._self, this._then);

  final SessionMessage _self;
  final $Res Function(SessionMessage) _then;

/// Create a copy of SessionMessage
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? role = null,Object? content = null,Object? timestamp = freezed,Object? quickReplies = freezed,Object? products = freezed,Object? searchType = freezed,}) {
  return _then(_self.copyWith(
role: null == role ? _self.role : role // ignore: cast_nullable_to_non_nullable
as String,content: null == content ? _self.content : content // ignore: cast_nullable_to_non_nullable
as String,timestamp: freezed == timestamp ? _self.timestamp : timestamp // ignore: cast_nullable_to_non_nullable
as String?,quickReplies: freezed == quickReplies ? _self.quickReplies : quickReplies // ignore: cast_nullable_to_non_nullable
as List<String>?,products: freezed == products ? _self.products : products // ignore: cast_nullable_to_non_nullable
as List<Product>?,searchType: freezed == searchType ? _self.searchType : searchType // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

}


/// Adds pattern-matching-related methods to [SessionMessage].
extension SessionMessagePatterns on SessionMessage {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _SessionMessage value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _SessionMessage() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _SessionMessage value)  $default,){
final _that = this;
switch (_that) {
case _SessionMessage():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _SessionMessage value)?  $default,){
final _that = this;
switch (_that) {
case _SessionMessage() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String role,  String content,  String? timestamp, @JsonKey(name: 'quick_replies')  List<String>? quickReplies,  List<Product>? products, @JsonKey(name: 'search_type')  String? searchType)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _SessionMessage() when $default != null:
return $default(_that.role,_that.content,_that.timestamp,_that.quickReplies,_that.products,_that.searchType);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String role,  String content,  String? timestamp, @JsonKey(name: 'quick_replies')  List<String>? quickReplies,  List<Product>? products, @JsonKey(name: 'search_type')  String? searchType)  $default,) {final _that = this;
switch (_that) {
case _SessionMessage():
return $default(_that.role,_that.content,_that.timestamp,_that.quickReplies,_that.products,_that.searchType);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String role,  String content,  String? timestamp, @JsonKey(name: 'quick_replies')  List<String>? quickReplies,  List<Product>? products, @JsonKey(name: 'search_type')  String? searchType)?  $default,) {final _that = this;
switch (_that) {
case _SessionMessage() when $default != null:
return $default(_that.role,_that.content,_that.timestamp,_that.quickReplies,_that.products,_that.searchType);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _SessionMessage implements SessionMessage {
  const _SessionMessage({required this.role, required this.content, this.timestamp, @JsonKey(name: 'quick_replies') final  List<String>? quickReplies, final  List<Product>? products, @JsonKey(name: 'search_type') this.searchType}): _quickReplies = quickReplies,_products = products;
  factory _SessionMessage.fromJson(Map<String, dynamic> json) => _$SessionMessageFromJson(json);

@override final  String role;
@override final  String content;
@override final  String? timestamp;
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

@override@JsonKey(name: 'search_type') final  String? searchType;

/// Create a copy of SessionMessage
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$SessionMessageCopyWith<_SessionMessage> get copyWith => __$SessionMessageCopyWithImpl<_SessionMessage>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$SessionMessageToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _SessionMessage&&(identical(other.role, role) || other.role == role)&&(identical(other.content, content) || other.content == content)&&(identical(other.timestamp, timestamp) || other.timestamp == timestamp)&&const DeepCollectionEquality().equals(other._quickReplies, _quickReplies)&&const DeepCollectionEquality().equals(other._products, _products)&&(identical(other.searchType, searchType) || other.searchType == searchType));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,role,content,timestamp,const DeepCollectionEquality().hash(_quickReplies),const DeepCollectionEquality().hash(_products),searchType);

@override
String toString() {
  return 'SessionMessage(role: $role, content: $content, timestamp: $timestamp, quickReplies: $quickReplies, products: $products, searchType: $searchType)';
}


}

/// @nodoc
abstract mixin class _$SessionMessageCopyWith<$Res> implements $SessionMessageCopyWith<$Res> {
  factory _$SessionMessageCopyWith(_SessionMessage value, $Res Function(_SessionMessage) _then) = __$SessionMessageCopyWithImpl;
@override @useResult
$Res call({
 String role, String content, String? timestamp,@JsonKey(name: 'quick_replies') List<String>? quickReplies, List<Product>? products,@JsonKey(name: 'search_type') String? searchType
});




}
/// @nodoc
class __$SessionMessageCopyWithImpl<$Res>
    implements _$SessionMessageCopyWith<$Res> {
  __$SessionMessageCopyWithImpl(this._self, this._then);

  final _SessionMessage _self;
  final $Res Function(_SessionMessage) _then;

/// Create a copy of SessionMessage
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? role = null,Object? content = null,Object? timestamp = freezed,Object? quickReplies = freezed,Object? products = freezed,Object? searchType = freezed,}) {
  return _then(_SessionMessage(
role: null == role ? _self.role : role // ignore: cast_nullable_to_non_nullable
as String,content: null == content ? _self.content : content // ignore: cast_nullable_to_non_nullable
as String,timestamp: freezed == timestamp ? _self.timestamp : timestamp // ignore: cast_nullable_to_non_nullable
as String?,quickReplies: freezed == quickReplies ? _self._quickReplies : quickReplies // ignore: cast_nullable_to_non_nullable
as List<String>?,products: freezed == products ? _self._products : products // ignore: cast_nullable_to_non_nullable
as List<Product>?,searchType: freezed == searchType ? _self.searchType : searchType // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}


}

// dart format on
