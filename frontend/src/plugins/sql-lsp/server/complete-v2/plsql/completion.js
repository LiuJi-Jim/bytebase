import antlr4 from "antlr4";
// import { PlSqlLexer, PlSqlParserListener } from "./parser";
import PlSqlLexer from "./parser/PlSqlLexer";
import PlSqlParserListener from "./parser/PlSqlParserListener";

const possibleIdentifierPrefix = /[\w]$/;
const lineSeparator = /\n|\r|\r\n/g;

export function findCursorTokenIndex(tokenStream, cursor) {
  // NOTE: cursor positions are 0-based
  const cursorCol = cursor.column;
  console.log("cursorCol", cursorCol);
  // console.log('tokenStream.tokens.length', tokenStream.tokens.length)
  for (let i = 0; i < tokenStream.tokens.length; i++) {
    const t = tokenStream.get(i);

    // console.log('t', t);

    const tokenStartCol = t.column;
    const tokenEndCol = tokenStartCol + t.text.length;
    const tokenStartLine = t.line;
    const tokenEndLine =
      t.type !== PlSqlLexer.SPACES || !t.text
        ? tokenStartLine
        : tokenStartLine + (t.text.match(lineSeparator)?.length || 0);

    // console.log('tokenStartCol', tokenStartCol);
    // console.log('tokenEndCol', tokenEndCol);
    // console.log('tokenStartLine', tokenStartLine);
    // console.log('tokenEndLine', tokenEndLine);

    // NOTE: tokenEndCol makes sense only of tokenStartLine === tokenEndLine
    if (
      tokenEndLine > cursor.line ||
      (tokenStartLine === cursor.line && tokenEndCol > cursorCol)
    ) {
      if (
        i > 0 &&
        tokenStartLine === cursor.line &&
        tokenStartCol === cursorCol &&
        possibleIdentifierPrefix.test(tokenStream.get(i - 1).text)
      ) {
        return i - 1;
      } else if (tokenStream.get(i).type === PlSqlLexer.SPACES) {
        return i + 1;
      } else return i;
    }
  }
  return undefined;
}

class FromClauseListener extends PlSqlParserListener {
  constructor(tokenIndex) {
    super();
    this.fromClause = [];
    this.suggestion = [];
    this.tokenIndex = tokenIndex;
  }

  enterSelect_list_elements(ctx) {
    // console.log("enterSelect_list_elements", ctx);
    // console.log(ctx.getText());
    // console.log(ctx.start, ctx.stop);
    if (
      ctx.start.tokenIndex <= this.tokenIndex &&
      this.tokenIndex <= ctx.stop.tokenIndex
    ) {
      // const text = ctx.parser._input.getText({
      //   start: ctx.start,
      //   stop: ctx.stop,
      // });
      // console.log("text", text);
      const tokens = [];
      let periods = 0;

      for (let i = ctx.start.tokenIndex; i <= ctx.stop.tokenIndex; i++) {
        const token = ctx.parser._input.get(i);
        // console.log("token", token);
        if (token.channel !== antlr4.Lexer.DEFAULT_TOKEN_CHANNEL) {
          continue;
        }
        if (token.type === PlSqlLexer.PERIOD) {
          periods = periods + 1;
          continue;
        }
        if (token.type === PlSqlLexer.DELIMITED_ID) {
          tokens.push(token.text.slice(1, -1));
        } else {
          tokens.push(token.text.toUpperCase());
        }
      }

      while (periods >= tokens.length) {
        tokens.push("");
      }
      // console.log("tokens", tokens);
      if (tokens.length === 1) {
        this.suggestion.push({
          field: tokens[0],
        });
      } else if (tokens.length === 2) {
        this.suggestion.push({
          field: tokens[1],
          table: tokens[0],
        });
      } else if (tokens.length === 3) {
        this.suggestion.push({
          field: tokens[2],
          table: tokens[1],
          schema: tokens[0],
        });
      }
    }
  }

  enterFrom_clause(ctx) {
    // TODO: handle sub-query and join.
    // console.log("enterFrom_clause", ctx);
    // console.log(ctx.table_ref_list());
    // console.log(ctx.table_ref_list().table_ref_list());
    // console.log(ctx.getText());
    for (let i = 0; i < ctx.table_ref_list().length; i++) {
      const tableRef = ctx.table_ref_list().table_ref(i);
      const tableItem = {
        alias: "",
        table: "",
        schema: "",
      };
      const tableViewName = tableRef
        .table_ref_aux()
        .table_ref_aux_internal()
        .dml_table_expression_clause()
        ?.tableview_name();
      if (tableViewName) {
        if (tableViewName.id_expression()) {
          tableItem.schema = normalizeIdentifier(tableViewName.identifier());
          tableItem.table = normalizeIDExpression(
            tableViewName.id_expression()
          );
        } else if (tableViewName.identifier()) {
          tableItem.table = normalizeIdentifier(tableViewName.identifier());
        }
      }
      if (tableRef.table_ref_aux().table_alias()) {
        tableItem.alias = normalizeTableAlias(
          tableRef.table_ref_aux().table_alias()
        );
      }
      if (
        tableItem.table !== "" ||
        tableItem.alias !== "" ||
        tableItem.schema !== ""
      ) {
        this.fromClause.push(tableItem);
      }
      // console.log("tableItem", tableItem);
    }
  }

  insteadOfAlias(ctx) {
    for (let i = 0; i < this.suggestion.length; i++) {
      for (let j = 0; j < this.fromClause.length; j++) {
        if (this.suggestion[i].table === this.fromClause[j].alias) {
          this.suggestion[i].table = this.fromClause[j].table;
          this.suggestion[i].schema = this.fromClause[j].schema;
          break;
        }
      }
    }
  }
}

export function normalizeIDExpression(idExpression) {
  if (!idExpression) return undefined;
  if (idExpression.regular_id()) {
    return idExpression.regular_id().getText().toUpperCase();
  }
  if (idExpression.DELIMITED_ID()) {
    // trim the double quotes
    return idExpression.DELIMITED_ID().getText().slice(1, -1);
  }
}

export function normalizeIdentifier(identifier) {
  return normalizeIDExpression(identifier.id_expression());
}

export function normalizeTableAlias(tableAlias) {
  if (!tableAlias) return undefined;
  if (tableAlias.identifier()) {
    return normalizeIdentifier(tableAlias.identifier());
  }
  if (tableAlias.quoted_string()) {
    return tableAlias.quoted_string().getText().slice(1, -1);
  }
}

export function completion(tree, completionTokenIndex) {
  const listener = new FromClauseListener(completionTokenIndex);
  // console.log("start walk", tree, listener);
  antlr4.tree.ParseTreeWalker.DEFAULT.walk(listener, tree);
  listener.insteadOfAlias();
  // console.log("listener", listener);
  return { suggestion: listener.suggestion, fromClause: listener.fromClause };
}
