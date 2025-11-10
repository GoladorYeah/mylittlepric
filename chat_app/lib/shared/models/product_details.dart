import 'package:freezed_annotation/freezed_annotation.dart';

part 'product_details.freezed.dart';
part 'product_details.g.dart';

@freezed
class ProductOffer with _$ProductOffer {
  const factory ProductOffer({
    required String merchant,
    required String price, required String link, String? logo,
    @JsonKey(name: 'extracted_price') double? extractedPrice,
    String? currency,
    String? title,
    String? availability,
    String? shipping,
    @JsonKey(name: 'shipping_extracted') double? shippingExtracted,
    String? total,
    @JsonKey(name: 'extracted_total') double? extractedTotal,
    double? rating,
    int? reviews,
    @JsonKey(name: 'payment_methods') String? paymentMethods,
    String? tag,
    @JsonKey(name: 'details_and_offers') List<String>? detailsAndOffers,
    @JsonKey(name: 'monthly_payment_duration') int? monthlyPaymentDuration,
    @JsonKey(name: 'down_payment') String? downPayment,
  }) = _ProductOffer;

  factory ProductOffer.fromJson(Map<String, dynamic> json) => _$ProductOfferFromJson(json);
}

@freezed
class Specification with _$Specification {
  const factory Specification({
    required String title,
    required String value,
  }) = _Specification;

  factory Specification.fromJson(Map<String, dynamic> json) => _$SpecificationFromJson(json);
}

@freezed
class RatingBreakdown with _$RatingBreakdown {
  const factory RatingBreakdown({
    required int stars,
    required int amount,
  }) = _RatingBreakdown;

  factory RatingBreakdown.fromJson(Map<String, dynamic> json) => _$RatingBreakdownFromJson(json);
}

@freezed
class ProductDetailsResponse with _$ProductDetailsResponse {
  const factory ProductDetailsResponse({
    required String type,
    required String title,
    required String price,
    // variants would need more complex modeling
    required List<ProductOffer> offers, double? rating,
    int? reviews,
    String? description,
    List<String>? images,
    List<Specification>? specifications,
    // videos and more_options could be added if needed
    @JsonKey(name: 'rating_breakdown') List<RatingBreakdown>? ratingBreakdown,
  }) = _ProductDetailsResponse;

  factory ProductDetailsResponse.fromJson(Map<String, dynamic> json) =>
      _$ProductDetailsResponseFromJson(json);
}
