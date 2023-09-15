import antlr4 from 'antlr4';

export default class PlSqlParserBase extends antlr4.Parser {
  _isVersion10: boolean;
  _isVersion12: boolean;

  constructor(input: antlr4.CommonTokenStream) {
    super(input);
    this._isVersion10 = false;
    this._isVersion12 = true;
  }

  isVersion10() {
    return this._isVersion10;
  }

  isVersion12() {
    return this._isVersion12;
  }

  setVersion10(value: boolean) {
    this._isVersion10 = value;
  }

  setVersion12(value: boolean) {
    this._isVersion12 = value;
  }
}
