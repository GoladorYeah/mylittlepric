import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import '../../../shared/models/saved_search.dart';
import '../../../shared/widgets/widgets.dart';
import '../providers/session_provider.dart';
import '../providers/session_state.dart';
import '../widgets/widgets.dart';

/// History screen showing saved searches
class HistoryScreen extends ConsumerStatefulWidget {
  const HistoryScreen({super.key});

  @override
  ConsumerState<HistoryScreen> createState() => _HistoryScreenState();
}

class _HistoryScreenState extends ConsumerState<HistoryScreen> {
  final _searchController = TextEditingController();
  String? _selectedCategory;
  List<SavedSearch> _filteredSearches = [];

  @override
  void initState() {
    super.initState();
    _searchController.addListener(_filterSearches);
  }

  @override
  void dispose() {
    _searchController.dispose();
    super.dispose();
  }

  void _filterSearches() {
    final sessionNotifier = ref.read(sessionProvider.notifier);
    final searches = ref.read(sessionProvider).searches;

    List<SavedSearch> filtered = searches;

    // Filter by search query
    if (_searchController.text.isNotEmpty) {
      filtered = sessionNotifier.searchInHistory(_searchController.text);
    }

    // Filter by category
    if (_selectedCategory != null && _selectedCategory!.isNotEmpty) {
      filtered = filtered.where((s) => s.category == _selectedCategory).toList();
    }

    setState(() {
      _filteredSearches = filtered;
    });
  }

  Future<void> _deleteSearch(String searchId) async {
    // Show confirmation dialog
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Delete Search'),
        content: const Text('Are you sure you want to delete this search?'),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            style: FilledButton.styleFrom(
              backgroundColor: Theme.of(context).colorScheme.error,
            ),
            child: const Text('Delete'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      await ref.read(sessionProvider.notifier).deleteSearch(searchId);
      if (mounted) {
        SuccessSnackBar.show(context, 'Search deleted successfully');
        _filterSearches();
      }
    }
  }

  Future<void> _clearAllSearches() async {
    // Show confirmation dialog
    final confirmed = await showDialog<bool>(
      context: context,
      builder: (context) => AlertDialog(
        title: const Text('Clear All History'),
        content: const Text(
          'Are you sure you want to delete all search history? This action cannot be undone.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context, false),
            child: const Text('Cancel'),
          ),
          FilledButton(
            onPressed: () => Navigator.pop(context, true),
            style: FilledButton.styleFrom(
              backgroundColor: Theme.of(context).colorScheme.error,
            ),
            child: const Text('Clear All'),
          ),
        ],
      ),
    );

    if (confirmed == true && mounted) {
      await ref.read(sessionProvider.notifier).clearAllSearches();
      if (mounted) {
        SuccessSnackBar.show(context, 'All searches cleared');
        setState(() {
          _filteredSearches = [];
          _selectedCategory = null;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final sessionState = ref.watch(sessionProvider);
    final categories = ref.watch(uniqueCategoriesProvider);

    // Update filtered searches when data changes
    if (_filteredSearches.isEmpty && sessionState.searches.isNotEmpty) {
      _filteredSearches = sessionState.searches;
    }

    return Scaffold(
      appBar: AppBar(
        title: const Text('Search History'),
        actions: [
          if (sessionState.hasSearches)
            IconButton(
              icon: const Icon(Icons.delete_sweep),
              onPressed: _clearAllSearches,
              tooltip: 'Clear all history',
            ),
        ],
      ),
      body: Column(
        children: [
          // Search bar and filters
          Padding(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                // Search field
                TextField(
                  controller: _searchController,
                  decoration: InputDecoration(
                    hintText: 'Search in history...',
                    prefixIcon: const Icon(Icons.search),
                    suffixIcon: _searchController.text.isNotEmpty
                        ? IconButton(
                            icon: const Icon(Icons.clear),
                            onPressed: () {
                              _searchController.clear();
                              _filterSearches();
                            },
                          )
                        : null,
                    border: OutlineInputBorder(
                      borderRadius: BorderRadius.circular(12),
                    ),
                  ),
                ),
                const SizedBox(height: 12),
                // Category filter chips
                if (categories.isNotEmpty)
                  SizedBox(
                    height: 40,
                    child: ListView(
                      scrollDirection: Axis.horizontal,
                      children: [
                        // All categories chip
                        Padding(
                          padding: const EdgeInsets.only(right: 8),
                          child: FilterChip(
                            label: const Text('All'),
                            selected: _selectedCategory == null,
                            onSelected: (selected) {
                              setState(() {
                                _selectedCategory = null;
                              });
                              _filterSearches();
                            },
                          ),
                        ),
                        // Category chips
                        ...categories.map((category) {
                          return Padding(
                            padding: const EdgeInsets.only(right: 8),
                            child: FilterChip(
                              label: Text(category),
                              selected: _selectedCategory == category,
                              onSelected: (selected) {
                                setState(() {
                                  _selectedCategory = selected ? category : null;
                                });
                                _filterSearches();
                              },
                            ),
                          );
                        }).toList(),
                      ],
                    ),
                  ),
              ],
            ),
          ),
          // Content
          Expanded(
            child: _buildContent(sessionState),
          ),
        ],
      ),
      floatingActionButton: FloatingActionButton.extended(
        onPressed: () => context.go('/chat'),
        icon: const Icon(Icons.add),
        label: const Text('New Search'),
      ),
    );
  }

  Widget _buildContent(SessionState state) {
    // Loading state
    if (state.isLoading && state.searches.isEmpty) {
      return const ListShimmer(itemCount: 5);
    }

    // Error state
    if (state.error != null && state.searches.isEmpty) {
      return CustomErrorWidget(
        message: state.error!,
        onRetry: () => ref.read(sessionProvider.notifier).refreshSearches(),
      );
    }

    // Empty state
    if (_filteredSearches.isEmpty) {
      if (_searchController.text.isNotEmpty || _selectedCategory != null) {
        // No results for current filter
        return const EmptyStateWidget(
          title: 'No Results',
          message: 'No searches match your current filters.',
          icon: Icons.search_off,
        );
      } else {
        // No searches at all
        return EmptyStateWidget(
          title: 'No Search History',
          message: 'Your search history will appear here.',
          icon: Icons.history,
          action: FilledButton.icon(
            onPressed: () => context.go('/chat'),
            icon: const Icon(Icons.search),
            label: const Text('Start Searching'),
          ),
        );
      }
    }

    // List of searches
    return RefreshIndicator(
      onRefresh: () => ref.read(sessionProvider.notifier).refreshSearches(),
      child: ListView.builder(
        itemCount: _filteredSearches.length + (state.hasMore ? 1 : 0),
        itemBuilder: (context, index) {
          // Load more indicator
          if (index == _filteredSearches.length) {
            if (state.isLoading) {
              return const Padding(
                padding: EdgeInsets.all(16),
                child: Center(child: CircularProgressIndicator()),
              );
            } else {
              // Load more button
              return Padding(
                padding: const EdgeInsets.all(16),
                child: Center(
                  child: OutlinedButton(
                    onPressed: () {
                      ref.read(sessionProvider.notifier).loadMoreSearches();
                    },
                    child: const Text('Load More'),
                  ),
                ),
              );
            }
          }

          // Session card
          final search = _filteredSearches[index];
          return SessionCard(
            search: search,
            onDelete: () => _deleteSearch(search.id),
          );
        },
      ),
    );
  }
}
