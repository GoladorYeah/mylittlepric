import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../../../shared/models/product.dart';
import '../../../shared/models/product_details.dart';
import '../../config/app_config.dart';
import '../dio_client.dart';

/// API service for product operations
class ProductApiService {
  final DioClient _client;

  ProductApiService(this._client);

  /// Get product details by page token
  Future<ProductDetails> getProductDetails({
    required String pageToken,
    String? country,
    String? language,
  }) async {
    final response = await _client.post(
      AppConfig.productDetailsEndpoint,
      data: {
        'page_token': pageToken,
        if (country != null) 'country': country,
        if (language != null) 'language': language,
      },
    );

    return ProductDetails.fromJson(response.data as Map<String, dynamic>);
  }

  /// Search products (if backend supports direct product search)
  Future<List<Product>> searchProducts({
    required String query,
    String? country,
    String? language,
    int? limit,
  }) async {
    final response = await _client.get(
      '/api/products/search',
      queryParameters: {
        'query': query,
        if (country != null) 'country': country,
        if (language != null) 'language': language,
        if (limit != null) 'limit': limit,
      },
    );

    final data = response.data as List<dynamic>;
    return data
        .map((json) => Product.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  /// Get product by ID (if backend supports)
  Future<Product> getProduct(String productId) async {
    final response = await _client.get('/api/products/$productId');
    return Product.fromJson(response.data as Map<String, dynamic>);
  }

  /// Get similar products
  Future<List<Product>> getSimilarProducts({
    required String productId,
    int? limit,
  }) async {
    final response = await _client.get(
      '/api/products/$productId/similar',
      queryParameters: {
        if (limit != null) 'limit': limit,
      },
    );

    final data = response.data as List<dynamic>;
    return data
        .map((json) => Product.fromJson(json as Map<String, dynamic>))
        .toList();
  }

  /// Track product view
  Future<void> trackProductView(String productId) async {
    await _client.post(
      '/api/products/$productId/view',
    );
  }

  /// Track product click
  Future<void> trackProductClick({
    required String productId,
    required String source,
  }) async {
    await _client.post(
      '/api/products/$productId/click',
      data: {
        'source': source,
      },
    );
  }
}

/// Provider for ProductApiService
final productApiServiceProvider = Provider<ProductApiService>((ref) {
  final client = ref.watch(dioClientProvider);
  return ProductApiService(client);
});
