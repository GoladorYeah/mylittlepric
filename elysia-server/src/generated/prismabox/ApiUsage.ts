import { Type } from "@sinclair/typebox";

import { __transformDate__ } from "./__transformDate__";

import { __nullable__ } from "./__nullable__";

export const ApiUsagePlain = Type.Object(
  {
    id: Type.String(),
    apiName: Type.String(),
    keyIndex: Type.Integer(),
    requestType: __nullable__(Type.String()),
    responseTimeMs: __nullable__(Type.Integer()),
    success: Type.Boolean(),
    errorMessage: __nullable__(Type.String()),
    createdAt: Type.Date(),
  },
  { additionalProperties: false },
);

export const ApiUsageRelations = Type.Object(
  {},
  { additionalProperties: false },
);

export const ApiUsageWhere = Type.Partial(
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
          apiName: Type.String(),
          keyIndex: Type.Integer(),
          requestType: Type.String(),
          responseTimeMs: Type.Integer(),
          success: Type.Boolean(),
          errorMessage: Type.String(),
          createdAt: Type.Date(),
        },
        { additionalProperties: false },
      ),
    { $id: "ApiUsage" },
  ),
);

export const ApiUsageWhereUnique = Type.Recursive(
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
              apiName: Type.String(),
              keyIndex: Type.Integer(),
              requestType: Type.String(),
              responseTimeMs: Type.Integer(),
              success: Type.Boolean(),
              errorMessage: Type.String(),
              createdAt: Type.Date(),
            },
            { additionalProperties: false },
          ),
        ),
      ],
      { additionalProperties: false },
    ),
  { $id: "ApiUsage" },
);

export const ApiUsageSelect = Type.Partial(
  Type.Object(
    {
      id: Type.Boolean(),
      apiName: Type.Boolean(),
      keyIndex: Type.Boolean(),
      requestType: Type.Boolean(),
      responseTimeMs: Type.Boolean(),
      success: Type.Boolean(),
      errorMessage: Type.Boolean(),
      createdAt: Type.Boolean(),
      _count: Type.Boolean(),
    },
    { additionalProperties: false },
  ),
);

export const ApiUsageInclude = Type.Partial(
  Type.Object({ _count: Type.Boolean() }, { additionalProperties: false }),
);

export const ApiUsageOrderBy = Type.Partial(
  Type.Object(
    {
      id: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      apiName: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      keyIndex: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      requestType: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      responseTimeMs: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      success: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      errorMessage: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      createdAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
    },
    { additionalProperties: false },
  ),
);

export const ApiUsage = Type.Composite([ApiUsagePlain, ApiUsageRelations], {
  additionalProperties: false,
});
