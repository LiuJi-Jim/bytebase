import { Expr as CELExpr } from "@/types/proto/google/api/expr/v1alpha1/syntax";
import type {
  Factor,
  NumberFactor,
  StringFactor,
  Operator,
  CollectionOperator,
  CompareOperator,
  StringOperator,
  EqualityOperator,
  SimpleExpr,
  CollectionExpr,
  ConditionExpr,
  ConditionGroupExpr,
  CompareExpr,
  EqualityExpr,
  StringExpr,
} from "../types";
import {
  LogicalOperatorList,
  isEqualityOperator,
  isConditionGroupExpr,
  isCollectionOperator,
  isLogicalOperator,
  isCompareOperator,
  isStringOperator,
  isNumberFactor,
  isStringFactor,
} from "../types";

// For simplify UI implementation, the "root" condition need to be a group.
export const wrapAsGroup = (expr: SimpleExpr): ConditionGroupExpr => {
  if (isConditionGroupExpr(expr)) return expr;
  return {
    operator: "_&&_",
    args: [expr],
  };
};

// Convert common expr to simple expr
export const resolveCELExpr = (expr: CELExpr): SimpleExpr => {
  const dfs = (expr: CELExpr): ConditionGroupExpr | ConditionExpr => {
    const { callExpr } = expr;
    if (!callExpr) {
      return emptySimpleExpr();
      // throw new Error(`unsupported expr "${JSON.stringify(expr)}"`);
    }

    const { args } = callExpr;
    const operator = callExpr.function as Operator;
    if (isLogicalOperator(operator)) {
      const group: ConditionGroupExpr = {
        operator,
        args: args.map(dfs),
      };
      return group;
    }
    if (isEqualityOperator(operator)) {
      return resolveEqualityExpr(expr);
    }
    if (isCompareOperator(operator)) {
      return resolveCompareExpr(expr);
    }
    if (isStringOperator(operator)) {
      return resolveStringExpr(expr);
    }
    if (isCollectionOperator(operator)) {
      return resolveCollectionExpr(expr);
    }
    throw new Error(`unsupported expr "${JSON.stringify(expr)}"`);
  };
  return dfs(expr);
};

const resolveEqualityExpr = (expr: CELExpr): EqualityExpr => {
  const operator = expr.callExpr!.function as EqualityOperator;
  const [factorExpr, valueExpr] = expr.callExpr!.args;
  const factor = factorExpr.identExpr!.name;
  if (isNumberFactor(factor)) {
    return {
      operator,
      args: [
        combineResourceNameAndResourceId(factor),
        valueExpr.constExpr!.int64Value! ?? 0,
      ],
    };
  }
  if (isStringFactor(factor)) {
    return {
      operator,
      args: [
        combineResourceNameAndResourceId(factor),
        valueExpr.constExpr!.stringValue! ?? "",
      ],
    };
  }
  throw new Error(`cannot resolve expr ${JSON.stringify(expr)}`);
};

const resolveCompareExpr = (expr: CELExpr): CompareExpr => {
  const operator = expr.callExpr!.function as CompareOperator;
  const [factor, value] = expr.callExpr!.args;
  return {
    operator,
    args: [
      factor.identExpr!.name as NumberFactor,
      value.constExpr!.int64Value!,
    ],
  };
};

const resolveStringExpr = (expr: CELExpr): StringExpr => {
  const operator = expr.callExpr!.function as StringOperator;
  const factor = expr.callExpr!.target!;
  const value = expr.callExpr!.args[0];
  return {
    operator,
    args: [
      factor.identExpr!.name as StringFactor,
      value.constExpr!.stringValue!,
    ],
  };
};

const resolveCollectionExpr = (expr: CELExpr): CollectionExpr => {
  const operator = expr.callExpr!.function as CollectionOperator;
  const [factorExpr, valuesExpr] = expr.callExpr!.args;
  const factor = factorExpr.identExpr!.name;
  if (isNumberFactor(factor)) {
    return {
      operator,
      args: [
        combineResourceNameAndResourceId(factor),
        valuesExpr.listExpr?.elements?.map(
          (constant) => constant.constExpr?.int64Value ?? 0
        ) ?? [],
      ],
    };
  }
  if (isStringFactor(factor)) {
    return {
      operator,
      args: [
        combineResourceNameAndResourceId(factor),
        valuesExpr.listExpr?.elements?.map(
          (constant) => constant.constExpr?.stringValue ?? ""
        ) ?? [],
      ],
    };
  }

  throw new Error(`cannot resolve expr ${JSON.stringify(expr)}`);
};

const emptySimpleExpr = (): ConditionGroupExpr => {
  return {
    operator: LogicalOperatorList[0],
    args: [],
  };
};

const combineResourceNameAndResourceId = (factor: string): Factor => {
  if (factor === "project_id" || factor === "project_name") {
    return "project";
  }
  if (factor === "environment_id" || factor === "environment_name") {
    return "environment";
  }
  return factor as Factor;
};
