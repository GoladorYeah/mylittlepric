// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'product.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_Product _$ProductFromJson(Map<String, dynamic> json) => _Product(
  position: (json['position'] as num).toInt(),
  title: json['title'] as String,
  link: json['link'] as String,
  productLink: json['product_link'] as String,
  productId: json['product_id'] as String,
  serpapiProductApi: json['serpapi_product_api'] as String,
  serpapiProductApiComparative:
      json['serpapi_product_api_comparative'] as String?,
  source: json['source'] as String,
  price: json['price'] as String,
  extractedPrice: (json['extracted_price'] as num).toDouble(),
  rating: (json['rating'] as num?)?.toDouble(),
  reviews: (json['reviews'] as num?)?.toInt(),
  thumbnail: json['thumbnail'] as String,
  delivery: json['delivery'] as String?,
  tag: json['tag'] as String?,
  extensions: (json['extensions'] as List<dynamic>?)
      ?.map((e) => e as String)
      .toList(),
  currency: json['currency'] as String?,
  pageToken: json['page_token'] as String?,
  relevanceScore: (json['relevance_score'] as num?)?.toDouble(),
);

Map<String, dynamic> _$ProductToJson(_Product instance) => <String, dynamic>{
  'position': instance.position,
  'title': instance.title,
  'link': instance.link,
  'product_link': instance.productLink,
  'product_id': instance.productId,
  'serpapi_product_api': instance.serpapiProductApi,
  'serpapi_product_api_comparative': instance.serpapiProductApiComparative,
  'source': instance.source,
  'price': instance.price,
  'extracted_price': instance.extractedPrice,
  'rating': instance.rating,
  'reviews': instance.reviews,
  'thumbnail': instance.thumbnail,
  'delivery': instance.delivery,
  'tag': instance.tag,
  'extensions': instance.extensions,
  'currency': instance.currency,
  'page_token': instance.pageToken,
  'relevance_score': instance.relevanceScore,
};
