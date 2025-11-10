// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'session_state.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;
/// @nodoc
mixin _$SessionState {

/// List of saved searches
 List<SavedSearch> get searches;/// Whether searches are being loaded
 bool get isLoading;/// Whether more searches can be loaded
 bool get hasMore;/// Current page offset for pagination
 int get offset;/// Number of items per page
 int get limit;/// Error message
 String? get error;/// Whether a search is being deleted
 bool get isDeleting;
/// Create a copy of SessionState
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$SessionStateCopyWith<SessionState> get copyWith => _$SessionStateCopyWithImpl<SessionState>(this as SessionState, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is SessionState&&const DeepCollectionEquality().equals(other.searches, searches)&&(identical(other.isLoading, isLoading) || other.isLoading == isLoading)&&(identical(other.hasMore, hasMore) || other.hasMore == hasMore)&&(identical(other.offset, offset) || other.offset == offset)&&(identical(other.limit, limit) || other.limit == limit)&&(identical(other.error, error) || other.error == error)&&(identical(other.isDeleting, isDeleting) || other.isDeleting == isDeleting));
}


@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(searches),isLoading,hasMore,offset,limit,error,isDeleting);

@override
String toString() {
  return 'SessionState(searches: $searches, isLoading: $isLoading, hasMore: $hasMore, offset: $offset, limit: $limit, error: $error, isDeleting: $isDeleting)';
}


}

/// @nodoc
abstract mixin class $SessionStateCopyWith<$Res>  {
  factory $SessionStateCopyWith(SessionState value, $Res Function(SessionState) _then) = _$SessionStateCopyWithImpl;
@useResult
$Res call({
 List<SavedSearch> searches, bool isLoading, bool hasMore, int offset, int limit, String? error, bool isDeleting
});




}
/// @nodoc
class _$SessionStateCopyWithImpl<$Res>
    implements $SessionStateCopyWith<$Res> {
  _$SessionStateCopyWithImpl(this._self, this._then);

  final SessionState _self;
  final $Res Function(SessionState) _then;

/// Create a copy of SessionState
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? searches = null,Object? isLoading = null,Object? hasMore = null,Object? offset = null,Object? limit = null,Object? error = freezed,Object? isDeleting = null,}) {
  return _then(_self.copyWith(
searches: null == searches ? _self.searches : searches // ignore: cast_nullable_to_non_nullable
as List<SavedSearch>,isLoading: null == isLoading ? _self.isLoading : isLoading // ignore: cast_nullable_to_non_nullable
as bool,hasMore: null == hasMore ? _self.hasMore : hasMore // ignore: cast_nullable_to_non_nullable
as bool,offset: null == offset ? _self.offset : offset // ignore: cast_nullable_to_non_nullable
as int,limit: null == limit ? _self.limit : limit // ignore: cast_nullable_to_non_nullable
as int,error: freezed == error ? _self.error : error // ignore: cast_nullable_to_non_nullable
as String?,isDeleting: null == isDeleting ? _self.isDeleting : isDeleting // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}

}


/// Adds pattern-matching-related methods to [SessionState].
extension SessionStatePatterns on SessionState {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _SessionState value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _SessionState() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _SessionState value)  $default,){
final _that = this;
switch (_that) {
case _SessionState():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _SessionState value)?  $default,){
final _that = this;
switch (_that) {
case _SessionState() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( List<SavedSearch> searches,  bool isLoading,  bool hasMore,  int offset,  int limit,  String? error,  bool isDeleting)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _SessionState() when $default != null:
return $default(_that.searches,_that.isLoading,_that.hasMore,_that.offset,_that.limit,_that.error,_that.isDeleting);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( List<SavedSearch> searches,  bool isLoading,  bool hasMore,  int offset,  int limit,  String? error,  bool isDeleting)  $default,) {final _that = this;
switch (_that) {
case _SessionState():
return $default(_that.searches,_that.isLoading,_that.hasMore,_that.offset,_that.limit,_that.error,_that.isDeleting);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( List<SavedSearch> searches,  bool isLoading,  bool hasMore,  int offset,  int limit,  String? error,  bool isDeleting)?  $default,) {final _that = this;
switch (_that) {
case _SessionState() when $default != null:
return $default(_that.searches,_that.isLoading,_that.hasMore,_that.offset,_that.limit,_that.error,_that.isDeleting);case _:
  return null;

}
}

}

/// @nodoc


class _SessionState extends SessionState {
  const _SessionState({final  List<SavedSearch> searches = const [], this.isLoading = false, this.hasMore = true, this.offset = 0, this.limit = 20, this.error, this.isDeleting = false}): _searches = searches,super._();
  

/// List of saved searches
 final  List<SavedSearch> _searches;
/// List of saved searches
@override@JsonKey() List<SavedSearch> get searches {
  if (_searches is EqualUnmodifiableListView) return _searches;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_searches);
}

/// Whether searches are being loaded
@override@JsonKey() final  bool isLoading;
/// Whether more searches can be loaded
@override@JsonKey() final  bool hasMore;
/// Current page offset for pagination
@override@JsonKey() final  int offset;
/// Number of items per page
@override@JsonKey() final  int limit;
/// Error message
@override final  String? error;
/// Whether a search is being deleted
@override@JsonKey() final  bool isDeleting;

/// Create a copy of SessionState
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$SessionStateCopyWith<_SessionState> get copyWith => __$SessionStateCopyWithImpl<_SessionState>(this, _$identity);



@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _SessionState&&const DeepCollectionEquality().equals(other._searches, _searches)&&(identical(other.isLoading, isLoading) || other.isLoading == isLoading)&&(identical(other.hasMore, hasMore) || other.hasMore == hasMore)&&(identical(other.offset, offset) || other.offset == offset)&&(identical(other.limit, limit) || other.limit == limit)&&(identical(other.error, error) || other.error == error)&&(identical(other.isDeleting, isDeleting) || other.isDeleting == isDeleting));
}


@override
int get hashCode => Object.hash(runtimeType,const DeepCollectionEquality().hash(_searches),isLoading,hasMore,offset,limit,error,isDeleting);

@override
String toString() {
  return 'SessionState(searches: $searches, isLoading: $isLoading, hasMore: $hasMore, offset: $offset, limit: $limit, error: $error, isDeleting: $isDeleting)';
}


}

/// @nodoc
abstract mixin class _$SessionStateCopyWith<$Res> implements $SessionStateCopyWith<$Res> {
  factory _$SessionStateCopyWith(_SessionState value, $Res Function(_SessionState) _then) = __$SessionStateCopyWithImpl;
@override @useResult
$Res call({
 List<SavedSearch> searches, bool isLoading, bool hasMore, int offset, int limit, String? error, bool isDeleting
});




}
/// @nodoc
class __$SessionStateCopyWithImpl<$Res>
    implements _$SessionStateCopyWith<$Res> {
  __$SessionStateCopyWithImpl(this._self, this._then);

  final _SessionState _self;
  final $Res Function(_SessionState) _then;

/// Create a copy of SessionState
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? searches = null,Object? isLoading = null,Object? hasMore = null,Object? offset = null,Object? limit = null,Object? error = freezed,Object? isDeleting = null,}) {
  return _then(_SessionState(
searches: null == searches ? _self._searches : searches // ignore: cast_nullable_to_non_nullable
as List<SavedSearch>,isLoading: null == isLoading ? _self.isLoading : isLoading // ignore: cast_nullable_to_non_nullable
as bool,hasMore: null == hasMore ? _self.hasMore : hasMore // ignore: cast_nullable_to_non_nullable
as bool,offset: null == offset ? _self.offset : offset // ignore: cast_nullable_to_non_nullable
as int,limit: null == limit ? _self.limit : limit // ignore: cast_nullable_to_non_nullable
as int,error: freezed == error ? _self.error : error // ignore: cast_nullable_to_non_nullable
as String?,isDeleting: null == isDeleting ? _self.isDeleting : isDeleting // ignore: cast_nullable_to_non_nullable
as bool,
  ));
}


}

// dart format on
