// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'product.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$Product {

 int get position; String get title; String get link;@JsonKey(name: 'product_link') String get productLink;@JsonKey(name: 'product_id') String get productId;@JsonKey(name: 'serpapi_product_api') String get serpapiProductApi; String get source; String get price;@JsonKey(name: 'extracted_price') double get extractedPrice; String get thumbnail;@JsonKey(name: 'serpapi_product_api_comparative') String? get serpapiProductApiComparative; double? get rating; int? get reviews; String? get delivery; String? get tag; List<String>? get extensions; String? get currency;@JsonKey(name: 'page_token') String? get pageToken;@JsonKey(name: 'relevance_score') double? get relevanceScore;
/// Create a copy of Product
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ProductCopyWith<Product> get copyWith => _$ProductCopyWithImpl<Product>(this as Product, _$identity);

  /// Serializes this Product to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Product&&(identical(other.position, position) || other.position == position)&&(identical(other.title, title) || other.title == title)&&(identical(other.link, link) || other.link == link)&&(identical(other.productLink, productLink) || other.productLink == productLink)&&(identical(other.productId, productId) || other.productId == productId)&&(identical(other.serpapiProductApi, serpapiProductApi) || other.serpapiProductApi == serpapiProductApi)&&(identical(other.source, source) || other.source == source)&&(identical(other.price, price) || other.price == price)&&(identical(other.extractedPrice, extractedPrice) || other.extractedPrice == extractedPrice)&&(identical(other.thumbnail, thumbnail) || other.thumbnail == thumbnail)&&(identical(other.serpapiProductApiComparative, serpapiProductApiComparative) || other.serpapiProductApiComparative == serpapiProductApiComparative)&&(identical(other.rating, rating) || other.rating == rating)&&(identical(other.reviews, reviews) || other.reviews == reviews)&&(identical(other.delivery, delivery) || other.delivery == delivery)&&(identical(other.tag, tag) || other.tag == tag)&&const DeepCollectionEquality().equals(other.extensions, extensions)&&(identical(other.currency, currency) || other.currency == currency)&&(identical(other.pageToken, pageToken) || other.pageToken == pageToken)&&(identical(other.relevanceScore, relevanceScore) || other.relevanceScore == relevanceScore));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hashAll([runtimeType,position,title,link,productLink,productId,serpapiProductApi,source,price,extractedPrice,thumbnail,serpapiProductApiComparative,rating,reviews,delivery,tag,const DeepCollectionEquality().hash(extensions),currency,pageToken,relevanceScore]);

@override
String toString() {
  return 'Product(position: $position, title: $title, link: $link, productLink: $productLink, productId: $productId, serpapiProductApi: $serpapiProductApi, source: $source, price: $price, extractedPrice: $extractedPrice, thumbnail: $thumbnail, serpapiProductApiComparative: $serpapiProductApiComparative, rating: $rating, reviews: $reviews, delivery: $delivery, tag: $tag, extensions: $extensions, currency: $currency, pageToken: $pageToken, relevanceScore: $relevanceScore)';
}


}

/// @nodoc
abstract mixin class $ProductCopyWith<$Res>  {
  factory $ProductCopyWith(Product value, $Res Function(Product) _then) = _$ProductCopyWithImpl;
@useResult
$Res call({
 int position, String title, String link,@JsonKey(name: 'product_link') String productLink,@JsonKey(name: 'product_id') String productId,@JsonKey(name: 'serpapi_product_api') String serpapiProductApi, String source, String price,@JsonKey(name: 'extracted_price') double extractedPrice, String thumbnail,@JsonKey(name: 'serpapi_product_api_comparative') String? serpapiProductApiComparative, double? rating, int? reviews, String? delivery, String? tag, List<String>? extensions, String? currency,@JsonKey(name: 'page_token') String? pageToken,@JsonKey(name: 'relevance_score') double? relevanceScore
});




}
/// @nodoc
class _$ProductCopyWithImpl<$Res>
    implements $ProductCopyWith<$Res> {
  _$ProductCopyWithImpl(this._self, this._then);

  final Product _self;
  final $Res Function(Product) _then;

/// Create a copy of Product
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? position = null,Object? title = null,Object? link = null,Object? productLink = null,Object? productId = null,Object? serpapiProductApi = null,Object? source = null,Object? price = null,Object? extractedPrice = null,Object? thumbnail = null,Object? serpapiProductApiComparative = freezed,Object? rating = freezed,Object? reviews = freezed,Object? delivery = freezed,Object? tag = freezed,Object? extensions = freezed,Object? currency = freezed,Object? pageToken = freezed,Object? relevanceScore = freezed,}) {
  return _then(_self.copyWith(
position: null == position ? _self.position : position // ignore: cast_nullable_to_non_nullable
as int,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,link: null == link ? _self.link : link // ignore: cast_nullable_to_non_nullable
as String,productLink: null == productLink ? _self.productLink : productLink // ignore: cast_nullable_to_non_nullable
as String,productId: null == productId ? _self.productId : productId // ignore: cast_nullable_to_non_nullable
as String,serpapiProductApi: null == serpapiProductApi ? _self.serpapiProductApi : serpapiProductApi // ignore: cast_nullable_to_non_nullable
as String,source: null == source ? _self.source : source // ignore: cast_nullable_to_non_nullable
as String,price: null == price ? _self.price : price // ignore: cast_nullable_to_non_nullable
as String,extractedPrice: null == extractedPrice ? _self.extractedPrice : extractedPrice // ignore: cast_nullable_to_non_nullable
as double,thumbnail: null == thumbnail ? _self.thumbnail : thumbnail // ignore: cast_nullable_to_non_nullable
as String,serpapiProductApiComparative: freezed == serpapiProductApiComparative ? _self.serpapiProductApiComparative : serpapiProductApiComparative // ignore: cast_nullable_to_non_nullable
as String?,rating: freezed == rating ? _self.rating : rating // ignore: cast_nullable_to_non_nullable
as double?,reviews: freezed == reviews ? _self.reviews : reviews // ignore: cast_nullable_to_non_nullable
as int?,delivery: freezed == delivery ? _self.delivery : delivery // ignore: cast_nullable_to_non_nullable
as String?,tag: freezed == tag ? _self.tag : tag // ignore: cast_nullable_to_non_nullable
as String?,extensions: freezed == extensions ? _self.extensions : extensions // ignore: cast_nullable_to_non_nullable
as List<String>?,currency: freezed == currency ? _self.currency : currency // ignore: cast_nullable_to_non_nullable
as String?,pageToken: freezed == pageToken ? _self.pageToken : pageToken // ignore: cast_nullable_to_non_nullable
as String?,relevanceScore: freezed == relevanceScore ? _self.relevanceScore : relevanceScore // ignore: cast_nullable_to_non_nullable
as double?,
  ));
}

}


/// Adds pattern-matching-related methods to [Product].
extension ProductPatterns on Product {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Product value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Product() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Product value)  $default,){
final _that = this;
switch (_that) {
case _Product():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Product value)?  $default,){
final _that = this;
switch (_that) {
case _Product() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( int position,  String title,  String link, @JsonKey(name: 'product_link')  String productLink, @JsonKey(name: 'product_id')  String productId, @JsonKey(name: 'serpapi_product_api')  String serpapiProductApi,  String source,  String price, @JsonKey(name: 'extracted_price')  double extractedPrice,  String thumbnail, @JsonKey(name: 'serpapi_product_api_comparative')  String? serpapiProductApiComparative,  double? rating,  int? reviews,  String? delivery,  String? tag,  List<String>? extensions,  String? currency, @JsonKey(name: 'page_token')  String? pageToken, @JsonKey(name: 'relevance_score')  double? relevanceScore)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Product() when $default != null:
return $default(_that.position,_that.title,_that.link,_that.productLink,_that.productId,_that.serpapiProductApi,_that.source,_that.price,_that.extractedPrice,_that.thumbnail,_that.serpapiProductApiComparative,_that.rating,_that.reviews,_that.delivery,_that.tag,_that.extensions,_that.currency,_that.pageToken,_that.relevanceScore);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( int position,  String title,  String link, @JsonKey(name: 'product_link')  String productLink, @JsonKey(name: 'product_id')  String productId, @JsonKey(name: 'serpapi_product_api')  String serpapiProductApi,  String source,  String price, @JsonKey(name: 'extracted_price')  double extractedPrice,  String thumbnail, @JsonKey(name: 'serpapi_product_api_comparative')  String? serpapiProductApiComparative,  double? rating,  int? reviews,  String? delivery,  String? tag,  List<String>? extensions,  String? currency, @JsonKey(name: 'page_token')  String? pageToken, @JsonKey(name: 'relevance_score')  double? relevanceScore)  $default,) {final _that = this;
switch (_that) {
case _Product():
return $default(_that.position,_that.title,_that.link,_that.productLink,_that.productId,_that.serpapiProductApi,_that.source,_that.price,_that.extractedPrice,_that.thumbnail,_that.serpapiProductApiComparative,_that.rating,_that.reviews,_that.delivery,_that.tag,_that.extensions,_that.currency,_that.pageToken,_that.relevanceScore);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( int position,  String title,  String link, @JsonKey(name: 'product_link')  String productLink, @JsonKey(name: 'product_id')  String productId, @JsonKey(name: 'serpapi_product_api')  String serpapiProductApi,  String source,  String price, @JsonKey(name: 'extracted_price')  double extractedPrice,  String thumbnail, @JsonKey(name: 'serpapi_product_api_comparative')  String? serpapiProductApiComparative,  double? rating,  int? reviews,  String? delivery,  String? tag,  List<String>? extensions,  String? currency, @JsonKey(name: 'page_token')  String? pageToken, @JsonKey(name: 'relevance_score')  double? relevanceScore)?  $default,) {final _that = this;
switch (_that) {
case _Product() when $default != null:
return $default(_that.position,_that.title,_that.link,_that.productLink,_that.productId,_that.serpapiProductApi,_that.source,_that.price,_that.extractedPrice,_that.thumbnail,_that.serpapiProductApiComparative,_that.rating,_that.reviews,_that.delivery,_that.tag,_that.extensions,_that.currency,_that.pageToken,_that.relevanceScore);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Product implements Product {
  const _Product({required this.position, required this.title, required this.link, @JsonKey(name: 'product_link') required this.productLink, @JsonKey(name: 'product_id') required this.productId, @JsonKey(name: 'serpapi_product_api') required this.serpapiProductApi, required this.source, required this.price, @JsonKey(name: 'extracted_price') required this.extractedPrice, required this.thumbnail, @JsonKey(name: 'serpapi_product_api_comparative') this.serpapiProductApiComparative, this.rating, this.reviews, this.delivery, this.tag, final  List<String>? extensions, this.currency, @JsonKey(name: 'page_token') this.pageToken, @JsonKey(name: 'relevance_score') this.relevanceScore}): _extensions = extensions;
  factory _Product.fromJson(Map<String, dynamic> json) => _$ProductFromJson(json);

@override final  int position;
@override final  String title;
@override final  String link;
@override@JsonKey(name: 'product_link') final  String productLink;
@override@JsonKey(name: 'product_id') final  String productId;
@override@JsonKey(name: 'serpapi_product_api') final  String serpapiProductApi;
@override final  String source;
@override final  String price;
@override@JsonKey(name: 'extracted_price') final  double extractedPrice;
@override final  String thumbnail;
@override@JsonKey(name: 'serpapi_product_api_comparative') final  String? serpapiProductApiComparative;
@override final  double? rating;
@override final  int? reviews;
@override final  String? delivery;
@override final  String? tag;
 final  List<String>? _extensions;
@override List<String>? get extensions {
  final value = _extensions;
  if (value == null) return null;
  if (_extensions is EqualUnmodifiableListView) return _extensions;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(value);
}

@override final  String? currency;
@override@JsonKey(name: 'page_token') final  String? pageToken;
@override@JsonKey(name: 'relevance_score') final  double? relevanceScore;

/// Create a copy of Product
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ProductCopyWith<_Product> get copyWith => __$ProductCopyWithImpl<_Product>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ProductToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Product&&(identical(other.position, position) || other.position == position)&&(identical(other.title, title) || other.title == title)&&(identical(other.link, link) || other.link == link)&&(identical(other.productLink, productLink) || other.productLink == productLink)&&(identical(other.productId, productId) || other.productId == productId)&&(identical(other.serpapiProductApi, serpapiProductApi) || other.serpapiProductApi == serpapiProductApi)&&(identical(other.source, source) || other.source == source)&&(identical(other.price, price) || other.price == price)&&(identical(other.extractedPrice, extractedPrice) || other.extractedPrice == extractedPrice)&&(identical(other.thumbnail, thumbnail) || other.thumbnail == thumbnail)&&(identical(other.serpapiProductApiComparative, serpapiProductApiComparative) || other.serpapiProductApiComparative == serpapiProductApiComparative)&&(identical(other.rating, rating) || other.rating == rating)&&(identical(other.reviews, reviews) || other.reviews == reviews)&&(identical(other.delivery, delivery) || other.delivery == delivery)&&(identical(other.tag, tag) || other.tag == tag)&&const DeepCollectionEquality().equals(other._extensions, _extensions)&&(identical(other.currency, currency) || other.currency == currency)&&(identical(other.pageToken, pageToken) || other.pageToken == pageToken)&&(identical(other.relevanceScore, relevanceScore) || other.relevanceScore == relevanceScore));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hashAll([runtimeType,position,title,link,productLink,productId,serpapiProductApi,source,price,extractedPrice,thumbnail,serpapiProductApiComparative,rating,reviews,delivery,tag,const DeepCollectionEquality().hash(_extensions),currency,pageToken,relevanceScore]);

@override
String toString() {
  return 'Product(position: $position, title: $title, link: $link, productLink: $productLink, productId: $productId, serpapiProductApi: $serpapiProductApi, source: $source, price: $price, extractedPrice: $extractedPrice, thumbnail: $thumbnail, serpapiProductApiComparative: $serpapiProductApiComparative, rating: $rating, reviews: $reviews, delivery: $delivery, tag: $tag, extensions: $extensions, currency: $currency, pageToken: $pageToken, relevanceScore: $relevanceScore)';
}


}

/// @nodoc
abstract mixin class _$ProductCopyWith<$Res> implements $ProductCopyWith<$Res> {
  factory _$ProductCopyWith(_Product value, $Res Function(_Product) _then) = __$ProductCopyWithImpl;
@override @useResult
$Res call({
 int position, String title, String link,@JsonKey(name: 'product_link') String productLink,@JsonKey(name: 'product_id') String productId,@JsonKey(name: 'serpapi_product_api') String serpapiProductApi, String source, String price,@JsonKey(name: 'extracted_price') double extractedPrice, String thumbnail,@JsonKey(name: 'serpapi_product_api_comparative') String? serpapiProductApiComparative, double? rating, int? reviews, String? delivery, String? tag, List<String>? extensions, String? currency,@JsonKey(name: 'page_token') String? pageToken,@JsonKey(name: 'relevance_score') double? relevanceScore
});




}
/// @nodoc
class __$ProductCopyWithImpl<$Res>
    implements _$ProductCopyWith<$Res> {
  __$ProductCopyWithImpl(this._self, this._then);

  final _Product _self;
  final $Res Function(_Product) _then;

/// Create a copy of Product
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? position = null,Object? title = null,Object? link = null,Object? productLink = null,Object? productId = null,Object? serpapiProductApi = null,Object? source = null,Object? price = null,Object? extractedPrice = null,Object? thumbnail = null,Object? serpapiProductApiComparative = freezed,Object? rating = freezed,Object? reviews = freezed,Object? delivery = freezed,Object? tag = freezed,Object? extensions = freezed,Object? currency = freezed,Object? pageToken = freezed,Object? relevanceScore = freezed,}) {
  return _then(_Product(
position: null == position ? _self.position : position // ignore: cast_nullable_to_non_nullable
as int,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,link: null == link ? _self.link : link // ignore: cast_nullable_to_non_nullable
as String,productLink: null == productLink ? _self.productLink : productLink // ignore: cast_nullable_to_non_nullable
as String,productId: null == productId ? _self.productId : productId // ignore: cast_nullable_to_non_nullable
as String,serpapiProductApi: null == serpapiProductApi ? _self.serpapiProductApi : serpapiProductApi // ignore: cast_nullable_to_non_nullable
as String,source: null == source ? _self.source : source // ignore: cast_nullable_to_non_nullable
as String,price: null == price ? _self.price : price // ignore: cast_nullable_to_non_nullable
as String,extractedPrice: null == extractedPrice ? _self.extractedPrice : extractedPrice // ignore: cast_nullable_to_non_nullable
as double,thumbnail: null == thumbnail ? _self.thumbnail : thumbnail // ignore: cast_nullable_to_non_nullable
as String,serpapiProductApiComparative: freezed == serpapiProductApiComparative ? _self.serpapiProductApiComparative : serpapiProductApiComparative // ignore: cast_nullable_to_non_nullable
as String?,rating: freezed == rating ? _self.rating : rating // ignore: cast_nullable_to_non_nullable
as double?,reviews: freezed == reviews ? _self.reviews : reviews // ignore: cast_nullable_to_non_nullable
as int?,delivery: freezed == delivery ? _self.delivery : delivery // ignore: cast_nullable_to_non_nullable
as String?,tag: freezed == tag ? _self.tag : tag // ignore: cast_nullable_to_non_nullable
as String?,extensions: freezed == extensions ? _self._extensions : extensions // ignore: cast_nullable_to_non_nullable
as List<String>?,currency: freezed == currency ? _self.currency : currency // ignore: cast_nullable_to_non_nullable
as String?,pageToken: freezed == pageToken ? _self.pageToken : pageToken // ignore: cast_nullable_to_non_nullable
as String?,relevanceScore: freezed == relevanceScore ? _self.relevanceScore : relevanceScore // ignore: cast_nullable_to_non_nullable
as double?,
  ));
}


}

// dart format on
