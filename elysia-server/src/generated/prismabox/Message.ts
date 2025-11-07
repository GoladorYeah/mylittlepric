import { Type } from "@sinclair/typebox";

import { __transformDate__ } from "./__transformDate__";

import { __nullable__ } from "./__nullable__";

export const MessagePlain = Type.Object(
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
);

export const MessageRelations = Type.Object(
  {
    session: Type.Object(
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
  },
  { additionalProperties: false },
);

export const MessageWhere = Type.Partial(
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
          role: Type.String(),
          content: Type.String(),
          responseType: Type.String(),
          quickReplies: Type.Any(),
          metadata: Type.Any(),
          products: Type.Any(),
          searchInfo: Type.Any(),
          createdAt: Type.Date(),
        },
        { additionalProperties: false },
      ),
    { $id: "Message" },
  ),
);

export const MessageWhereUnique = Type.Recursive(
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
              role: Type.String(),
              content: Type.String(),
              responseType: Type.String(),
              quickReplies: Type.Any(),
              metadata: Type.Any(),
              products: Type.Any(),
              searchInfo: Type.Any(),
              createdAt: Type.Date(),
            },
            { additionalProperties: false },
          ),
        ),
      ],
      { additionalProperties: false },
    ),
  { $id: "Message" },
);

export const MessageSelect = Type.Partial(
  Type.Object(
    {
      id: Type.Boolean(),
      sessionId: Type.Boolean(),
      role: Type.Boolean(),
      content: Type.Boolean(),
      responseType: Type.Boolean(),
      quickReplies: Type.Boolean(),
      metadata: Type.Boolean(),
      products: Type.Boolean(),
      searchInfo: Type.Boolean(),
      createdAt: Type.Boolean(),
      session: Type.Boolean(),
      _count: Type.Boolean(),
    },
    { additionalProperties: false },
  ),
);

export const MessageInclude = Type.Partial(
  Type.Object(
    { session: Type.Boolean(), _count: Type.Boolean() },
    { additionalProperties: false },
  ),
);

export const MessageOrderBy = Type.Partial(
  Type.Object(
    {
      id: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      sessionId: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      role: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      content: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      responseType: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      quickReplies: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      metadata: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      products: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      searchInfo: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      createdAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
    },
    { additionalProperties: false },
  ),
);

export const Message = Type.Composite([MessagePlain, MessageRelations], {
  additionalProperties: false,
});
