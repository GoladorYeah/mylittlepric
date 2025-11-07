import { Type } from "@sinclair/typebox";

import { __transformDate__ } from "./__transformDate__";

import { __nullable__ } from "./__nullable__";

export const ChatSessionPlain = Type.Object(
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
);

export const ChatSessionRelations = Type.Object(
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
    messages: Type.Array(
      Type.Object(
        {
          id: Type.String(),
          sessionId: Type.String(),
          role: Type.String(),
          content: Type.String(),
          responseType: __nullable__(Type.String()),
          quickReplies: __nullable__(Type.Any()),
          metadata: __nullable__(Type.Any()),
          products: __nullable__(Type.Any()),
          searchInfo: __nullable__(Type.Any()),
          createdAt: Type.Date(),
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
    searchQueries: Type.Array(
      Type.Object(
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
      ),
      { additionalProperties: false },
    ),
  },
  { additionalProperties: false },
);

export const ChatSessionWhere = Type.Partial(
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
          userId: Type.String(),
          countryCode: Type.String(),
          languageCode: Type.String(),
          currency: Type.String(),
          messageCount: Type.Integer(),
          searchState: Type.Any(),
          cycleState: Type.Any(),
          conversationContext: Type.Any(),
          createdAt: Type.Date(),
          updatedAt: Type.Date(),
          expiresAt: Type.Date(),
        },
        { additionalProperties: false },
      ),
    { $id: "ChatSession" },
  ),
);

export const ChatSessionWhereUnique = Type.Recursive(
  (Self) =>
    Type.Intersect(
      [
        Type.Partial(
          Type.Object(
            { id: Type.String(), sessionId: Type.String() },
            { additionalProperties: false },
          ),
          { additionalProperties: false },
        ),
        Type.Union(
          [
            Type.Object({ id: Type.String() }),
            Type.Object({ sessionId: Type.String() }),
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
              sessionId: Type.String(),
              userId: Type.String(),
              countryCode: Type.String(),
              languageCode: Type.String(),
              currency: Type.String(),
              messageCount: Type.Integer(),
              searchState: Type.Any(),
              cycleState: Type.Any(),
              conversationContext: Type.Any(),
              createdAt: Type.Date(),
              updatedAt: Type.Date(),
              expiresAt: Type.Date(),
            },
            { additionalProperties: false },
          ),
        ),
      ],
      { additionalProperties: false },
    ),
  { $id: "ChatSession" },
);

export const ChatSessionSelect = Type.Partial(
  Type.Object(
    {
      id: Type.Boolean(),
      sessionId: Type.Boolean(),
      userId: Type.Boolean(),
      countryCode: Type.Boolean(),
      languageCode: Type.Boolean(),
      currency: Type.Boolean(),
      messageCount: Type.Boolean(),
      searchState: Type.Boolean(),
      cycleState: Type.Boolean(),
      conversationContext: Type.Boolean(),
      createdAt: Type.Boolean(),
      updatedAt: Type.Boolean(),
      expiresAt: Type.Boolean(),
      user: Type.Boolean(),
      messages: Type.Boolean(),
      searchHistory: Type.Boolean(),
      searchQueries: Type.Boolean(),
      _count: Type.Boolean(),
    },
    { additionalProperties: false },
  ),
);

export const ChatSessionInclude = Type.Partial(
  Type.Object(
    {
      user: Type.Boolean(),
      messages: Type.Boolean(),
      searchHistory: Type.Boolean(),
      searchQueries: Type.Boolean(),
      _count: Type.Boolean(),
    },
    { additionalProperties: false },
  ),
);

export const ChatSessionOrderBy = Type.Partial(
  Type.Object(
    {
      id: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      sessionId: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      userId: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
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
      messageCount: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      searchState: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      cycleState: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      conversationContext: Type.Union(
        [Type.Literal("asc"), Type.Literal("desc")],
        { additionalProperties: false },
      ),
      createdAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      updatedAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      expiresAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
    },
    { additionalProperties: false },
  ),
);

export const ChatSession = Type.Composite(
  [ChatSessionPlain, ChatSessionRelations],
  { additionalProperties: false },
);
