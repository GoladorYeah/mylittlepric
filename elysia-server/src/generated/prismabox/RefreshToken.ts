import { Type } from "@sinclair/typebox";

import { __transformDate__ } from "./__transformDate__";

import { __nullable__ } from "./__nullable__";

export const RefreshTokenPlain = Type.Object(
  {
    id: Type.String(),
    userId: Type.String(),
    tokenHash: Type.String(),
    expiresAt: Type.Date(),
    createdAt: Type.Date(),
    revokedAt: __nullable__(Type.Date()),
  },
  { additionalProperties: false },
);

export const RefreshTokenRelations = Type.Object(
  {
    user: Type.Object(
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
  },
  { additionalProperties: false },
);

export const RefreshTokenWhere = Type.Partial(
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
          tokenHash: Type.String(),
          expiresAt: Type.Date(),
          createdAt: Type.Date(),
          revokedAt: Type.Date(),
        },
        { additionalProperties: false },
      ),
    { $id: "RefreshToken" },
  ),
);

export const RefreshTokenWhereUnique = Type.Recursive(
  (Self) =>
    Type.Intersect(
      [
        Type.Partial(
          Type.Object(
            { id: Type.String(), tokenHash: Type.String() },
            { additionalProperties: false },
          ),
          { additionalProperties: false },
        ),
        Type.Union(
          [
            Type.Object({ id: Type.String() }),
            Type.Object({ tokenHash: Type.String() }),
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
              userId: Type.String(),
              tokenHash: Type.String(),
              expiresAt: Type.Date(),
              createdAt: Type.Date(),
              revokedAt: Type.Date(),
            },
            { additionalProperties: false },
          ),
        ),
      ],
      { additionalProperties: false },
    ),
  { $id: "RefreshToken" },
);

export const RefreshTokenSelect = Type.Partial(
  Type.Object(
    {
      id: Type.Boolean(),
      userId: Type.Boolean(),
      tokenHash: Type.Boolean(),
      expiresAt: Type.Boolean(),
      createdAt: Type.Boolean(),
      revokedAt: Type.Boolean(),
      user: Type.Boolean(),
      _count: Type.Boolean(),
    },
    { additionalProperties: false },
  ),
);

export const RefreshTokenInclude = Type.Partial(
  Type.Object(
    { user: Type.Boolean(), _count: Type.Boolean() },
    { additionalProperties: false },
  ),
);

export const RefreshTokenOrderBy = Type.Partial(
  Type.Object(
    {
      id: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      userId: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      tokenHash: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      expiresAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      createdAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      revokedAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
    },
    { additionalProperties: false },
  ),
);

export const RefreshToken = Type.Composite(
  [RefreshTokenPlain, RefreshTokenRelations],
  { additionalProperties: false },
);
