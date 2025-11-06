# Backend Services Audit Report

## Executive Summary

This audit examined 14 service files in `/backend/internal/services/`. The codebase shows good structural patterns with dependency injection and proper service separation, but has several areas for improvement in error handling, logging, performance, and code quality.

---

## 1. CODE STRUCTURE & PATTERNS

### Services Examined:
1. ✅ **auth_service.go** - Well-structured, good patterns
2. ✅ **cache.go** - Basic structure, could be improved
3. ✅ **context_extractor.go** - Good service design
4. ✅ **context_optimizer.go** - Stateless service, good patterns
5. ⚠️ **embedding.go** - Creates own Redis client
6. ⚠️ **gemini.go** - Mixed old and new approaches
7. ✅ **google_oauth.go** - Clean, focused service
8. ✅ **grounding_strategy.go** - Stateless, good patterns
9. ✅ **prompt_manager.go** - Simple, effective
10. ✅ **response_schemas.go** - Helper functions only
11. ✅ **search_history_service.go** - Well-structured
12. ⚠️ **serp.go** - Good logic but verbose
13. ⚠️ **session.go** - Dual storage (Redis + DB) adds complexity
14. ✅ **universal_prompt_manager.go** - Well-designed

### Issues Found:

#### 1.1 Inconsistent Dependency Injection Pattern
**File:** `cache.go`
**Line:** 23-45
**Issue:** CacheService creates its own Redis client instead of receiving injected instance
```go
func NewCacheService(cfg *config.Config, embedding *EmbeddingService) *CacheService {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisURL,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})
	// ...
}
```
**Impact:** Creates duplicate Redis connections, not following DI pattern
**Recommendation:** Accept injected Redis client like other services

#### 1.2 Context as Field vs Parameter Mismatch
**Files:** Multiple (auth_service.go:38, cache.go:34, embedding.go:32, etc.)
**Issue:** Services store `context.Context` as field initialized with `context.Background()`
**Lines Examples:**
- auth_service.go:38: `ctx: context.Background()`
- cache.go:34: `ctx: context.Background()`
- embedding.go:32: `ctx: context.Background()`
- gemini.go:55: `ctx: context.Background()`

**Problem:** Creating background contexts at initialization time is inefficient and unconventional. Context should be passed per-request or per-operation.
**Recommendation:** Use context parameters in method signatures instead of storing as field

#### 1.3 Duplicated Service Logic
**Files:** `gemini.go` and `universal_prompt_manager.go`
**Issue:** Both services manage prompts and context building
**Lines:** 
- gemini.go:107-250 (ProcessMessageWithContext - OLD)
- gemini.go:491-829 (ProcessWithUniversalPrompt - NEW)

This shows incomplete refactoring from old to new prompt system.

---

## 2. OUTDATED CODE

#### 2.1 Commented-out Code Blocks
**File:** `gemini.go`
**Lines:** 116-119 (rotateClient call)
```go
// УБРАЛИ rotateClient() отсюда!
// if err := g.rotateClient(); err != nil {
// 	return nil, -1, err
// }
```
**Lines:** 280-282 (Smart strategy code)
```go
// Smart strategy code preserved for future reference:
// decision := g.groundingStrategy.ShouldUseGrounding(userMessage, history, category)
// return decision.UseGrounding
```
**Issue:** Dead code cluttering the codebase
**Action:** Remove commented code and use git history if needed

#### 2.2 Old Method Still Active
**File:** `gemini.go`
**Method:** `ProcessMessageWithContext()` (lines 107-250)
**Issue:** This method appears to be replaced by `ProcessWithUniversalPrompt()` (lines 491-829) but both still exist
**Finding:** The new method includes sophisticated retry logic and fallback handling that the old one lacks
**Recommendation:** Complete migration and remove old method

#### 2.3 Unused Function Parameter
**File:** `serp.go`
**Method:** `getMaxProducts()` (lines 462-473)
**Issue:** Method exists but appears to be unused - `validateRelevance()` takes hardcoded `maxProducts := 10` instead of using this method
**Lines:** serp.go:177 hardcodes 10, method would allow dynamic values

---

## 3. ERROR HANDLING

### Critical Issues:

#### 3.1 Ignored Error Returns
**File:** `embedding.go`
**Lines:** 45, 91, 97, 98
```go
// Line 45-46
jsonData, _ := json.Marshal(e.categoryEmbeddings)  // IGNORING ERROR
e.redis.Set(e.ctx, key, jsonData, 0)

// Line 91-92  
json.Unmarshal(cached, &embedding)  // IGNORING ERROR

// Line 97-98
jsonData, _ := json.Marshal(embedding)  // IGNORING ERROR
```
**Impact:** Silent failures that could cause incorrect behavior
**Action:** Proper error handling and logging

#### 3.2 Silent Error Suppression
**File:** `cache.go`
**Lines:** 55-56
```go
if err := json.Unmarshal(data, &cards); err == nil {
	return cards, nil
	// ERROR IS SILENTLY IGNORED - continues to next option
}
```
**Issue:** JSON unmarshal error at line 55 is silently ignored, similar pattern at line 67

#### 3.3 Weak Error Messages
**File:** `context_extractor.go`
**Lines:** 89-92
```go
if err != nil {
	fmt.Printf("⚠️ Failed to extract preferences: %v\n", err)
	return currentPreferences, err  // Returns error but logs first
}
```
**Issue:** Inconsistent - sometimes errors are returned, sometimes only logged

#### 3.4 Missing Error Context
**File:** `session.go`
**Lines:** 319
```go
json.Unmarshal([]byte(msgData), &msg)
if err != nil {
	continue  // SILENTLY SKIPS BAD MESSAGE
}
```
**Issue:** Could hide data corruption issues

---

## 4. LOGGING

### Major Issues:

#### 4.1 Over-reliance on fmt.Printf()
**File:** All service files
**Issue:** Extensive use of fmt.Printf() instead of structured logging framework

**Examples:**
- gemini.go: Lines 59, 90, 116, 193, 328, 405, 577, 603, 628, 673, etc. (50+ printf calls)
- serp.go: Lines 57, 64, 91, 111, 113, 140, 146, 152, 187, etc. (30+ printf calls)
- session.go: Lines 85, 95, 97
- embedding.go: Lines 45, 46, 48

**Problem:** 
- No log levels (debug/info/warn/error)
- Can't filter by severity
- Timestamps not standardized
- Can't control verbosity
- Goes to stdout not proper logger

**Recommendation:** Implement structured logging (e.g., using slog, logrus, or zap)

#### 4.2 No Logging in Critical Paths
**File:** `cache.go`
**Issue:** Cache operations have no logging, making debugging difficult
- GetSearchResults() - line 47-72 (no logs on miss/hit)
- deduplicateProducts() - line 85-106 (no logging of dedup count)
- SetSearchResults() - line 74-83 (no logging of cache set)

#### 4.3 Debug-Level Info at Production Level
**File:** `google_oauth.go`
**Lines:** 75
```go
fmt.Printf("🔍 Google tokeninfo response: %s\n", string(bodyBytes))  // LOGS SENSITIVE DATA
```
**Security Issue:** Printing entire OAuth token response to stdout

---

## 5. PERFORMANCE ISSUES

### High Priority:

#### 5.1 O(n²) Deduplication Algorithm
**File:** `cache.go`
**Lines:** 85-106
**Function:** `deduplicateProducts()`

```go
func (c *CacheService) deduplicateProducts(cards []models.ProductCard) []models.ProductCard {
	unique := []models.ProductCard{cards[0]}
	
	for i := 1; i < len(cards); i++ {
		isDuplicate := false
		for j := range unique {
			// O(n) inner loop for each card
			if c.embedding.AreDuplicateProducts(cards[i].Name, unique[j].Name, 0.95) {
				isDuplicate = true
				break
			}
		}
		if !isDuplicate {
			unique = append(unique, cards[i])
		}
	}
	return unique
}
```

**Analysis:** 
- For 10 products: 10×10 = 100 comparisons
- Each comparison calls embedding similarity (API call to Gemini!)
- Total: O(n²) with expensive operations inside

**Impact:** Severe performance degradation on large result sets
**Recommendation:** Use hash-based dedup or batch similarity checks

#### 5.2 Inefficient Redis Scan in Cache Lookup
**File:** `embedding.go`
**Lines:** 130-154
**Function:** `FindSimilarCachedQuery()`

```go
func (e *EmbeddingService) FindSimilarCachedQuery(query string, threshold float32) string {
	queryEmbedding := e.GetQueryEmbedding(query)
	pattern := "cache:search:*"
	iter := e.redis.Scan(e.ctx, 0, pattern, 100).Iterator()
	
	for iter.Next(e.ctx) {
		cacheKey := iter.Val()
		cachedQuery := cacheKey[len("cache:search:"):]
		cachedEmbedding := e.GetQueryEmbedding(cachedQuery)
		// Compares every cached query
	}
}
```

**Analysis:**
- SCAN with pattern "cache:search:*" scans ALL Redis keys
- For each key, extracts and gets embedding (more API calls)
- No pagination or limit
- Could scan thousands of keys on large deployments

**Impact:** Slow response times, excessive API calls to Gemini
**Recommendation:** Use indexed Redis data structure or limit scan results

#### 5.3 Inefficient JSON Parsing in Error Recovery
**File:** `gemini.go`
**Lines:** 717-817
**Issue:** Multiple json.Unmarshal() calls during error recovery
- First parse to check if error (line 776)
- Then extract specific fields
- Could parse same JSON 5-10 times

**Recommendation:** Parse once, then work with parsed object

#### 5.4 Context Creation Per Retry
**File:** `gemini.go`
**Lines:** 420
```go
ctx, cancel := context.WithTimeout(g.ctx, 30*time.Second)
resp, err := client.Models.GenerateContent(ctx, ...)
cancel()
```

**Issue:** New context created for each retry attempt, not reused

---

## 6. SECURITY CONCERNS

### High Priority:

#### 6.1 Sensitive Data Logging
**File:** `google_oauth.go`
**Line:** 75
```go
fmt.Printf("🔍 Google tokeninfo response: %s\n", string(bodyBytes))
```
**Risk:** Entire OAuth token response printed to stdout/logs
**Contains:** ID tokens, user information
**Action:** Remove or use debug-only logging

#### 6.2 No Input Validation
**Files:** Multiple
**Examples:**
- `auth_service.go`: No validation of email format (line 45)
- `google_oauth.go`: No validation of idToken length/format (line 48)
- `serp.go`: No validation of query string (line 36)
- `session.go`: No validation of sessionID format (line 67)

**Impact:** Potential for injection attacks or invalid operations
**Recommendation:** Add input validation helpers

#### 6.3 Weak Token Validation
**File:** `auth_service.go`
**Line:** 87-88
```go
user, err := s.getUserByProviderID("google", googleUser.Sub)
if err != nil && !errors.Is(err, redis.Nil) {
	return nil, fmt.Errorf("failed to get user: %w", err)
}
```

**Issue:** Confusing error handling - checks for redis.Nil specifically
Better to check explicit error types

#### 6.4 Potential for Race Conditions
**File:** `session.go`
**Lines:** 67-101
**Issue:** GetSession() has race condition between Redis check and DB fallback
```go
// Redis might have session
data, err := s.redis.Get(s.ctx, key).Bytes()
// ...
// But between Redis check and DB restore, session might be modified
session, err := s.getSessionFromDB(sessionID)
```

---

## 7. DEPENDENCIES & IMPORTS

### Issues Found:

#### 7.1 Circular Dependency Risk
**Services:** GeminiService → EmbeddingService
**Files:** 
- gemini.go line 28: `embedding *EmbeddingService`
- cache.go line 19: `embedding *EmbeddingService`
- context_optimizer.go line 29: `embedding *EmbeddingService`

**Finding:** EmbeddingService is a common dependency, but no issues currently
**Status:** ✅ No actual circular dependency, but tightly coupled

#### 7.2 Unused Imports
**File:** `response_schemas.go`
**Line:** 4
```go
import "google.golang.org/genai"
```
**Issue:** Only uses genai.Schema, genai.TypeObject, etc. - all from genai package
**Finding:** Imports are necessary

#### 7.3 Missing Interface Definitions
**Services:** Most services
**Issue:** Services are concrete types, not behind interfaces
**Example:** GeminiService used directly, not IGeminiService interface
**Impact:** Harder to mock for testing, tight coupling

**Recommendation:** Consider interface-based design for testability

---

## 8. SPECIFIC FINDINGS BY FILE

### auth_service.go
✅ **Strengths:**
- Good error wrapping (lines 47, 56, 72, etc.)
- Proper use of private helper methods
- Thread-safe token storage

⚠️ **Issues:**
- Line 116, 144: Uses fmt.Printf for non-critical errors instead of logging
- Line 252-256: userExists() method doesn't validate email format
- No input validation on passwords (minimum length, complexity)

---

### cache.go
✅ **Strengths:**
- Clean interface
- Handles deduplication

❌ **Critical Issues:**
- Lines 23-28: Creates own Redis client (DI violation)
- Lines 85-106: O(n²) deduplication with expensive operations inside
- Line 55: Silent JSON unmarshal error
- No logging on cache operations

---

### context_extractor.go
✅ **Strengths:**
- Well-separated concerns
- Good extraction patterns

⚠️ **Issues:**
- Lines 89-92: Logs but still returns error (inconsistent pattern)
- Lines 205-206: Simple keyword-based extraction, not AI-powered despite having Gemini client
- Line 341-347: Duplicate contains() helper (should use slices package)

---

### context_optimizer.go
✅ **Strengths:**
- Stateless design
- Good pattern matching for context decisions
- Handles multiple languages

⚠️ **Issues:**
- Multiple hardcoded keyword lists (lines 84-107, 124-136, etc.) - should be in config
- Lines 48-79: Many similar methods doing pattern matching - could be consolidated

---

### embedding.go
❌ **Critical Issues:**
- Lines 24, 45, 46, 91, 97, 98: Multiple ignored errors
- Lines 130-154: SCAN pattern over all Redis keys - huge performance issue
- Line 48: json.Unmarshal error ignored silently
- No logging at all

⚠️ **Issues:**
- Line 35: Calls loadCategoryEmbeddings() in init - blocking operation
- Thread safety: categoryEmbeddings map accessed without locking (line 116)

---

### gemini.go
❌ **Critical Issues:**
- Lines 116-119, 280-282: Commented-out code blocks
- Lines 107-250: Old ProcessMessageWithContext() still exists alongside new method
- Lines 717-817: Multiple json.Unmarshal() calls in error recovery
- 50+ fmt.Printf() calls for logging

⚠️ **Issues:**
- Line 420: New context created per retry
- Line 270-278: Always enables grounding (comment says "smart strategy was too conservative")
- Lines 932-1045: Complex JSON repair logic - hard to maintain

---

### google_oauth.go
❌ **Critical Issues:**
- Line 75: Logs entire OAuth token response (security risk)
- Line 82: No validation of Client ID format

✅ **Strengths:**
- Clean service design
- Proper error handling

---

### grounding_strategy.go
✅ **Strengths:**
- Well-designed decision-making logic
- Good use of embedding-based similarity
- No obvious issues

⚠️ **Minor:**
- Lines 79-81: Multiple API calls to get embedding for same concept
- Could cache "fresh_info", "product_concept", etc. embeddings

---

### prompt_manager.go
✅ **Strengths:**
- Simple, effective design
- Proper mutex usage
- Good error handling

⚠️ **Minor Issues:**
- Line 37: Error logged with fmt.Printf instead of proper logger
- Hard-coded file paths (line 25-28) should be in config

---

### response_schemas.go
✅ **Strengths:**
- Clean helper functions
- Well-structured schemas
- Good documentation

---

### search_history_service.go
✅ **Strengths:**
- Proper database operations
- Good error handling
- Parameterized queries (prevents SQL injection)

⚠️ **Minor Issues:**
- No logging on CRUD operations
- Line 319: Silently skips bad JSON messages without logging

---

### serp.go
❌ **Issues:**
- Lines 55-67: Repeated logging for debugging (should be cleaner)
- Lines 164-213: Overly simplistic relevance logic (just takes top 10)
- Lines 462-473: getMaxProducts() method unused
- 30+ fmt.Printf() calls
- Line 176-212: Comment says "НОВАЯ ЛОГИКА" (NEW LOGIC) but doesn't match comment at line 214

⚠️ **Issues:**
- Lines 281-288: Brand list hardcoded (should be configurable)
- Performance: Multiple string operations in tight loops

---

### session.go
✅ **Strengths:**
- Dual storage (Redis + DB) provides resilience
- Good separation of concerns
- Proper error handling

❌ **Critical Issues:**
- Line 319: json.Unmarshal error silently ignored in GetMessages()
- Lines 67-101: Potential race condition between Redis and DB

⚠️ **Issues:**
- Line 291: Context format string has inconsistent formatting
- No logging on session operations

---

### universal_prompt_manager.go
✅ **Strengths:**
- Well-designed context management
- Good separation of full/compact/minimal contexts
- Proper mutex usage

⚠️ **Issues:**
- Lines 47, 55: panic() on file load - better to return error
- Lines 292-310: extractGroups()/extractSubgroups() return empty lists (not implemented)
- No logging

---

## SUMMARY TABLE

| Category | Files | Count | Severity |
|----------|-------|-------|----------|
| Error Handling | embedding.go, cache.go, session.go | 6 | HIGH |
| Logging (missing/excessive) | All files | 50+ | MEDIUM |
| Performance | cache.go, embedding.go, gemini.go | 3 | HIGH |
| Security | google_oauth.go, auth_service.go | 2 | MEDIUM |
| Code Quality | gemini.go, serp.go | 2 | MEDIUM |
| Dependencies | cache.go | 1 | MEDIUM |
| Unused Code | gemini.go, serp.go | 3 | LOW |

---

## RECOMMENDATIONS (Priority Order)

### 🔴 CRITICAL (Fix Immediately)
1. **embedding.go:** Fix error handling (lines 45, 91, 97, 98) - Silent failures could cause data issues
2. **embedding.go:** Replace SCAN pattern with indexed lookup (lines 130-154) - Performance issue
3. **cache.go:** Fix O(n²) deduplication (lines 85-106) - Can cause timeouts
4. **session.go:** Fix silent JSON unmarshal error (line 319) - Data loss risk
5. **google_oauth.go:** Remove sensitive token logging (line 75) - Security risk

### 🟡 HIGH (Fix Soon)
1. **All files:** Replace fmt.Printf with structured logging framework
2. **cache.go:** Inject Redis client instead of creating new one
3. **gemini.go:** Remove commented-out code (lines 116-119, 280-282)
4. **gemini.go:** Complete migration from old ProcessMessageWithContext()
5. **Add input validation** to all service methods

### 🟢 MEDIUM (Improve Quality)
1. Remove unused getMaxProducts() in serp.go
2. Extract hardcoded keyword lists to configuration
3. Create service interfaces for better testability
4. Move context from service field to method parameters
5. Implement proper database transaction handling
6. Add comprehensive error context in error messages

### 💡 NICE TO HAVE
1. Add metrics/tracing for performance monitoring
2. Implement caching for embedding similarity checks
3. Add rate limiting to API calls
4. Create helper functions for common patterns

