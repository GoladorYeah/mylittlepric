// GENERATED CODE - DO NOT MODIFY BY HAND
// coverage:ignore-file
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of 'product_details.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

// dart format off
T _$identity<T>(T value) => value;

/// @nodoc
mixin _$ProductOffer {

 String get merchant; String get price; String get link; String? get logo;@JsonKey(name: 'extracted_price') double? get extractedPrice; String? get currency; String? get title; String? get availability; String? get shipping;@JsonKey(name: 'shipping_extracted') double? get shippingExtracted; String? get total;@JsonKey(name: 'extracted_total') double? get extractedTotal; double? get rating; int? get reviews;@JsonKey(name: 'payment_methods') String? get paymentMethods; String? get tag;@JsonKey(name: 'details_and_offers') List<String>? get detailsAndOffers;@JsonKey(name: 'monthly_payment_duration') int? get monthlyPaymentDuration;@JsonKey(name: 'down_payment') String? get downPayment;
/// Create a copy of ProductOffer
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ProductOfferCopyWith<ProductOffer> get copyWith => _$ProductOfferCopyWithImpl<ProductOffer>(this as ProductOffer, _$identity);

  /// Serializes this ProductOffer to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ProductOffer&&(identical(other.merchant, merchant) || other.merchant == merchant)&&(identical(other.price, price) || other.price == price)&&(identical(other.link, link) || other.link == link)&&(identical(other.logo, logo) || other.logo == logo)&&(identical(other.extractedPrice, extractedPrice) || other.extractedPrice == extractedPrice)&&(identical(other.currency, currency) || other.currency == currency)&&(identical(other.title, title) || other.title == title)&&(identical(other.availability, availability) || other.availability == availability)&&(identical(other.shipping, shipping) || other.shipping == shipping)&&(identical(other.shippingExtracted, shippingExtracted) || other.shippingExtracted == shippingExtracted)&&(identical(other.total, total) || other.total == total)&&(identical(other.extractedTotal, extractedTotal) || other.extractedTotal == extractedTotal)&&(identical(other.rating, rating) || other.rating == rating)&&(identical(other.reviews, reviews) || other.reviews == reviews)&&(identical(other.paymentMethods, paymentMethods) || other.paymentMethods == paymentMethods)&&(identical(other.tag, tag) || other.tag == tag)&&const DeepCollectionEquality().equals(other.detailsAndOffers, detailsAndOffers)&&(identical(other.monthlyPaymentDuration, monthlyPaymentDuration) || other.monthlyPaymentDuration == monthlyPaymentDuration)&&(identical(other.downPayment, downPayment) || other.downPayment == downPayment));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hashAll([runtimeType,merchant,price,link,logo,extractedPrice,currency,title,availability,shipping,shippingExtracted,total,extractedTotal,rating,reviews,paymentMethods,tag,const DeepCollectionEquality().hash(detailsAndOffers),monthlyPaymentDuration,downPayment]);

@override
String toString() {
  return 'ProductOffer(merchant: $merchant, price: $price, link: $link, logo: $logo, extractedPrice: $extractedPrice, currency: $currency, title: $title, availability: $availability, shipping: $shipping, shippingExtracted: $shippingExtracted, total: $total, extractedTotal: $extractedTotal, rating: $rating, reviews: $reviews, paymentMethods: $paymentMethods, tag: $tag, detailsAndOffers: $detailsAndOffers, monthlyPaymentDuration: $monthlyPaymentDuration, downPayment: $downPayment)';
}


}

/// @nodoc
abstract mixin class $ProductOfferCopyWith<$Res>  {
  factory $ProductOfferCopyWith(ProductOffer value, $Res Function(ProductOffer) _then) = _$ProductOfferCopyWithImpl;
@useResult
$Res call({
 String merchant, String price, String link, String? logo,@JsonKey(name: 'extracted_price') double? extractedPrice, String? currency, String? title, String? availability, String? shipping,@JsonKey(name: 'shipping_extracted') double? shippingExtracted, String? total,@JsonKey(name: 'extracted_total') double? extractedTotal, double? rating, int? reviews,@JsonKey(name: 'payment_methods') String? paymentMethods, String? tag,@JsonKey(name: 'details_and_offers') List<String>? detailsAndOffers,@JsonKey(name: 'monthly_payment_duration') int? monthlyPaymentDuration,@JsonKey(name: 'down_payment') String? downPayment
});




}
/// @nodoc
class _$ProductOfferCopyWithImpl<$Res>
    implements $ProductOfferCopyWith<$Res> {
  _$ProductOfferCopyWithImpl(this._self, this._then);

  final ProductOffer _self;
  final $Res Function(ProductOffer) _then;

/// Create a copy of ProductOffer
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? merchant = null,Object? price = null,Object? link = null,Object? logo = freezed,Object? extractedPrice = freezed,Object? currency = freezed,Object? title = freezed,Object? availability = freezed,Object? shipping = freezed,Object? shippingExtracted = freezed,Object? total = freezed,Object? extractedTotal = freezed,Object? rating = freezed,Object? reviews = freezed,Object? paymentMethods = freezed,Object? tag = freezed,Object? detailsAndOffers = freezed,Object? monthlyPaymentDuration = freezed,Object? downPayment = freezed,}) {
  return _then(_self.copyWith(
merchant: null == merchant ? _self.merchant : merchant // ignore: cast_nullable_to_non_nullable
as String,price: null == price ? _self.price : price // ignore: cast_nullable_to_non_nullable
as String,link: null == link ? _self.link : link // ignore: cast_nullable_to_non_nullable
as String,logo: freezed == logo ? _self.logo : logo // ignore: cast_nullable_to_non_nullable
as String?,extractedPrice: freezed == extractedPrice ? _self.extractedPrice : extractedPrice // ignore: cast_nullable_to_non_nullable
as double?,currency: freezed == currency ? _self.currency : currency // ignore: cast_nullable_to_non_nullable
as String?,title: freezed == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String?,availability: freezed == availability ? _self.availability : availability // ignore: cast_nullable_to_non_nullable
as String?,shipping: freezed == shipping ? _self.shipping : shipping // ignore: cast_nullable_to_non_nullable
as String?,shippingExtracted: freezed == shippingExtracted ? _self.shippingExtracted : shippingExtracted // ignore: cast_nullable_to_non_nullable
as double?,total: freezed == total ? _self.total : total // ignore: cast_nullable_to_non_nullable
as String?,extractedTotal: freezed == extractedTotal ? _self.extractedTotal : extractedTotal // ignore: cast_nullable_to_non_nullable
as double?,rating: freezed == rating ? _self.rating : rating // ignore: cast_nullable_to_non_nullable
as double?,reviews: freezed == reviews ? _self.reviews : reviews // ignore: cast_nullable_to_non_nullable
as int?,paymentMethods: freezed == paymentMethods ? _self.paymentMethods : paymentMethods // ignore: cast_nullable_to_non_nullable
as String?,tag: freezed == tag ? _self.tag : tag // ignore: cast_nullable_to_non_nullable
as String?,detailsAndOffers: freezed == detailsAndOffers ? _self.detailsAndOffers : detailsAndOffers // ignore: cast_nullable_to_non_nullable
as List<String>?,monthlyPaymentDuration: freezed == monthlyPaymentDuration ? _self.monthlyPaymentDuration : monthlyPaymentDuration // ignore: cast_nullable_to_non_nullable
as int?,downPayment: freezed == downPayment ? _self.downPayment : downPayment // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}

}


/// Adds pattern-matching-related methods to [ProductOffer].
extension ProductOfferPatterns on ProductOffer {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ProductOffer value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ProductOffer() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ProductOffer value)  $default,){
final _that = this;
switch (_that) {
case _ProductOffer():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ProductOffer value)?  $default,){
final _that = this;
switch (_that) {
case _ProductOffer() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String merchant,  String price,  String link,  String? logo, @JsonKey(name: 'extracted_price')  double? extractedPrice,  String? currency,  String? title,  String? availability,  String? shipping, @JsonKey(name: 'shipping_extracted')  double? shippingExtracted,  String? total, @JsonKey(name: 'extracted_total')  double? extractedTotal,  double? rating,  int? reviews, @JsonKey(name: 'payment_methods')  String? paymentMethods,  String? tag, @JsonKey(name: 'details_and_offers')  List<String>? detailsAndOffers, @JsonKey(name: 'monthly_payment_duration')  int? monthlyPaymentDuration, @JsonKey(name: 'down_payment')  String? downPayment)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ProductOffer() when $default != null:
return $default(_that.merchant,_that.price,_that.link,_that.logo,_that.extractedPrice,_that.currency,_that.title,_that.availability,_that.shipping,_that.shippingExtracted,_that.total,_that.extractedTotal,_that.rating,_that.reviews,_that.paymentMethods,_that.tag,_that.detailsAndOffers,_that.monthlyPaymentDuration,_that.downPayment);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String merchant,  String price,  String link,  String? logo, @JsonKey(name: 'extracted_price')  double? extractedPrice,  String? currency,  String? title,  String? availability,  String? shipping, @JsonKey(name: 'shipping_extracted')  double? shippingExtracted,  String? total, @JsonKey(name: 'extracted_total')  double? extractedTotal,  double? rating,  int? reviews, @JsonKey(name: 'payment_methods')  String? paymentMethods,  String? tag, @JsonKey(name: 'details_and_offers')  List<String>? detailsAndOffers, @JsonKey(name: 'monthly_payment_duration')  int? monthlyPaymentDuration, @JsonKey(name: 'down_payment')  String? downPayment)  $default,) {final _that = this;
switch (_that) {
case _ProductOffer():
return $default(_that.merchant,_that.price,_that.link,_that.logo,_that.extractedPrice,_that.currency,_that.title,_that.availability,_that.shipping,_that.shippingExtracted,_that.total,_that.extractedTotal,_that.rating,_that.reviews,_that.paymentMethods,_that.tag,_that.detailsAndOffers,_that.monthlyPaymentDuration,_that.downPayment);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String merchant,  String price,  String link,  String? logo, @JsonKey(name: 'extracted_price')  double? extractedPrice,  String? currency,  String? title,  String? availability,  String? shipping, @JsonKey(name: 'shipping_extracted')  double? shippingExtracted,  String? total, @JsonKey(name: 'extracted_total')  double? extractedTotal,  double? rating,  int? reviews, @JsonKey(name: 'payment_methods')  String? paymentMethods,  String? tag, @JsonKey(name: 'details_and_offers')  List<String>? detailsAndOffers, @JsonKey(name: 'monthly_payment_duration')  int? monthlyPaymentDuration, @JsonKey(name: 'down_payment')  String? downPayment)?  $default,) {final _that = this;
switch (_that) {
case _ProductOffer() when $default != null:
return $default(_that.merchant,_that.price,_that.link,_that.logo,_that.extractedPrice,_that.currency,_that.title,_that.availability,_that.shipping,_that.shippingExtracted,_that.total,_that.extractedTotal,_that.rating,_that.reviews,_that.paymentMethods,_that.tag,_that.detailsAndOffers,_that.monthlyPaymentDuration,_that.downPayment);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ProductOffer implements ProductOffer {
  const _ProductOffer({required this.merchant, required this.price, required this.link, this.logo, @JsonKey(name: 'extracted_price') this.extractedPrice, this.currency, this.title, this.availability, this.shipping, @JsonKey(name: 'shipping_extracted') this.shippingExtracted, this.total, @JsonKey(name: 'extracted_total') this.extractedTotal, this.rating, this.reviews, @JsonKey(name: 'payment_methods') this.paymentMethods, this.tag, @JsonKey(name: 'details_and_offers') final  List<String>? detailsAndOffers, @JsonKey(name: 'monthly_payment_duration') this.monthlyPaymentDuration, @JsonKey(name: 'down_payment') this.downPayment}): _detailsAndOffers = detailsAndOffers;
  factory _ProductOffer.fromJson(Map<String, dynamic> json) => _$ProductOfferFromJson(json);

@override final  String merchant;
@override final  String price;
@override final  String link;
@override final  String? logo;
@override@JsonKey(name: 'extracted_price') final  double? extractedPrice;
@override final  String? currency;
@override final  String? title;
@override final  String? availability;
@override final  String? shipping;
@override@JsonKey(name: 'shipping_extracted') final  double? shippingExtracted;
@override final  String? total;
@override@JsonKey(name: 'extracted_total') final  double? extractedTotal;
@override final  double? rating;
@override final  int? reviews;
@override@JsonKey(name: 'payment_methods') final  String? paymentMethods;
@override final  String? tag;
 final  List<String>? _detailsAndOffers;
@override@JsonKey(name: 'details_and_offers') List<String>? get detailsAndOffers {
  final value = _detailsAndOffers;
  if (value == null) return null;
  if (_detailsAndOffers is EqualUnmodifiableListView) return _detailsAndOffers;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(value);
}

@override@JsonKey(name: 'monthly_payment_duration') final  int? monthlyPaymentDuration;
@override@JsonKey(name: 'down_payment') final  String? downPayment;

/// Create a copy of ProductOffer
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ProductOfferCopyWith<_ProductOffer> get copyWith => __$ProductOfferCopyWithImpl<_ProductOffer>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ProductOfferToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ProductOffer&&(identical(other.merchant, merchant) || other.merchant == merchant)&&(identical(other.price, price) || other.price == price)&&(identical(other.link, link) || other.link == link)&&(identical(other.logo, logo) || other.logo == logo)&&(identical(other.extractedPrice, extractedPrice) || other.extractedPrice == extractedPrice)&&(identical(other.currency, currency) || other.currency == currency)&&(identical(other.title, title) || other.title == title)&&(identical(other.availability, availability) || other.availability == availability)&&(identical(other.shipping, shipping) || other.shipping == shipping)&&(identical(other.shippingExtracted, shippingExtracted) || other.shippingExtracted == shippingExtracted)&&(identical(other.total, total) || other.total == total)&&(identical(other.extractedTotal, extractedTotal) || other.extractedTotal == extractedTotal)&&(identical(other.rating, rating) || other.rating == rating)&&(identical(other.reviews, reviews) || other.reviews == reviews)&&(identical(other.paymentMethods, paymentMethods) || other.paymentMethods == paymentMethods)&&(identical(other.tag, tag) || other.tag == tag)&&const DeepCollectionEquality().equals(other._detailsAndOffers, _detailsAndOffers)&&(identical(other.monthlyPaymentDuration, monthlyPaymentDuration) || other.monthlyPaymentDuration == monthlyPaymentDuration)&&(identical(other.downPayment, downPayment) || other.downPayment == downPayment));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hashAll([runtimeType,merchant,price,link,logo,extractedPrice,currency,title,availability,shipping,shippingExtracted,total,extractedTotal,rating,reviews,paymentMethods,tag,const DeepCollectionEquality().hash(_detailsAndOffers),monthlyPaymentDuration,downPayment]);

@override
String toString() {
  return 'ProductOffer(merchant: $merchant, price: $price, link: $link, logo: $logo, extractedPrice: $extractedPrice, currency: $currency, title: $title, availability: $availability, shipping: $shipping, shippingExtracted: $shippingExtracted, total: $total, extractedTotal: $extractedTotal, rating: $rating, reviews: $reviews, paymentMethods: $paymentMethods, tag: $tag, detailsAndOffers: $detailsAndOffers, monthlyPaymentDuration: $monthlyPaymentDuration, downPayment: $downPayment)';
}


}

/// @nodoc
abstract mixin class _$ProductOfferCopyWith<$Res> implements $ProductOfferCopyWith<$Res> {
  factory _$ProductOfferCopyWith(_ProductOffer value, $Res Function(_ProductOffer) _then) = __$ProductOfferCopyWithImpl;
@override @useResult
$Res call({
 String merchant, String price, String link, String? logo,@JsonKey(name: 'extracted_price') double? extractedPrice, String? currency, String? title, String? availability, String? shipping,@JsonKey(name: 'shipping_extracted') double? shippingExtracted, String? total,@JsonKey(name: 'extracted_total') double? extractedTotal, double? rating, int? reviews,@JsonKey(name: 'payment_methods') String? paymentMethods, String? tag,@JsonKey(name: 'details_and_offers') List<String>? detailsAndOffers,@JsonKey(name: 'monthly_payment_duration') int? monthlyPaymentDuration,@JsonKey(name: 'down_payment') String? downPayment
});




}
/// @nodoc
class __$ProductOfferCopyWithImpl<$Res>
    implements _$ProductOfferCopyWith<$Res> {
  __$ProductOfferCopyWithImpl(this._self, this._then);

  final _ProductOffer _self;
  final $Res Function(_ProductOffer) _then;

/// Create a copy of ProductOffer
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? merchant = null,Object? price = null,Object? link = null,Object? logo = freezed,Object? extractedPrice = freezed,Object? currency = freezed,Object? title = freezed,Object? availability = freezed,Object? shipping = freezed,Object? shippingExtracted = freezed,Object? total = freezed,Object? extractedTotal = freezed,Object? rating = freezed,Object? reviews = freezed,Object? paymentMethods = freezed,Object? tag = freezed,Object? detailsAndOffers = freezed,Object? monthlyPaymentDuration = freezed,Object? downPayment = freezed,}) {
  return _then(_ProductOffer(
merchant: null == merchant ? _self.merchant : merchant // ignore: cast_nullable_to_non_nullable
as String,price: null == price ? _self.price : price // ignore: cast_nullable_to_non_nullable
as String,link: null == link ? _self.link : link // ignore: cast_nullable_to_non_nullable
as String,logo: freezed == logo ? _self.logo : logo // ignore: cast_nullable_to_non_nullable
as String?,extractedPrice: freezed == extractedPrice ? _self.extractedPrice : extractedPrice // ignore: cast_nullable_to_non_nullable
as double?,currency: freezed == currency ? _self.currency : currency // ignore: cast_nullable_to_non_nullable
as String?,title: freezed == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String?,availability: freezed == availability ? _self.availability : availability // ignore: cast_nullable_to_non_nullable
as String?,shipping: freezed == shipping ? _self.shipping : shipping // ignore: cast_nullable_to_non_nullable
as String?,shippingExtracted: freezed == shippingExtracted ? _self.shippingExtracted : shippingExtracted // ignore: cast_nullable_to_non_nullable
as double?,total: freezed == total ? _self.total : total // ignore: cast_nullable_to_non_nullable
as String?,extractedTotal: freezed == extractedTotal ? _self.extractedTotal : extractedTotal // ignore: cast_nullable_to_non_nullable
as double?,rating: freezed == rating ? _self.rating : rating // ignore: cast_nullable_to_non_nullable
as double?,reviews: freezed == reviews ? _self.reviews : reviews // ignore: cast_nullable_to_non_nullable
as int?,paymentMethods: freezed == paymentMethods ? _self.paymentMethods : paymentMethods // ignore: cast_nullable_to_non_nullable
as String?,tag: freezed == tag ? _self.tag : tag // ignore: cast_nullable_to_non_nullable
as String?,detailsAndOffers: freezed == detailsAndOffers ? _self._detailsAndOffers : detailsAndOffers // ignore: cast_nullable_to_non_nullable
as List<String>?,monthlyPaymentDuration: freezed == monthlyPaymentDuration ? _self.monthlyPaymentDuration : monthlyPaymentDuration // ignore: cast_nullable_to_non_nullable
as int?,downPayment: freezed == downPayment ? _self.downPayment : downPayment // ignore: cast_nullable_to_non_nullable
as String?,
  ));
}


}


/// @nodoc
mixin _$Specification {

 String get title; String get value;
/// Create a copy of Specification
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$SpecificationCopyWith<Specification> get copyWith => _$SpecificationCopyWithImpl<Specification>(this as Specification, _$identity);

  /// Serializes this Specification to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is Specification&&(identical(other.title, title) || other.title == title)&&(identical(other.value, value) || other.value == value));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,title,value);

@override
String toString() {
  return 'Specification(title: $title, value: $value)';
}


}

/// @nodoc
abstract mixin class $SpecificationCopyWith<$Res>  {
  factory $SpecificationCopyWith(Specification value, $Res Function(Specification) _then) = _$SpecificationCopyWithImpl;
@useResult
$Res call({
 String title, String value
});




}
/// @nodoc
class _$SpecificationCopyWithImpl<$Res>
    implements $SpecificationCopyWith<$Res> {
  _$SpecificationCopyWithImpl(this._self, this._then);

  final Specification _self;
  final $Res Function(Specification) _then;

/// Create a copy of Specification
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? title = null,Object? value = null,}) {
  return _then(_self.copyWith(
title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,value: null == value ? _self.value : value // ignore: cast_nullable_to_non_nullable
as String,
  ));
}

}


/// Adds pattern-matching-related methods to [Specification].
extension SpecificationPatterns on Specification {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _Specification value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _Specification() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _Specification value)  $default,){
final _that = this;
switch (_that) {
case _Specification():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _Specification value)?  $default,){
final _that = this;
switch (_that) {
case _Specification() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String title,  String value)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _Specification() when $default != null:
return $default(_that.title,_that.value);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String title,  String value)  $default,) {final _that = this;
switch (_that) {
case _Specification():
return $default(_that.title,_that.value);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String title,  String value)?  $default,) {final _that = this;
switch (_that) {
case _Specification() when $default != null:
return $default(_that.title,_that.value);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _Specification implements Specification {
  const _Specification({required this.title, required this.value});
  factory _Specification.fromJson(Map<String, dynamic> json) => _$SpecificationFromJson(json);

@override final  String title;
@override final  String value;

/// Create a copy of Specification
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$SpecificationCopyWith<_Specification> get copyWith => __$SpecificationCopyWithImpl<_Specification>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$SpecificationToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _Specification&&(identical(other.title, title) || other.title == title)&&(identical(other.value, value) || other.value == value));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,title,value);

@override
String toString() {
  return 'Specification(title: $title, value: $value)';
}


}

/// @nodoc
abstract mixin class _$SpecificationCopyWith<$Res> implements $SpecificationCopyWith<$Res> {
  factory _$SpecificationCopyWith(_Specification value, $Res Function(_Specification) _then) = __$SpecificationCopyWithImpl;
@override @useResult
$Res call({
 String title, String value
});




}
/// @nodoc
class __$SpecificationCopyWithImpl<$Res>
    implements _$SpecificationCopyWith<$Res> {
  __$SpecificationCopyWithImpl(this._self, this._then);

  final _Specification _self;
  final $Res Function(_Specification) _then;

/// Create a copy of Specification
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? title = null,Object? value = null,}) {
  return _then(_Specification(
title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,value: null == value ? _self.value : value // ignore: cast_nullable_to_non_nullable
as String,
  ));
}


}


/// @nodoc
mixin _$RatingBreakdown {

 int get stars; int get amount;
/// Create a copy of RatingBreakdown
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$RatingBreakdownCopyWith<RatingBreakdown> get copyWith => _$RatingBreakdownCopyWithImpl<RatingBreakdown>(this as RatingBreakdown, _$identity);

  /// Serializes this RatingBreakdown to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is RatingBreakdown&&(identical(other.stars, stars) || other.stars == stars)&&(identical(other.amount, amount) || other.amount == amount));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,stars,amount);

@override
String toString() {
  return 'RatingBreakdown(stars: $stars, amount: $amount)';
}


}

/// @nodoc
abstract mixin class $RatingBreakdownCopyWith<$Res>  {
  factory $RatingBreakdownCopyWith(RatingBreakdown value, $Res Function(RatingBreakdown) _then) = _$RatingBreakdownCopyWithImpl;
@useResult
$Res call({
 int stars, int amount
});




}
/// @nodoc
class _$RatingBreakdownCopyWithImpl<$Res>
    implements $RatingBreakdownCopyWith<$Res> {
  _$RatingBreakdownCopyWithImpl(this._self, this._then);

  final RatingBreakdown _self;
  final $Res Function(RatingBreakdown) _then;

/// Create a copy of RatingBreakdown
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? stars = null,Object? amount = null,}) {
  return _then(_self.copyWith(
stars: null == stars ? _self.stars : stars // ignore: cast_nullable_to_non_nullable
as int,amount: null == amount ? _self.amount : amount // ignore: cast_nullable_to_non_nullable
as int,
  ));
}

}


/// Adds pattern-matching-related methods to [RatingBreakdown].
extension RatingBreakdownPatterns on RatingBreakdown {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _RatingBreakdown value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _RatingBreakdown() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _RatingBreakdown value)  $default,){
final _that = this;
switch (_that) {
case _RatingBreakdown():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _RatingBreakdown value)?  $default,){
final _that = this;
switch (_that) {
case _RatingBreakdown() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( int stars,  int amount)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _RatingBreakdown() when $default != null:
return $default(_that.stars,_that.amount);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( int stars,  int amount)  $default,) {final _that = this;
switch (_that) {
case _RatingBreakdown():
return $default(_that.stars,_that.amount);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( int stars,  int amount)?  $default,) {final _that = this;
switch (_that) {
case _RatingBreakdown() when $default != null:
return $default(_that.stars,_that.amount);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _RatingBreakdown implements RatingBreakdown {
  const _RatingBreakdown({required this.stars, required this.amount});
  factory _RatingBreakdown.fromJson(Map<String, dynamic> json) => _$RatingBreakdownFromJson(json);

@override final  int stars;
@override final  int amount;

/// Create a copy of RatingBreakdown
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$RatingBreakdownCopyWith<_RatingBreakdown> get copyWith => __$RatingBreakdownCopyWithImpl<_RatingBreakdown>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$RatingBreakdownToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _RatingBreakdown&&(identical(other.stars, stars) || other.stars == stars)&&(identical(other.amount, amount) || other.amount == amount));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,stars,amount);

@override
String toString() {
  return 'RatingBreakdown(stars: $stars, amount: $amount)';
}


}

/// @nodoc
abstract mixin class _$RatingBreakdownCopyWith<$Res> implements $RatingBreakdownCopyWith<$Res> {
  factory _$RatingBreakdownCopyWith(_RatingBreakdown value, $Res Function(_RatingBreakdown) _then) = __$RatingBreakdownCopyWithImpl;
@override @useResult
$Res call({
 int stars, int amount
});




}
/// @nodoc
class __$RatingBreakdownCopyWithImpl<$Res>
    implements _$RatingBreakdownCopyWith<$Res> {
  __$RatingBreakdownCopyWithImpl(this._self, this._then);

  final _RatingBreakdown _self;
  final $Res Function(_RatingBreakdown) _then;

/// Create a copy of RatingBreakdown
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? stars = null,Object? amount = null,}) {
  return _then(_RatingBreakdown(
stars: null == stars ? _self.stars : stars // ignore: cast_nullable_to_non_nullable
as int,amount: null == amount ? _self.amount : amount // ignore: cast_nullable_to_non_nullable
as int,
  ));
}


}


/// @nodoc
mixin _$ProductDetailsResponse {

 String get type; String get title; String get price;// variants would need more complex modeling
 List<ProductOffer> get offers; double? get rating; int? get reviews; String? get description; List<String>? get images; List<Specification>? get specifications;// videos and more_options could be added if needed
@JsonKey(name: 'rating_breakdown') List<RatingBreakdown>? get ratingBreakdown;
/// Create a copy of ProductDetailsResponse
/// with the given fields replaced by the non-null parameter values.
@JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
$ProductDetailsResponseCopyWith<ProductDetailsResponse> get copyWith => _$ProductDetailsResponseCopyWithImpl<ProductDetailsResponse>(this as ProductDetailsResponse, _$identity);

  /// Serializes this ProductDetailsResponse to a JSON map.
  Map<String, dynamic> toJson();


@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is ProductDetailsResponse&&(identical(other.type, type) || other.type == type)&&(identical(other.title, title) || other.title == title)&&(identical(other.price, price) || other.price == price)&&const DeepCollectionEquality().equals(other.offers, offers)&&(identical(other.rating, rating) || other.rating == rating)&&(identical(other.reviews, reviews) || other.reviews == reviews)&&(identical(other.description, description) || other.description == description)&&const DeepCollectionEquality().equals(other.images, images)&&const DeepCollectionEquality().equals(other.specifications, specifications)&&const DeepCollectionEquality().equals(other.ratingBreakdown, ratingBreakdown));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,type,title,price,const DeepCollectionEquality().hash(offers),rating,reviews,description,const DeepCollectionEquality().hash(images),const DeepCollectionEquality().hash(specifications),const DeepCollectionEquality().hash(ratingBreakdown));

@override
String toString() {
  return 'ProductDetailsResponse(type: $type, title: $title, price: $price, offers: $offers, rating: $rating, reviews: $reviews, description: $description, images: $images, specifications: $specifications, ratingBreakdown: $ratingBreakdown)';
}


}

/// @nodoc
abstract mixin class $ProductDetailsResponseCopyWith<$Res>  {
  factory $ProductDetailsResponseCopyWith(ProductDetailsResponse value, $Res Function(ProductDetailsResponse) _then) = _$ProductDetailsResponseCopyWithImpl;
@useResult
$Res call({
 String type, String title, String price, List<ProductOffer> offers, double? rating, int? reviews, String? description, List<String>? images, List<Specification>? specifications,@JsonKey(name: 'rating_breakdown') List<RatingBreakdown>? ratingBreakdown
});




}
/// @nodoc
class _$ProductDetailsResponseCopyWithImpl<$Res>
    implements $ProductDetailsResponseCopyWith<$Res> {
  _$ProductDetailsResponseCopyWithImpl(this._self, this._then);

  final ProductDetailsResponse _self;
  final $Res Function(ProductDetailsResponse) _then;

/// Create a copy of ProductDetailsResponse
/// with the given fields replaced by the non-null parameter values.
@pragma('vm:prefer-inline') @override $Res call({Object? type = null,Object? title = null,Object? price = null,Object? offers = null,Object? rating = freezed,Object? reviews = freezed,Object? description = freezed,Object? images = freezed,Object? specifications = freezed,Object? ratingBreakdown = freezed,}) {
  return _then(_self.copyWith(
type: null == type ? _self.type : type // ignore: cast_nullable_to_non_nullable
as String,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,price: null == price ? _self.price : price // ignore: cast_nullable_to_non_nullable
as String,offers: null == offers ? _self.offers : offers // ignore: cast_nullable_to_non_nullable
as List<ProductOffer>,rating: freezed == rating ? _self.rating : rating // ignore: cast_nullable_to_non_nullable
as double?,reviews: freezed == reviews ? _self.reviews : reviews // ignore: cast_nullable_to_non_nullable
as int?,description: freezed == description ? _self.description : description // ignore: cast_nullable_to_non_nullable
as String?,images: freezed == images ? _self.images : images // ignore: cast_nullable_to_non_nullable
as List<String>?,specifications: freezed == specifications ? _self.specifications : specifications // ignore: cast_nullable_to_non_nullable
as List<Specification>?,ratingBreakdown: freezed == ratingBreakdown ? _self.ratingBreakdown : ratingBreakdown // ignore: cast_nullable_to_non_nullable
as List<RatingBreakdown>?,
  ));
}

}


/// Adds pattern-matching-related methods to [ProductDetailsResponse].
extension ProductDetailsResponsePatterns on ProductDetailsResponse {
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

@optionalTypeArgs TResult maybeMap<TResult extends Object?>(TResult Function( _ProductDetailsResponse value)?  $default,{required TResult orElse(),}){
final _that = this;
switch (_that) {
case _ProductDetailsResponse() when $default != null:
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

@optionalTypeArgs TResult map<TResult extends Object?>(TResult Function( _ProductDetailsResponse value)  $default,){
final _that = this;
switch (_that) {
case _ProductDetailsResponse():
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

@optionalTypeArgs TResult? mapOrNull<TResult extends Object?>(TResult? Function( _ProductDetailsResponse value)?  $default,){
final _that = this;
switch (_that) {
case _ProductDetailsResponse() when $default != null:
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

@optionalTypeArgs TResult maybeWhen<TResult extends Object?>(TResult Function( String type,  String title,  String price,  List<ProductOffer> offers,  double? rating,  int? reviews,  String? description,  List<String>? images,  List<Specification>? specifications, @JsonKey(name: 'rating_breakdown')  List<RatingBreakdown>? ratingBreakdown)?  $default,{required TResult orElse(),}) {final _that = this;
switch (_that) {
case _ProductDetailsResponse() when $default != null:
return $default(_that.type,_that.title,_that.price,_that.offers,_that.rating,_that.reviews,_that.description,_that.images,_that.specifications,_that.ratingBreakdown);case _:
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

@optionalTypeArgs TResult when<TResult extends Object?>(TResult Function( String type,  String title,  String price,  List<ProductOffer> offers,  double? rating,  int? reviews,  String? description,  List<String>? images,  List<Specification>? specifications, @JsonKey(name: 'rating_breakdown')  List<RatingBreakdown>? ratingBreakdown)  $default,) {final _that = this;
switch (_that) {
case _ProductDetailsResponse():
return $default(_that.type,_that.title,_that.price,_that.offers,_that.rating,_that.reviews,_that.description,_that.images,_that.specifications,_that.ratingBreakdown);case _:
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

@optionalTypeArgs TResult? whenOrNull<TResult extends Object?>(TResult? Function( String type,  String title,  String price,  List<ProductOffer> offers,  double? rating,  int? reviews,  String? description,  List<String>? images,  List<Specification>? specifications, @JsonKey(name: 'rating_breakdown')  List<RatingBreakdown>? ratingBreakdown)?  $default,) {final _that = this;
switch (_that) {
case _ProductDetailsResponse() when $default != null:
return $default(_that.type,_that.title,_that.price,_that.offers,_that.rating,_that.reviews,_that.description,_that.images,_that.specifications,_that.ratingBreakdown);case _:
  return null;

}
}

}

/// @nodoc
@JsonSerializable()

class _ProductDetailsResponse implements ProductDetailsResponse {
  const _ProductDetailsResponse({required this.type, required this.title, required this.price, required final  List<ProductOffer> offers, this.rating, this.reviews, this.description, final  List<String>? images, final  List<Specification>? specifications, @JsonKey(name: 'rating_breakdown') final  List<RatingBreakdown>? ratingBreakdown}): _offers = offers,_images = images,_specifications = specifications,_ratingBreakdown = ratingBreakdown;
  factory _ProductDetailsResponse.fromJson(Map<String, dynamic> json) => _$ProductDetailsResponseFromJson(json);

@override final  String type;
@override final  String title;
@override final  String price;
// variants would need more complex modeling
 final  List<ProductOffer> _offers;
// variants would need more complex modeling
@override List<ProductOffer> get offers {
  if (_offers is EqualUnmodifiableListView) return _offers;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(_offers);
}

@override final  double? rating;
@override final  int? reviews;
@override final  String? description;
 final  List<String>? _images;
@override List<String>? get images {
  final value = _images;
  if (value == null) return null;
  if (_images is EqualUnmodifiableListView) return _images;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(value);
}

 final  List<Specification>? _specifications;
@override List<Specification>? get specifications {
  final value = _specifications;
  if (value == null) return null;
  if (_specifications is EqualUnmodifiableListView) return _specifications;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(value);
}

// videos and more_options could be added if needed
 final  List<RatingBreakdown>? _ratingBreakdown;
// videos and more_options could be added if needed
@override@JsonKey(name: 'rating_breakdown') List<RatingBreakdown>? get ratingBreakdown {
  final value = _ratingBreakdown;
  if (value == null) return null;
  if (_ratingBreakdown is EqualUnmodifiableListView) return _ratingBreakdown;
  // ignore: implicit_dynamic_type
  return EqualUnmodifiableListView(value);
}


/// Create a copy of ProductDetailsResponse
/// with the given fields replaced by the non-null parameter values.
@override @JsonKey(includeFromJson: false, includeToJson: false)
@pragma('vm:prefer-inline')
_$ProductDetailsResponseCopyWith<_ProductDetailsResponse> get copyWith => __$ProductDetailsResponseCopyWithImpl<_ProductDetailsResponse>(this, _$identity);

@override
Map<String, dynamic> toJson() {
  return _$ProductDetailsResponseToJson(this, );
}

@override
bool operator ==(Object other) {
  return identical(this, other) || (other.runtimeType == runtimeType&&other is _ProductDetailsResponse&&(identical(other.type, type) || other.type == type)&&(identical(other.title, title) || other.title == title)&&(identical(other.price, price) || other.price == price)&&const DeepCollectionEquality().equals(other._offers, _offers)&&(identical(other.rating, rating) || other.rating == rating)&&(identical(other.reviews, reviews) || other.reviews == reviews)&&(identical(other.description, description) || other.description == description)&&const DeepCollectionEquality().equals(other._images, _images)&&const DeepCollectionEquality().equals(other._specifications, _specifications)&&const DeepCollectionEquality().equals(other._ratingBreakdown, _ratingBreakdown));
}

@JsonKey(includeFromJson: false, includeToJson: false)
@override
int get hashCode => Object.hash(runtimeType,type,title,price,const DeepCollectionEquality().hash(_offers),rating,reviews,description,const DeepCollectionEquality().hash(_images),const DeepCollectionEquality().hash(_specifications),const DeepCollectionEquality().hash(_ratingBreakdown));

@override
String toString() {
  return 'ProductDetailsResponse(type: $type, title: $title, price: $price, offers: $offers, rating: $rating, reviews: $reviews, description: $description, images: $images, specifications: $specifications, ratingBreakdown: $ratingBreakdown)';
}


}

/// @nodoc
abstract mixin class _$ProductDetailsResponseCopyWith<$Res> implements $ProductDetailsResponseCopyWith<$Res> {
  factory _$ProductDetailsResponseCopyWith(_ProductDetailsResponse value, $Res Function(_ProductDetailsResponse) _then) = __$ProductDetailsResponseCopyWithImpl;
@override @useResult
$Res call({
 String type, String title, String price, List<ProductOffer> offers, double? rating, int? reviews, String? description, List<String>? images, List<Specification>? specifications,@JsonKey(name: 'rating_breakdown') List<RatingBreakdown>? ratingBreakdown
});




}
/// @nodoc
class __$ProductDetailsResponseCopyWithImpl<$Res>
    implements _$ProductDetailsResponseCopyWith<$Res> {
  __$ProductDetailsResponseCopyWithImpl(this._self, this._then);

  final _ProductDetailsResponse _self;
  final $Res Function(_ProductDetailsResponse) _then;

/// Create a copy of ProductDetailsResponse
/// with the given fields replaced by the non-null parameter values.
@override @pragma('vm:prefer-inline') $Res call({Object? type = null,Object? title = null,Object? price = null,Object? offers = null,Object? rating = freezed,Object? reviews = freezed,Object? description = freezed,Object? images = freezed,Object? specifications = freezed,Object? ratingBreakdown = freezed,}) {
  return _then(_ProductDetailsResponse(
type: null == type ? _self.type : type // ignore: cast_nullable_to_non_nullable
as String,title: null == title ? _self.title : title // ignore: cast_nullable_to_non_nullable
as String,price: null == price ? _self.price : price // ignore: cast_nullable_to_non_nullable
as String,offers: null == offers ? _self._offers : offers // ignore: cast_nullable_to_non_nullable
as List<ProductOffer>,rating: freezed == rating ? _self.rating : rating // ignore: cast_nullable_to_non_nullable
as double?,reviews: freezed == reviews ? _self.reviews : reviews // ignore: cast_nullable_to_non_nullable
as int?,description: freezed == description ? _self.description : description // ignore: cast_nullable_to_non_nullable
as String?,images: freezed == images ? _self._images : images // ignore: cast_nullable_to_non_nullable
as List<String>?,specifications: freezed == specifications ? _self._specifications : specifications // ignore: cast_nullable_to_non_nullable
as List<Specification>?,ratingBreakdown: freezed == ratingBreakdown ? _self._ratingBreakdown : ratingBreakdown // ignore: cast_nullable_to_non_nullable
as List<RatingBreakdown>?,
  ));
}


}

// dart format on
