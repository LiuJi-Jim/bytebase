import DefaultErrorStrategy from "../../3rd-party/antlr4/error/DefaultErrorStrategy";
import IntervalSet from "../../3rd-party/antlr4/misc/IntervalSet";
// import { PlSqlLexer, PlSqlParser } from "./parser";
import PlSqlLexer from "./parser/PlSqlLexer";
import PlSqlParser from "./parser/PlSqlParser";

export class PlSqlErrorStrategy extends DefaultErrorStrategy {
  singleTokenDeletion() {
    return null;
  }

  getErrorRecoverySet(recognizer) {
    const defaultRecoverySet = super.getErrorRecoverySet(recognizer);
    const currentRuleIndex = recognizer._ctx.ruleIndex;
    if (currentRuleIndex === PlSqlParser.RULE_query_block) {
      const plsqlFieldFollowSet = new IntervalSet();
      plsqlFieldFollowSet.addOne(PlSqlLexer.COMMA);
      plsqlFieldFollowSet.addOne(PlSqlLexer.FROM);
      const intersection = getIntersection(
        defaultRecoverySet,
        plsqlFieldFollowSet
      );
      if (intersection.length > 0) return intersection;
    }
    return defaultRecoverySet;
  }
}

function getIntersection(set1, set2) {
  const intersection = new IntervalSet();
  if (set1 === null || set2 === null) return intersection;
  if (set1.intervals === null || set2.intervals === null) return intersection;
  for (const interval of set1.intervals) {
    if (interval === null) continue;
    for (let i = interval.start; i <= interval.stop; i++) {
      if (set2.contains(i)) {
        intersection.addOne(i);
      }
    }
  }
  return intersection;
}
