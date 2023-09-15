import antlr4 from "antlr4";
import { completion, findCursorTokenIndex } from "./completion";
// import { PlSqlLexer, PlSqlParser } from "./parser";
import PlSqlLexer from "./parser/PlSqlLexer";
import PlSqlParser from "./parser/PlSqlParser";

export const parse = (input) => {
  const inputCharStream = antlr4.CharStreams.fromString(input);
  const lexer = new PlSqlLexer(inputCharStream);
  const tokenStream = new antlr4.CommonTokenStream(lexer);
  const parser = new PlSqlParser(tokenStream);
  try {
    const tree = parser.sql_script();
    return tree;
  } catch (err) {
    console.error("err", err);
    return undefined;
  }
};

export const complete = (inputString, caretIndex) => {
  console.log("input text:", inputString);
  console.log("caret index:", caretIndex);
  console.log("parsing...");

  console.time("parse");
  let tree = null;
  const inputCharStream = antlr4.CharStreams.fromString(inputString);
  const lexer = new PlSqlLexer(inputCharStream);
  const tokenStream = new antlr4.CommonTokenStream(lexer);
  const parser = new PlSqlParser(tokenStream);

  try {
    tree = parser.sql_script();
  } catch (e) {
    console.error("error:", e);
  } finally {
    console.timeEnd("parse");
  }
  // console.log("tree:", tree);

  const completionTokenIndex = findCursorTokenIndex(tokenStream, {
    line: 1,
    column: caretIndex,
  });
  console.log("completionTokenIndex:", completionTokenIndex);

  if (completionTokenIndex === undefined) {
    console.log("no completion token found");
    return;
  }

  const suggestion = completion(tree, completionTokenIndex);
  console.log("suggestion:", JSON.stringify(suggestion, null, "  "));

  return suggestion;
};
