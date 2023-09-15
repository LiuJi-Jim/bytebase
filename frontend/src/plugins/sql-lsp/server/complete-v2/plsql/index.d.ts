export declare type FromClause = {
  alias: string;
  table: string;
  schema: string;
};

export declare type Suggestion = {
  field: string;
  table?: string;
  schema?: string;
};

export declare function parse(input: string): unknown;

export declare function complete(
  inputString: string,
  caretIndex: number
): {
  suggestion: Suggestion[];
  fromClause: FromClause[];
};
