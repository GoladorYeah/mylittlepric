// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'chat_state.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$ChatState {

/// Current session ID
 String? get sessionId;/// List of chat messages
 List<ChatMessage> get messages;/// WebSocket connection state
 WebSocketState get wsState;/// Whether a message is being sent
 bool get isSending;/// Whether messages are being loaded
 bool get isLoading;/// Error message
 String? get error;/// Quick replies from the assistant
 List<String> get quickReplies;/// Whether the assistant is typing
 bool get isTyping;
/// Create a copy of ChatState
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ChatStateCopyWith<ChatState> get copyWith => _$ChatStateCopyWithImpl<ChatState>(this as ChatState, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ChatState&&(identical(other.sessionId, sessionId) || other.sessionId == sessionId)&&const DeepCollectionEquality().equals(other.messages, messages)&&(identical(other.wsState, wsState) || other.wsState == wsState)&&(identical(other.isSending, isSending) || other.isSending == isSending)&&(identical(other.isLoading, isLoading) || other.isLoading == isLoading)&&(identical(other.error, error) || other.error == error)&&const DeepCollectionEquality().equals(other.quickReplies, quickReplies)&&(identical(other.isTyping, isTyping) || other.isTyping == isTyping));
}


@override
int get hashCode => Object.hash(runtimeType,sessionId,const DeepCollectionEquality().hash(messages),wsState,isSending,isLoading,error,const DeepCollectionEquality().hash(quickReplies),isTyping);

@override
String toString() {
  return 'ChatState(sessionId: $sessionId, messages: $messages, wsState: $wsState, isSending: $isSending, isLoading: $isLoading, error: $error, quickReplies: $quickReplies, isTyping: $isTyping)';
}


}

/// @nodoc
abstract mixin class $ChatStateCopyWith<$Res>  {
  factory $ChatStateCopyWith(ChatState value, $Res Function(ChatState) _then) = _$ChatStateCopyWithImpl;
@useResult
$Res call({
 String? sessionId, List<ChatMessage> messages, WebSocketState wsState, bool isSending, bool isLoading, String? error, List<String> quickReplies, bool isTyping
});




}
/// @nodoc
class _$ChatStateCopyWithImpl<$Res>
    implements $ChatStateCopyWith<$Res> {
  _$ChatStateCopyWithImpl(this._self, this._then);

  final ChatState _self;
  final $Res Function(ChatState) _then;

/// Create a copy of ChatState
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? sessionId = freezed,Object? messages = null,Object? wsState = null,Object? isSending = null,Object? isLoading = null,Object? error = freezed,Object? quickReplies = null,Object? isTyping = null,}) {
  return _then(_self.copyWith(
sessionId: freezed == sessionId ? _self.sessionId : sessionId // ignore: cast_nullable_to_non_nullable
as String?,messages: null == messages ? _self.messages : messages // ignore: cast_nullable_to_non_nullable
as List<ChatMessage>,wsState: null == wsState ? _self.wsState : wsState // ignore: cast_nullable_to_non_nullable
as WebSocketState,isSending: null == isSending ? _self.isSending : isSending // ignore: cast_nullable_to_non_nullable
as bool,isLoading: null == isLoading ? _self.isLoading : isLoading // ignore: cast_nullable_to_non_nullable
as bool,error: freezed == error ? _self.error : error // ignore: cast_nullable_to_non_nullable
as String?,quickReplies: null == quickReplies ? _self.quickReplies : quickReplies // ignore: cast_nullable_to_non_nullable
as List<String>,isTyping: null == isTyping ? _self.isTyping : isTyping // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}

}


/// Adds pattern-matching-related methods to [ChatState].
extension ChatStatePatterns on ChatState {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ChatState value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ChatState() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ChatState value)  $default,){
final _that = this;
switch (_that) {
case _ChatState():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ChatState value)?  $default,){
final _that = this;
switch (_that) {
case _ChatState() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String? sessionId,  List<ChatMessage> messages,  WebSocketState wsState,  bool isSending,  bool isLoading,  String? error,  List<String> quickReplies,  bool isTyping)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ChatState() when $default != null:
return $default(_that.sessionId,_that.messages,_that.wsState,_that.isSending,_that.isLoading,_that.error,_that.quickReplies,_that.isTyping);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String? sessionId,  List<ChatMessage> messages,  WebSocketState wsState,  bool isSending,  bool isLoading,  String? error,  List<String> quickReplies,  bool isTyping)  $default,) {final _that = this;
switch (_that) {
case _ChatState():
return $default(_that.sessionId,_that.messages,_that.wsState,_that.isSending,_that.isLoading,_that.error,_that.quickReplies,_that.isTyping);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String? sessionId,  List<ChatMessage> messages,  WebSocketState wsState,  bool isSending,  bool isLoading,  String? error,  List<String> quickReplies,  bool isTyping)?  $default,) {final _that = this;
switch (_that) {
case _ChatState() when $default != null:
return $default(_that.sessionId,_that.messages,_that.wsState,_that.isSending,_that.isLoading,_that.error,_that.quickReplies,_that.isTyping);case _:
  return null;

}
}

}

/// @nodoc


class _ChatState extends ChatState {
  const _ChatState({this.sessionId, final  List<ChatMessage> messages = const [], this.wsState = WebSocketState.disconnected, this.isSending = false, this.isLoading = false, this.error, final  List<String> quickReplies = const [], this.isTyping = false}): _messages = messages,_quickReplies = quickReplies,super._();
  

/// Current session ID
@override final  String? sessionId;
/// List of chat messages
 final  List<ChatMessage> _messages;
/// List of chat messages
@override@JsonKey() List<ChatMessage> get messages {
  if (_messages is EqualUnmodifiableListView) return _messages;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_messages);
}

/// WebSocket connection state
@override@JsonKey() final  WebSocketState wsState;
/// Whether a message is being sent
@override@JsonKey() final  bool isSending;
/// Whether messages are being loaded
@override@JsonKey() final  bool isLoading;
/// Error message
@override final  String? error;
/// Quick replies from the assistant
 final  List<String> _quickReplies;
/// Quick replies from the assistant
@override@JsonKey() List<String> get quickReplies {
  if (_quickReplies is EqualUnmodifiableListView) return _quickReplies;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_quickReplies);
}

/// Whether the assistant is typing
@override@JsonKey() final  bool isTyping;

/// Create a copy of ChatState
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ChatStateCopyWith<_ChatState> get copyWith => __$ChatStateCopyWithImpl<_ChatState>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ChatState&&(identical(other.sessionId, sessionId) || other.sessionId == sessionId)&&const DeepCollectionEquality().equals(other._messages, _messages)&&(identical(other.wsState, wsState) || other.wsState == wsState)&&(identical(other.isSending, isSending) || other.isSending == isSending)&&(identical(other.isLoading, isLoading) || other.isLoading == isLoading)&&(identical(other.error, error) || other.error == error)&&const DeepCollectionEquality().equals(other._quickReplies, _quickReplies)&&(identical(other.isTyping, isTyping) || other.isTyping == isTyping));
}


@override
int get hashCode => Object.hash(runtimeType,sessionId,const DeepCollectionEquality().hash(_messages),wsState,isSending,isLoading,error,const DeepCollectionEquality().hash(_quickReplies),isTyping);

@override
String toString() {
  return 'ChatState(sessionId: $sessionId, messages: $messages, wsState: $wsState, isSending: $isSending, isLoading: $isLoading, error: $error, quickReplies: $quickReplies, isTyping: $isTyping)';
}


}

/// @nodoc
abstract mixin class _$ChatStateCopyWith<$Res> implements $ChatStateCopyWith<$Res> {
  factory _$ChatStateCopyWith(_ChatState value, $Res Function(_ChatState) _then) = __$ChatStateCopyWithImpl;
@override @useResult
$Res call({
 String? sessionId, List<ChatMessage> messages, WebSocketState wsState, bool isSending, bool isLoading, String? error, List<String> quickReplies, bool isTyping
});




}
/// @nodoc
class __$ChatStateCopyWithImpl<$Res>
    implements _$ChatStateCopyWith<$Res> {
  __$ChatStateCopyWithImpl(this._self, this._then);

  final _ChatState _self;
  final $Res Function(_ChatState) _then;

/// Create a copy of ChatState
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? sessionId = freezed,Object? messages = null,Object? wsState = null,Object? isSending = null,Object? isLoading = null,Object? error = freezed,Object? quickReplies = null,Object? isTyping = null,}) {
  return _then(_ChatState(
sessionId: freezed == sessionId ? _self.sessionId : sessionId // ignore: cast_nullable_to_non_nullable
as String?,messages: null == messages ? _self._messages : messages // ignore: cast_nullable_to_non_nullable
as List<ChatMessage>,wsState: null == wsState ? _self.wsState : wsState // ignore: cast_nullable_to_non_nullable
as WebSocketState,isSending: null == isSending ? _self.isSending : isSending // ignore: cast_nullable_to_non_nullable
as bool,isLoading: null == isLoading ? _self.isLoading : isLoading // ignore: cast_nullable_to_non_nullable
as bool,error: freezed == error ? _self.error : error // ignore: cast_nullable_to_non_nullable
as String?,quickReplies: null == quickReplies ? _self._quickReplies : quickReplies // ignore: cast_nullable_to_non_nullable
as List<String>,isTyping: null == isTyping ? _self.isTyping : isTyping // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}


}

// dart format on
