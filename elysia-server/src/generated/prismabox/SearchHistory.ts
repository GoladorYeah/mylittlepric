import { Type } from "@sinclair/typebox";

import { __transformDate__ } from "./__transformDate__";

import { __nullable__ } from "./__nullable__";

export const SearchHistoryPlain = Type.Object(
  {
    id: Type.String(),
    userId: __nullable__(Type.String()),
    sessionId: __nullable__(Type.String()),
    searchQuery: Type.String(),
    optimizedQuery: __nullable__(Type.String()),
    searchType: Type.String(),
    category: __nullable__(Type.String()),
    countryCode: Type.String(),
    languageCode: Type.String(),
    currency: Type.String(),
    resultCount: Type.Integer(),
    productsFound: __nullable__(Type.Any()),
    clickedProductId: __nullable__(Type.String()),
    createdAt: Type.Date(),
    expiresAt: __nullable__(Type.Date()),
  },
  { additionalProperties: false },
);

export const SearchHistoryRelations = Type.Object(
  {
    user: __nullable__(
      Type.Object(
        {
          id: Type.String(),
          email: Type.String(),
          passwordHash: Type.String(),
          fullName: __nullable__(Type.String()),
          googleId: __nullable__(Type.String()),
          avatarUrl: __nullable__(Type.String()),
          createdAt: Type.Date(),
          updatedAt: Type.Date(),
          lastLoginAt: __nullable__(Type.Date()),
        },
        { additionalProperties: false },
      ),
    ),
    session: __nullable__(
      Type.Object(
        {
          id: Type.String(),
          sessionId: Type.String(),
          userId: __nullable__(Type.String()),
          countryCode: Type.String(),
          languageCode: Type.String(),
          currency: Type.String(),
          messageCount: Type.Integer(),
          searchState: __nullable__(Type.Any()),
          cycleState: __nullable__(Type.Any()),
          conversationContext: __nullable__(Type.Any()),
          createdAt: Type.Date(),
          updatedAt: Type.Date(),
          expiresAt: Type.Date(),
        },
        { additionalProperties: false },
      ),
    ),
  },
  { additionalProperties: false },
);

export const SearchHistoryWhere = Type.Partial(
  Type.Recursive(
    (Self) =>
      Type.Object(
        {
          AND: Type.Union([
            Self,
            Type.Array(Self, { additionalProperties: false }),
          ]),
          NOT: Type.Union([
            Self,
            Type.Array(Self, { additionalProperties: false }),
          ]),
          OR: Type.Array(Self, { additionalProperties: false }),
          id: Type.String(),
          userId: Type.String(),
          sessionId: Type.String(),
          searchQuery: Type.String(),
          optimizedQuery: Type.String(),
          searchType: Type.String(),
          category: Type.String(),
          countryCode: Type.String(),
          languageCode: Type.String(),
          currency: Type.String(),
          resultCount: Type.Integer(),
          productsFound: Type.Any(),
          clickedProductId: Type.String(),
          createdAt: Type.Date(),
          expiresAt: Type.Date(),
        },
        { additionalProperties: false },
      ),
    { $id: "SearchHistory" },
  ),
);

export const SearchHistoryWhereUnique = Type.Recursive(
  (Self) =>
    Type.Intersect(
      [
        Type.Partial(
          Type.Object({ id: Type.String() }, { additionalProperties: false }),
          { additionalProperties: false },
        ),
        Type.Union([Type.Object({ id: Type.String() })], {
          additionalProperties: false,
        }),
        Type.Partial(
          Type.Object({
            AND: Type.Union([
              Self,
              Type.Array(Self, { additionalProperties: false }),
            ]),
            NOT: Type.Union([
              Self,
              Type.Array(Self, { additionalProperties: false }),
            ]),
            OR: Type.Array(Self, { additionalProperties: false }),
          }),
          { additionalProperties: false },
        ),
        Type.Partial(
          Type.Object(
            {
              id: Type.String(),
              userId: Type.String(),
              sessionId: Type.String(),
              searchQuery: Type.String(),
              optimizedQuery: Type.String(),
              searchType: Type.String(),
              category: Type.String(),
              countryCode: Type.String(),
              languageCode: Type.String(),
              currency: Type.String(),
              resultCount: Type.Integer(),
              productsFound: Type.Any(),
              clickedProductId: Type.String(),
              createdAt: Type.Date(),
              expiresAt: Type.Date(),
            },
            { additionalProperties: false },
          ),
        ),
      ],
      { additionalProperties: false },
    ),
  { $id: "SearchHistory" },
);

export const SearchHistorySelect = Type.Partial(
  Type.Object(
    {
      id: Type.Boolean(),
      userId: Type.Boolean(),
      sessionId: Type.Boolean(),
      searchQuery: Type.Boolean(),
      optimizedQuery: Type.Boolean(),
      searchType: Type.Boolean(),
      category: Type.Boolean(),
      countryCode: Type.Boolean(),
      languageCode: Type.Boolean(),
      currency: Type.Boolean(),
      resultCount: Type.Boolean(),
      productsFound: Type.Boolean(),
      clickedProductId: Type.Boolean(),
      createdAt: Type.Boolean(),
      expiresAt: Type.Boolean(),
      user: Type.Boolean(),
      session: Type.Boolean(),
      _count: Type.Boolean(),
    },
    { additionalProperties: false },
  ),
);

export const SearchHistoryInclude = Type.Partial(
  Type.Object(
    { user: Type.Boolean(), session: Type.Boolean(), _count: Type.Boolean() },
    { additionalProperties: false },
  ),
);

export const SearchHistoryOrderBy = Type.Partial(
  Type.Object(
    {
      id: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      userId: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      sessionId: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      searchQuery: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      optimizedQuery: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      searchType: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      category: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      countryCode: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      languageCode: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      currency: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      resultCount: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      productsFound: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      clickedProductId: Type.Union(
        [Type.Literal("asc"), Type.Literal("desc")],
        { additionalProperties: false },
      ),
      createdAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      expiresAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
    },
    { additionalProperties: false },
  ),
);

export const SearchHistory = Type.Composite(
  [SearchHistoryPlain, SearchHistoryRelations],
  { additionalProperties: false },
);
