import 'package:freezed_annotation/freezed_annotation.dart';

part 'product.freezed.dart';
part 'product.g.dart';

@freezed
class Product with _$Product {
  const factory Product({
    required int position,
    required String title,
    required String link,
    @JsonKey(name: 'product_link') required String productLink,
    @JsonKey(name: 'product_id') required String productId,
    @JsonKey(name: 'serpapi_product_api') required String serpapiProductApi,
    required String source, required String price, @JsonKey(name: 'extracted_price') required double extractedPrice, required String thumbnail, @JsonKey(name: 'serpapi_product_api_comparative') String? serpapiProductApiComparative,
    double? rating,
    int? reviews,
    String? delivery,
    String? tag,
    List<String>? extensions,
    String? currency,
    @JsonKey(name: 'page_token') String? pageToken,
    @JsonKey(name: 'relevance_score') double? relevanceScore,
  }) = _Product;

  factory Product.fromJson(Map<String, dynamic> json) => _$ProductFromJson(json);
}
