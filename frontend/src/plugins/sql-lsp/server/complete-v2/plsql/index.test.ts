import { expect, it } from "vitest";
import { parse, complete } from "./index";

const CARET_STR = "<|>";

const prepareTestCase = (input: string) => {
  const offset = Math.max(input.indexOf(CARET_STR), 0);
  const content = input.replace(CARET_STR, "");
  return { content, offset };
};

it("Hello World", async () => {
  const tree = parse("SELECT 'HELLO WORLD';");
  expect(tree).toBeDefined();
});

it("First Blood", async () => {
  const { content, offset } = prepareTestCase(
    "select u.name<|> from sys.user u"
  );
  const result = complete(content, offset);

  expect(result?.suggestion.length).toBe(1);
  const sug = result?.suggestion[0];
  expect(sug.field).toBe("name");
  expect(sug.table).toBe("user");
  expect(sug.schema).toBe("sys");

  expect(result?.fromClause.length).toBe(1);
});
