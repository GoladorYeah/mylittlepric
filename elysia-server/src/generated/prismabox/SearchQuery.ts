import { Type } from "@sinclair/typebox";

import { __transformDate__ } from "./__transformDate__";

import { __nullable__ } from "./__nullable__";

export const SearchQueryPlain = Type.Object(
  {
    id: Type.String(),
    sessionId: __nullable__(Type.String()),
    originalQuery: Type.String(),
    optimizedQuery: Type.String(),
    searchType: Type.String(),
    countryCode: Type.String(),
    resultCount: Type.Integer(),
    createdAt: Type.Date(),
  },
  { additionalProperties: false },
);

export const SearchQueryRelations = Type.Object(
  {
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

export const SearchQueryWhere = Type.Partial(
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
          sessionId: Type.String(),
          originalQuery: Type.String(),
          optimizedQuery: Type.String(),
          searchType: Type.String(),
          countryCode: Type.String(),
          resultCount: Type.Integer(),
          createdAt: Type.Date(),
        },
        { additionalProperties: false },
      ),
    { $id: "SearchQuery" },
  ),
);

export const SearchQueryWhereUnique = Type.Recursive(
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
              sessionId: Type.String(),
              originalQuery: Type.String(),
              optimizedQuery: Type.String(),
              searchType: Type.String(),
              countryCode: Type.String(),
              resultCount: Type.Integer(),
              createdAt: Type.Date(),
            },
            { additionalProperties: false },
          ),
        ),
      ],
      { additionalProperties: false },
    ),
  { $id: "SearchQuery" },
);

export const SearchQuerySelect = Type.Partial(
  Type.Object(
    {
      id: Type.Boolean(),
      sessionId: Type.Boolean(),
      originalQuery: Type.Boolean(),
      optimizedQuery: Type.Boolean(),
      searchType: Type.Boolean(),
      countryCode: Type.Boolean(),
      resultCount: Type.Boolean(),
      createdAt: Type.Boolean(),
      session: Type.Boolean(),
      _count: Type.Boolean(),
    },
    { additionalProperties: false },
  ),
);

export const SearchQueryInclude = Type.Partial(
  Type.Object(
    { session: Type.Boolean(), _count: Type.Boolean() },
    { additionalProperties: false },
  ),
);

export const SearchQueryOrderBy = Type.Partial(
  Type.Object(
    {
      id: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      sessionId: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      originalQuery: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      optimizedQuery: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      searchType: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      countryCode: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      resultCount: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      createdAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
    },
    { additionalProperties: false },
  ),
);

export const SearchQuery = Type.Composite(
  [SearchQueryPlain, SearchQueryRelations],
  { additionalProperties: false },
);
