// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'product_details.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_ProductOffer _$ProductOfferFromJson(Map<String, dynamic> json) =>
    _ProductOffer(
      merchant: json['merchant'] as String,
      price: json['price'] as String,
      link: json['link'] as String,
      logo: json['logo'] as String?,
      extractedPrice: (json['extracted_price'] as num?)?.toDouble(),
      currency: json['currency'] as String?,
      title: json['title'] as String?,
      availability: json['availability'] as String?,
      shipping: json['shipping'] as String?,
      shippingExtracted: (json['shipping_extracted'] as num?)?.toDouble(),
      total: json['total'] as String?,
      extractedTotal: (json['extracted_total'] as num?)?.toDouble(),
      rating: (json['rating'] as num?)?.toDouble(),
      reviews: (json['reviews'] as num?)?.toInt(),
      paymentMethods: json['payment_methods'] as String?,
      tag: json['tag'] as String?,
      detailsAndOffers: (json['details_and_offers'] as List<dynamic>?)
          ?.map((e) => e as String)
          .toList(),
      monthlyPaymentDuration: (json['monthly_payment_duration'] as num?)
          ?.toInt(),
      downPayment: json['down_payment'] as String?,
    );

Map<String, dynamic> _$ProductOfferToJson(_ProductOffer instance) =>
    <String, dynamic>{
      'merchant': instance.merchant,
      'price': instance.price,
      'link': instance.link,
      'logo': instance.logo,
      'extracted_price': instance.extractedPrice,
      'currency': instance.currency,
      'title': instance.title,
      'availability': instance.availability,
      'shipping': instance.shipping,
      'shipping_extracted': instance.shippingExtracted,
      'total': instance.total,
      'extracted_total': instance.extractedTotal,
      'rating': instance.rating,
      'reviews': instance.reviews,
      'payment_methods': instance.paymentMethods,
      'tag': instance.tag,
      'details_and_offers': instance.detailsAndOffers,
      'monthly_payment_duration': instance.monthlyPaymentDuration,
      'down_payment': instance.downPayment,
    };

_Specification _$SpecificationFromJson(Map<String, dynamic> json) =>
    _Specification(
      title: json['title'] as String,
      value: json['value'] as String,
    );

Map<String, dynamic> _$SpecificationToJson(_Specification instance) =>
    <String, dynamic>{'title': instance.title, 'value': instance.value};

_RatingBreakdown _$RatingBreakdownFromJson(Map<String, dynamic> json) =>
    _RatingBreakdown(
      stars: (json['stars'] as num).toInt(),
      amount: (json['amount'] as num).toInt(),
    );

Map<String, dynamic> _$RatingBreakdownToJson(_RatingBreakdown instance) =>
    <String, dynamic>{'stars': instance.stars, 'amount': instance.amount};

_ProductDetailsResponse _$ProductDetailsResponseFromJson(
  Map<String, dynamic> json,
) => _ProductDetailsResponse(
  type: json['type'] as String,
  title: json['title'] as String,
  price: json['price'] as String,
  offers: (json['offers'] as List<dynamic>)
      .map((e) => ProductOffer.fromJson(e as Map<String, dynamic>))
      .toList(),
  rating: (json['rating'] as num?)?.toDouble(),
  reviews: (json['reviews'] as num?)?.toInt(),
  description: json['description'] as String?,
  images: (json['images'] as List<dynamic>?)?.map((e) => e as String).toList(),
  specifications: (json['specifications'] as List<dynamic>?)
      ?.map((e) => Specification.fromJson(e as Map<String, dynamic>))
      .toList(),
  ratingBreakdown: (json['rating_breakdown'] as List<dynamic>?)
      ?.map((e) => RatingBreakdown.fromJson(e as Map<String, dynamic>))
      .toList(),
);

Map<String, dynamic> _$ProductDetailsResponseToJson(
  _ProductDetailsResponse instance,
) => <String, dynamic>{
  'type': instance.type,
  'title': instance.title,
  'price': instance.price,
  'offers': instance.offers,
  'rating': instance.rating,
  'reviews': instance.reviews,
  'description': instance.description,
  'images': instance.images,
  'specifications': instance.specifications,
  'rating_breakdown': instance.ratingBreakdown,
};
