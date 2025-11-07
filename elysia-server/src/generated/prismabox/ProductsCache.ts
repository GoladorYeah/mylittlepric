import { Type } from "@sinclair/typebox";

import { __transformDate__ } from "./__transformDate__";

import { __nullable__ } from "./__nullable__";

export const ProductsCachePlain = Type.Object(
  {
    id: Type.String(),
    pageToken: Type.String(),
    productName: Type.String(),
    countryCode: Type.String(),
    priceInfo: Type.Any(),
    merchantInfo: __nullable__(Type.Any()),
    productDetails: __nullable__(Type.Any()),
    imageUrl: __nullable__(Type.String()),
    link: __nullable__(Type.String()),
    rating: __nullable__(Type.Number()),
    fetchCount: Type.Integer(),
    lastFetchedAt: Type.Date(),
    createdAt: Type.Date(),
  },
  { additionalProperties: false },
);

export const ProductsCacheRelations = Type.Object(
  {},
  { additionalProperties: false },
);

export const ProductsCacheWhere = Type.Partial(
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
          pageToken: Type.String(),
          productName: Type.String(),
          countryCode: Type.String(),
          priceInfo: Type.Any(),
          merchantInfo: Type.Any(),
          productDetails: Type.Any(),
          imageUrl: Type.String(),
          link: Type.String(),
          rating: Type.Number(),
          fetchCount: Type.Integer(),
          lastFetchedAt: Type.Date(),
          createdAt: Type.Date(),
        },
        { additionalProperties: false },
      ),
    { $id: "ProductsCache" },
  ),
);

export const ProductsCacheWhereUnique = Type.Recursive(
  (Self) =>
    Type.Intersect(
      [
        Type.Partial(
          Type.Object(
            { id: Type.String(), pageToken: Type.String() },
            { additionalProperties: false },
          ),
          { additionalProperties: false },
        ),
        Type.Union(
          [
            Type.Object({ id: Type.String() }),
            Type.Object({ pageToken: Type.String() }),
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
              pageToken: Type.String(),
              productName: Type.String(),
              countryCode: Type.String(),
              priceInfo: Type.Any(),
              merchantInfo: Type.Any(),
              productDetails: Type.Any(),
              imageUrl: Type.String(),
              link: Type.String(),
              rating: Type.Number(),
              fetchCount: Type.Integer(),
              lastFetchedAt: Type.Date(),
              createdAt: Type.Date(),
            },
            { additionalProperties: false },
          ),
        ),
      ],
      { additionalProperties: false },
    ),
  { $id: "ProductsCache" },
);

export const ProductsCacheSelect = Type.Partial(
  Type.Object(
    {
      id: Type.Boolean(),
      pageToken: Type.Boolean(),
      productName: Type.Boolean(),
      countryCode: Type.Boolean(),
      priceInfo: Type.Boolean(),
      merchantInfo: Type.Boolean(),
      productDetails: Type.Boolean(),
      imageUrl: Type.Boolean(),
      link: Type.Boolean(),
      rating: Type.Boolean(),
      fetchCount: Type.Boolean(),
      lastFetchedAt: Type.Boolean(),
      createdAt: Type.Boolean(),
      _count: Type.Boolean(),
    },
    { additionalProperties: false },
  ),
);

export const ProductsCacheInclude = Type.Partial(
  Type.Object({ _count: Type.Boolean() }, { additionalProperties: false }),
);

export const ProductsCacheOrderBy = Type.Partial(
  Type.Object(
    {
      id: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      pageToken: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      productName: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      countryCode: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      priceInfo: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      merchantInfo: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      productDetails: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      imageUrl: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      link: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      rating: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      fetchCount: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      lastFetchedAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
      createdAt: Type.Union([Type.Literal("asc"), Type.Literal("desc")], {
        additionalProperties: false,
      }),
    },
    { additionalProperties: false },
  ),
);

export const ProductsCache = Type.Composite(
  [ProductsCachePlain, ProductsCacheRelations],
  { additionalProperties: false },
);
