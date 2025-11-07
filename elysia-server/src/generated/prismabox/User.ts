import { Type } from "@sinclair/typebox";

import { __transformDate__ } from "./__transformDate__";

import { __nullable__ } from "./__nullable__";

export const UserPlain = Type.Object(
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
);

export const UserRelations = Type.Object(
  {
    refreshTokens: Type.Array(
      Type.Object(
        {
          id: Type.String(),
          userId: Type.String(),
          tokenHash: Type.String(),
          expiresAt: Type.Date(),
          createdAt: Type.Date(),
          revokedAt: __nullable__(Type.Date()),
        },
        { additionalProperties: false },
      ),
      { additionalProperties: false },
    ),
    chatSessions: Type.Array(
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
      { additionalProperties: false },
    ),
    searchHistory: Type.Array(
      Type.Object(
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
      ),
      { additionalProperties: false },
    ),
  },
  { additionalProperties: false },
);

export const UserWhere = Type.Partial(
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
          email: Type.String(),
          passwordHash: Type.String(),
          fullName: Type.String(),
          googleId: Type.String(),
          avatarUrl: Type.String(),
          createdAt: Type.Date(),
          updatedAt: Type.Date(),
          lastLoginAt: Type.Date(),
        },
        { additionalProperties: false },
      ),
    { $id: "User" },
  ),
);

export const UserWhereUnique = Type.Recursive(
  (Self) =>
    Type.Intersect(
      [
        Type.Partial(
          Type.Object(
            {
              id: Type.String(),
              email: Type.String(),
              googleId: Type.String(),
            },
            { additionalProperties: false },
          ),
          { additionalProperties: false },
        ),
        Type.Union(
          [
            Type.Object({ id: Type.String() }),
            Type.Object({ email: Type.String() }),
            Type.Object({ googleId: Type.String() }),
          ],
          { additionalProperties: false },
        ),
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
              email: Type.String(),
              passwordHash: Type.String(),
              fullName: Type.String(),
              googleId: Type.String(),
              avatarUrl: Type.String(),
              createdAt: Type.Date(),
              updatedAt: Type.Date(),
              lastLoginAt: Type.Date(),
            },
            { additionalProperties: false },
          ),
        ),
      ],
      { additionalProperties: false },
    ),
  { $id: "User" },
);

export const UserSelect = Type.Partial(
  Type.Object(
    {
      id: Type.Boolean(),
      email: Type.Boolean(),
      passwordHash: Type.Boolean(),
      fullName: Type.Boolean(),
      googleId: Type.Boolean(),
      avatarUrl: Type.Boolean(),
      createdAt: Type.Boolean(),
      updatedAt: Type.Boolean(),
      lastLoginAt: Type.Boolean(),
      refreshTokens: Type.Boolean(),
      chatSessions: Type.Boolean(),
      searchHistory: Type.Boolean(),
      _count: Type.Boolean(),
    },
    { additionalProperties: false },
  ),
);

export const UserInclude = Type.Partial(
  Type.Object(
    {
      refreshTokens: Type.Boolean(),
      chatSessions: Type.Boolean(),
      searchHistory: Type.Boolean(),
      _count: Type.Boolean(),
    },
    { additionalProperties: false },
  ),
);

export const UserOrderBy = Type.Partial(
  Type.Object(
    {
      id: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      email: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      passwordHash: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      fullName: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      googleId: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      avatarUrl: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      createdAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      updatedAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      lastLoginAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
    },
    { additionalProperties: false },
  ),
);

export const User = Type.Composite([UserPlain, UserRelations], {
  additionalProperties: false,
});
